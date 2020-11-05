import { DataQuery, SelectableValue } from '@grafana/data';
import { AwsDataSourceJsonData, AwsDataSourceSecureJsonData } from 'common/types';

// Matches https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/query.go#L3
export enum QueryType {
  ListAssetModels = 'ListAssetModels',
  ListAssets = 'ListAssets',
  DescribeAsset = 'DescribeAsset',
  PropertyValue = 'PropertyValue',
  PropertyValueHistory = 'PropertyValueHistory',
  PropertyAggregate = 'PropertyAggregate',
}

export enum SiteWiseQuality {
  ANY = 'ANY',
  GOOD = 'GOOD',
  BAD = 'BAD',
  UNCERTAIN = 'UNCERTAIN',
}

export enum SiteWiseTimeOrder {
  ASCENDING = 'ASCENDING',
  DESCENDING = 'DESCENDING',
}

export enum SiteWiseResolution {
  Auto = 'AUTO', // or missing!
  Min = '1m',
  Hour = '1h',
  Day = '1d',
}

export enum AggregateType {
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

  // Although these are not required everywhere, many queries use them
  assetId?: string;
  propertyId?: string;
}

export interface SitewiseNextQuery extends SitewiseQuery {
  /**
   * The next token should never be saved in the JSON model, however some queries
   * will require multiple pages in order to fulfil the requests
   */
  nextToken?: string;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface ListAssetModelsQuery extends SitewiseQuery {
  queryType: QueryType.ListAssetModels;
}

export function isListAssetModelsQuery(q?: SitewiseQuery): q is ListAssetModelsQuery {
  return q?.queryType === QueryType.ListAssetModels;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface ListAssetsQuery extends SitewiseQuery {
  queryType: QueryType.ListAssets;
  modelId?: string;
  filter: 'ALL' | 'TOP_LEVEL';
}

export function isListAssetsQuery(q?: SitewiseQuery): q is ListAssetsQuery {
  return q?.queryType === QueryType.ListAssets;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface DescribeAssetQuery extends SitewiseQuery {
  queryType: QueryType.DescribeAsset;
}

export function isDescribeAssetQuery(q?: SitewiseQuery): q is ListAssetModelsQuery {
  return q?.queryType === QueryType.DescribeAsset;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValue.html}
 * {@link https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/property.go#L15}
 */
export interface AssetPropertyValueQuery extends SitewiseQuery {
  queryType: QueryType.PropertyValue;
}

export function isAssetPropertyValueQuery(q?: SitewiseQuery): q is AssetPropertyValueQuery {
  return q?.queryType === QueryType.PropertyAggregate;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValueHistory.html}
 */
export interface AssetPropertyValueHistoryQuery extends SitewiseQuery {
  queryType: QueryType.PropertyValueHistory;

  quality?: SiteWiseQuality;
  timeOrdering?: SiteWiseTimeOrder;
}

export function isAssetPropertyValueHistoryQuery(q?: SitewiseQuery): q is AssetPropertyValueHistoryQuery {
  return q?.queryType === QueryType.PropertyValueHistory;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html}
 */
export interface AssetPropertyAggregatesQuery extends SitewiseQuery {
  queryType: QueryType.PropertyAggregate;

  resolution?: SiteWiseResolution;
  aggregates: AggregateType[]; // at least one

  quality?: SiteWiseQuality;
  timeOrdering?: SiteWiseTimeOrder;
}

export function isAssetPropertyAggregatesQuery(q?: SitewiseQuery): q is AssetPropertyAggregatesQuery {
  return q?.queryType === QueryType.PropertyAggregate;
}

export function isPropertyQueryType(queryType?: QueryType): boolean {
  return (
    queryType === QueryType.PropertyAggregate ||
    queryType === QueryType.PropertyValue ||
    queryType === QueryType.PropertyValueHistory
  );
}

// matches native sitewise API with capitals
export interface AssetPropertyInfo extends SelectableValue<string> {
  Alias?: string;
  DataType: string;
  Id: string;
  Name: string;
  Unit: string;

  // Filled in for selectable values
  value: string;
  label: string;
}

// Processed form DescribeAssetResult frame
export interface AssetInfo {
  name: string; // string
  id: string; // string
  arn: string; // string
  model_id: string;
  properties: AssetPropertyInfo[];
}

/**
 * Metadata attached to DataFrame results
 */
export interface SitewiseCustomMeta {
  nextToken?: string;

  // Show the aggregate value actually used
  resolution?: string;
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
