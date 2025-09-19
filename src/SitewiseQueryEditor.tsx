import React, { useState, useCallback } from 'react';
import { Space, CodeEditor } from '@grafana/ui';
import { EditorRows, QueryEditorMode, InlineSelect } from '@grafana/plugin-ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';
import { SqlQueryBuilder } from 'components/query/sql-query-builder/SqlQueryBuilder';
import { defaultSitewiseQueryState, SitewiseQueryState } from 'components/query/sql-query-builder/types';
import { regionOptions, type Region } from 'regions';
import { SitewiseCompletionProvider } from 'language/autoComplete';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;
// Uncomment the following code when Builder mode is ready
// const editorModeOptions: Array<SelectableValue<QueryEditorMode | 'sql'>> = [
//   { label: 'Builder', value: QueryEditorMode.Builder },
//   { label: 'SQL Builder', value: 'sql' }, // custom option
//   { label: 'Code', value: QueryEditorMode.Code },
// ];

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange, onRunQuery, datasource } = props;
  // Uncomment the following code when Builder mode is ready
  // const [showConfirmation, setShowConfirmation] = useState(false);
  // const [pendingMode, setPendingMode] = useState<SelectableValue<QueryEditorMode | 'sql'> | null>(null);
  const [builderState, setBuilderState] = useState(query.sqlQueryState || defaultSitewiseQueryState);

  // Hardcoded to Builder mode til code is ready
  // Add the setEditorMode when Builder mode is ready
  const [editorMode] = useState<QueryEditorMode | 'sql'>(query.editorMode || QueryEditorMode.Builder);

  const handleQueryChange = useCallback(
    (updatedState: SitewiseQueryState) => {
      setBuilderState(updatedState);
      onChange({
        ...query,
        rawSQL: updatedState.rawSQL,
        sqlQueryState: updatedState,
      });
    },
    [query, onChange]
  );
  // Uncomment the following code when Builder mode is ready
  // const onEditorModeChange = (sel: SelectableValue<QueryEditorMode | 'sql'>, skipConfirmation = false) => {
  //   const newEditorMode = sel.value;
  //   console.log(query)
  //   if (!newEditorMode) {
  //     return;
  //   }
  //   if (!skipConfirmation && editorMode === QueryEditorMode.Code && newEditorMode === 'sql') {
  //     setPendingMode(sel);
  //     setShowConfirmation(true);
  //     return;
  //   }
  //   const newQuery = { ...query };
  //   if (newEditorMode === QueryEditorMode.Code || newEditorMode === 'sql') {
  //     newQuery.queryType = QueryType.ExecuteQuery;
  //     newQuery.clientCache = false;
  //     newQuery.rawSQL = newQuery.rawSQL || datasource.defaultQuery;
  //   }
  //   setEditorMode(newEditorMode);
  //   onChange({
  //     ...newQuery,
  //     editorMode: newEditorMode as QueryEditorMode
  //   });
  // };

  // const handleConfirmModeChange = () => {
  //   if (pendingMode) {
  //     query.rawSQL = builderState.rawSQL;
  //     onEditorModeChange(pendingMode, true);
  //   }
  //   setShowConfirmation(false);
  //   setPendingMode(null);
  // };

  // const handleCancelModeChange = () => {
  //   setShowConfirmation(false);
  //   setPendingMode(null);
  // };

  const onRegionChange = (sel: SelectableValue<Region>) => {
    onChange({ ...query, region: sel.value });
  };

  const selectedRegionOption = regionOptions.find((option) => {
    return option.value === query.region;
  });

  return (
    <>
      <QueryEditorHeader<DataSource, SitewiseQuery, SitewiseOptions>
        {...props}
        enableRunButton
        // Uncomment the following code when Builder mode is ready
        // extraHeaderElementRight={
        //   <InlineSelect
        //     label="Mode"
        //     options={editorModeOptions}
        //     value={editorModeOptions.find((opt) => opt.value === editorMode) || editorModeOptions[0]}
        //     onChange={(sel) => onEditorModeChange(sel)}
        //     menuPlacement="auto"
        //     isSearchable={false}
        //     width={14}
        //   />
        // }
        extraHeaderElementLeft={
          editorMode === QueryEditorMode.Code ? (
            <InlineSelect
              label="AWS Region"
              options={regionOptions}
              value={selectedRegionOption}
              onChange={onRegionChange}
              backspaceRemovesValue
              allowCustomValue
              isClearable
              menuPlacement="auto"
            />
          ) : undefined
        }
      />
      <Space v={0.5} />
      <EditorRows>
        {editorMode === QueryEditorMode.Code && (
          <CodeEditor
            language="sql"
            showLineNumbers
            showMiniMap={false}
            value={query.rawSQL || datasource.defaultQuery}
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
        )}
        {editorMode === 'sql' && <SqlQueryBuilder builderState={builderState} onChange={handleQueryChange} />}
        {editorMode === QueryEditorMode.Builder && (
          <VisualQueryBuilder datasource={datasource} query={query} onChange={onChange} onRunQuery={onRunQuery} />
        )}
      </EditorRows>

      {/* Confirmation Dialog */}
      {/* Uncomment the following code when Builder mode is ready */}
      {/* <ConfirmDialog
      isOpen={showConfirmation}
      title="Switch to SQL Builder"
      message="Are you sure to switch to sql builder mode? You will lose the changes done in code editor mode."
      onConfirm={handleConfirmModeChange}
      onCancel={handleCancelModeChange}
    /> */}
    </>
  );
}
