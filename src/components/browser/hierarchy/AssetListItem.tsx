import React, { FunctionComponent } from 'react';
import { AssetInfo } from '../../../types';
import { AssetSummary } from '../../../queryResponseTypes';
import { stylesFactory, useTheme } from '@grafana/ui';
import { LinkButton } from '@grafana/ui';
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
      border: 1px solid blue;
    `,
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
      onClick={() => {
        onInspect ? onInspect(asset.id) : onSelect(asset.id);
      }}
      actions={
        <LinkButton
          variant="primary"
          target="_blank"
          rel="noopener"
          onClick={() => onSelect(asset.id)}
          icon="external-link-alt"
        >
          SELECT
        </LinkButton>
      }
    />
  );
};
