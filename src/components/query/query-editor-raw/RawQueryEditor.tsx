import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { SitewiseCompletionProvider } from 'language/autoComplete';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function RawQueryEditor(props: Props) {
  const query = props.query as SitewiseQuery;
  const defaultQuery = 'select $__selectAll from raw_time_series where $__unixEpochFilter(event_timestamp)';

  query.clientCache = false;
  query.rawSQL = query.rawSQL || defaultQuery;

  const onChange = (query: SitewiseQuery) => {
    props.onChange(query);
  };

  return (
    <CodeEditor
      aria-label="SQL"
      language="sql"
      monacoOptions={{ automaticLayout: true, minimap: { enabled: false } }}
      value={query.rawSQL}
      onSave={(text) => onChange({ ...query, rawSQL: text })}
      onBlur={(text) => onChange({ ...query, rawSQL: text })}
      onBeforeEditorMount={(monaco) => {
        if (SitewiseCompletionProvider.monaco === null) {
          SitewiseCompletionProvider.monaco = monaco;
          monaco.languages.registerCompletionItemProvider('sql', SitewiseCompletionProvider);
        }
      }}
      height={'45vh'}
    />
  );
}
