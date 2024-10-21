import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions, SqlQuery } from 'types';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export const firstLabelWith = 20;

export function RawQueryEditor(props: Props) {
  const query = props.query as SqlQuery;

  const onChange = (query: SqlQuery) => {
    console.log(query.queryStatement);
    console.log(props);
    props.onChange(query);
    props.onRunQuery();
  };

  return (
    <CodeEditor
      aria-label="SQL"
      language="sql"
      value={query.queryStatement || ''}
      onSave={(text) => onChange({ ...query, queryStatement: text })}
      onBlur={(text) => onChange({ ...query, queryStatement: text })}
      height={'45vh'}
    />
  );
}
