import { QueryType, SitewiseQuery } from '../types';
import { migrateQuery } from './migrateQuery';

describe('migrateQuery()', () => {
  it('should return the same query by reference equality if there are no migrations', () => {
    const query: SitewiseQuery = { refId: 'a', queryType: QueryType.PropertyAggregate };
    const migratedQuery = migrateQuery(query);
    expect(query).toBe(migratedQuery);
  });

  it('should migrate assetId to assetIds if assetIds does not exist', () => {
    const query: SitewiseQuery = { refId: 'a', queryType: QueryType.PropertyAggregate, assetId: 'asset-id' };
    const migratedQuery = migrateQuery(query);
    expect(query).not.toBe(migratedQuery);
    expect(migratedQuery.assetIds).toEqual(expect.arrayContaining(['asset-id']));
  });
});
