import { DataQuery } from '@grafana/data';
import { AwsDataSourceJsonData, AwsDataSourceSecureJsonData } from 'common/types';

export enum QueryType {
  Builder = 'builder',
  Samples = 'samples',
  Raw = 'raw',
}

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

export interface SitewiseQuery extends DataQuery {
  // When specified, use this rather than the default for macros
  database?: string;
  table?: string;
  measure?: string;

  // The rendered query
  rawQuery?: string;

  // Not a real parameter...
  // nextToken?: string;
}

export interface SitewiseOptions extends AwsDataSourceJsonData {
  defaultDatabase?: string;
  defaultTable?: string;
  defaultMeasure?: string;
}

export interface SitewiseSecureJsonData extends AwsDataSourceSecureJsonData {
  // nothing for now
}
