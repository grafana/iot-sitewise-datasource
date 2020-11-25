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
  tab: 'Modal' | 'Hierarchy';
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
  state: State = { isOpen: false, tab: 'Hierarchy' };

  async componentDidMount() {
    const { assetId, region } = this.props;
    const cache = this.props.datasource.getCache(region);
    const asset = assetId ? await cache.getAssetInfo(assetId) : undefined;
    this.setState({ cache, asset });
  }

  async componentDidUpdate(oldProps: Props) {

    let update: State = {...this.state}
    let shouldUpdate = false;

    if (this.props.region !== oldProps.region) {
      shouldUpdate = true;
      update.cache = this.props.datasource.getCache(this.props.region);
      // this.setState({ cache });
    }
    if (this.props.assetId !== oldProps.assetId) {
      shouldUpdate = true;
      const { cache } = this.state;
      const { assetId } = this.props;
      // Asset changed from the parent... reset state
      update.asset = assetId ? await cache!.getAssetInfo(assetId) : undefined;
      // this.setState({ asset });
    }

    if (shouldUpdate) {
      this.setState(update)
    }
  }

  onSelectAsset = (assetId?: string) => {
    this.props.onAssetChanged(assetId);
    this.setState({ isOpen: false });
  };

  renderBody() {
    const { cache, tab, asset } = this.state;
    if (!cache) {
      return (
        <div>
          <Spinner />
          Loading...
        </div>
      );
    }

    switch (tab) {
      case 'Hierarchy':
        return <BrowseHierarchy cache={cache} asset={asset} onAssetSelected={this.onSelectAsset} />;
      case 'Modal':
        return <BrowseModels cache={cache} asset={asset} onAssetChanged={this.onSelectAsset} />;
    }
  }

  render() {
    const { isOpen, tab } = this.state;

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
                <Tab
                  css
                  label={'Hierarchy'}
                  active={'Hierarchy' === tab}
                  onChangeTab={() => this.setState({ tab: 'Hierarchy' })}
                />
                <Tab
                  css
                  label={'By Model'}
                  active={'Modal' === tab}
                  onChangeTab={() => this.setState({ tab: 'Modal' })}
                />
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
