import React from 'react';
import { QueryType, ListAssetsQuery, SqlQuery } from 'types';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';
import { PropertyQueryEditor } from './PropertyQueryEditor';
import { SqlQueryEditor } from './SqlQueryEditor';
import { Props } from './QueryEditor';

export const QueryField = ({ query, ...props }: Props) => {
  if (!query.queryType) {
    return null;
  }

  switch (query.queryType) {
    case QueryType.ListAssetModels:
      return null; // nothing required
    case QueryType.ListAssets:
      return <ListAssetsQueryEditor {...props} query={query as ListAssetsQuery} />;
    case QueryType.ListAssociatedAssets:
    case QueryType.PropertyValue:
    case QueryType.PropertyInterpolated:
    case QueryType.PropertyAggregate:
    case QueryType.PropertyValueHistory:
      return <PropertyQueryEditor {...props} query={query} />;
    case QueryType.SQL:
      return <SqlQueryEditor {...props} query={query as SqlQuery} />;
    default:
      return <div>Missing UI for query type: {query.queryType}</div>;
  }
};
