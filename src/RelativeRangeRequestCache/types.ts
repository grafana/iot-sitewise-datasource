import { DataFrame } from '@grafana/data';
import { AssetPropertyAggregatesQuery, AssetPropertyValueHistoryQuery, ListAssetsQuery, ListAssociatedAssetsQuery, QueryType, SitewiseQuery } from 'types';

export const TIME_SERIES_QUERY_TYPES = new Set<QueryType>([
  QueryType.PropertyAggregate,
  QueryType.PropertyInterpolated,
  QueryType.PropertyValue,
  QueryType.PropertyValueHistory,
]);

export interface CachedQueryInfo {
  query: Pick<SitewiseQueriesUnion, 'queryType' | 'timeOrdering' | 'lastObservation'>;
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
