import { css } from '@emotion/css';
import { type SelectableValue } from '@grafana/data';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { LinkButton, Select, Input, Icon } from '@grafana/ui';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { getAssetProperty, getDefaultAggregate } from 'queryInfo';
import {
  isAssetPropertyAggregatesQuery,
  isAssetPropertyValueHistoryQuery,
  isListAssociatedAssetsQuery,
  isAssetPropertyInterpolatedQuery,
  shouldShowOptionsRow,
  type AssetInfo,
  type ListAssociatedAssetsQuery,
  type SitewiseQuery,
} from 'types';
import { DEFAULT_REGION } from '../../../regions';
import { AssetBrowser } from '../../browser/AssetBrowser';
import { aggReg } from './AggregationSettings/AggregatePicker';
import { AggregationSettings } from './AggregationSettings/AggregationSettings';
import { InterpolatedResolutionSettings } from './InterpolatedResolutionSettings';
import { QueryOptions } from './QueryOptions';
import type { SitewiseQueryEditorProps } from './types';

const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89abAB][0-9a-f]{3}-[0-9a-f]{12}$/i;
const ALL_HIERARCHIES = '*';

export const PropertyQueryEditor = ({ query, datasource, onChange }: SitewiseQueryEditorProps) => {
  const [isLoading, setIsLoading] = useState(false);
  const [assetId, setAssetId] = useState<string | undefined>(query.assetIds?.[0]);
  const [asset, setAsset] = useState<AssetInfo | undefined>(undefined);
  const [assets, setAssets] = useState<Array<SelectableValue<string>>>([]);
  const [assetProperties, setAssetProperties] = useState<Array<{ id: string; name: string }>>([]);

  const cache = useMemo(() => datasource.getCache(query.region), [datasource, query.region]);

  const onAliasChange = useCallback(
    (evt: React.SyntheticEvent<HTMLInputElement>) => {
      onChange({ ...query, propertyAlias: evt.currentTarget.value });
    },
    [onChange, query]
  );

  const onSetAssetId = useCallback(
    (assetId?: string) => {
      if (!assetId) {
        setAssetId(undefined);
        onChange({ ...query, assetIds: undefined });
      } else {
        const validId = uuidRegex.test(assetId) || assetId.startsWith('externalId:') || assetId.startsWith('$');
        const assetIds = validId ? [assetId] : [`externalId:${assetId}`];
        setAssetId(assetIds[0]);
        onChange({ ...query, assetIds });
      }
    },
    [onChange, query]
  );

  const onAssetChange = useCallback(
    (sel: SelectableValue<string> | Array<SelectableValue<string>>) => {
      const assetIds: Set<string> = new Set();
      if (Array.isArray(sel)) {
        sel.forEach((s) => {
          if (s.value) {
            assetIds.add(s.value);
          }
        });
      } else if (sel.value) {
        assetIds.add(sel.value);
      }

      const newQuery =
        Array.isArray(sel) && sel.length === 0
          ? {
              ...query,
              assetIds: [],
              propertyId: undefined,
            }
          : {
              ...query,
              assetIds: [...assetIds],
            };

      onChange(newQuery);
    },
    [onChange, query]
  );

  const onPropertyChange = useCallback(
    (sel: SelectableValue<string>) => {
      const update = { ...query, propertyId: sel.value! } satisfies SitewiseQuery;

      // Make sure the selected aggregates are actually supported
      if (isAssetPropertyAggregatesQuery(update)) {
        if (update.propertyId) {
          const info = getAssetProperty(asset, update.propertyId);

          if (!update.aggregates) {
            update.aggregates = [];
          }

          if (info) {
            update.aggregates = update.aggregates.filter((a) => aggReg.get(a).isValid(info));
          }

          if (!update.aggregates.length) {
            update.aggregates = [getDefaultAggregate(info)];
          }
        }
      }

      onChange(update);
    },
    [onChange, query]
  );

  const onSetPropertyId = useCallback(
    (propertyId?: string) => {
      onChange({ ...query, propertyId });
    },
    [onChange, query]
  );

  const onHierarchyIdChange = useCallback(
    (sel: SelectableValue<string>) => {
      const update = { ...query };
      if (isListAssociatedAssetsQuery(update)) {
        if (sel.value === ALL_HIERARCHIES) {
          delete update.hierarchyId;
          update.loadAllChildren = true;
        } else if (sel.value && sel.value.length) {
          update.hierarchyId = sel.value;
          update.loadAllChildren = false;
        } else {
          delete update.hierarchyId;
          update.loadAllChildren = false;
        }
      }
      onChange(update);
    },
    [onChange, query]
  );

  const onSetHierarchyId = useCallback(
    (hierarchyId?: string) => {
      // FIXME: query is being casted to `any` as hierarchy does not exist on the base query type
      onChange({ ...(query as any), hierarchyId });
    },
    [onChange, query]
  );

  const onLastObservationChange = useCallback(() => {
    onChange({ ...query, lastObservation: !query.lastObservation });
  }, [onChange, query]);

  const onFlattenL4eChange = useCallback(() => {
    onChange({ ...query, flattenL4e: !query.flattenL4e });
  }, [onChange, query]);

  const renderAssociatedAsset = (query: ListAssociatedAssetsQuery) => {
    const hierarchies: Array<SelectableValue<string>> = [
      { value: '', label: '** Parent **' },
      { value: ALL_HIERARCHIES, label: '** All **' },
    ];
    if (asset) {
      hierarchies.push(...asset.hierarchy);
    }

    let current = hierarchies.find((v) => v.value === query.hierarchyId);
    if (!current) {
      if (query.hierarchyId) {
        current = { value: query.hierarchyId, label: 'ID: ' + query.hierarchyId };
        hierarchies.push(current);
      } else {
        current = query.loadAllChildren ? hierarchies[1] /* all */ : hierarchies[0]; // parent
      }
    }

    return (
      <EditorField label="Asset Hierarchy" htmlFor="assetHierarchy">
        <Select
          id="assetHierarchy"
          aria-label="Asset Hierarchy"
          isLoading={isLoading}
          options={hierarchies}
          value={current}
          onChange={onHierarchyIdChange}
          placeholder="Select..."
          allowCustomValue={true}
          backspaceRemovesValue={true}
          isClearable={true}
          isSearchable={true}
          onCreateOption={onSetHierarchyId}
          formatCreateLabel={(txt) => `Hierarchy Id: ${txt}`}
          menuPlacement="auto"
        />
      </EditorField>
    );
  };

  useEffect(() => {
    setIsLoading(true);
    const assetId = query.assetIds?.[0];

    Promise.allSettled([
      assetId && cache.getAssetInfo(assetId).then(setAsset),
      assetId &&
        cache.listAssetProperties(assetId).then((assetProperties) => {
          setAssetProperties(assetProperties?.map(({ id, name }) => ({ id, name })) ?? []);
        }),
      cache.getAssetPickerOptions().then(setAssets),
    ])
      .catch(console.error)
      .finally(() => setIsLoading(false));
  }, [query]);

  let current = assets.filter((a) => (a.value ? query.assetIds?.includes(a.value) : false));
  if (current.length === 0 && query.assetIds?.length) {
    if (isLoading) {
      current = query.assetIds.map((id) => ({ label: 'loading...', value: id }));
    } else {
      current = query.assetIds.map((id) => ({ label: `ID: ${assetId}`, value: id }));
    }
  }

  const isAssociatedAssets = isListAssociatedAssetsQuery(query);
  const showProp = Boolean(!isAssociatedAssets && (query.propertyId || query.assetIds));
  const assetPropertyOptions = showProp
    ? assetProperties.map(({ id, name }) => ({ id, name, value: id, label: name }))
    : [];

  const showQuality = Boolean(
    query.propertyId ||
      (query.propertyAlias && isAssetPropertyAggregatesQuery(query)) ||
      isAssetPropertyValueHistoryQuery(query) ||
      isAssetPropertyInterpolatedQuery(query)
  );

  const showOptionsRow = shouldShowOptionsRow(query, showProp);

  let currentAssetPropertyOption = assetPropertyOptions.find((p) => p.id === query.propertyId);
  if (!currentAssetPropertyOption && query.propertyId) {
    currentAssetPropertyOption = {
      id: query.propertyId,
      name: 'ID: ' + query.propertyId,
      value: query.propertyId,
      label: 'ID: ' + query.propertyId,
    };
  }

  return (
    <>
      {!isAssociatedAssets && (
        <EditorRow>
          <EditorField label="Property Alias" tooltip={<QueryTooltip />} tooltipInteractive htmlFor="alias" width={80}>
            <Input
              id="alias"
              aria-label="Property alias"
              value={query.propertyAlias}
              onChange={onAliasChange}
              placeholder="optional alias that identifies the property, such as an OPC-UA server data stream path"
            />
          </EditorField>
        </EditorRow>
      )}

      {(!Boolean(query.propertyAlias) || isAssociatedAssets) && (
        <>
          <EditorRow>
            <EditorFieldGroup>
              <EditorField label="Asset" tooltip={<AssetTooltip />} tooltipInteractive htmlFor="asset" width={30}>
                <Select
                  id="asset"
                  inputId="asset"
                  aria-label="Asset"
                  isMulti={true}
                  key={query.region ?? DEFAULT_REGION}
                  isLoading={isLoading}
                  options={assets}
                  value={current}
                  onChange={(sel) => onAssetChange(sel)}
                  placeholder="Select an asset"
                  allowCustomValue={true}
                  isClearable={true}
                  isSearchable={true}
                  onCreateOption={onSetAssetId}
                  formatCreateLabel={(txt) => `Asset ID: ${txt}`}
                  menuPlacement="auto"
                />
              </EditorField>

              <div className={styles.exploreContainer}>
                <AssetBrowser
                  datasource={datasource}
                  region={query.region}
                  assetId={query.assetIds?.[0]}
                  onAssetChanged={onSetAssetId}
                />
              </div>
            </EditorFieldGroup>
          </EditorRow>

          {showProp && (
            <EditorRow>
              <EditorFieldGroup>
                <EditorField label="Property" htmlFor="property" width={30}>
                  <Select
                    id="property"
                    inputId="property"
                    aria-label="Property"
                    isLoading={isLoading}
                    options={assetPropertyOptions}
                    value={currentAssetPropertyOption ?? null}
                    onChange={onPropertyChange}
                    placeholder="Select a property"
                    allowCustomValue
                    isSearchable
                    onCreateOption={onSetPropertyId}
                    formatCreateLabel={(txt) => `Property ID: ${txt}`}
                    menuPlacement="auto"
                  />
                </EditorField>
              </EditorFieldGroup>

              <EditorFieldGroup>
                {showQuality && isAssetPropertyAggregatesQuery(query) && (
                  <AggregationSettings query={query} property={getAssetProperty(asset)} onChange={onChange} />
                )}

                {isAssetPropertyInterpolatedQuery(query) && (
                  <InterpolatedResolutionSettings query={query} onChange={onChange} />
                )}
              </EditorFieldGroup>
            </EditorRow>
          )}
        </>
      )}

      {query.propertyAlias && isAssetPropertyAggregatesQuery(query) && (
        <EditorRow>
          <AggregationSettings query={query} property={getAssetProperty(asset)} onChange={onChange} />
        </EditorRow>
      )}

      {query.propertyAlias && isAssetPropertyInterpolatedQuery(query) && (
        <EditorRow>
          <InterpolatedResolutionSettings query={query} onChange={onChange} />
        </EditorRow>
      )}

      {isAssociatedAssets && (
        <EditorRow>
          <EditorFieldGroup>{renderAssociatedAsset(query as ListAssociatedAssetsQuery)}</EditorFieldGroup>
        </EditorRow>
      )}

      {showOptionsRow && (
        <EditorRow>
          <EditorFieldGroup>
            <QueryOptions
              query={query}
              datasource={datasource}
              onChange={onChange}
              showProp={showProp}
              showQuality={!!(query.propertyId || query.propertyAlias)}
              onLastObservationChange={onLastObservationChange}
              onFlattenL4eChange={onFlattenL4eChange}
            />
          </EditorFieldGroup>
        </EditorRow>
      )}
    </>
  );
};

const AssetTooltip = () => (
  <div>
    Set the asset ID. It can be either the actual ID in UUID format, or else "externalId:" followed by the external ID,
    if it has one.
    <LinkButton
      href="https://docs.aws.amazon.com/iot-sitewise/latest/userguide/object-ids.html#external-ids"
      target="_blank"
    >
      API Docs <Icon name="external-link-alt" />
    </LinkButton>
  </div>
);

const QueryTooltip = () => (
  <div>
    Setting an alias for an asset property. <br />
    <LinkButton
      href="https://docs.aws.amazon.com/iot-sitewise/latest/userguide/connect-data-streams.html"
      target="_blank"
    >
      API Docs <Icon name="external-link-alt" />
    </LinkButton>
  </div>
);

const styles = {
  exploreContainer: css({ display: 'flex', alignItems: 'flex-end' }),
};
