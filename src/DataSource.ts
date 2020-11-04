import { DataSourceInstanceSettings, ScopedVars, DataQueryResponse, DataQueryRequest, DataFrame } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { SitewiseCache } from 'sitewiseCache';

import { SitewiseQuery, SitewiseOptions, SitewiseCustomMeta } from './types';
import { Observable } from 'rxjs';
import { getRequestLooper, MultiRequestTracker } from 'requestLooper';
import { appendMatchingFrames } from 'appendFrames';

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

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: SitewiseQuery): boolean {
    return !!query.queryType;
  }

  getQueryDisplayText(query: SitewiseQuery): string {
    return 'TODO: ' + JSON.stringify(query);
  }

  applyTemplateVariables(query: SitewiseQuery, scopedVars: ScopedVars): SitewiseQuery {
    // if (!query.rawQuery) {
    //   return query;
    // }

    // const templateSrv = getTemplateSrv();
    // return {
    //   ...query,
    //   database: templateSrv.replace(query.database || '', scopedVars),
    //   table: templateSrv.replace(query.table || '', scopedVars),
    //   measure: templateSrv.replace(query.measure || '', scopedVars),
    //   rawQuery: templateSrv.replace(query.rawQuery), // DO NOT include scopedVars! it uses $__interval_ms!!!!!
    // };
    return query;
  }

  runQuery(query: SitewiseQuery, maxDataPoints?: number): Observable<DataQueryResponse> {
    // @ts-ignore
    return this.query({ targets: [query], requestId: `iot.${counter++}`, maxDataPoints });
  }

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    return getRequestLooper(request, {
      // Check for a "nextToken" in the response
      // NOTE:  this well behave weird if multiple requests are in one query!
      getNextQuery: (rsp: DataQueryResponse) => {
        if (rsp.data?.length) {
          const first = rsp.data[0] as DataFrame;
          const meta = first.meta?.custom as SitewiseCustomMeta;
          if (meta && meta.nextToken) {
            const query = request.targets.find(t => t.refId === first.refId);
            if (query) {
              return {
                ...query,
                nextToken: meta.nextToken,
              };
            }
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
