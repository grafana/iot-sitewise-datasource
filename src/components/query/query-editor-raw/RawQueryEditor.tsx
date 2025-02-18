import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { SitewiseCompletionProvider } from 'language/autoComplete';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function RawQueryEditor(props: Props) {
  const defaultQuery = 'select $__selectAll from raw_time_series where $__unixEpochFilter(event_timestamp)';
  const { onChange, query } = props;
  query.rawSQL = query.rawSQL || defaultQuery;

  return (
    <CodeEditor
      language="sql"
      showLineNumbers
      showMiniMap={false}
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
      height={'200px'}
    />
  );
}
