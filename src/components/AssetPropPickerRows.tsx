import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { SitewiseQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetInfo } from 'queryResponseTypes';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

interface State {
  showId?: boolean;
  asset?: AssetInfo;
  options: Array<SelectableValue<string>>;
  loading: boolean;
}

export class AssetPropPickerRows extends PureComponent<Props, State> {
  state: State = {
    options: [],
    loading: true,
  };

  async componentDidMount() {
    const { query, datasource } = this.props;
    const update: State = {
      loading: false,
    } as State;

    const cache = datasource.getCache(query.region);
    if (query.assetId) {
      try {
        update.asset = await cache.getAssetInfo(query.assetId);
      } catch (err) {
        console.warn('error reading assets', err);
      }
    }
    update.options = await cache.getAssetPickerOptions();

    this.setState(update);
  }

  onAssetChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId: sel.value! });
    onRunQuery();
  };

  onSetAssetId = (assetId: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    const { loading, options, asset } = this.state;
    let current = options.find(v => v.value === query.assetId);
    if (!current && query.assetId) {
      if (loading) {
        current = { label: 'loading...', value: query.assetId };
      } else if (asset) {
        current = { label: asset.name, description: query.assetId, value: query.assetId };
      } else {
        current = { label: `Unknown: ${query.assetId}`, value: query.assetId };
      }
    }

    const showProp = query.propertyId || asset?.properties;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Asset" labelWidth={10} grow={true}>
            <Select
              isLoading={loading}
              options={options}
              value={current}
              onChange={this.onAssetChange}
              placeholder="Select an asset"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              onCreateOption={this.onSetAssetId}
              formatCreateLabel={txt => `Asset ID: ${txt}`}
            />
          </InlineField>
        </div>
        {showProp && (
          <div className="gf-form">
            <InlineField label="Property" labelWidth={10} grow={true}>
              <div>TODO: property picker</div>
            </InlineField>
          </div>
        )}
      </>
    );
  }
}
