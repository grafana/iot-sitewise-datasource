import defaults from 'lodash/defaults';
import React, { useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from 'DataSource';
import { SitewiseQuery, SitewiseOptions, QueryType } from 'types';
import { InlineField, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWiseQueryTypes, changeQueryType } from 'queryInfo';
import { standardRegionOptions } from 'regions';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/experimental';
import { config } from '@grafana/runtime';
import { QueryField } from './QueryField';
import { QueryToolTip } from './QueryToolTip';

export type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

const queryDefaults: Partial<SitewiseQuery> = {
  maxPageAggregations: 1,
};

export const firstLabelWith = 20;

export const QueryEditor = (props: Props) => {
  const { datasource, query: baseQuery, onRunQuery = () => { }, onChange = () => { } } = props;
  const newFormStylingEnabled = config.featureToggles.awsDatasourcesNewFormStyling;

  const query = defaults(baseQuery, queryDefaults);

  const defaultRegion: SelectableValue<string> = {
    label: `Default`,
    description: datasource.options?.defaultRegion,
    value: undefined,
  };

  const regions = query.region ? [defaultRegion, ...standardRegionOptions] : standardRegionOptions;
  const currentQueryType = siteWiseQueryTypes.find((v) => v.value === query.queryType);

  const onQueryTypeChange = useCallback((sel: SelectableValue<QueryType>) => {
    // hack to use QueryEditor as VariableQueryEditor
    onChange(changeQueryType(query, sel as QueryTypeInfo));
    onRunQuery();
  }, [onChange, onRunQuery, query]);

  const onRegionChange = useCallback((sel: SelectableValue<string>) => {
    onChange({ ...query, assetId: undefined, propertyId: undefined, region: sel.value });
    onRunQuery();
  }, [onChange, onRunQuery, query]);

  return (
    <>
      {newFormStylingEnabled ? (
        <>
          <EditorRows>
            <EditorRow>
              <EditorFieldGroup>
                <EditorField label="Query type" tooltip={currentQueryType && <QueryToolTip {...currentQueryType} />} tooltipInteractive width={30}>
                  <Select
                    options={siteWiseQueryTypes}
                    value={currentQueryType}
                    onChange={onQueryTypeChange}
                    placeholder="Select query type"
                    menuPlacement="bottom"
                  />
                </EditorField>
                <EditorField label="Region" width={15}>
                  <Select
                    options={regions}
                    value={standardRegionOptions.find((v) => v.value === query.region) || defaultRegion}
                    onChange={onRegionChange}
                    backspaceRemovesValue={true}
                    allowCustomValue={true}
                    isClearable={true}
                    menuPlacement="bottom"
                  />
                </EditorField>
              </EditorFieldGroup>
            </EditorRow>
            <QueryField {...props} query={query} />
          </EditorRows>
        </>
      ) : <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={firstLabelWith} grow={true} tooltip={currentQueryType && <QueryToolTip {...currentQueryType} />} interactive>
            <Select
              options={siteWiseQueryTypes}
              value={currentQueryType}
              onChange={onQueryTypeChange}
              placeholder="Select query type"
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Region" labelWidth={14}>
            <Select
              width={18}
              options={regions}
              value={standardRegionOptions.find((v) => v.value === query.region) || defaultRegion}
              onChange={onRegionChange}
              backspaceRemovesValue={true}
              allowCustomValue={true}
              isClearable={true}
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
        <QueryField {...props} query={query} />
      </>}
    </>
  );
}
