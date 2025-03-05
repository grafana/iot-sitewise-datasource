import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from './types';
import { MetadataInspector } from 'components/MetadataInspector';
import { ConfigEditor } from 'components/ConfigEditor';
import { SitewiseQueryEditor } from 'SitewiseQueryEditor';

export const plugin = new DataSourcePlugin<DataSource, SitewiseQuery, SitewiseOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setMetadataInspector(MetadataInspector)
  .setQueryEditor(SitewiseQueryEditor);
