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
  // returns string that will be shown in the panel header when the panel is collapsed
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
    } else if (query.propertyAlias) {
      txt += ' / ' + query.propertyAlias;
    }
    return txt;
  }

  /**
   * Supports template variables for region, asset and property
   */
  applyTemplateVariables(query: SitewiseQuery, scopedVars: ScopedVars): SitewiseQuery {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      propertyAlias: templateSrv.replace(query.propertyAlias, scopedVars),
      region: templateSrv.replace(query.region || '', scopedVars),
      propertyId: templateSrv.replace(query.propertyId || '', scopedVars),
      assetIds: query.assetIds?.flatMap((assetId) => templateSrv.replace(assetId, scopedVars, 'csv').split(',')) ?? [],
    };
  }

  runQuery(query: SitewiseQuery, maxDataPoints?: number): Observable<DataQueryResponse> {
    // @ts-ignore
    return this.query({ targets: [query], requestId: `iot.${counter++}`, maxDataPoints });
  }

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    /*
    I think we want to do something like this: (from twinmaker)
    if (this.grafanaLiveEnabled) {
      return super.query(options);
    } 
    or possibly something like this: (from grafana-plugin-examples)
      const observables = request.targets.map((target, index) => {
      const query = defaults(target, DEFAULT_QUERY);
      this.grafanaLiveEnabled = true;

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: `my-ws/custom-${query.lowerLimit}-${query.upperLimit}`, // this will allow each new query to create a new connection
          data: {
            ...query,
          },
        },
      });
    });

    return merge(...observables);


    In general I remain unclear about the relationship between grafana live and a streaming toggle

    Another TODO: remove streaming toggle from non-stream-able query types like list assets, otherwise they just repeat the same data over and over again
    */
    return super.query(request);
    // return getRequestLooper(request, {
    //   // Check for a "nextToken" in the response
    //   getNextQueries: (rsp: DataQueryResponse) => {
    //     if (rsp.data?.length) {
    //       const next: SitewiseNextQuery[] = [];
    //       for (const frame of rsp.data as DataFrame[]) {
    //         const meta = frame.meta?.custom as SitewiseCustomMeta;
    //         if (meta && meta.nextToken) {
    //           const query = request.targets.find((t) => t.refId === frame.refId);
    //           if (query) {
    //             const existingNextQuery = next.find((v) => v.refId === frame.refId);
    //             if (existingNextQuery) {
    //               if (existingNextQuery.nextToken !== meta.nextToken && meta.entryId && meta.nextToken) {
    //                 if (!existingNextQuery.nextTokens) {
    //                   existingNextQuery.nextTokens = {};
    //                 }
    //                 existingNextQuery.nextTokens[meta.entryId] = meta.nextToken;
    //               }
    //             } else {
    //               next.push({
    //                 ...query,
    //                 nextToken: meta.nextToken,
    //                 nextTokens: { ...(meta.entryId && meta.nextToken ? { [meta.entryId]: meta.nextToken } : {}) },
    //               });
    //             }
    //           }
    //         }
    //       }
    //       if (next.length) {
    //         return next;
    //       }
    //     }
    //     return undefined;
    //   },
    //   /**
    //    * The original request
    //    */
    //   query: (request: DataQueryRequest<SitewiseQuery>) => {
    //     return super.query(request);
    //   },
    //   /**
    //    * Process the results
    //    */
    //   process: (t: MultiRequestTracker, data: DataFrame[], isLast: boolean) => {
    //     if (t.data) {
    //       // append rows to fields with the same structure
    //       t.data = appendMatchingFrames(t.data, data);
    //     } else {
    //       t.data = data; // hang on to the results from the last query
    //     }
    //     return t.data;
    //   },
    //   /**
    //    * Callback that gets executed when unsubscribed
    //    */
    //   onCancel: (tracker: MultiRequestTracker) => {},
    // });
  }
}

let counter = 1000;
