import type { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';
import type { SelectableValue } from '@grafana/data';
import { QueryEditorMode } from '@grafana/plugin-ui';
import type { DataQuery } from '@grafana/schema';
import type { Region } from './regions';
import { SitewiseQueryState } from 'components/query/sql-query-builder/types';

// Matches https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/query.go#L3
export enum QueryType {
  ListAssetModels = 'ListAssetModels',
  ListAssets = 'ListAssets',
  ListAssociatedAssets = 'ListAssociatedAssets',
  ListAssetProperties = 'ListAssetProperties',
  DescribeAsset = 'DescribeAsset',
  PropertyValue = 'PropertyValue',
  PropertyValueHistory = 'PropertyValueHistory',
  PropertyAggregate = 'PropertyAggregate',
  PropertyInterpolated = 'PropertyInterpolated',
  ListTimeSeries = 'ListTimeSeries',
  ExecuteQuery = 'ExecuteQuery',
}

export enum SiteWiseQuality {
  ANY = 'ANY',
  GOOD = 'GOOD',
  BAD = 'BAD',
  UNCERTAIN = 'UNCERTAIN',
}

export enum SiteWiseResponseFormat {
  Table = 'table',
  TimeSeries = 'timeseries',
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
  FifteenMin = '15m',
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
  region?: Region; // aws region string
  responseFormat?: SiteWiseResponseFormat;

  // QueryEditor
  editorMode?: QueryEditorMode;

  // RawQueryEditor
  rawSQL?: string;

  // SQL Query Builder
  sqlQueryState?: SitewiseQueryState;

  /** @deprecated -- this is migrated to assetIds */
  assetId?: string;
  // One or more assets to filter -- when multiple, they should share the same properties, the batch API will be called
  assetIds?: string[];
  /** @deprecated -- this is migrated to propertyIds */
  propertyId?: string;
  // One or more properties to fetch data
  propertyIds?: string[];
  /** @deprecated -- this is migrated to propertyAlias */
  propertyAlias?: string;
  // One or more properties to fetch data
  propertyAliases?: string[];
  quality?: SiteWiseQuality;
  resolution?: SiteWiseResolution;
  lastObservation?: boolean;
  flattenL4e?: boolean;
  maxPageAggregations?: number;
  clientCache?: boolean;
}

export interface SitewiseNextQuery extends SitewiseQuery {
  /**
   * The next token should never be saved in the JSON model, however some queries
   * will require multiple pages in order to fulfil the requests
   */
  nextToken?: string;
  nextTokens?: Record<string, string>;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListAssetModels.html}
 */
export interface ListAssetModelsQuery extends SitewiseQuery {
  queryType: QueryType.ListAssetModels;
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
  loadAllChildren?: boolean; // When passed, we will loop through all associated hierarchies, and return children from all.
  hierarchyId?: string; // if empty and loadAllChildren is false, will list the parents
}

export function isListAssociatedAssetsQuery(q?: SitewiseQuery): q is ListAssociatedAssetsQuery {
  return q?.queryType === QueryType.ListAssociatedAssets;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValue.html}
 * {@link https://github.com/grafana/iot-sitewise-datasource/blob/main/pkg/models/property.go#L15}
 */
export interface AssetPropertyValueQuery extends SitewiseQuery {
  queryType: QueryType.PropertyValue;

  flattenL4e?: boolean;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyValueHistory.html}
 */
export interface AssetPropertyValueHistoryQuery extends SitewiseQuery {
  queryType: QueryType.PropertyValueHistory;

  timeOrdering?: SiteWiseTimeOrder;
  flattenL4e?: boolean;
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
}

export function isAssetPropertyInterpolatedQuery(q?: SitewiseQuery): q is AssetPropertyInterpolatedQuery {
  return q?.queryType === QueryType.PropertyInterpolated;
}

/**
 * {@link https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_ListTimeSeries.html}
 */

export interface ListTimeSeriesQuery extends SitewiseQuery {
  queryType: QueryType.ListTimeSeries;
  aliasPrefix?: string;
  assetId?: string;
  timeSeriesType?: 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL';
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
  return queryType === QueryType.PropertyAggregate || queryType === QueryType.PropertyValueHistory;
}

export function shouldShowOptionsRow(query: SitewiseQuery, showProp: boolean): boolean {
  const shouldShowLastObservedSwitch =
    shouldShowLastObserved(query.queryType) && !query.propertyAliases?.length && showProp;
  const shouldShowWithPropertyAlias =
    // shouldn't show the row when querying associated assets with property alias, otherwise show it every time property alias is set
    query.propertyAliases?.length && !isListAssociatedAssetsQuery(query);
  return !!(query.propertyIds?.length || shouldShowWithPropertyAlias || shouldShowLastObservedSwitch);
}

export function shouldShowL4eOptions(queryType?: QueryType): boolean {
  return queryType === QueryType.PropertyValue || queryType === QueryType.PropertyValueHistory;
}

export function shouldShowQualityAndOrderComponent(queryType?: QueryType): boolean {
  return queryType !== QueryType.PropertyValue;
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
export interface SitewiseCustomMetadata {
  nextToken?: string;
  entryId?: string;
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
