import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { ListAssetsQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';

type Props = SitewiseQueryEditorProps<ListAssetsQuery>;

const filters = [
  {
    label: 'Top Level',
    value: 'TOP_LEVEL',
    description: 'The list includes only top-level assets in the asset hierarchy tree',
  },
  { label: 'All', value: 'ALL', description: 'The list includes all assets for a given asset model ID' },
];

export class ListAssetsQueryEditor extends PureComponent<Props> {
  onAssetModelIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetModelId: sel.value! });
    onRunQuery();
  };

  onFilterChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: sel.value as 'ALL' | 'TOP_LEVEL' });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    const modelIds: Array<SelectableValue<string>> = [];

    if (query.assetModelId) {
      modelIds.push({
        label: query.assetModelId,
        value: query.assetModelId,
      });
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Model ID" labelWidth={10} grow={true}>
            <Select
              options={modelIds}
              value={modelIds.find(v => v.value === query.assetModelId) || undefined}
              onChange={this.onAssetModelIdChange}
              placeholder="Select an asset model id"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              formatCreateLabel={txt => `Model ID: ${txt}`}
            />
          </InlineField>
        </div>
        <div className="gf-form">
          <InlineField label="Filter" labelWidth={10} grow={true}>
            <Select
              options={filters}
              value={filters.find(v => v.value === query.filter) || filters[0]}
              onChange={this.onFilterChange}
              placeholder="Select a property"
            />
          </InlineField>
        </div>
      </>
    );
  }
}
