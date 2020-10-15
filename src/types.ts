import { DataQuery } from '@grafana/data';
import { AwsDataSourceJsonData, AwsDataSourceSecureJsonData } from 'common/types';

// Matches https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/query.go#L3
export enum QueryType {
  ListAssetModels = 'ListAssetModels',
  ListAssets = 'ListAssets',
  PropertyValue = 'PropertyValue',
  PropertyValueHistory = 'PropertyValueHistory',
  PropertyAggregate = 'PropertyAggregate',
}

export enum SiteWiseQualities {
  GOOD = 'GOOD',
  BAD = 'BAD',
  UNCERTAIN = 'UNCERTAIN',
}

export enum SiteWiseResolution {
  Min = '1m',
  Hour = '1h',
  Day = '1d',
  Auto = 'AUTO', // or missing!
}

export enum AggregateTypes {
  AVERAGE = 'AVERAGE',
  COUNT = 'COUNT',
  MAXIMUM = 'MAXIMUM',
  MINIMUM = 'MINIMUM',
  SUM = 'SUM',
  STANDARD_DEVIATION = 'STANDARD_DEVIATION',
}

export interface SitewiseQuery extends DataQuery {
  queryType: QueryType;
  region?: string; // aws region string
}

export interface SitewiseQueryWithPages {
  /**
   * The next token should never be saved in the JSON model, however some queries
   * will require multiple pages in order to fulfil the requests
   */
  nextToken?: string;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface ListAssetModelsQuery extends SitewiseQuery, SitewiseQueryWithPages {
  queryType: QueryType.ListAssetModels;
}

export function isListAssetModelsQuery(q?: SitewiseQuery): q is ListAssetModelsQuery {
  return q?.queryType === QueryType.ListAssetModels;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface ListAssetsQuery extends SitewiseQuery, SitewiseQueryWithPages {
  queryType: QueryType.ListAssets;
  assetModelId: string;
  filter: 'ALL' | 'TOP_LEVEL';
}

export function isListAssetsQuery(q?: SitewiseQuery): q is ListAssetsQuery {
  return q?.queryType === QueryType.ListAssets;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValue.html}
 * {@link https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/property.go#L15}
 */
export interface AssetPropertyValueQuery extends SitewiseQuery {
  queryType: QueryType.PropertyValue;

  assetId: string;
  propertyId: string;
  // NOTE: 'propertyAlias' is not supported, but the UI should be able to show aliases
}

export function isAssetPropertyValueQuery(q?: SitewiseQuery): q is AssetPropertyValueQuery {
  return q?.queryType === QueryType.PropertyAggregate;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValueHistory.html}
 */
export interface AssetPropertyValueHistoryQuery extends SitewiseQuery, SitewiseQueryWithPages {
  queryType: QueryType.PropertyValueHistory;

  assetId: string;
  propertyId: string;
  qualities?: SiteWiseQualities[]; // Docs say "Fixed number of 1 item.????" does that mean only one?
  timeOrdering?: 'ASCENDING' | 'DESCENDING';
}

export function isAssetPropertyValueHistoryQuery(q?: SitewiseQuery): q is AssetPropertyValueHistoryQuery {
  return q?.queryType === QueryType.PropertyValueHistory;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html}
 */
export interface AssetPropertyAggregatesQuery extends SitewiseQuery, SitewiseQueryWithPages {
  queryType: QueryType.PropertyAggregate;

  assetId: string;
  propertyId: string;
  resolution?: SiteWiseResolution;
  aggregateTypes: AggregateTypes[]; // at least one
  qualities?: SiteWiseQualities[];
  timeOrdering?: 'ASCENDING' | 'DESCENDING';
}

export function isAssetPropertyAggregatesQuery(q?: SitewiseQuery): q is AssetPropertyAggregatesQuery {
  return q?.queryType === QueryType.PropertyAggregate;
}

/**
 * Metadata attached to DataFrame results
 */
export interface SitewiseCustomMeta {
  queryId: string;
  nextToken?: string;
  hasSeries?: boolean;

  executionStartTime?: number; // The backend clock
  executionFinishTime?: number; // The backend clock

  fetchStartTime?: number; // The frontend clock
  fetchEndTime?: number; // The frontend clock
  fetchTime?: number; // The frontend clock

  // when multiple queries exist we keep track of each request
  subs?: SitewiseCustomMeta[];
}

/**
 * Global datasource options
 */
export interface SitewiseOptions extends AwsDataSourceJsonData {
  // nothing for now
}

export interface SitewiseSecureJsonData extends AwsDataSourceSecureJsonData {
  // nothing for now
}
