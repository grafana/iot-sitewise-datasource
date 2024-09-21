import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface SitewiseRawQuery extends DataQuery {
  queryText?: string;
}

export const DEFAULT_QUERY: Partial<SitewiseRawQuery> = {};

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  data_points: DataPoint[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface RawDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
