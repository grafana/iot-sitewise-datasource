import React, { Component } from 'react';
import { assetSummaryToAssetInfo, SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { SelectableValue } from '@grafana/data';
import { Input, Select } from '@grafana/ui';
import { AssetHierarchyList } from './hierarchy/AssetHierarchyList';
import { AssetListItem } from './hierarchy/AssetListItem';
import { AssetList } from './hierarchy/AssetList';

export interface State {
  assets: Array<SelectableValue<string>>;
  asset?: AssetInfo;
  parents?: AssetInfo[];
  search?: string;
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

  onSearchChange = (event: React.FormEvent<HTMLInputElement>) => {
    this.setState({ search: event.currentTarget.value });
  };

  renderHierarchies = () => {
    const { asset, search } = this.state;
    if (!asset) {
      return;
    }
    if (!asset.hierarchy.length) {
      return <h6>No hierarchies found for asset.</h6>;
    }

    return (
      <>
        <h5> Asset Hierarchies: </h5>
        <div style={{ height: '40vh', overflow: 'auto' }}>
          <Input value={search} onChange={this.onSearchChange} placeholder="search..." />
          <br />

          {asset.hierarchy.map((h) => {
            return (
              <AssetHierarchyList
                key={h.value}
                hierarchy={{ name: h.label, id: h.value }}
                asset={asset}
                search={search}
                cache={this.props.cache}
                onInspect={this.onSetAssetId}
                onSelect={this.onAssetSelected}
              />
            );
          })}
        </div>
      </>
    );
  };

  renderParents = () => {
    const { parents } = this.state;

    return (
      <>
        <h5> Parents: </h5>
        {parents ? (
          <AssetList
            listInfo={{ id: '', name: 'Show Parents', description: '' }}
            assets={parents}
            onSelect={this.onAssetSelected}
            onInspect={this.onSetAssetId}
            open={false}
          />
        ) : (
          <h6>No parents found.</h6>
        )}
      </>
    );
  };

  render() {
    const { asset, assets } = this.state;

    let current = asset ? assets.find((v) => v.value === asset.id) : undefined;
    if (!current && asset) {
      current = { label: asset.name, value: asset.id, description: asset.arn };
    }

    return (
      <div style={{ height: '60vh' }}>
        {asset ? (
          <>
            <this.renderParents />
            <p />
            <h5> Selected Asset: </h5>
            <AssetListItem asset={asset} onSelect={() => this.onAssetSelected(asset?.id)} />
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
            formatCreateLabel={(txt) => `Asset ID: ${txt}`}
            menuPlacement="auto"
          />
        )}
        <br />
        {this.renderHierarchies()}
      </div>
    );
  }
}
