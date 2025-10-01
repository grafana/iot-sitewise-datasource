import { validateQuery } from './validateQuery';
import { SitewiseQueryState, defaultSitewiseQueryState } from '../types';

describe('validateQuery', () => {
  it('should return no errors for a fully valid query', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'asset_id' }],
      whereConditions: [{ column: 'asset_id', operator: '=', value: '25' }],
      limit: 100,
      groupByTags: ['region'],
      orderByFields: [{ column: 'timestamp', direction: 'DESC' }],
    };

    const errors = validateQuery(query);
    expect(errors).toEqual([]);
  });

  it('should return error if no selectFields', () => {
    const invalidQuery: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [],
    };

    const errors = validateQuery(invalidQuery);
    expect(errors).toContainEqual({
      error: 'At least one column must be selected in the SELECT clause.',
      type: 'select',
    });
  });

  it('should return error if selectFields has no valid column', () => {
    const invalidQuery: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: '' }],
    };

    const errors = validateQuery(invalidQuery);
    expect(errors).toContainEqual({
      error: 'At least one column must be selected in the SELECT clause.',
      type: 'select',
    });
  });

  it('should return error if selectedAssetModel is missing', () => {
    const invalidQuery: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: '',
      selectFields: [{ column: 'asset_id' }],
    };

    const errors = validateQuery(invalidQuery);
    expect(errors).toContainEqual({
      error: 'A source (e.g., asset model or table) must be specified in the FROM clause.',
      type: 'from',
    });
  });

  it('should return error if a WHERE condition is missing operator or value', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'asset_id' }],
      whereConditions: [{ column: 'asset_name', operator: '', value: '' }],
    };

    const errors = validateQuery(query);
    expect(errors).toContainEqual({
      error: 'Each WHERE condition must include both an operator and a value when a column is selected.',
      type: 'where',
    });
  });

  it('should ignore WHERE condition if column is not provided', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'asset_id' }],
      whereConditions: [{ column: '', operator: '', value: '' }],
      orderByFields: [{ column: 'timestamp', direction: 'DESC' }],
    };

    const errors = validateQuery(query);
    expect(errors).toEqual([]);
  });

  // LIMIT clause validations
  it('should return error if limit is NaN', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      limit: NaN,
    };

    const errors = validateQuery(query);
    expect(errors).toContainEqual({ error: 'Limit must be a valid number.', type: 'limit' });
  });

  it('should return error if limit is zero or negative', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      limit: 0,
    };

    const errors = validateQuery(query);
    expect(errors).toContainEqual({ error: 'Limit must be greater than 0.', type: 'limit' });
  });

  it('should return error if limit exceeds 100000', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      limit: 100001,
    };

    const errors = validateQuery(query);
    expect(errors).toContainEqual({ error: 'Limit must not exceed 100,000 rows.', type: 'limit' });
  });

  // GROUP BY TAGS validation
  it('should return error if groupByTags has empty values', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      groupByTags: ['region', ''],
    };

    const errors = validateQuery(query);
    expect(errors).toContainEqual({ error: 'Group by tags must not contain empty values.', type: 'group' });
  });
});
