import React, { FunctionComponent } from 'react';
import { AssetInfo } from '../../../types';
import { AssetSummary } from '../../../queryResponseTypes';
import { Button } from '@grafana/ui';
import { Card } from 'common/Card';

export interface Props {
  asset: AssetInfo | AssetSummary;
  onSelect: (assetId: string) => void;
  onInspect?: (assetId: string) => void;
}

export const AssetListItem: FunctionComponent<Props> = ({ asset, onInspect, onSelect }) => {
  return (
    <Card
      title={asset.name}
      description={asset.id}
      onClick={() => onSelect(asset.id)}
      actions={
        <>
          {onInspect && (
            <Button
              variant="secondary"
              onClick={(event: React.MouseEvent) => {
                event.stopPropagation();
                onInspect(asset.id);
              }}
              icon="folder"
            >
              BROWSE
            </Button>
          )}

          <Button variant="primary" icon="check">
            SELECT
          </Button>
        </>
      }
    />
  );
};
