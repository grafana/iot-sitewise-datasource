import React, { PureComponent } from 'react';
import { DataFrameView, SelectableValue } from '@grafana/data';
import { ListAssetsQuery } from 'types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetModelSummary } from 'queryResponseTypes';
import { firstLabelWith } from './QueryEditor';

type Props = SitewiseQueryEditorProps<ListAssetsQuery>;

interface State {
  models?: DataFrameView<AssetModelSummary>;
}

const filters = [
  {
    label: 'Top Level',
    value: 'TOP_LEVEL',
    description: 'The list includes only top-level assets in the asset hierarchy tree',
  },
  { label: 'All', value: 'ALL', description: 'The list includes all assets for a given asset model ID' },
];

export class ListAssetsQueryEditor extends PureComponent<Props, State> {
  state: State = {};

  async componentDidMount() {
    const { query } = this.props;
    const cache = this.props.datasource.getCache(query.region);
    const models = await cache.getModels();
    this.setState({ models });
  }

  onAssetModelIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, modelId: sel.value! });
    onRunQuery();
  };

  onFilterChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: sel.value as 'ALL' | 'TOP_LEVEL' });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    const { models } = this.state;
    const modelIds = models
      ? models.map((m) => ({
          value: m.id,
          label: m.name,
          description: m.description,
        }))
      : [];
    let currentModel = modelIds.find((m) => m.value === query.modelId);
    if (query.modelId && !currentModel) {
      currentModel = {
        value: query.modelId,
        label: 'Model ID: ' + query.modelId,
        description: '',
      };
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Model ID" labelWidth={firstLabelWith} grow={true}>
            <Select
              isLoading={!models}
              options={modelIds}
              value={currentModel}
              onChange={this.onAssetModelIdChange}
              placeholder="Select an asset model id"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              formatCreateLabel={(txt) => `Model ID: ${txt}`}
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
        <div className="gf-form">
          <InlineField label="Filter" labelWidth={firstLabelWith} grow={true}>
            <Select
              options={filters}
              value={filters.find((v) => v.value === query.filter) || filters[0]}
              onChange={this.onFilterChange}
              placeholder="Select a property"
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
      </>
    );
  }
}
