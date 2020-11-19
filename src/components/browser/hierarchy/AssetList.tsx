import { AssetInfo } from '../../../types';
import { AssetSummary } from '../../../queryResponseTypes';
import React, { FunctionComponent } from 'react';
import { GrafanaTheme } from '@grafana/data';
import { CollapsableSection, Label, Spinner, styleMixins, stylesFactory, useTheme } from '@grafana/ui';
import { css } from 'emotion';
import { AssetListItem } from './AssetListItem';

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    container: css`
      width: 100%;
      height: auto;
    `,
    listItem: css`
      ${styleMixins.listItem(theme)}
    `,
  };
});

export interface ListInfo {
  name?: string;
  id?: string;
  description?: string;
}

export interface Props {
  listInfo: ListInfo;
  assets?: Array<AssetInfo | AssetSummary>;
  onSelect: (assetId: string) => void;
  onInspect?: (assetId: string) => void;
}

export const AssetList: FunctionComponent<Props> = ({ listInfo, assets, onSelect, onInspect }) => {
  const theme = useTheme();
  const style = getStyles(theme);

  const label = ((<Label description={listInfo.description}>{listInfo.name}</Label>) as unknown) as string;

  const renderChildren = () => {
    if (!assets) {
      return (
        <>
          <Spinner />
          Loading assets...
        </>
      );
    }

    return (
      <ul key={listInfo.id}>
        {assets.map(c => {
          return (
            <li key={c.name} className={style.listItem}>
              <AssetListItem asset={c} onInspect={onInspect} onSelect={onSelect} />
            </li>
          );
        })}
      </ul>
    );
  };

  return (
    <div className={style.container}>
      <CollapsableSection label={label} isOpen={false}>
        {renderChildren()}
      </CollapsableSection>
    </div>
  );
};
