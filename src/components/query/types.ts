import { DataSource } from 'DataSource';
import { SitewiseQuery } from 'types';

export interface SitewiseQueryEditorProps<TQuery extends SitewiseQuery = SitewiseQuery> {
  datasource: DataSource;
  query: TQuery;
  onRunQuery: () => void;
  onChange: (value: TQuery) => void;
}
