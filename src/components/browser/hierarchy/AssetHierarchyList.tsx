import React, { FunctionComponent, useState, useEffect } from 'react';
import { css } from 'emotion';
import { DataFrameView, GrafanaTheme } from '@grafana/data';
import { AssetSummary } from '../../../queryResponseTypes';
import { CollapsableSection, Label, Spinner, styleMixins, stylesFactory, useTheme } from '@grafana/ui';
import { AssetHierarchyNode } from './AssetHierarchyNode';
import { AssetInfo } from '../../../types';
import { SitewiseCache } from '../../../sitewiseCache';

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
}

export interface Props {
  asset?: AssetInfo | AssetSummary;
  hierarchy: HierarchyInfo;
  children?: DataFrameView<AssetSummary>;
  cache?: SitewiseCache;
  onSelect: (assetId: string) => void;
  onInspect: (assetId: string) => void;
}

const hierarchyLabel = (info: HierarchyInfo) => {
  return <Label description={info.id}>{info.name}</Label>;
};

export const AssetHierarchy: FunctionComponent<Props> = ({
  asset,
  hierarchy,
  children,
  cache,
  onSelect,
  onInspect,
}) => {
  const [currentChildren, setChildren] = useState<DataFrameView<AssetSummary> | undefined>(children);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const theme = useTheme();
  const style = getStyles(theme);

  const label = (hierarchyLabel(hierarchy) as unknown) as string;

  const renderChildren = () => {
    return currentChildren?.toArray().map(c => {
      return (
        <li key={c.name} className={style.listItem}>
          <AssetHierarchyNode asset={c} onInspect={onInspect} onSelect={onSelect} />
        </li>
      );
    });
  };

  useEffect(() => {
    // try to load children is none passed in
    if (!currentChildren && asset && cache) {
      setIsLoading(true);
      const fetchData = async () => {
        const results = await cache.getAssociatedAssets(asset.id, hierarchy.id);
        setChildren(results);
        setIsLoading(false);
      };
      fetchData();
    }
  }, [currentChildren, asset, cache, hierarchy.id]);

  return (
    <div className={style.container}>
      <CollapsableSection label={label} isOpen={false}>
        {isLoading ? (
          <div>
            <Spinner /> Loading children...{' '}
          </div>
        ) : (
          <ul>{renderChildren()}</ul>
        )}
      </CollapsableSection>
    </div>
  );
};

// export class AssetHierarchyList extends PureComponent<Props> {
//   renderChildren = () => {
//     return this.props.children?.map(c => {
//       return(
//         <li key={c.name}>
//           <AssetHierarchyNode asset={c} />
//         </li>
//       );
//
//     });
//   };
//
//   render() {
//     const { hierarchy } = this.props;
//
//     const label = hierarchyLabel(hierarchy) as unknown as string;
//
//     return (
//       <div style={{ height: '60vh' }}>
//         <CollapsableSection label={label} isOpen={false}>
//           <ul>
//             {this.renderChildren()}
//           </ul>
//         </CollapsableSection>
//       </div>
//     );
//   }
// }
