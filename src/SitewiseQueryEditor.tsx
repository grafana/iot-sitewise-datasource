import React from 'react';
import { InlineField, Select, Space } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { QueryEditorHeader } from '@grafana/aws-sdk';
import { EditorRows, QueryEditorMode, QueryEditorModeToggle } from '@grafana/plugin-ui';
import { SitewiseQuery, SitewiseOptions, QueryType } from 'types';
import { DataSource } from 'SitewiseDataSource';
import { RawQueryEditor } from 'components/query/query-editor-raw/RawQueryEditor';
import { VisualQueryBuilder } from 'components/query/visual-query-builder/VisualQueryBuilder';
import { standardRegionOptions } from 'regions';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export function SitewiseQueryEditor(props: Props) {
  const { query, onChange, onRunQuery } = props;
  const editorMode = query.editorMode || QueryEditorMode.Builder;

  const onEditorModeChange = (newEditorMode: QueryEditorMode) => {
    const newQuery = { ...query };
    if (newEditorMode === QueryEditorMode.Code) {
      newQuery.queryType = QueryType.ExecuteQuery;
      newQuery.clientCache = false;
    }
    onChange({ ...newQuery, editorMode: newEditorMode });
  };

  const onRegionChange = (sel: SelectableValue<string>) => {
    onChange({ ...query, region: sel.value });
  };

  return (
    <>
      <QueryEditorHeader<DataSource, SitewiseQuery, SitewiseOptions>
        {...props}
        enableRunButton
        extraHeaderElementRight={<QueryEditorModeToggle mode={editorMode!} onChange={onEditorModeChange} />}
        extraHeaderElementLeft={
          editorMode == QueryEditorMode.Code ? (
            <InlineField label="AWS Region">
              <Select
                options={standardRegionOptions}
                value={
                  standardRegionOptions.find((v) => v.value === query.region) || props.datasource.options.defaultRegion
                }
                onChange={onRegionChange}
                backspaceRemovesValue
                allowCustomValue
                isClearable
                menuPlacement="auto"
              />
            </InlineField>
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
