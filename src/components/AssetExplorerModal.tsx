import React, { Component } from 'react';
import { Button, Icon, IconName, Modal } from '@grafana/ui';
import { IconSize } from '@grafana/ui/types/icon';
import { SitewiseQueryEditorProps } from './types';
import { SitewiseQuery } from 'types';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

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
  state: State = { isOpen: false };

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
          Explore
        </Button>
        <Modal
          title={<ModalHeader title="Browse for assets" iconName="search" size="xxl" />}
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
