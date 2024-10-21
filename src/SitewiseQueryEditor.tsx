import React, { useCallback, useState } from 'react';
import { ConfirmModal, Space } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { reportInteraction } from '@grafana/runtime';
import { EditorRows, QueryEditorMode, QueryEditorModeToggle } from '@grafana/experimental';
import { SitewiseQuery, SitewiseOptions, QueryType, SqlQuery } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { RawQueryEditor } from 'components/query/query-editor-raw/RawQueryEditor';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange, onRunQuery } = props;
  const editorMode = query.editorMode || QueryEditorMode.Builder;

  const onEditorModeChange = useCallback(
    (newEditorMode: QueryEditorMode) => {
      reportInteraction('grafana_sitewise_editor_mode_clicked', {
        newEditor: newEditorMode,
        previousEditor: query.editorMode ?? '',
      });

      if (newEditorMode === QueryEditorMode.Code) {
        query.queryType = QueryType.ExecuteQuery;
      }
      changeEditorMode(query, newEditorMode, onChange);
    },
    [onChange, query]
  );

  const onChangeInternal = (query: SitewiseQuery) => {
    onChange(query);
  };

  const [parseModalOpen, setParseModalOpen] = useState(false);

  function changeEditorMode(
    query: SitewiseQuery,
    editorMode: QueryEditorMode,
    onChange: (query: SitewiseQuery) => void
  ) {
    onChange({ ...query, editorMode });
  }

  function onRunQueryInternal() {
    onRunQuery();
  }

  return (
    <>
      <ConfirmModal
        isOpen={parseModalOpen}
        title="Query parsing"
        body="There were errors while trying to parse the query. Continuing to visual builder may lose some parts of the query."
        confirmText="Continue"
        onConfirm={() => {
          onChange({ ...query, editorMode: QueryEditorMode.Builder });
          setParseModalOpen(false);
        }}
        onDismiss={() => setParseModalOpen(false)}
      />
      <QueryEditorHeader<DataSource, SitewiseQuery, SitewiseOptions>
        {...props}
        enableRunButton
        extraHeaderElementRight={<QueryEditorModeToggle mode={editorMode!} onChange={onEditorModeChange} />}
      />
      <Space v={0.5} />
      <EditorRows>
        {editorMode === QueryEditorMode.Code && (
          <RawQueryEditor
            {...props}
            datasource={props.datasource}
            query={query as SqlQuery}
            onChange={onChangeInternal}
            onRunQuery={onRunQueryInternal}
          />
        )}
        {editorMode === QueryEditorMode.Builder && (
          <VisualQueryBuilder
            datasource={props.datasource}
            query={query}
            onChange={onChangeInternal}
            onRunQuery={onRunQueryInternal}
          />
        )}
      </EditorRows>
    </>
  );
}
