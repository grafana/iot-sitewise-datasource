import { DataSource } from 'DataSource';
import { SitewiseQuery } from 'types';

export interface SitewiseQueryEditorProps<TQuery extends SitewiseQuery = SitewiseQuery> {
  datasource: DataSource;
  query: TQuery;
  onChange: (value: TQuery) => void;
}
