import { DataSourcePlugin } from '@grafana/data';
import { SitewiseDatasource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { SitewiseQuery, SitewiseDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<SitewiseDatasource, SitewiseQuery, SitewiseDataSourceOptions>(
  SitewiseDatasource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
