import { SelectableValue } from '@grafana/data';

export interface SqlQueryBuilderProps {
  builderState: any;
  onChange: (query: any) => void;
}

export interface SelectField {
  column: string;
  aggregation?: string;
  alias?: string;
  functionArg?: string;
  functionArgValue?: string;
  functionArgValue2?: string;
}

export interface WhereCondition {
  column: string;
  operator: string;
  value: string;
  value2?: string; // For BETWEEN operator
  logicalOperator?: 'AND' | 'OR';
  operator2?: string; // For BETWEEN operator
}

export interface OrderByField {
  column: string;
  direction: 'ASC' | 'DESC';
}

export type AggregationFunction = 'COUNT' | 'SUM' | 'AVG' | 'MAX' | 'MIN' | '';

export interface HavingCondition {
  aggregation: AggregationFunction;
  column: string;
  operator: '=' | '!=' | '>' | '<' | '>=' | '<=';
  value: string;
  logicalOperator?: 'AND' | 'OR';
}

export interface SitewiseQueryState {
  selectedAssetModel?: string;
  selectFields: SelectField[];
  whereConditions: WhereCondition[];
  groupByTags: string[];
  orderByFields: OrderByField[];
  havingConditions: HavingCondition[];
  limit?: number;
  rawSQL: string;
}

export const defaultSitewiseQueryState: SitewiseQueryState = {
  selectedAssetModel: '',
  selectFields: [{ column: '', aggregation: '', alias: '' }],
  whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
  groupByTags: [],
  orderByFields: [{ column: '', direction: 'ASC' }],
  havingConditions: [{ aggregation: 'COUNT', column: '', operator: '=', value: '', logicalOperator: 'AND' }],
  limit: 100,
  rawSQL: '',
};

export interface AssetProperty {
  id: string;
  name: string;
  dataType: 'STRING' | 'INTEGER' | 'DOUBLE' | 'BOOLEAN' | 'TIMESTAMP';
  alias?: string;
}

export interface AssetModel {
  id: string;
  name: string;
  properties: AssetProperty[];
}

