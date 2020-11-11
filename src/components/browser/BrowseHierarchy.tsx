import React, { Component } from 'react';
import { AssetInfo } from '../../types';
import { SitewiseCache } from 'sitewiseCache';
import { DataFrameView } from '@grafana/data';
import { AssetSummary } from 'queryResponseTypes';
import { AssetHierarchy } from './AssetHierarchy';

export interface Props {
  cache: SitewiseCache;
  asset?: AssetInfo; // The incoming value
  onAssetChanged: (assetId?: string) => void;
}

interface State {
  modelId?: string;
  hierarchy: Array<DataFrameView<AssetSummary>>;
}

export class BrowseHierarchy extends Component<Props, State> {
  state: State = { hierarchy: [] };

  async componentDidMount() {
    const { asset, cache } = this.props;
    if (asset != null) {
      console.log('TODO... find the tree...');
    }
    const topLevel = await cache.getTopLevelAssets();
    const hierarchy = [topLevel];
    this.setState({ hierarchy });
  }

  render() {
    // const { hierarchy } = this.state;
    // if (!hierarchy.length) {
    //   return (
    //     <div>
    //       <Spinner />
    //       Loading hierarchy...
    //     </div>
    //   );
    // }

    return (
      // <div style={{ height: '60vh' }}>
      //   {hierarchy.map((level, idx) => {
      //     if (idx === hierarchy.length - 1) {
      //       return <div key={idx}>SHOW EACH?</div>;
      //     }
      //     return <div key={idx}>SELECT for level... {idx}</div>;
      //   })}
      // </div>
      <AssetHierarchy
        cache={this.props.cache}
        currentAsset={this.props.asset}
        onAssetChanged={assetId => console.log(assetId)}
      />
    );
  }
}
