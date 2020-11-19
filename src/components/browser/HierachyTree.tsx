import React, { Component } from 'react';
import { assetSummaryToAssetInfo, SitewiseCache } from '../../sitewiseCache';
import { AssetInfo } from '../../types';
import { SelectableValue } from '@grafana/data';
import { Button, Label, Select } from '@grafana/ui';
import { AssetHierarchy } from './hierarchy/AssetHierarchyList';

const UNSET_VAL = { value: undefined, description: undefined };

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

export class HierachyTree extends Component<Props, State> {
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

  onParentChanged = async (sel: SelectableValue<string>) => {
    await this.setSelectedAssetInfo(sel.value);
  };

  onAssetSelected = async (assetId?: string) => {
    const { onAssetSelected } = this.props;
    if (assetId) {
      onAssetSelected(assetId);
    }
  };

  render() {
    const { asset, assets, parents } = this.state;

    let current = asset ? assets.find(v => v.value === asset.id) : undefined;
    if (!current && asset) {
      current = { label: asset.name, value: asset.id, description: asset.arn };
    }

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
      <div style={{ height: '60vh', overflow: 'auto' }}>
        <Button name="copy" size="md" variant="secondary" onClick={_ => this.onAssetSelected(asset?.id)}>
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
        {parents && parents.length > 0 ? (
          <>
            <Label description="asset parent to select">Parents:</Label>
            <Select
              options={parentVals}
              value={UNSET_VAL}
              onChange={this.onParentChanged}
              backspaceRemovesValue={true}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </>
        ) : (
          undefined
        )}
        <p />

        {asset ? (
          <ul>
            {asset.hierarchy.map(h => {
              return (
                <li key={h.label}>
                  <AssetHierarchy
                    hierarchy={{ name: h.label, id: h.value }}
                    asset={asset}
                    cache={this.props.cache}
                    onInspect={this.onSetAssetId}
                    onSelect={this.onAssetSelected}
                  />
                </li>
              );
            })}
          </ul>
        ) : (
          undefined
        )}
      </div>
    );
  }
}
