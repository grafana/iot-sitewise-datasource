import { DataFrame } from '@grafana/data';
import { AssetPropertyAggregatesQuery, AssetPropertyValueHistoryQuery, ListAssetsQuery, ListAssociatedAssetsQuery, QueryType, SitewiseQuery } from 'types';

const TIME_SERIES_QUERY_TYPES = new Set<QueryType>([
  QueryType.PropertyAggregate,
  QueryType.PropertyInterpolated,
  QueryType.PropertyValue,
  QueryType.PropertyValueHistory,
]);

export function isTimeSeriesQueryType(queryType: QueryType) {
  return TIME_SERIES_QUERY_TYPES.has(queryType);
}

const TIME_ORDERING_QUERY_TYPES = new Set<QueryType>([
  QueryType.PropertyAggregate,
  QueryType.PropertyValueHistory,
]);

export function isTimeOrderingQueryType(queryType: QueryType) {
  return TIME_ORDERING_QUERY_TYPES.has(queryType);
}

export interface CachedQueryInfo {
  query: SitewiseQueriesUnion;
  dataFrame: DataFrame;
}

// Union of all SiteWise queries variants
export type SitewiseQueriesUnion = SitewiseQuery
  & Partial<Pick<AssetPropertyAggregatesQuery, 'aggregates'>>
  & Partial<Pick<AssetPropertyValueHistoryQuery, 'timeOrdering'>>
  & Partial<Pick<ListAssociatedAssetsQuery, 'loadAllChildren'>>
  & Partial<Pick<ListAssociatedAssetsQuery, 'hierarchyId'>>
  & Partial<Pick<ListAssetsQuery, 'modelId'>>
  & Partial<Pick<ListAssetsQuery, 'filter'>>;