export const queryReferenceViews: AssetModel[] = [
  {
    id: 'asset',
    name: 'asset',
    properties: [
      { id: 'asset_id', name: 'asset_id', dataType: 'DOUBLE' },
      { id: 'asset_name', name: 'asset_name', dataType: 'STRING' },
      { id: 'asset_description', name: 'asset_description', dataType: 'STRING' },
      { id: 'asset_model_id', name: 'asset_model_id', dataType: 'STRING' },
      { id: 'parent_asset_id', name: 'parent_asset_id', dataType: 'STRING' },
      { id: 'asset_external_id', name: 'asset_external_id', dataType: 'STRING' },
      { id: 'asset_model_external_id', name: 'asset_model_external_id', dataType: 'STRING' },
    ],
  },
  {
    id: 'asset_property',
    name: 'asset_property',
    properties: [
      { id: 'asset_id', name: 'asset_id', dataType: 'STRING' },
      { id: 'property_id', name: 'property_id', dataType: 'STRING' },
      { id: 'property_name', name: 'property_name', dataType: 'STRING' },
      { id: 'property_alias', name: 'property_alias', dataType: 'STRING' },
      { id: 'property_external_id', name: 'property_external_id', dataType: 'STRING' },
      { id: 'asset_composite_model_id', name: 'asset_composite_model_id', dataType: 'STRING' },
      { id: 'property_type', name: 'property_type', dataType: 'STRING' },
      { id: 'property_data_type', name: 'property_data_type', dataType: 'STRING' },
      { id: 'int_attribute_value', name: 'int_attribute_value', dataType: 'INTEGER' },
      { id: 'double_attribute_value', name: 'double_attribute_value', dataType: 'DOUBLE' },
      { id: 'boolean_attribute_value', name: 'boolean_attribute_value', dataType: 'BOOLEAN' },
      { id: 'string_attribute_value', name: 'string_attribute_value', dataType: 'STRING' },
    ],
  },
  {
    id: 'raw_time_series',
    name: 'raw_time_series',
    properties: [
      { id: 'property_id', name: 'property_id', dataType: 'STRING' },
      { id: 'asset_id', name: 'asset_id', dataType: 'STRING' },
      { id: 'property_alias', name: 'property_alias', dataType: 'STRING' },
      { id: 'event_timestamp', name: 'event_timestamp', dataType: 'TIMESTAMP' },
      { id: 'quality', name: 'quality', dataType: 'STRING' },
      { id: 'boolean_value', name: 'boolean_value', dataType: 'BOOLEAN' },
      { id: 'int_value', name: 'int_value', dataType: 'INTEGER' },
      { id: 'double_value', name: 'double_value', dataType: 'DOUBLE' },
      { id: 'string_value', name: 'string_value', dataType: 'STRING' },
    ],
  },
  {
    id: 'latest_value_time_series',
    name: 'latest_value_time_series',
    properties: [
      { id: 'property_id', name: 'property_id', dataType: 'STRING' },
      { id: 'asset_id', name: 'asset_id', dataType: 'STRING' },
      { id: 'property_alias', name: 'property_alias', dataType: 'STRING' },
      { id: 'event_timestamp', name: 'event_timestamp', dataType: 'TIMESTAMP' },
      { id: 'quality', name: 'quality', dataType: 'STRING' },
      { id: 'boolean_value', name: 'boolean_value', dataType: 'BOOLEAN' },
      { id: 'int_value', name: 'int_value', dataType: 'INTEGER' },
      { id: 'double_value', name: 'double_value', dataType: 'DOUBLE' },
      { id: 'string_value', name: 'string_value', dataType: 'STRING' },
    ],
  },
  {
    id: 'precomputed_aggregates',
    name: 'precomputed_aggregates',
    properties: [
      { id: 'asset_id', name: 'asset_id', dataType: 'STRING' },
      { id: 'property_id', name: 'property_id', dataType: 'STRING' },
      { id: 'property_alias', name: 'property_alias', dataType: 'STRING' },
      { id: 'event_timestamp', name: 'event_timestamp', dataType: 'TIMESTAMP' },
      { id: 'quality', name: 'quality', dataType: 'STRING' },
      { id: 'resolution', name: 'resolution', dataType: 'STRING' },
      { id: 'sum_value', name: 'sum_value', dataType: 'DOUBLE' },
      { id: 'count_value', name: 'count_value', dataType: 'INTEGER' },
      { id: 'average_value', name: 'average_value', dataType: 'DOUBLE' },
      { id: 'maximum_value', name: 'maximum_value', dataType: 'DOUBLE' },
      { id: 'minimum_value', name: 'minimum_value', dataType: 'DOUBLE' },
      { id: 'stdev_value', name: 'stdev_value', dataType: 'DOUBLE' },
    ],
  },
];

export const timeIntervalProperty: AssetProperty = {
  id: 'timeInterval',
  name: 'Time Interval',
  dataType: 'STRING',
};

export const whereOperators = [
  { label: '>', value: '>' },
  { label: '<', value: '<' },
  { label: '>=', value: '>=' },
  { label: '<=', value: '<=' },
  { label: '=', value: '=' },
  { label: '!=', value: '!=' },
  { label: 'LIKE', value: 'LIKE' },
  { label: 'IN', value: 'IN' },
  { label: 'BETWEEN', value: 'BETWEEN' },
  { label: 'IS NULL', value: 'IS NULL' },
  { label: 'IS NOT NULL', value: 'IS NOT NULL' },
];

export const timeIntervals: Array<SelectableValue<string>> = [
  { label: '1s', value: '1s' },
  { label: '10s', value: '10s' },
  { label: '1m', value: '1m' },
  { label: '5m', value: '5m' },
  { label: '10m', value: '10m' },
  { label: '15m', value: '15m' },
  { label: '1h', value: '1h' },
];

export interface FunctionOption {
  label: string;
  value: string;
  items?: FunctionOption[]; // children
}

