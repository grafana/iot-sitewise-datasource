import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataFrame,
  DataQueryRequest,
  ScopedVars,
  QueryResultMetaStat,
} from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable } from 'rxjs';

import { SitewiseQuery, SitewiseOptions, SitewiseCustomMeta } from './types';
import { getRequestLooper, MultiRequestTracker } from 'requestLooper';

export class DataSource extends DataSourceWithBackend<SitewiseQuery, SitewiseOptions> {
  // Easy access for QueryEditor
  options: SitewiseOptions;

  constructor(instanceSettings: DataSourceInstanceSettings<SitewiseOptions>) {
    super(instanceSettings);
    this.options = instanceSettings.jsonData;
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: SitewiseQuery): boolean {
    return !!query.rawQuery;
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

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    let tracker: SitewiseCustomMeta | undefined = undefined;
    let queryId: string | undefined = undefined;
    return getRequestLooper(request, {
      // Check for a "nextToken" in the response
      getNextQuery: (rsp: DataQueryResponse) => {
        if (rsp.data?.length) {
          const first = rsp.data[0] as DataFrame;
          const meta = first.meta?.custom as SitewiseCustomMeta;
          if (meta && meta.nextToken) {
            queryId = meta.queryId;

            return {
              refId: first.refId,
              rawQuery: first.meta?.executedQueryString,
              nextToken: meta.nextToken,
            } as SitewiseQuery;
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
        const meta = data[0]?.meta?.custom as SitewiseCustomMeta;
        if (!meta) {
          return data; // NOOP
        }
        // Single request
        meta.fetchStartTime = t.fetchStartTime;
        meta.fetchEndTime = t.fetchEndTime;
        meta.fetchTime = t.fetchEndTime! - t.fetchStartTime!;

        if (tracker) {
          // Additional request
          if (!tracker.subs?.length) {
            const { subs, nextToken, queryId, ...rest } = tracker;
            (rest as any).requestNumber = 1;
            tracker.subs?.push(rest as SitewiseCustomMeta);
          }
          for (const m of tracker.subs!) {
            delete m.nextToken; // not useful in the
          }
          delete (meta as any).queryId;
          (meta as any).requestNumber = tracker.subs!.length + 1;

          tracker.subs!.push(meta);
          tracker.fetchEndTime = t.fetchEndTime;
          tracker.fetchTime = t.fetchEndTime! - tracker.fetchStartTime!;
          tracker.executionFinishTime = meta.executionFinishTime;

          data[0].meta!.custom = tracker;
        } else {
          // First request
          tracker = {
            ...t,
            ...meta,
            subs: [],
          } as SitewiseCustomMeta;
        }

        // Calculate stats
        if (isLast && tracker.executionStartTime && tracker.executionFinishTime) {
          delete tracker.nextToken;

          const tsTime = tracker.executionFinishTime - tracker.executionStartTime;
          if (tsTime > 0) {
            const stats: QueryResultMetaStat[] = [];
            if (tracker.subs && tracker.subs.length) {
              stats.push({
                displayName: 'HTTP request count',
                value: tracker.subs.length,
                unit: 'none',
              });
            }
            stats.push({
              displayName: 'Execution time (Grafana server ⇆ Timestream)',
              value: tsTime,
              unit: 'ms',
              decimals: 2,
            });
            if (tracker.fetchStartTime) {
              tracker.fetchEndTime = Date.now();
              const dsTime = tracker.fetchEndTime - tracker.fetchStartTime;
              tracker.fetchTime = dsTime - tsTime;
              if (dsTime > tsTime) {
                stats.push({
                  displayName: 'Fetch time (Browser ⇆ Grafana server w/o Timestream)',
                  value: tracker.fetchTime,
                  unit: 'ms',
                  decimals: 2,
                });
                stats.push({
                  displayName: 'Fetch overhead',
                  value: (tracker.fetchTime / dsTime) * 100,
                  unit: 'percent', // 0 - 100
                });
              }
            }
            data[0].meta!.stats = stats;
          }
        }
        return data;
      },

      /**
       * Callback that gets executed when unsubscribed
       */
      onCancel: (tracker: MultiRequestTracker) => {
        if (queryId) {
          console.log('Cancelling running timestream query');

          // tracker.killed = true;
          this.postResource(`cancel`, {
            queryId,
          })
            .then(v => {
              console.log('Timestream query Canceled:', v);
            })
            .catch(err => {
              err.isHandled = true; // avoid the popup
              console.log('error killing', err);
            });
        }
      },
    });
  }
}

export function getNextTokenMeta(rsp: DataQueryResponse): SitewiseCustomMeta | undefined {
  if (rsp.data?.length) {
    const first = rsp.data[0] as DataFrame;
    const meta = first.meta?.custom as SitewiseCustomMeta;
    if (meta && meta.nextToken) {
      return meta;
    }
  }
  return undefined;
}
