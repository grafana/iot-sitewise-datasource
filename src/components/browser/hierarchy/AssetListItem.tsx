import React, { FunctionComponent, useState } from 'react';
import { css } from 'emotion';
import { AssetInfo } from '../../../types';
import { AssetSummary } from '../../../queryResponseTypes';
import { Button, Icon, styleMixins, stylesFactory, useTheme } from '@grafana/ui';
import { GrafanaTheme } from '@grafana/data';

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    assetRow: css`
      width: 100%;
      display: flex;
      flex-direction: row;
      margin: 0 ${theme.spacing.xs} ${theme.spacing.xs} 0;
      vertical-align: middle
      &:hover {
        border: 1px solid green;
        background: ${styleMixins.hoverColor(theme.colors.bg2, theme)};
      }
    `,
    description: css`
      label: Label-description;
      color: ${theme.colors.formDescription};
      font-size: ${theme.typography.size.sm};
      font-weight: ${theme.typography.weight.regular};
      margin-top: ${theme.spacing.xxs};
      display: block;
    `,
    infoIcon: css`
      margin-right: ${theme.spacing.sm};
    `,
    assetTitle: css`
      margin-right: ${theme.spacing.md};
    `,
    buttons: css`
      margin-top: ${theme.spacing.sm};
      margin-left: auto;
      margin-right: ${theme.spacing.md};
    `,
  };
});

export interface Props {
  asset: AssetInfo | AssetSummary;
  onSelect: (assetId: string) => void;
  onInspect?: (assetId: string) => void;
}

export const AssetListItem: FunctionComponent<Props> = ({ asset, onInspect, onSelect }) => {
  const theme = useTheme();
  const style = getStyles(theme);

  const [isEntered, setEntered] = useState<boolean>(false);

  const onComponentMouseEntered = (_: any) => {
    setEntered(true);
  };
  const onComponentMouseLeave = (_: any) => {
    setEntered(false);
  };

  return (
    <div className={style.assetRow} onMouseEnter={onComponentMouseEntered} onMouseLeave={onComponentMouseLeave}>
      <Icon name="info-circle" size="md" className={style.infoIcon} />

      <div className={style.assetTitle}>
        <h4>{asset.name}</h4>
        <div className={style.description}>{asset.id}</div>
      </div>

      <div className={style.buttons} hidden={!isEntered}>
        <Button
          icon="arrow-up"
          size="md"
          variant="link"
          hidden={!onInspect}
          onClick={_ => {
            if (onInspect) {
              onInspect(asset.id);
            }
          }}
        >
          {' '}
          Inspect{' '}
        </Button>
        <Button icon="save" size="md" variant="primary" onClick={_ => onSelect(asset.id)}>
          {' '}
          Select{' '}
        </Button>
      </div>
    </div>
  );
};
