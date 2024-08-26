import { SitewiseQuery } from '../types';

const migrateAssetId = (query: SitewiseQuery): SitewiseQuery => {
  if (query.assetId && !query.assetIds) {
    return { ...query, assetIds: [query.assetId] };
  }

  return query;
};

export const migrateQuery = (query: SitewiseQuery): SitewiseQuery => {
  let migratedQuery = migrateAssetId(query);

  return migratedQuery;
};
