import React, { FunctionComponent } from 'react';
import { AssetInfo } from '../../../types';
import { AssetSummary } from '../../../queryResponseTypes';
import { Button, stylesFactory, useTheme } from '@grafana/ui';
import { Card } from 'common/Card';
import { GrafanaTheme } from '@grafana/data';
import { css } from 'emotion';

export interface Props {
  asset: AssetInfo | AssetSummary;
  onSelect: (assetId: string) => void;
  onInspect?: (assetId: string) => void;
  current?: boolean;
}

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    current: css`
      border: 1px solid ${theme.colors.formInputBorderHover};
    `, // $panel-editor-viz-item-border-hover;
  };
});

export const AssetListItem: FunctionComponent<Props> = ({ asset, current, onInspect, onSelect }) => {
  const theme = useTheme();
  const style = getStyles(theme);

  return (
    <Card
      className={current ? style.current : undefined}
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
