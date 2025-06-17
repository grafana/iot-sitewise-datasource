import React from 'react';
import { Space } from '@grafana/ui';
import { InlineSelect, EditorRows, QueryEditorMode } from '@grafana/plugin-ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { SitewiseQuery, SitewiseOptions } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { RawQueryEditor } from 'components/query/query-editor-raw/RawQueryEditor';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';
import { regionOptions, type Region } from 'regions';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange, onRunQuery } = props;

  // Hardcoded to Builder mode til code is ready
  const editorMode = query.editorMode || QueryEditorMode.Builder;

  // Uncomment the following code when Builder mode is ready
  // const onEditorModeChange = (newEditorMode: QueryEditorMode) => {
  //   const newQuery = { ...query };
  //   if (newEditorMode === QueryEditorMode.Code) {
  //     newQuery.queryType = QueryType.ExecuteQuery;
  //     newQuery.clientCache = false;
  //     newQuery.rawSQL = newQuery.rawSQL || props.datasource.defaultQuery;
  //   }
  //   onChange({ ...newQuery, editorMode: newEditorMode });
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
        // extraHeaderElementRight={<QueryEditorModeToggle mode={editorMode!} onChange={onEditorModeChange} />}
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
          <RawQueryEditor
            {...props}
            datasource={props.datasource}
            query={query}
            onChange={onChange}
            onRunQuery={onRunQuery}
          />
        )}
        {editorMode === QueryEditorMode.Builder && (
          <VisualQueryBuilder datasource={props.datasource} query={query} onChange={onChange} onRunQuery={onRunQuery} />
        )}
      </EditorRows>
    </>
  );
}
