import { DataFrameView, SelectableValue } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { ListAssetsQuery, ListAssociatedAssetsQuery, QueryType } from 'types';
import { AssetModelSummary, AssetSummary, DescribeAssetResult } from './queryResponseTypes';
import { AssetInfo, AssetPropertyInfo } from './types';
import { map } from 'rxjs/operators';
import { getTemplateSrv } from '@grafana/runtime';
import { useEffect, useState } from 'react';

/**
 * Keep a different cache for each region
 */
export class SitewiseCache {
  private models?: DataFrameView<AssetModelSummary>;
  private assetsById = new Map<string, AssetInfo>();
  private topLevelAssets?: DataFrameView<AssetSummary>;
  private assetPropertiesByAssetId = new Map<string, DataFrameView<{ id: string; name: string }>>();

  constructor(
    private ds: DataSource,
    private region: string
  ) {}

  async getAssetInfo(id: string): Promise<AssetInfo | undefined> {
    const v = this.assetsById.get(id);
    if (v) {
      return Promise.resolve(v);
    }

    return this.ds
      .runQuery(
        {
          refId: 'getAssetInfo',
          queryType: QueryType.DescribeAsset,
          assetId: id,
          region: this.region,
        },
        1000
      )
      .pipe(
        map((res) => {
          if (res.data.length) {
            const view = new DataFrameView<DescribeAssetResult>(res.data[0]);
            if (view && view.length) {
              const info = frameToAssetInfo(view.get(0));
              this.assetsById.set(id, info);
              return info;
            }
          }
          throw 'asset not found';
        })
      )
      .toPromise();
  }

  getAssetInfoSync(id: string): AssetInfo | undefined {
    const v = this.assetsById.get(id);
    if (v) {
      return v;
    }
    try {
      (async () => await this.getAssetInfo(id))();
    } catch {}
    return this.assetsById.get(id);
  }

  async listAssetProperties(assetId: string): Promise<DataFrameView<{ id: string; name: string }> | undefined> {
    const ap = this.assetPropertiesByAssetId.get(assetId);

    if (ap) {
      return ap;
    }

    return this.ds
      .runQuery({
        refId: 'listAssetProperties',
        queryType: QueryType.ListAssetProperties,
        assetId,
        region: this.region,
      })
      .pipe(
        map((res) => {
          if (res.data.length) {
            const assetProperties = new DataFrameView<{ id: string; name: string }>(res.data[0]);

            this.assetPropertiesByAssetId.set(assetId, assetProperties);

            return assetProperties;
          }

          throw 'asset properties not found';
        })
      )
      .toPromise();
  }

  async getModels(): Promise<DataFrameView<AssetModelSummary> | undefined> {
    if (this.models) {
      return Promise.resolve(this.models);
    }

    return this.ds
      .runQuery({
        refId: 'getModels',
        queryType: QueryType.ListAssetModels,
        region: this.region,
      })
      .pipe(
        map((res) => {
          if (res.data.length) {
            this.models = new DataFrameView<AssetModelSummary>(res.data[0]);
            return this.models;
          }
          throw 'no models found';
        })
      )
      .toPromise();
  }

  async getModelsOptions(): Promise<Array<SelectableValue<string>> | undefined> {
    const models = await this.getModels();
    if (!models) {
      return;
    }

    return models.toArray().map((model) => ({
      label: model.name,
      value: model.id,
      description: model.description,
    }));
  }

  // No cache for now
  async getAssetsOfType(modelId: string): Promise<DataFrameView<AssetSummary> | undefined> {
    const query: ListAssetsQuery = {
      refId: 'getAssetsOfType',
      queryType: QueryType.ListAssets,
      filter: 'ALL',
      modelId,
      region: this.region,
    };
    return this.ds
      .runQuery(query, 1000)
      .pipe(
        map((res) => {
          if (res.data.length) {
            this.topLevelAssets = new DataFrameView<AssetSummary>(res.data[0]);
            return this.topLevelAssets;
          }
          throw 'no assets found';
        })
      )
      .toPromise();
  }

