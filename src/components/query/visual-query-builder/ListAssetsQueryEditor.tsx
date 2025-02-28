import React from 'react';
import { SelectableValue } from '@grafana/data';
import { ListAssetsQuery } from 'types';
import { Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { useModelsOptions } from 'sitewiseCache';
import { useOptionsWithVariables } from 'common/useOptionsWithVariables';

const FILTERS = [
  {
    label: 'Top Level',
    value: 'TOP_LEVEL',
    description: 'The list includes only top-level assets in the asset hierarchy tree',
  },
  { label: 'All', value: 'ALL', description: 'The list includes all assets for a given asset model ID' },
];

export const ListAssetsQueryEditor = ({ query, datasource, onChange }: SitewiseQueryEditorProps<ListAssetsQuery>) => {
  const cache = datasource.getCache(query.region);
  const { isLoading, options } = useModelsOptions(cache);
  const modelId = useOptionsWithVariables({ current: query.modelId, options });

  const onAssetModelIdChange = (sel: SelectableValue<string>) => {
    onChange({ ...query, modelId: sel.value! });
  };

  const onFilterChange = (sel: SelectableValue<string>) => {
    onChange({ ...query, filter: sel.value as 'ALL' | 'TOP_LEVEL' });
  };

  return (
    <EditorRow>
      <EditorFieldGroup>
        <EditorField label="Model ID" htmlFor="model" width={30}>
          <Select
            inputId="model"
            aria-label="Model ID"
            isLoading={isLoading}
            options={modelId.options}
            value={modelId.current}
            onChange={onAssetModelIdChange}
            placeholder="Select an asset model id"
            allowCustomValue={true}
            isClearable={true}
            isSearchable={true}
            formatCreateLabel={(txt) => `Model ID: ${txt}`}
            menuPlacement="auto"
          />
        </EditorField>
        <EditorField label="Filter" htmlFor="filter" width={20}>
          <Select
            inputId="filter"
            aria-label="Filter"
            options={FILTERS}
            value={FILTERS.find((v) => v.value === query.filter) || FILTERS[0]}
            onChange={onFilterChange}
            placeholder="Select a filter"
            menuPlacement="auto"
          />
        </EditorField>
      </EditorFieldGroup>
    </EditorRow>
  );
};
