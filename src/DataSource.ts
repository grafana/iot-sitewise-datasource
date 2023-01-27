import {
  DataSourceInstanceSettings,
  ScopedVars,
  DataQueryResponse,
  DataQueryRequest,
  DataFrame,
  MetricFindValue,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { SitewiseCache } from 'sitewiseCache';

import { SitewiseQuery, SitewiseOptions, SitewiseCustomMeta, isPropertyQueryType, SitewiseNextQuery } from './types';
import { Observable } from 'rxjs';
import { getRequestLooper, MultiRequestTracker } from 'requestLooper';
import { appendMatchingFrames } from 'appendFrames';
import { frameToMetricFindValues } from 'utils';

export class DataSource extends DataSourceWithBackend<SitewiseQuery, SitewiseOptions> {
  // Easy access for QueryEditor
  readonly options: SitewiseOptions;
  private cache = new Map<string, SitewiseCache>();

  constructor(instanceSettings: DataSourceInstanceSettings<SitewiseOptions>) {
    super(instanceSettings);
    this.options = instanceSettings.jsonData;
  }

  /**
   * Get a region scoped cache
   */
  getCache(region?: string): SitewiseCache {
    if (!region || region === 'default') {
      region = this.options.defaultRegion || '';
    }
    let v = this.cache.get(region);
    if (!v) {
      v = new SitewiseCache(this, region);
      this.cache.set(region, v);
    }
    return v;
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  async metricFindQuery(query: SitewiseQuery, options: any): Promise<MetricFindValue[]> {
    const request = {
      targets: [
        {
          ...query,
          refId: 'metricFindQuery',
        },
      ],
      range: options.range,
      rangeRaw: options.rangeRaw,
    } as DataQueryRequest<SitewiseQuery>;

    let res: DataQueryResponse;

    try {
      res = await this.query(request).toPromise();
    } catch (err) {
      return Promise.reject(err);
    }

    if (!res || !res.data || res.data.length <= 0) {
      return [];
    }
    return frameToMetricFindValues(res.data[0] as DataFrame);
  }

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: SitewiseQuery): boolean {
    if (!query.queryType) {
      return false; // skip the query
    }
    // Migrate assetId to assetIDs (v1.6)
    if (query.assetId) {
      const ids = new Set<string>();
      ids.add(query.assetId);
      if (query.assetIds) {
        for (const id of query.assetIds) {
          ids.add(id);
        }
      }
      query.assetIds = Array.from(ids);
      delete query.assetId;
    }

    if (isPropertyQueryType(query.queryType)) {
      return Boolean((query.assetIds?.length && query.propertyId) || query.propertyAlias);
    }
    return true; // keep the query
  }

  getQueryDisplayText(query: SitewiseQuery): string {
    const cache = this.getCache(query.region);
    let txt: string = query.queryType;
    if (query.assetIds?.length) {
      const info = cache.getAssetInfoSync(query.assetIds[0]);
      if (!info) {
        return txt + ' / ' + query.assetIds.join('/');
      }
      txt += ' / ' + info.name;

      if (query.propertyId && info.properties) {
        const p = info.properties.find((v) => v.Id === query.propertyId);
        if (p) {
          txt += ' / ' + p.Name;
        } else {
          txt += ' / ' + query.propertyId;
        }
      }
    }
    return txt;
  }

  /**
   * Supports template variables for region, asset and property
   */
  applyTemplateVariables(query: SitewiseQuery, scopedVars: ScopedVars): SitewiseQuery {
    const templateSrv = getTemplateSrv();
    let assetIds: string[] = [];
    if (query.assetIds) {
      for (const id of query.assetIds) {
        const out = templateSrv.replace(id, scopedVars, 'json');
        try {
          const parsed: string | string[] = JSON.parse(out);
          if (Array.isArray(parsed)) {
            assetIds = assetIds.concat(parsed);
          } else {
            assetIds.push(out);
          }
        } catch (err) {
          assetIds.push(id);
        }
      }
    }
    return {
      ...query,
      propertyAlias: templateSrv.replace(query.propertyAlias, scopedVars),
      region: templateSrv.replace(query.region || '', scopedVars),
      propertyId: templateSrv.replace(query.propertyId || '', scopedVars),
      assetIds,
    };
    return query;
  }

  runQuery(query: SitewiseQuery, maxDataPoints?: number): Observable<DataQueryResponse> {
    // @ts-ignore
    return this.query({ targets: [query], requestId: `iot.${counter++}`, maxDataPoints });
  }

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    return getRequestLooper(request, {
      // Check for a "nextToken" in the response
      getNextQueries: (rsp: DataQueryResponse) => {
        if (rsp.data?.length) {
          const next: SitewiseNextQuery[] = [];
          for (const frame of rsp.data as DataFrame[]) {
            const meta = frame.meta?.custom as SitewiseCustomMeta;
            if (meta && meta.nextToken) {
              const query = request.targets.find((t) => t.refId === frame.refId);
              if (query) {
                next.push({
                  ...query,
                  nextToken: meta.nextToken,
                });
              }
            }
          }
          if (next.length) {
            return next;
          }
        }
        return undefined;
      },

      /**
       * The original request
       */
      query: (request: DataQueryRequest<SitewiseQuery>) => {
        return super.query(request);
      },

      /**
       * Process the results
       */
      process: (t: MultiRequestTracker, data: DataFrame[], isLast: boolean) => {
        if (t.data) {
          // append rows to fields with the same structure
          t.data = appendMatchingFrames(t.data, data);
        } else {
          t.data = data; // hang on to the results from the last query
        }
        return t.data;
      },

      /**
       * Callback that gets executed when unsubscribed
       */
      onCancel: (tracker: MultiRequestTracker) => {},
    });
  }
}

let counter = 1000;
