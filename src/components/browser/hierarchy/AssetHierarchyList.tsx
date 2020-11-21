import React, { FunctionComponent, useState, useEffect } from 'react';
import { css } from 'emotion';
import { GrafanaTheme } from '@grafana/data';
import { AssetSummary } from '../../../queryResponseTypes';
import { styleMixins, stylesFactory, useTheme } from '@grafana/ui';
import { AssetInfo } from '../../../types';
import { SitewiseCache } from '../../../sitewiseCache';
import { AssetList } from './AssetList';

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    container: css`
      width: 100%;
      height: 60vh;
    `,
    listItem: css`
      ${styleMixins.listItem(theme)}
    `,
  };
});

export interface HierarchyInfo {
  name?: string;
  id?: string;
  description?: string;
}

// either must have children injected, or have asset + cache
export interface Props {
  asset?: AssetInfo | AssetSummary;
  hierarchy: HierarchyInfo;
  children?: AssetSummary[];
  cache?: SitewiseCache;
  search?: string;
  onSelect: (assetId: string) => void;
  onInspect: (assetId: string) => void;
}

export const AssetHierarchyList: FunctionComponent<Props> = ({
  asset,
  hierarchy,
  children,
  cache,
  search,
  onSelect,
  onInspect,
}) => {
  const [currentChildren, setChildren] = useState<AssetSummary[] | undefined>(children);

  const theme = useTheme();
  const style = getStyles(theme);

  useEffect(() => {
    // try to load children if none passed in
    if (!currentChildren && asset && cache) {
      const fetchData = async () => {
        const results = await cache.getAssociatedAssets(asset.id, hierarchy.id);
        setChildren(results.toArray());
      };
      fetchData();
    }
  }, [currentChildren, asset, cache, hierarchy.id]);

  return (
    <div key={hierarchy.id} className={style.container}>
      <AssetList
        search={search}
        assets={currentChildren}
        listInfo={{ id: hierarchy.id, description: hierarchy.id, name: hierarchy.name }}
        onSelect={onSelect}
        onInspect={onInspect}
      />
    </div>
  );
};
