import React, { Component } from 'react';
import { SitewiseDatasource } from '../../DataSource';
import { Button, Icon, IconName, Modal } from '@grafana/ui';
import { IconSize } from '@grafana/ui/types/icon';

interface Props {
  datasource: SitewiseDatasource;
  isOpen?: boolean;
}

interface State {
  isOpen: boolean;
}

export const ModalHeader = ({
  title,
  iconName = 'exclamation-triangle',
  size = 'lg',
}: {
  title: string;
  iconName: IconName;
  size: IconSize;
}) => {
  return (
    <div className="modal-header-title">
      <Icon name={iconName} size={size} />
      <span className="p-l-1">{title}</span>
    </div>
  );
};

export class AssetExplorerModal extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { isOpen: props.isOpen || false };
  }

  render() {
    const { isOpen } = this.state;

    return (
      <>
        <Button
          variant="primary"
          size="md"
          icon="search-plus"
          onClick={event =>
            this.setState({ isOpen: true }, () => {
              console.log(this.state);
            })
          }
        >
          Asset Search
        </Button>
        <Modal
          title={<ModalHeader title="Search for assets" iconName="search" size="xxl" />}
          isOpen={isOpen}
          onDismiss={() => this.setState({ isOpen: false })}
        >
          <div>
            <h3>Search by Asset Model types, or Property attributes!!!</h3>
          </div>
        </Modal>
      </>
    );
  }
}
