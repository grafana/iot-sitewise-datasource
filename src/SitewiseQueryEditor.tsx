import { isEqual } from 'lodash';
import React, { useCallback, useState } from 'react';
import { ConfirmModal, Space } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { reportInteraction } from '@grafana/runtime';
import { EditorRows, QueryEditorMode, QueryEditorModeToggle } from '@grafana/experimental';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { RawQueryEditor } from 'components/query/query-editor-raw/RawQueryEditor';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange } = props;
  const editorMode = query.editorMode || QueryEditorMode.Builder;

  const onEditorModeChange = useCallback(
    (newEditorMode: QueryEditorMode) => {
      reportInteraction('grafana_sitewise_editor_mode_clicked', {
        newEditor: newEditorMode,
        previousEditor: query.editorMode ?? '',
        newQuery: !query.expression,
      });

      // if (newEditorMode === QueryEditorMode.Builder) {
      //     const result = buildVisualQueryFromString(query.expression || '');
      //     // If there are errors, give user a chance to decide if they want to go to builder as that can lose some data.
      //     if (result.errors.length) {
      //         setParseModalOpen(true);
      //         return;
      //     }
      // }
      changeEditorMode(query, newEditorMode, onChange);
    },
    [onChange, query] // , app
  );

  const onChangeInternal = (query: SitewiseQuery) => {
    if (!isEqual(query, props.query)) {
      // setDataIsStale(true);
    }
    onChange(query);
  };

  const [parseModalOpen, setParseModalOpen] = useState(false);

  const queryEditorModeDefaultLocalStorageKey = 'SitewiseQueryEditorModeDefault';

  function changeEditorMode(
    query: SitewiseQuery,
    editorMode: QueryEditorMode,
    onChange: (query: SitewiseQuery) => void
  ) {
    // If empty query store new mode as default
    if (query.expression === '') {
      window.localStorage.setItem(queryEditorModeDefaultLocalStorageKey, editorMode);
    }

    onChange({ ...query, editorMode });
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
        showAsyncQueryButtons
        extraHeaderElementRight={<QueryEditorModeToggle mode={editorMode!} onChange={onEditorModeChange} />}
        // cancel={props.datasource.cancel} TODO: Implement cancel
      />
      <Space v={0.5} />
      <EditorRows>
        {editorMode === QueryEditorMode.Code && (
          <RawQueryEditor
            {...props}
            datasource={props.datasource}
            query={query}
            onChange={onChangeInternal}
            onRunQuery={props.onRunQuery}
          /> // showExplain={explain} />
        )}
        {editorMode === QueryEditorMode.Builder && (
          <VisualQueryBuilder
            datasource={props.datasource}
            query={query}
            onChange={onChangeInternal}
            onRunQuery={props.onRunQuery}
            // showExplain={explain}
            // timeRange={timeRange}
          />
        )}
      </EditorRows>
    </>
  );
}
