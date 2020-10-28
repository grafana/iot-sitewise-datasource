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
} from './types';

export interface QueryTypeInfo extends SelectableValue<QueryType> {
  value: QueryType; // not optional
  defaultQuery: Partial<SitewiseQuery>;
}

export const siteWisteQueryTypes: QueryTypeInfo[] = [
  {
    label: 'Get property value aggregates',
    value: QueryType.PropertyAggregate,
    description: `Gets aggregated values for an asset property.`,
    defaultQuery: {
      resolution: SiteWiseResolution.Auto,
      aggregates: [AggregateType.AVERAGE],
    } as AssetPropertyAggregatesQuery,
  },
  {
    label: 'Get property value history',
    value: QueryType.PropertyValueHistory,
    description: `Gets the history of an asset property's value.`,
    defaultQuery: {} as AssetPropertyValueHistoryQuery,
  },
  {
    label: 'Get property value',
    value: QueryType.PropertyValue,
    description: `Gets an asset property's current value.`,
    defaultQuery: {} as AssetPropertyValueQuery,
  },
  {
    label: 'List assets',
    value: QueryType.ListAssets,
    description: 'Retrieves a paginated list of asset summaries.',
    defaultQuery: {
      filter: 'TOP_LEVEL',
    } as ListAssetsQuery,
  },
  {
    label: 'List asset models',
    value: QueryType.ListAssetModels,
    description: 'Retrieves this list of all asset models',
    defaultQuery: {} as ListAssetModelsQuery,
    keys: [],
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
