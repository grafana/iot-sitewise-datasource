import {
  SitewiseQueryState,
  mockAssetModels,
  HavingCondition,
  SelectField,
  WhereCondition,
  isFunctionOfType,
} from '../types';

/**
 * Wraps a value in single quotes if it's a non-variable string.
 *
 * @param val - Any value to be quoted.
 * @returns Quoted string or original value.
 */
const quote = (val: any): string | undefined =>
  typeof val === 'string' && val.trim() !== '' && !val.startsWith('$') ? `'${val}'` : val;

/**
 * Constructs the `SELECT` clause of the SQL query using field metadata and asset model properties.
 *
 * @param fields - List of columns to select.
 * @param properties - Asset model properties to resolve field names.
 * @returns SQL SELECT clause string.
 */
const buildSelectClause = (fields: SelectField[], properties: any[]): string => {
  const clauses = fields
    .filter(({ column }) => column)
    .map((field) => {
      const base = properties.find((p) => p.id === field.column)?.name || field.column;
      const { aggregation, functionArg, functionArgValue, functionArgValue2, alias } = field || '';
      let expr = base;

      if (!aggregation) {
        return alias ? `${expr} AS "${alias}"` : expr;
      }
      // Handle different function types
      switch (true) {
        case isFunctionOfType(aggregation, 'date'):
          expr = `${aggregation}(${functionArg ?? '1d'}, ${functionArgValue ?? '0'}, ${base})`;
          break;
        case isFunctionOfType(aggregation, 'math') || isFunctionOfType(aggregation, 'coalesce'):
          expr = `${aggregation}(${base}, ${functionArgValue ?? '0'})`;
          break;
        case isFunctionOfType(aggregation, 'str'):
          expr =
            aggregation === 'STR_REPLACE'
              ? `${aggregation}(${base}, '${functionArgValue}', '${functionArgValue2}')`
              : `${aggregation}(${base}, ${functionArgValue}, ${functionArgValue2})`;
          break;
        case isFunctionOfType(aggregation, 'concat'):
          expr = `${aggregation}(${base},${functionArg})`;
          break;
        case isFunctionOfType(aggregation, 'cast'):
          expr = `CAST(${base} AS ${functionArg})`;
          break;
        case isFunctionOfType(aggregation, 'now'):
          expr = 'NOW()';
          break;
        default:
          expr = `${aggregation}(${base})`;
      }

      return alias ? `${expr} AS "${alias}"` : expr;
    });

  return `SELECT ${clauses.length ? clauses.join(', ') : '*'}`;
};

/**
 * Builds the `WHERE` clause from a list of conditions.
 *
 * Handles operators like `=`, `BETWEEN`, etc. Supports AND/OR chaining.
 *
 * @param conditions - Array of where conditions.
 * @returns SQL WHERE clause string or empty string.
 */
const buildWhereClause = (conditions: WhereCondition[] = []): string => {
  const parts = conditions
    .filter((c) => c.column && c.operator && c.value !== undefined && c.value !== null)
    .map((c, i, arr) => {
      const val1 = quote(c.value);
      const val2 = quote(c.value2);
      const condition =
        c.operator === 'BETWEEN' && c.value2
          ? `${c.column} ${c.operator} ${val1} ${c.operator2} ${val2}`
          : `${c.column} ${c.operator} ${val1}`;
      const logic = i < arr.length - 1 ? `${c.logicalOperator ?? 'AND'}` : '';
      return `${condition} ${logic}`;
    });

  return parts.length > 0 ? `WHERE ${parts.join(' ')}` : '';
};

/**
 * Builds the `GROUP BY` clause from selected tag columns.
 *
 * @param columns - List of  column names.
 * @returns SQL GROUP BY clause or empty string.
 */
const buildGroupByClause = (columns: string[] = []): string => {
  if (!columns.length) {
    return '';
  }
  const parts = [...columns];
  return `GROUP BY ${parts.join(', ')}`;
};

/**
 * Builds the `HAVING` clause for aggregated filtering after `GROUP BY`.
 *
 * Supports logical chaining using `AND`/`OR`.
 *
 * @param conditions - List of having conditions.
 * @returns SQL HAVING clause or empty string.
 */
const buildHavingClause = (conditions: HavingCondition[] = []): string => {
  const validConditions = conditions.filter((c) => c.column?.trim() && c.aggregation?.trim() && c.operator?.trim());
  if (!validConditions.length) {
    return '';
  }
  const parts = validConditions.map((c) => `${c.aggregation}(${c.column}) ${c.operator} ${Number(c.value)}`);
  return (
    'HAVING ' +
    parts.map((expr, i) => (i === 0 ? expr : `${validConditions[i - 1].logicalOperator ?? 'AND'} ${expr}`)).join(' ')
  );
};

/**
 * Builds the `ORDER BY` clause based on field and direction (ASC/DESC).
 *
 * @param fields - Array of order fields.
 * @returns SQL ORDER BY clause or empty string.
 */
const buildOrderByClause = (fields: Array<{ column: string; direction: string }> = []): string => {
  const parts = fields.filter((f) => f.column && f.direction).map((f) => `${f.column} ${f.direction}`);
  return parts.length ? `ORDER BY ${parts.join(', ')}` : '';
};

/**
 * Builds the `LIMIT` clause for restricting number of query results.
 *
 * @param limit - Number of rows to return.
 * @returns SQL LIMIT clause with default of 100 if not defined.
 */
const buildLimitClause = (limit: number | undefined): string => `LIMIT ${typeof limit === 'number' ? limit : 100}`;

/**
 * Generates the full SQL query preview string by composing all query parts:
 * SELECT, FROM, WHERE, GROUP BY, HAVING, ORDER BY, and LIMIT.
 *
 * @param queryState - Full query state object containing selected model and clauses.
 * @returns A full SQL query string or message prompting to select a model.
 */
export const generateQueryPreview = async (queryState: SitewiseQueryState): Promise<string> => {
  if (!queryState.selectedAssetModel) {
    return 'Select an asset model to build your query';
  }

  const model = mockAssetModels.find((m) => m.id === queryState.selectedAssetModel);
  const properties = model?.properties || [];

  const selectClause = buildSelectClause(queryState.selectFields ?? [], properties);
  const whereClause = buildWhereClause(queryState.whereConditions ?? []);
  const groupByClause = buildGroupByClause(queryState.groupByTags ?? []);
  const havingClause = buildHavingClause(queryState.havingConditions ?? []);
  const orderByClause = buildOrderByClause(queryState.orderByFields ?? []);
  const limitClause = buildLimitClause(queryState.limit);

  return [
    selectClause,
    `FROM ${queryState.selectedAssetModel}`,
    whereClause,
    groupByClause,
    havingClause,
    orderByClause,
    limitClause,
  ]
    .filter(Boolean)
    .join('\n');
};
