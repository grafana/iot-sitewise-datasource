import { DataFrame, DataFrameView, SelectableValue } from '@grafana/data';
import { DataSource } from 'DataSource';
import { ListAssetsQuery, QueryType } from 'types';
import { AssetModelSummary, AssetInfo, AssetSummary } from './queryResponseTypes';
import { map } from 'rxjs/operators';

export class SitewiseCache {
  private models?: DataFrameView<AssetModelSummary>;
  private assetsById = new Map<string, AssetInfo>();
  private topLevelAssets?: DataFrameView<AssetSummary>;

  constructor(private ds: DataSource) {
    console.log('DS', ds);
  }

  async getAssetInfo(id: string): Promise<AssetInfo> {
    const v = this.assetsById.get(id);
    if (v) {
      return Promise.resolve(v);
    }

    return this.ds
      .runQuery({
        refId: 'getAssetInfo',
        queryType: QueryType.DescribeAsset,
        assetId: id,
      })
      .pipe(
        map(res => {
          if (res.data.length) {
            const info = frameToAssetInfo(res.data[0]);
            this.assetsById.set(id, info);
            return info;
          }
          throw 'asset not found';
        })
      )
      .toPromise();
  }

  async getModels(): Promise<DataFrameView<AssetModelSummary>> {
    if (this.models) {
      return Promise.resolve(this.models);
    }

    return this.ds
      .runQuery({
        refId: 'getModels',
        queryType: QueryType.ListAssetModels,
      })
      .pipe(
        map(res => {
          if (res.data.length) {
            this.models = new DataFrameView<AssetModelSummary>(res.data[0]);
            return this.models;
          }
          throw 'no models found';
        })
      )
      .toPromise();
  }

  async getTopLevelAssets(): Promise<DataFrameView<AssetSummary>> {
    if (this.topLevelAssets) {
      return Promise.resolve(this.topLevelAssets);
    }
    const query: ListAssetsQuery = {
      refId: 'topLevelAssets',
      queryType: QueryType.ListAssets,
      filter: 'TOP_LEVEL',
    };
    return this.ds
      .runQuery(query)
      .pipe(
        map(res => {
          if (res.data.length) {
            this.topLevelAssets = new DataFrameView<AssetSummary>(res.data[0]);
            return this.topLevelAssets;
          }
          throw 'no assets found';
        })
      )
      .toPromise();
  }

  async getAssetPickerOptions(): Promise<Array<SelectableValue<string>>> {
    const options: Array<SelectableValue<string>> = [];
    const topLevel = await this.getTopLevelAssets();
    for (const asset of topLevel) {
      options.push({
        label: asset.name,
        description: asset.arn,
      });
    }

    // Also add recent values
    for (const asset of this.assetsById.values()) {
      options.push({
        label: asset.name,
        description: asset.id,
      });
    }
    return options;
  }
}

export function frameToAssetInfo(frame: DataFrame): AssetInfo {
  console.log('TODO', frame);
  return {} as AssetInfo;
}
