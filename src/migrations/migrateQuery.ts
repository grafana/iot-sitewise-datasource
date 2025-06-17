/* eslint-disable @typescript-eslint/no-deprecated */

import { SitewiseQuery } from '../types';

const migrateAssetProperty = (query: SitewiseQuery): SitewiseQuery => {
  if (query.assetId && !query.assetIds) {
    return { ...query, assetIds: [query.assetId] };
  }

  if (query.propertyId && !query.propertyIds) {
    return { ...query, propertyIds: [query.propertyId] };
  }

  if (query.propertyAlias && !query.propertyAliases) {
    return { ...query, propertyAliases: [query.propertyAlias] };
  }

  return query;
};

export const migrateQuery = (query: SitewiseQuery): SitewiseQuery => {
  const migratedQuery = migrateAssetProperty(query);

  return migratedQuery;
};