  async getAssociatedAssets(assetId: string, hierarchyId?: string): Promise<DataFrameView<AssetSummary> | undefined> {
    const query: ListAssociatedAssetsQuery = {
      queryType: QueryType.ListAssociatedAssets,
      refId: 'associatedAssets',
      assetId: assetId,
      hierarchyId: hierarchyId,
      region: this.region,
    };

    return this.ds
      .runQuery(query, 1000)
      .pipe(
        map((res) => {
          if (res.data.length) {
            return new DataFrameView<AssetSummary>(res.data[0]);
          } else {
            throw 'no asset hierarchy found';
          }
        })
      )
      .toPromise();
  }

  async getTopLevelAssets(): Promise<DataFrameView<AssetSummary> | undefined> {
    if (this.topLevelAssets) {
      return Promise.resolve(this.topLevelAssets);
    }
    const query: ListAssetsQuery = {
      refId: 'topLevelAssets',
      queryType: QueryType.ListAssets,
      filter: 'TOP_LEVEL',
      region: this.region,
    };
    return this.ds
      .runQuery(query, 1000)
      .pipe(
        map((res) => {
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
    const options = getTemplateVariableOptions();
    try {
      const topLevel = (await this.getTopLevelAssets()) || [];
      for (const asset of topLevel) {
        options.push({
          label: asset.name,
          value: asset.id,
          description: asset.arn,
        });
      }
    } catch (err) {
      console.log('Error reading top level assets', err);
    }

    return options;
  }
}

export function frameToAssetInfo(res: DescribeAssetResult): AssetInfo {
  let properties: AssetPropertyInfo[] = [];
  let hierarchy: AssetPropertyInfo[] = [];

  try {
    properties = JSON.parse(res.properties);
    hierarchy = JSON.parse(res.hierarchies); // has Id, Name
  } catch (e) {
    console.log(res.properties, res.hierarchies);
    console.error('Error parsing JSON:', e);
    throw 'Could not parse returned JSON';
  }

  for (const p of properties) {
    p.value = p.Id;
    p.label = p.Name;

    if (p.Unit) {
      p.label += ' (' + p.Unit + ')';
    }

    if (p.DataType) {
      p.description = p.DataType;
      if (p.Alias) {
        p.description += ' // ' + p.Alias;
      }
    }
  }
  const options: AssetPropertyInfo[] = getTemplateSrv()
    .getVariables()
    .map((variable) => {
      const name = '${' + variable.name + '}';
      return {
        Id: name,
        Name: name,
        DataType: 'string',
        Unit: '',
        label: name,
        value: name,
        icon: 'arrow-right',
      };
    });

  const { hierarchies: _, ...rest } = res;

  return {
    ...rest,
    properties: [...options, ...properties],
    hierarchy: hierarchy.map((v) => ({
      label: v.Name,
      value: v.Id,
    })),
  };
}

export function assetSummaryToAssetInfo(res?: DataFrameView<AssetSummary>): AssetInfo[] {
  const results: AssetInfo[] = [];

  if (!res) {
    return results;
  }

  for (const info of res.toArray()) {
    const hierarchy: AssetPropertyInfo[] = JSON.parse(info.hierarchies); // has Id, Name
    const properties: AssetPropertyInfo[] = [];
    results.push({
      ...info,
      properties,
      hierarchy: hierarchy.map((v) => ({
        label: v.Name,
        value: v.Id,
      })),
    });
  }

  return results;
}

const getTemplateVariableOptions = (): Array<SelectableValue<string>> => {
  return getTemplateSrv()
    .getVariables()
    .map((variable) => ({
      label: '${' + (variable.label ?? variable.name) + '}',
      value: '${' + variable.name + '}',
      icon: 'arrow-right',
    }));
};

export const useModelsOptions = (cache: SitewiseCache): { isLoading: boolean; options: SelectableValue[] } => {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [options, setOptions] = useState<SelectableValue[]>([]);

  useEffect(() => {
    cache
      .getModelsOptions()
      .then((options) => {
        setIsLoading(false);
        setOptions(options || []);
      })
      .catch(() => {
        setIsLoading(false);
      });
  }, [cache]);

  return { isLoading, options };
};
