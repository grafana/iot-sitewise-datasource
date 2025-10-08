import { generateQueryPreview } from './queryGenerator';
import { AggregationFunction, SitewiseQueryState } from '../types';

describe('generateQueryPreview', () => {
  it('returns message when asset model is not selected', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: '',
      selectFields: [],
      whereConditions: [],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toBe('');
  });

  it('generates minimal SELECT query', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [{ column: 'prop-1', aggregation: '', alias: '' }],
      whereConditions: [],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('SELECT');
    expect(preview).toContain('FROM model-1');
  });

  it('includes WHERE clause with logical operator', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [{ column: 'prop-1', aggregation: '', alias: '' }],
      whereConditions: [
        { column: 'prop-1', operator: '=', value: '100', logicalOperator: 'AND' },
        { column: 'prop-2', operator: '!=', value: '200' },
      ],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview.replace(/\s+/g, ' ').trim()).toContain(
      "SELECT prop-1 FROM model-1 WHERE prop-1 = '100' AND prop-2 != '200' LIMIT 100"
    );
  });

  it('includes aggregation and alias in SELECT', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [
        { column: 'prop-1', aggregation: 'AVG', alias: 'avg1' },
        { column: 'prop-2', aggregation: 'MAX', alias: '' },
      ],
      whereConditions: [],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('AVG(');
    expect(preview).toContain('AS "avg1"');
    expect(preview).toContain('MAX(');
  });

  it('includes GROUP BY and ORDER BY clauses', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [{ column: 'prop-1', aggregation: '', alias: '' }],
      whereConditions: [],
      groupByTags: ['prop-1'],
      havingConditions: [],
      orderByFields: [{ column: 'prop-1', direction: 'DESC' }],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('GROUP BY prop-1');
    expect(preview).toContain('ORDER BY prop-1 DESC');
  });

  it('includes LIMIT clause', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [{ column: 'prop-1', aggregation: '', alias: '' }],
      whereConditions: [],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      limit: 10,
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('LIMIT 10');
  });

  it('handles CAST, NOW and DATE_BIN functions', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [
        {
          column: 'prop-1',
          aggregation: 'CAST',
          functionArg: 'DOUBLE',
        },
        {
          column: 'prop-2',
          aggregation: 'NOW',
          functionArg: '',
        },
        {
          column: 'prop-3',
          aggregation: 'DATE_ADD',
          functionArg: '1d',
          functionArgValue: '0',
        },
      ],
      whereConditions: [],
      groupByTags: [],
      havingConditions: [],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('CAST(prop-1 AS DOUBLE)');
    expect(preview).toContain('NOW()');
    expect(preview).toContain('DATE_ADD(1d, 0, prop-3)');
  });

  it('includes simple HAVING clause with aggregation', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-1',
      selectFields: [{ column: 'prop-1', aggregation: 'COUNT', alias: '' }],
      whereConditions: [],
      groupByTags: ['prop-1'],
      havingConditions: [{ aggregation: 'COUNT', column: 'prop-1', operator: '>', value: '5', logicalOperator: 'AND' }],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('HAVING COUNT(prop-1) > 5');
  });

  it('includes multiple HAVING conditions with logical operators', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-2',
      selectFields: [
        { column: 'prop-1', aggregation: 'SUM', alias: '' },
        { column: 'prop-2', aggregation: 'AVG', alias: '' },
      ],
      whereConditions: [],
      groupByTags: ['prop-2'],
      havingConditions: [
        { aggregation: 'SUM', column: 'prop-1', operator: '>=', value: '100', logicalOperator: 'AND' },
        { aggregation: 'AVG', column: 'prop-2', operator: '<', value: '50', logicalOperator: 'OR' },
        { aggregation: 'COUNT', column: 'prop-2', operator: '=', value: '10', logicalOperator: 'AND' },
      ],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('HAVING SUM(prop-1) >= 100 AND AVG(prop-2) < 50 OR COUNT(prop-2) = 10');
  });

  it('skips invalid or empty HAVING conditions', async () => {
    const query: SitewiseQueryState = {
      selectedAssetModel: 'model-3',
      selectFields: [{ column: 'prop-1', aggregation: 'MAX', alias: '' }],
      whereConditions: [],
      groupByTags: ['prop-1'],
      havingConditions: [
        { aggregation: 'MAX', column: '', operator: '=', value: '10', logicalOperator: 'AND' },
        {
          aggregation: '' as unknown as AggregationFunction,
          column: 'prop-1',
          operator: '=',
          value: '20',
          logicalOperator: 'AND',
        },
        { aggregation: 'MAX', column: 'prop-1', operator: '=', value: '20', logicalOperator: 'AND' },
      ],
      orderByFields: [],
      rawSQL: '',
    };

    const preview = await generateQueryPreview(query);
    expect(preview).toContain('HAVING MAX(prop-1) = 20');
    expect(preview).not.toContain('MAX()');
    expect(preview).not.toContain("= '10'");
  });
});
