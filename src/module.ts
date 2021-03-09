import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { SitewiseQuery, SitewiseOptions } from './types';
import { MetaInspector } from 'components/MetaInspector';
import { ConfigEditor } from 'components/ConfigEditor';
import { QueryEditor } from 'components/query/QueryEditor';

export const plugin = new DataSourcePlugin<DataSource, SitewiseQuery, SitewiseOptions>(DataSource)
  .setConfigEditor(ConfigEditor as any) // HACK since typename was added in 7.5 https://github.com/grafana/grafana/pull/31326/files#diff-c58fc1a09e9b9b17e5f45efbfb646273e69145f7687facb134440da4edafc745R552
  .setMetadataInspector(MetaInspector)
  .setQueryEditor(QueryEditor);
