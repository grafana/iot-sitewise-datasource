import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';
import {
  AggregateType,
  AssetProperty,
  Property,
  DescribeAssetResponse,
  Resolution,
  AggregatedValue,
  AssetModelSummary,
} from 'aws-sdk/clients/iotsitewise';
import { Dispatch, SetStateAction } from 'react';

export type Setter<T> = Dispatch<SetStateAction<T>>;

export type SitewiseAsset = DescribeAssetResponse;
export type SitewiseProperty = Property;
export type SitewiseAggregateValue = AggregatedValue;
export type SitewiseModelSummary = AssetModelSummary;

export enum SitewiseQueryType {
  Default = 'Default',
  Aggregate = 'Aggregate',
  PropertyHistory = 'PropertyHistory',
  DisableStreaming = 'DisableStreaming',
}

export interface SitewiseQuery extends DataQuery {
  // The Sitewise asset that owns the property to-query
  // In practice, this should be a truncated plugin-internal model
  asset: SitewiseAsset;
  // The Sitewise property to query
  // In practice, this should be a truncated plugin-internal model
  property: AssetProperty;
  // Optional AWS region
  // The data source should be able to aggregate data streams from multiple AWS regions,
  // at the user discretion
  // The default AWS region is part of the data source config
  region?: string;
  // Configures whether or not the data stream should update periodically
  streaming: boolean;
  // The interval in which to fetch new data
  streamingInterval: number;
  // Sitewise aggregate type
  // See: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html
  aggregateTypes: SitewiseAggregateType[];
  // Sitewise data resolution
  // See: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_GetAssetPropertyAggregates.html
  dataResolution: Resolution;

  models: Array<SelectableValue<SitewiseModelSummary>>;
  model: SitewiseModelSummary;

  // Overriding with Sitewise specific flavor
  queryType: SitewiseQueryType;
}

export type SitewiseAggregateType = AggregateType | 'NONE' | string;
export type SitewiseResolution = '1m' | '1h' | '1d' | 'ALL' | string;

/**
 * These are options configured for each DataSource instance
 */
export interface SitewiseDataSourceOptions extends DataSourceJsonData {
  defaultRegion: string;
  accessKeyId: string;
  secretAccessKey: string;
  sessionToken: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {}
