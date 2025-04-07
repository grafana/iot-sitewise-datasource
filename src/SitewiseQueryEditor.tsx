import React from 'react';
import { Space } from '@grafana/ui';
import { InlineSelect } from '@grafana/plugin-ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { EditorRows, QueryEditorMode, QueryEditorModeToggle } from '@grafana/plugin-ui';
import { SitewiseQuery, SitewiseOptions, QueryType } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { RawQueryEditor } from 'components/query/query-editor-raw/RawQueryEditor';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';
import { regionOptions, type Region } from 'regions';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange, onRunQuery } = props;
  const editorMode = query.editorMode || QueryEditorMode.Builder;

  const onEditorModeChange = (newEditorMode: QueryEditorMode) => {
    const newQuery = { ...query };
    if (newEditorMode === QueryEditorMode.Code) {
      newQuery.queryType = QueryType.ExecuteQuery;
      newQuery.clientCache = false;
      newQuery.rawSQL = newQuery.rawSQL || props.datasource.defaultQuery;
    }
    onChange({ ...newQuery, editorMode: newEditorMode });
  };

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
        extraHeaderElementRight={<QueryEditorModeToggle mode={editorMode!} onChange={onEditorModeChange} />}
        extraHeaderElementLeft={
          editorMode == QueryEditorMode.Code ? (
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
