import { DataSourceInstanceSettings, ScopedVars, DataQueryResponse } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { SitewiseCache } from 'sitewiseCache';

import { SitewiseQuery, SitewiseOptions } from './types';
import { Observable } from 'rxjs';

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
  getCache(region?: string) {
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

  runQuery(query: SitewiseQuery): Observable<DataQueryResponse> {
    // @ts-ignore
    return this.query({ targets: [query], requestId: `iot.${counter++}` });
  }
}

let counter = 1000;
