import React, { Component } from 'react';
import { SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { DataFrameView, SelectableValue } from '@grafana/data';
import { AssetSummary } from '../../queryResponseTypes';
import { Select } from '@grafana/ui';

export interface HierarchyState {
  childAssets?: DataFrameView<AssetSummary>;
  assets: Array<SelectableValue<string>>;
  currentAsset?: AssetInfo;
  currentHierarchy?: string;
}

export interface Props {
  cache: SitewiseCache;
  currentAsset?: AssetInfo; // The incoming value
  onAssetChanged: (assetId?: string) => void;
}

export class AssetHierarchy extends Component<Props, HierarchyState> {
  state: HierarchyState = { assets: [] };

  async componentDidMount() {
    const { currentAsset, cache } = this.props;

    const update: HierarchyState = {
      ...this.state,
      currentAsset: currentAsset,
      assets: await cache.getAssetPickerOptions(),
    };
    this.setState(update);
  }

  onAssetChange = async (sel: SelectableValue<string>) => {
    if (sel.value) {
      const update: HierarchyState = {
        ...this.state,
        currentAsset: await this.props.cache.getAssetInfo(sel.value),
      };

      this.setState(update);
    }
  };

  onSetAssetId = async (assetId?: string) => {
    if (assetId) {
      this.setState({ ...this.state, currentAsset: await this.props.cache.getAssetInfo(assetId) });
    }
  };

  onChildAssetChange = async (sel: SelectableValue<string>) => {
    const { cache } = this.props;

    if (sel.value) {
      this.setState({
        ...this.state,
        currentAsset: await cache.getAssetInfo(sel.value),
        assets: await cache.getAssetPickerOptions(),
        currentHierarchy: undefined,
        childAssets: undefined,
      });
    }
  };

  onHierarchyChange = async (sel: SelectableValue<string>) => {
    const { currentAsset } = this.state;
    const { cache } = this.props;

    if (sel.value && currentAsset) {
      this.setState({
        ...this.state,
        currentHierarchy: sel.value,
        childAssets: await cache.getAssociatedAssets(currentAsset.id, sel.value),
      });
    }
  };

  render() {
    const { currentAsset, assets, childAssets, currentHierarchy } = this.state;

    let current = currentAsset ? assets.find(v => v.value === currentAsset.id) : undefined;
    if (!current && currentAsset) {
      current = { label: currentAsset.name, value: currentAsset.id, description: currentAsset.arn };
    }

    let childOptions = childAssets
      ? childAssets.map(asset => {
          return {
            label: asset.name,
            value: asset.id,
            description: asset.arn,
          };
        })
      : [];

    let childVal = { value: undefined, description: undefined };

    let hierachyVal =
      currentAsset && currentHierarchy
        ? currentAsset.hierarchy.find(value => value.value === currentHierarchy)
        : { value: undefined, description: undefined };

    return (
      <div style={{ height: '60vh' }}>
        <h4>Selected Asset:</h4>
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

        <h4>Hierarchies:</h4>
        <Select
          options={currentAsset?.hierarchy}
          value={hierachyVal}
          onChange={this.onHierarchyChange}
          backspaceRemovesValue={true}
          isSearchable={true}
          menuPlacement="bottom"
        />

        {/* TODO: Add parent drill up */}

        <h4>Children:</h4>
        <Select options={childOptions} isSearchable={true} onChange={this.onChildAssetChange} value={childVal} />
      </div>
    );
  }
}
