// TODO? this file should be auto-generated!

// Mapped from DataFrame result
export interface AssetModelSummary {
  name: string; // string
  description: string; // string
  id: string; // string
  arn: string; // string
  error?: string; // *string
  state: string; // string
  creation_date: number; // time.Time
  last_update: number; // time.Time
}

// Mapped from DataFrame result
export interface AssetSummary {
  name: string; // string
  id: string; // string
  model_id: string; // string
  arn: string; // string
  creation_date: number; // time.Time
  last_update: number; // time.Time
  state: string; // string
  error?: string; // *string
  hierarchies: string; // string
}

// Mapped from DataFrame result
export interface DescribeAssetResult {
  name: string; // string
  id: string; // string
  arn: string; // string
  model_id: string; // string
  state: string; // string
  error?: string; // *string
  creation_date: number; // time.Time
  last_update: number; // time.Time
  hierarchies: string; // string
  properties: string; // string
}

// Processed form DescribeAssetResult
export interface AssetInfo {
  name: string; // string
  id: string; // string
  arn: string; // string
  model_id: string;
  properties: AssetPropertyInfo[];
}

// Mapped from DataFrame result
export interface AssetPropertyInfo {
  Alias?: string;
  DataType: string;
  Id: string;
  Name: string;
  Unit: string;
}
