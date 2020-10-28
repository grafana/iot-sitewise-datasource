import React, { Component } from 'react';
import { Button, Icon, Modal } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { SitewiseQuery } from 'types';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

interface State {
  isOpen: boolean;
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
  state: State = { isOpen: false };

  render() {
    const { isOpen } = this.state;

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
            <h3>Search by Asset Model types, or Property attributes!!!</h3>
          </div>
        </Modal>
      </>
    );
  }
}
