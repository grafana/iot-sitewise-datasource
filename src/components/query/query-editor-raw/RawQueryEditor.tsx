import React from 'react';
import { CodeEditor } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { SitewiseCompletionProvider } from 'language/autoComplete';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function RawQueryEditor(props: Props) {
  const { onChange, query } = props;

  return (
    <CodeEditor
      language="sql"
      showLineNumbers
      showMiniMap={false}
      value={query.rawSQL || props.datasource.defaultQuery}
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
