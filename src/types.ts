import { DataQuery, SelectableValue } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

// Matches https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/query.go#L3
export enum QueryType {
  ListAssetModels = 'ListAssetModels',
  ListAssets = 'ListAssets',
  ListAssociatedAssets = 'ListAssociatedAssets',
  DescribeAsset = 'DescribeAsset',
  PropertyValue = 'PropertyValue',
  PropertyValueHistory = 'PropertyValueHistory',
  PropertyAggregate = 'PropertyAggregate',
  PropertyInterpolated = 'PropertyInterpolated',
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
  Sec = '1s',
  TenSec = '10s',
  Min = '1m',
  TenMin = '10m',
  Hour = '1h',
  TenHour = '10h',
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
  propertyAlias?: string;
  quality?: SiteWiseQuality;
  resolution?: SiteWiseResolution;
  lastObservation?: boolean;
  maxPageAggregations?: number;
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

// https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssociatedAssets.html

export interface ListAssociatedAssetsQuery extends SitewiseQuery {
  queryType: QueryType.ListAssociatedAssets;
  hierarchyId?: string; // if empty, will list the parents
}

export function isListAssociatedAssetsQuery(q?: SitewiseQuery): q is ListAssociatedAssetsQuery {
  return q?.queryType === QueryType.ListAssociatedAssets;
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

  timeOrdering?: SiteWiseTimeOrder;
}

export function isAssetPropertyAggregatesQuery(q?: SitewiseQuery): q is AssetPropertyAggregatesQuery {
  return q?.queryType === QueryType.PropertyAggregate;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetInterpolatedAssetPropertyValues.html}
 */
export interface AssetPropertyInterpolatedQuery extends SitewiseQuery {
  queryType: QueryType.PropertyInterpolated;
  timeOrdering?: SiteWiseTimeOrder;
}

export function isAssetPropertyInterpolatedQuery(q?: SitewiseQuery): q is AssetPropertyInterpolatedQuery {
  return q?.queryType === QueryType.PropertyInterpolated;
}

export function isPropertyQueryType(queryType?: QueryType): boolean {
  return (
    queryType === QueryType.PropertyAggregate ||
    queryType === QueryType.PropertyValue ||
    queryType === QueryType.PropertyValueHistory ||
    queryType === QueryType.PropertyInterpolated
  );
}

export function shouldShowLastObserved(queryType?: QueryType): boolean {
  return (
    queryType === QueryType.PropertyAggregate ||
    queryType === QueryType.PropertyValueHistory
  );
}

// matches native sitewise API with capitals
export interface AssetPropertyInfo extends SelectableValue<string> {
  Id: string;
  Name: string;
  Alias?: string;
  DataType: string;
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
  hierarchy: Array<SelectableValue<string>>; // Id is value
}

/**
 * Metadata attached to DataFrame results
 */
export interface SitewiseCustomMeta {
  nextToken?: string;

  resolution?: string;

  aggregates?: string[];
}

/**
 * Global datasource options
 */
export interface SitewiseOptions extends AwsAuthDataSourceJsonData {
  // nothing for now
  edgeAuthMode?: string;
  edgeAuthUser?: string;
}

export interface SitewiseSecureJsonData extends AwsAuthDataSourceSecureJsonData {
  // nothing for now
  edgeAuthPass?: string;
  cert?: string;
}
