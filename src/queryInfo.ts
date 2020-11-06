import { SelectableValue } from '@grafana/data';
import {
  QueryType,
  SitewiseQuery,
  AggregateType,
  AssetPropertyAggregatesQuery,
  ListAssetsQuery,
  ListAssetModelsQuery,
  AssetPropertyValueQuery,
  AssetPropertyValueHistoryQuery,
  SiteWiseResolution,
  AssetInfo,
  AssetPropertyInfo,
} from './types';

export interface QueryTypeInfo extends SelectableValue<QueryType> {
  value: QueryType; // not optional
  defaultQuery: Partial<SitewiseQuery>;
  helpURL: string;
}

export const siteWisteQueryTypes: QueryTypeInfo[] = [
  {
    label: 'Get property value aggregates',
    value: QueryType.PropertyAggregate,
    description: `Gets aggregated values for an asset property.`,
    defaultQuery: {
      resolution: SiteWiseResolution.Auto,
      aggregates: [AggregateType.AVERAGE],
      timeOrdering: 'ASCENDING',
    } as AssetPropertyAggregatesQuery,
    helpURL: 'https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html',
  },
  {
    label: 'Get property value history',
    value: QueryType.PropertyValueHistory,
    description: `Gets the history of an asset property's value.`,
    defaultQuery: {
      timeOrdering: 'ASCENDING',
    } as AssetPropertyValueHistoryQuery,
    helpURL: 'https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html',
  },
  {
    label: 'Get property value',
    value: QueryType.PropertyValue,
    description: `Gets an asset property's current value.`,
    defaultQuery: {} as AssetPropertyValueQuery,
    helpURL: 'https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html',
  },
  {
    label: 'List assets',
    value: QueryType.ListAssets,
    description: 'Retrieves a paginated list of asset summaries.',
    defaultQuery: {
      filter: 'TOP_LEVEL',
    } as ListAssetsQuery,
    helpURL: 'https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html',
  },
  {
    label: 'List asset models',
    value: QueryType.ListAssetModels,
    description: 'Retrieves this list of all asset models',
    defaultQuery: {} as ListAssetModelsQuery,
    helpURL: 'https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html',
  },
];

export function changeQueryType(q: SitewiseQuery, info: QueryTypeInfo): SitewiseQuery {
  if (q.queryType === info.value) {
    return q; // no change;
  }
  const copy = {
    ...info.defaultQuery,
    ...q,
    queryType: info.value,
  };

  // TODO: for each query type, remove the unused fields

  console.log('CHANGE', q, copy);

  return copy;
}

export function getAssetProperty(asset?: AssetInfo, propId?: string): AssetPropertyInfo | undefined {
  if (!asset?.properties || !propId) {
    return undefined;
  }
  return asset.properties.find(p => p.Id === propId);
}

export function getDefaultAggregate(prop?: AssetPropertyInfo): AggregateType {
  if (prop?.DataType === 'STRING') {
    return AggregateType.COUNT;
  }
  return AggregateType.AVERAGE;
}
