import React, { Component } from 'react';
import { Select, Spinner } from '@grafana/ui';
import { AssetInfo } from '../../types';
import { SitewiseCache } from 'sitewiseCache';
import { DataFrameView, SelectableValue } from '@grafana/data';
import { AssetModelSummary, AssetSummary } from 'queryResponseTypes';
import { AssetList } from './hierarchy/AssetList';

export interface Props {
  cache: SitewiseCache;
  asset?: AssetInfo; // The incoming value
  onAssetChanged: (assetId?: string) => void;
}

interface State {
  modelId?: string;
  models?: DataFrameView<AssetModelSummary>;
  assets?: DataFrameView<AssetSummary>;
}

export class BrowseModels extends Component<Props, State> {
  state: State = {};

  async componentDidMount() {
    const { asset, cache } = this.props;
    const update: State = {
      models: await cache.getModels(),
    };
    update.modelId = asset?.model_id ?? update.models?.get(0).id;
    update.assets = await cache.getAssetsOfType(update.modelId!);
    this.setState(update);
  }

  onModelIdChange = async (sel: SelectableValue<string>) => {
    const modelId = sel.value;
    const assets = modelId ? await this.props.cache.getAssetsOfType(modelId) : undefined;
    this.setState({ modelId, assets });
  };

  onAssetChanged = async (assetId?: string) => {
    if (assetId) {
      this.props.onAssetChanged(assetId);
    }
  };

  render() {
    const { models, assets, modelId } = this.state;
    if (!models) {
      return (
        <div>
          <Spinner />
          Loading models...
        </div>
      );
    }
    const modelOptions = models.map(m => ({
      value: m.id,
      label: m.name,
      description: m.description,
    }));

    const selectedModel = modelOptions.find(v => v.value === modelId);

    return (
      <>
        <div style={{ height: '60vh' }}>
          <h4>Model:</h4>
          <Select
            options={modelOptions}
            value={selectedModel || {}}
            onChange={this.onModelIdChange}
            backspaceRemovesValue={true}
            isSearchable={true}
            menuPlacement="bottom"
          />
          <br />
          <br />
          <h4>Assets:</h4>
          {selectedModel && assets ? (
            <AssetList assets={assets.toArray()} onSelect={this.onAssetChanged} />
          ) : (
            <>
              <p />
              <h6>No assets found.</h6>
            </>
          )}
        </div>
      </>
    );
  }
}
