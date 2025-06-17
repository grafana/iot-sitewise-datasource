import defaults from 'lodash/defaults';
import React, { useEffect } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions, QueryType, ListAssetsQuery, ListTimeSeriesQuery } from 'types';
import { Icon, LinkButton, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWiseQueryTypes, changeQueryType } from 'queryInfo';
import { regionOptions, type Region } from 'regions';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';
import { PropertyQueryEditor } from './PropertyQueryEditor';
import { migrateQuery } from '../../../migrations/migrateQuery';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/plugin-ui';
import { ClientCacheRow } from './ClientCacheRow';
import { ListTimeSeriesQueryEditorFunction } from './ListTimeSeriesQueryEditor';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

const queryDefaults: Partial<SitewiseQuery> = {
  maxPageAggregations: 1,
  flattenL4e: true,
  clientCache: true,
};

export function VisualQueryBuilder(props: Props) {
  const { datasource } = props;
  const query = defaults(props.query, queryDefaults);

  useEffect(() => {
    const migratedQuery = migrateQuery(query);

    if (query !== migratedQuery) {
      props.onChange(migratedQuery);
    }
    // eslint-disable-next-line @typescript-eslint/no-deprecated, react-hooks/exhaustive-deps
  }, [query.assetId, query.propertyId, query.propertyAlias]);

  const defaultRegion: SelectableValue<Region> = {
    label: `Default`,
    description: datasource.options?.defaultRegion,
    value: undefined,
  };
  const regions = query.region ? [defaultRegion, ...regionOptions] : regionOptions;
  const currentQueryType = siteWiseQueryTypes.find((v) => v.value === query.queryType);

  const onQueryTypeChange = (sel: SelectableValue<QueryType>) => {
    const { onChange, query } = props;
    // hack to use QueryEditor as VariableQueryEditor
    onChange(changeQueryType(query, sel as QueryTypeInfo));
  };

  const onRegionChange = (sel: SelectableValue<Region>) => {
    const { onChange, query } = props;
    onChange({ ...query, assetId: undefined, propertyId: undefined, region: sel.value });
  };

  const onClientCacheChange = (evt: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    onChange({ ...query, clientCache: evt.currentTarget.checked });
  };

  const renderQuery = (query: SitewiseQuery) => {
    if (!query.queryType) {
      return;
    }
    switch (query.queryType) {
      case QueryType.ExecuteQuery:
      case QueryType.ListAssetModels:
        return null; // nothing required
      case QueryType.ListAssets:
        return <ListAssetsQueryEditor {...props} query={query as ListAssetsQuery} />;
      case QueryType.ListTimeSeries:
        return <ListTimeSeriesQueryEditorFunction {...props} query={query as ListTimeSeriesQuery} />;
      case QueryType.ListAssociatedAssets:
      case QueryType.PropertyValue:
      case QueryType.PropertyInterpolated:
      case QueryType.PropertyAggregate:
      case QueryType.PropertyValueHistory:
        return <PropertyQueryEditor {...props} />;
    }
    return <div>Missing UI for query type: {query.queryType}</div>;
  };

  const queryTooltip = currentQueryType ? (
    <div>
      {currentQueryType.description} <br />
      <LinkButton href={currentQueryType.helpURL} target="_blank">
        API Docs <Icon name="external-link-alt" />
      </LinkButton>
    </div>
  ) : undefined;

  const clientCacheRow = (
    <ClientCacheRow
      clientCache={query.clientCache}
      onClientCacheChange={onClientCacheChange}
      queryRefId={query.refId}
    />
  );

  return (
    <>
      <EditorRows>
        <EditorRow>
          <EditorFieldGroup>
            <EditorField htmlFor="query" label="Query type" tooltip={queryTooltip} tooltipInteractive width={30}>
              <Select
                id="query"
                aria-label="Query type"
                options={siteWiseQueryTypes}
                value={currentQueryType}
                onChange={onQueryTypeChange}
                placeholder="Select query type"
                menuPlacement="auto"
              />
            </EditorField>
            <EditorField label="Region" width={15}>
              <Select
                options={regions}
                value={regionOptions.find((v) => v.value === query.region) || defaultRegion}
                onChange={onRegionChange}
                backspaceRemovesValue={true}
                allowCustomValue={true}
                isClearable={true}
                menuPlacement="auto"
              />
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>
        {renderQuery(query)}
        {clientCacheRow}
      </EditorRows>
    </>
  );
}
