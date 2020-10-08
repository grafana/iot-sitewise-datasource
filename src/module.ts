import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { QueryEditor, ConfigEditor } from './components';
import { SitewiseQuery, SitewiseOptions } from './types';
import { MetaInspector } from 'components/MetaInspector';

export const plugin = new DataSourcePlugin<DataSource, SitewiseQuery, SitewiseOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setMetadataInspector(MetaInspector)
  .setQueryEditor(QueryEditor);
