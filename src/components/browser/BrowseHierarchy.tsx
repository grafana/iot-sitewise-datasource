import React, { Component } from 'react';
import { SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { DataFrameView, SelectableValue } from '@grafana/data';
import { AssetSummary } from '../../queryResponseTypes';
import { Button, Label, Select } from '@grafana/ui';

export interface State {
  childAssets?: DataFrameView<AssetSummary>;
  assets: Array<SelectableValue<string>>;
  asset?: AssetInfo;
  currentHierarchy?: string;
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

    const update: State = {
      ...this.state,
      asset: asset,
      assets: await cache.getAssetPickerOptions(),
    };
    this.setState(update);
  }

  onAssetChange = async (sel: SelectableValue<string>) => {
    if (sel.value) {
      const update: State = {
        ...this.state,
        asset: await this.props.cache.getAssetInfo(sel.value),
      };

      this.setState(update);
    }
  };

  onSetAssetId = async (assetId?: string) => {
    if (assetId) {
      this.setState({ ...this.state, asset: await this.props.cache.getAssetInfo(assetId) });
    }
  };

  onChildAssetChange = async (sel: SelectableValue<string>) => {
    const { cache } = this.props;

    if (sel.value) {
      this.setState({
        ...this.state,
        asset: await cache.getAssetInfo(sel.value),
        assets: await cache.getAssetPickerOptions(),
        currentHierarchy: undefined,
        childAssets: undefined,
      });
    }
  };

  onHierarchyChange = async (sel: SelectableValue<string>) => {
    const { asset } = this.state;
    const { cache } = this.props;

    if (sel.value && asset) {
      this.setState({
        ...this.state,
        currentHierarchy: sel.value,
        childAssets: await cache.getAssociatedAssets(asset.id, sel.value),
      });
    }
  };

  onAssetSelected = async (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    const { asset } = this.state;
    const { onAssetSelected } = this.props;
    if (asset) {
      onAssetSelected(asset.id);
    }
  };

  render() {
    const { asset, assets, childAssets, currentHierarchy } = this.state;

    let current = asset ? assets.find(v => v.value === asset.id) : undefined;
    if (!current && asset) {
      current = { label: asset.name, value: asset.id, description: asset.arn };
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
      asset && currentHierarchy
        ? asset.hierarchy.find(value => value.value === currentHierarchy)
        : { value: undefined, description: undefined };

    return (
      <div style={{ height: '60vh' }}>
        <Button name="copy" size="md" variant="secondary" onClick={this.onAssetSelected}>
          Select
        </Button>
        <br />
        <br />
        <Label description="asset to select">Asset:</Label>
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
        <br />
        <Label description="asset hierarchy to inspect">Hierarchies:</Label>
        <Select
          options={asset?.hierarchy}
          value={hierachyVal}
          onChange={this.onHierarchyChange}
          backspaceRemovesValue={true}
          isSearchable={true}
          menuPlacement="bottom"
        />
        {/* TODO: Add parent drill up */}
        <br />
        <Label description="child assets within the selected asset hierarchy">Children:</Label>
        <Select options={childOptions} isSearchable={true} onChange={this.onChildAssetChange} value={childVal} />
      </div>
    );
  }
}
