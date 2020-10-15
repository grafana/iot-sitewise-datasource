import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { SitewiseQuery, SitewiseOptions } from './types';

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
}
