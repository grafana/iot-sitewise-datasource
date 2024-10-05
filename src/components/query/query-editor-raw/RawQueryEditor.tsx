import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from 'types';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export const firstLabelWith = 20;

export function RawQueryEditor(props: Props) {
  const query = props.query;
  const updateQueryExpression = (queryText: string) => {
    query.expression = queryText;
  };

  const onChange = (query: SitewiseQuery) => {
    console.log(query.expression);
  };

  const onSqlChange = (text: string) => {
    console.log(text);
    updateQueryExpression(text);
  };

  return (
    <CodeEditor
      aria-label="SQL"
      language="sql"
      value={query.expression || ''}
      onSave={onSqlChange}
      onBlur={(text) => onChange({ ...query, expression: text })}
      height={'45vh'}
    />
  );
}
