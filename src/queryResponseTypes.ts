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

// Processed form
export interface AssetInfo {
  name: string; // string
  id: string; // string
  modelId: string;
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