export const allFunctions: FunctionOption[] = [
  {
    label: 'Aggregate',
    value: 'Aggregate',
    items: [
      { label: 'Raw Values', value: '' },
      { label: 'AVG', value: 'AVG' },
      { label: 'SUM', value: 'SUM' },
      { label: 'COUNT', value: 'COUNT' },
      { label: 'MAX', value: 'MAX' },
      { label: 'MIN', value: 'MIN' },
      { label: 'STDDEV', value: 'STDDEV' },
    ],
  },
  {
    label: 'String',
    value: 'String',
    items: [
      { label: 'LENGTH', value: 'LENGTH' },
      { label: 'CONCAT', value: 'CONCAT' },
      { label: 'SUBSTR', value: 'SUBSTR' },
      { label: 'UPPER', value: 'UPPER' },
      { label: 'LOWER', value: 'LOWER' },
      { label: 'TRIM', value: 'TRIM' },
      { label: 'LTRIM', value: 'LTRIM' },
      { label: 'RTRIM', value: 'RTRIM' },
      { label: 'STR_REPLACE', value: 'STR_REPLACE' },
      { label: 'COALESCE', value: 'COALESCE' },
      { label: 'Raw Values', value: 'Raw Values' },
    ],
  },
  {
    label: 'Math',
    value: 'Math',
    items: [
      { label: 'POWER', value: 'POWER' },
      { label: 'ROUND', value: 'ROUND' },
      { label: 'FLOOR', value: 'FLOOR' },
    ],
  },
  {
    label: 'DateTime',
    value: 'DateTime',
    items: [
      { label: 'NOW', value: 'NOW' },
      { label: 'DATE_ADD', value: 'DATE_ADD' },
      { label: 'DATE_SUB', value: 'DATE_SUB' },
      { label: 'TIMESTAMP_ADD', value: 'TIMESTAMP_ADD' },
      { label: 'TIMESTAMP_SUB', value: 'TIMESTAMP_SUB' },
      { label: 'CAST', value: 'CAST' },
      { label: 'TO_DATE', value: 'TO_DATE' },
      { label: 'TO_TIMESTAMP', value: 'TO_TIMESTAMP' },
      { label: 'TO_TIME', value: 'TO_TIME' },
    ],
  },
];

export const FUNCTION_ARGS = {
  DATE: ['DAY', 'MONTH', 'YEAR'],
  CAST: ['BOOLEAN', 'INTEGER', 'INT', 'TIMESTAMP'],
};

export const tooltipMessages: Record<string, string> = {
  FROM: 'Select the view to query data from.',
  SELECT: 'Choose the fields or columns you want to retrieve in the result.',
  WHERE: 'Filter rows based on specific conditions.',
  'GROUP BY': 'Select one or more columns to group your query by',
  'ORDER BY': 'Sort the result set by one or more columns.',
  HAVING: 'Filter groups based on aggregate conditions.',
  LIMIT: 'Restrict the number of records returned by the query.',
};

export const FUNCTION_TYPES = {
  date: ['DATE_ADD', 'DATE_SUB', 'TIMESTAMP_ADD', 'TIMESTAMP_SUB'],
  math: ['POWER', 'ROUND', 'COALESCE'],
  cast: ['CAST'],
  now: ['NOW'],
  str: ['STR_REPLACE', 'SUBSTR'],
  coalesce: ['COALESCE'],
  val: ['IS NULL', 'IS NOT NULL'],
  concat: ['CONCAT'],
} as const;

export type FunctionType = keyof typeof FUNCTION_TYPES;
export type FunctionName = (typeof FUNCTION_TYPES)[FunctionType][number];
export const isFunctionOfType = (fn?: string, ...types: FunctionType[]): fn is FunctionName => {
  if (!fn) {
    return false;
  }

  return types.some((type) => {
    const functions = FUNCTION_TYPES[type] as readonly string[];
    return functions.includes(fn);
  });
};

export interface ValidationError {
  type: string;
  error: string;
}
