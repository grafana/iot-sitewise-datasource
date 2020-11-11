import React, { Component } from 'react';
import { Button, Icon, Modal, Spinner, Tab, TabContent, TabsBar } from '@grafana/ui';
import { AssetInfo } from '../../types';
import { DataSource } from 'DataSource';
import { SitewiseCache } from 'sitewiseCache';
import { BrowseModels } from './BrowseModels';
import { BrowseHierarchy } from './BrowseHierarchy';

export interface Props {
  datasource: DataSource;
  assetId?: string; // The incoming value
  region?: string;
  onAssetChanged: (assetId?: string) => void;
}

interface State {
  isOpen: boolean;
  byModel: boolean;
  cache?: SitewiseCache;
  asset?: AssetInfo;
}

export const ModalHeader = () => {
  return (
    <div className="modal-header-title">
      <Icon name="folder-open" size="lg" />
      <span className="p-l-1">Asset Browser</span>
    </div>
  );
};

export class AssetBrowser extends Component<Props, State> {
  state: State = { isOpen: false, byModel: false };

  async componentDidMount() {
    const { assetId, region } = this.props;
    const cache = this.props.datasource.getCache(region);
    const asset = assetId ? await cache.getAssetInfo(assetId) : undefined;
    this.setState({ cache, asset });
  }

  async componentDidUpdate(oldProps: Props) {
    if (this.props.region !== oldProps.region) {
      const cache = this.props.datasource.getCache(this.props.region);
      this.setState({ cache });
    }
    if (this.props.assetId !== oldProps.assetId) {
      const { cache } = this.state;
      const { assetId } = this.props;
      // Asset changed from the parent... reset state
      const asset = assetId ? await cache!.getAssetInfo(assetId) : undefined;
      this.setState({ asset });
    }
  }

  onSelectAsset = (assetId?: string) => {
    this.props.onAssetChanged(assetId);
    this.setState({ isOpen: false });
  };

  renderBody() {
    const { cache, byModel, asset } = this.state;
    if (!cache) {
      return (
        <div>
          <Spinner />
          Loading...
        </div>
      );
    }
    if (byModel) {
      return <BrowseModels cache={cache} asset={asset} onAssetChanged={this.onSelectAsset} />;
    }
    return <BrowseHierarchy cache={cache} asset={asset} onAssetSelected={this.onSelectAsset} />;
  }

  render() {
    const { isOpen, byModel } = this.state;

    return (
      <>
        <Button
          variant="secondary"
          size="md"
          icon="folder-open"
          onClick={event =>
            this.setState({ isOpen: true }, () => {
              console.log(this.state);
            })
          }
        >
          Explore
        </Button>
        <Modal title={<ModalHeader />} isOpen={isOpen} onDismiss={() => this.setState({ isOpen: false })}>
          <div>
            <div>
              <TabsBar>
                <Tab css label={'Heiarchy'} active={!byModel} onChangeTab={() => this.setState({ byModel: false })} />
                <Tab css label={'By Model'} active={byModel} onChangeTab={() => this.setState({ byModel: true })} />
              </TabsBar>
              <TabContent style={{ maxHeight: '90vh' }}>
                <div>{this.renderBody()}</div>
              </TabContent>
            </div>
          </div>
        </Modal>
      </>
    );
  }
}
