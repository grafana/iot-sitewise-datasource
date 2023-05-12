import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { SitewiseQuery, SitewiseOptions } from './types';
import { MetaInspector } from 'components/MetaInspector';
import { ConfigEditor } from 'components/ConfigEditor';
import { QueryEditor } from 'components/query/QueryEditor';

export const plugin = new DataSourcePlugin<DataSource, SitewiseQuery, SitewiseOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setMetadataInspector(MetaInspector)
  .setQueryEditor(QueryEditor);

