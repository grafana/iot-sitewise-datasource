import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { SitewiseQuery, AssetInfo } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetExplorerModal } from './AssetExplorerModal';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

interface State {
  asset?: AssetInfo;
  assets: Array<SelectableValue<string>>;
  loading: boolean;
  openModal: boolean;
}

export class AssetPropPickerRows extends PureComponent<Props, State> {
  state: State = {
    assets: [],
    loading: true,
    openModal: false,
  };

  async updateInfo() {
    const { query, datasource } = this.props;
    const update: State = {
      loading: false,
    } as State;

    const cache = datasource.getCache(query.region);
    if (query?.assetId) {
      try {
        update.asset = await cache.getAssetInfo(query.assetId);
      } catch (err) {
        console.warn('error reading asset info', err);
      }
    }

    try {
      update.assets = await cache.getAssetPickerOptions();
    } catch (err) {
      console.warn('error getting options', err);
    }

    this.setState(update);
  }

  async componentDidMount() {
    this.updateInfo();
  }

  async componentDidUpdate(oldProps: Props) {
    const { query } = this.props;
    if (query?.assetId !== oldProps?.query?.assetId) {
      if (!query.assetId) {
        this.setState({ asset: undefined, loading: false });
      } else {
        this.setState({ loading: true });
        this.updateInfo();
      }
    }
  }

  onAssetChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId: sel.value! });
    onRunQuery();
  };

  onPropertyChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, propertyId: sel.value! });
    onRunQuery();
  };

  onSetAssetId = (assetId: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId });
    onRunQuery();
  };

  openAssetExplorer = () => {
    console.log('TODO!');
  };

  render() {
    const { query } = this.props;
    const { loading, asset, assets } = this.state;

    let current = query.assetId ? assets.find(v => v.value === query.assetId) : undefined;
    if (!current && query.assetId) {
      if (loading) {
        current = { label: 'loading...', value: query.assetId };
      } else if (asset) {
        current = { label: asset.name, description: query.assetId, value: query.assetId };
      } else {
        current = { label: `Unknown: ${query.assetId}`, value: query.assetId };
      }
    }

    const showProp = query.propertyId || asset;
    const properties = showProp ? (asset ? asset.properties : []) : [];

    return (
      <>
        <div className="gf-form">
          <InlineField label="Asset" labelWidth={10} grow={true}>
            <Select
              isLoading={loading}
              options={assets}
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
          <AssetExplorerModal {...this.props} />
        </div>
        {showProp && (
          <div className="gf-form">
            <InlineField label="Property" labelWidth={10} grow={true}>
              <Select
                isLoading={loading}
                options={properties}
                value={properties.find(p => p.value === query.propertyId)}
                onChange={this.onPropertyChange}
                placeholder="Select a property"
                isSearchable={true}
              />
            </InlineField>
          </div>
        )}
      </>
    );
  }
}
