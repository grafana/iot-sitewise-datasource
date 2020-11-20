import React, { Component } from 'react';
import { assetSummaryToAssetInfo, SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { SelectableValue } from '@grafana/data';
import { Select } from '@grafana/ui';
import { AssetHierarchyList } from './hierarchy/AssetHierarchyList';
import { AssetListItem } from './hierarchy/AssetListItem';

export interface State {
  assets: Array<SelectableValue<string>>;
  asset?: AssetInfo;
  parents?: AssetInfo[];
}

export interface Props {
  cache: SitewiseCache;
  asset?: AssetInfo; // The incoming value
  onAssetSelected: (assetId?: string) => void;
}

export class BrowseHierarchy extends Component<Props, State> {
  state: State = { assets: [] };

  async componentDidMount() {
    const { asset, cache } = this.props;

    const parentInfo = asset ? await this.getParentInfo(asset.id) : undefined;

    const update: State = {
      ...this.state,
      asset: asset,
      assets: await cache.getAssetPickerOptions(),
      parents: parentInfo,
    };
    this.setState(update);
  }

  getParentInfo = async (assetId: string): Promise<AssetInfo[]> => {
    const { cache } = this.props;

    const parentSummaries = await cache.getAssociatedAssets(assetId);
    return assetSummaryToAssetInfo(parentSummaries);
  };

  onSetAssetId = async (assetId?: string) => {
    await this.setSelectedAssetInfo(assetId);
  };

  setSelectedAssetInfo = async (assetId?: string) => {
    const { cache } = this.props;
    if (assetId) {
      this.setState({
        ...this.state,
        asset: await cache.getAssetInfo(assetId),
        assets: await cache.getAssetPickerOptions(),
        parents: await this.getParentInfo(assetId),
      });
    }
  };

  onAssetChange = async (sel: SelectableValue<string>) => {
    await this.setSelectedAssetInfo(sel.value);
  };

  onAssetSelected = async (assetId?: string) => {
    if (assetId) {
      this.props.onAssetSelected(assetId);
    }
  };

  renderParents() {
    const { asset, parents } = this.state;
    if (asset && parents?.length) {
      return (
        <Select
          options={parents.map(v => ({ label: v.name, value: v.id, description: v.id }))}
          onChange={this.onAssetChange}
          placeholder="Parent assets..."
          menuPlacement="bottom"
        />
      );
    }
    return;
  }

  renderHierarchies = () => {
    const { asset } = this.state;

    if (asset) {
      return (
        <ul>
          {asset.hierarchy.length ? (
            asset.hierarchy.map(h => {
              return (
                <li key={h.label}>
                  <AssetHierarchyList
                    hierarchy={{ name: h.label, id: h.value }}
                    asset={asset}
                    cache={this.props.cache}
                    onInspect={this.onSetAssetId}
                    onSelect={this.onAssetSelected}
                  />
                </li>
              );
            })
          ) : (
            <h6>No hierarchies found for asset.</h6>
          )}
        </ul>
      );
    }

    return <></>;
  };

  render() {
    const { asset, assets } = this.state;

    let current = asset ? assets.find(v => v.value === asset.id) : undefined;
    if (!current && asset) {
      current = { label: asset.name, value: asset.id, description: asset.arn };
    }

    return (
      <>
        {asset ? (
          <>
            <AssetListItem current={true} asset={asset} onSelect={() => this.onAssetSelected(asset?.id)} />
            {this.renderParents()}
          </>
        ) : (
          <Select
            options={assets}
            value={current}
            onChange={this.onAssetChange}
            placeholder="Select an asset"
            allowCustomValue={true}
            isClearable={true}
            isSearchable={true}
            onCreateOption={this.onSetAssetId}
            formatCreateLabel={txt => `Asset ID: ${txt}`}
            menuPlacement="bottom"
          />
        )}
        <br />
        <div style={{ height: '60vh', overflow: 'auto' }}>
          <h5> Asset Hierarchies: </h5>
          <this.renderHierarchies />
        </div>
      </>
    );
  }
}
