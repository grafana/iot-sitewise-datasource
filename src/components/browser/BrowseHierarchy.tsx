import React, { Component } from 'react';
import { assetSummaryToAssetInfo, SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { DataFrameView, SelectableValue } from '@grafana/data';
import { AssetSummary } from '../../queryResponseTypes';
import { Button, Label, Select } from '@grafana/ui';

const UNSET_VAL = { value: undefined, description: undefined };

export interface State {
  children?: DataFrameView<AssetSummary>;
  assets: Array<SelectableValue<string>>;
  asset?: AssetInfo;
  parents?: AssetInfo[];
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
    return [];
  };

  onSetAssetId = async (assetId?: string) => {
    if (assetId) {
      this.setState({ ...this.state, asset: await this.props.cache.getAssetInfo(assetId) });
    }
  };

  setSelectedAssetInfo = async (assetId?: string) => {
    const { cache } = this.props;
    if (assetId) {
      this.setState({
        ...this.state,
        asset: await cache.getAssetInfo(assetId),
        assets: await cache.getAssetPickerOptions(),
        parents: await this.getParentInfo(assetId),
        currentHierarchy: undefined,
        children: undefined,
      });
    }
  };

  onAssetChange = async (sel: SelectableValue<string>) => {
    await this.setSelectedAssetInfo(sel.value);
  };

  onChildAssetChange = async (sel: SelectableValue<string>) => {
    await this.setSelectedAssetInfo(sel.value);
  };

  onParentChanged = async (sel: SelectableValue<string>) => {
    await this.setSelectedAssetInfo(sel.value);
  };

  onHierarchyChange = async (sel: SelectableValue<string>) => {
    const { asset } = this.state;
    const { cache } = this.props;

    if (sel.value && asset) {
      this.setState({
        ...this.state,
        currentHierarchy: sel.value,
        children: await cache.getAssociatedAssets(asset.id, sel.value),
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

  onParentSelected = async (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    const { cache } = this.props;
    const { asset } = this.state;
    if (asset) {
      const parentSummary = await cache.getAssociatedAssets(asset.id);
      if (parentSummary.length === 1) {
        this.setState({ ...this.state, asset: assetSummaryToAssetInfo(parentSummary)[0] });
      }
    }
  };

  render() {
    const { asset, assets, parents, children, currentHierarchy } = this.state;

    let current = asset ? assets.find(v => v.value === asset.id) : undefined;
    if (!current && asset) {
      current = { label: asset.name, value: asset.id, description: asset.arn };
    }

    let childOptions = children
      ? children.map(asset => {
          return {
            label: asset.name,
            value: asset.id,
            description: asset.arn,
          };
        })
      : [];

    let hierachyVal =
      asset && currentHierarchy
        ? asset.hierarchy.find(value => value.value === currentHierarchy)
        : { value: undefined, description: undefined };

    let parentVals = parents
      ? parents.map(p => {
          return {
            label: p.name,
            value: p.id,
            description: p.arn,
          };
        })
      : [];

    return (
      <div style={{ height: '60vh' }}>
        <Button name="copy" size="md" variant="secondary" onClick={this.onAssetSelected}>
          Select
        </Button>
        <p />
        <p />
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
        <p />
        <Label description="asset hierarchy to inspect">Hierarchies:</Label>
        <Select
          options={asset?.hierarchy}
          value={hierachyVal}
          onChange={this.onHierarchyChange}
          backspaceRemovesValue={true}
          isSearchable={true}
          menuPlacement="bottom"
        />
        <p />
        <Label description="asset parent to select">Parents:</Label>
        <Select
          options={parentVals}
          value={UNSET_VAL}
          onChange={this.onParentChanged}
          backspaceRemovesValue={true}
          isSearchable={true}
          menuPlacement="bottom"
        />
        <p />
        <Label description="child assets within the selected asset hierarchy">Children:</Label>
        <Select options={childOptions} isSearchable={true} onChange={this.onChildAssetChange} value={UNSET_VAL} />
      </div>
    );
  }
}
