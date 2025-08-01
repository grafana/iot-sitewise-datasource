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
    expect(errors).toContain('At least one column must be selected in the SELECT clause.');
  });

  it('should return error if selectFields has no valid column', () => {
    const invalidQuery: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: '' }],
    };

    const errors = validateQuery(invalidQuery);
    expect(errors).toContain('At least one column must be selected in the SELECT clause.');
  });

  it('should return error if selectedAssetModel is missing', () => {
    const invalidQuery: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: '',
      selectFields: [{ column: 'asset_id' }],
    };

    const errors = validateQuery(invalidQuery);
    expect(errors).toContain('A source (e.g., asset model or table) must be specified in the FROM clause.');
  });

  it('should return error if a WHERE condition is missing operator or value', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'asset_id' }],
      whereConditions: [{ column: 'asset_name', operator: '', value: '' }],
    };

    const errors = validateQuery(query);
    expect(errors).toContain(
      'Each WHERE condition must include both an operator and a value when a column is selected.'
    );
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
    expect(errors).toContain('Limit must be a valid number.');
  });

  it('should return error if limit is zero or negative', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      limit: 0,
    };

    const errors = validateQuery(query);
    expect(errors).toContain('Limit must be greater than 0.');
  });

  it('should return error if limit exceeds 100000', () => {
    const query: SitewiseQueryState = {
      ...defaultSitewiseQueryState,
      selectedAssetModel: 'asset',
      selectFields: [{ column: 'id' }],
      limit: 100001,
    };

    const errors = validateQuery(query);
    expect(errors).toContain('Limit must not exceed 100,000 rows.');
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
    expect(errors).toContain('Group by tags must not contain empty values.');
  });
});
