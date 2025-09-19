// components/common/ConfirmDialog.tsx
import React from 'react';
import { Modal, Button, HorizontalGroup } from '@grafana/ui';

interface ConfirmDialogProps {
  isOpen: boolean;
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export const ConfirmDialog: React.FC<ConfirmDialogProps> = ({ isOpen, title, message, onConfirm, onCancel }) => {
  return (
    <Modal title={title} isOpen={isOpen} onDismiss={onCancel} closeOnEscape closeOnBackdropClick>
      <div style={{ padding: '16px 0' }}>
        <p>{message}</p>
      </div>
      <HorizontalGroup justify="flex-end">
        <Button variant="secondary" onClick={onCancel}>
          No
        </Button>
        <Button variant="destructive" onClick={onConfirm}>
          Yes
        </Button>
      </HorizontalGroup>
    </Modal>
  );
};
