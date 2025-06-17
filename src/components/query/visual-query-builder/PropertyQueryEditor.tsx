import { css } from '@emotion/css';
import { type SelectableValue } from '@grafana/data';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { LinkButton, Select, Icon } from '@grafana/ui';
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
import { getSelectableTemplateVariables } from '../../../variables';

const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89abAB][0-9a-f]{3}-[0-9a-f]{12}$/i;
const ALL_HIERARCHIES = '*';

export const PropertyQueryEditor = ({ query, datasource, onChange }: SitewiseQueryEditorProps) => {
  const [isLoading, setIsLoading] = useState(false);
  const [assetId, setAssetId] = useState<string | undefined>(query.assetIds?.[0]);
  const [asset, setAsset] = useState<AssetInfo | undefined>(undefined);
  const [assets, setAssets] = useState<Array<SelectableValue<string>>>([]);
  const [assetProperties, setAssetProperties] = useState<Array<SelectableValue<string>>>([]);
  const [propertyAliases, setPropertyAliases] = useState<Array<SelectableValue<string>>>([]);

  const cache = useMemo(() => datasource.getCache(query.region), [datasource, query.region]);

  const onAliasChange = useCallback(
    (sel: SelectableValue<string> | Array<SelectableValue<string>>) => {
      const propertyAliases: Set<string> = new Set();
      if (Array.isArray(sel)) {
        sel.forEach((s) => {
          if (s.value) {
            propertyAliases.add(s.value);
          }
        });
      } else if (sel.value) {
        propertyAliases.add(sel.value);
      }

      onChange({ ...query, propertyAliases: [...propertyAliases] });
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
          ? { ...query, assetIds: [], propertyIds: [] }
          : { ...query, assetIds: [...assetIds] };

      onChange(newQuery);
    },
    [onChange, query]
  );

  const onPropertyChange = useCallback(
    (sel: SelectableValue<string> | Array<SelectableValue<string>>) => {
      const propertyIds: Set<string> = new Set();
      if (Array.isArray(sel)) {
        sel.forEach((s) => {
          if (s.value) {
            propertyIds.add(s.value);
          }
        });
      } else if (sel.value) {
        propertyIds.add(sel.value);
      }

      const newQuery = {
        ...query,
        propertyIds: [...propertyIds],
      } satisfies SitewiseQuery;

      // Make sure the selected aggregates are actually supported
      if (isAssetPropertyAggregatesQuery(newQuery)) {
        newQuery.aggregates = newQuery.aggregates ?? [];
        newQuery.propertyIds.forEach((propertyId) => {
          const info = getAssetProperty(asset, propertyId);
          if (info) {
            newQuery.aggregates = newQuery.aggregates.filter((a) => aggReg.get(a).isValid(info));
          }
        });

        if (!newQuery.aggregates.length) {
          newQuery.aggregates = [getDefaultAggregate()];
        }
      }

      onChange(newQuery);
    },
    [onChange, query, asset]
  );

  const onSetPropertyAlias = useCallback(
    (propertyAlias?: string) => {
      if (!propertyAlias) {
        onChange({ ...query, propertyAliases: [] });
      } else if (query.propertyAliases) {
        onChange({ ...query, propertyAliases: [...query.propertyAliases, propertyAlias] });
      } else {
        onChange({ ...query, propertyAliases: [propertyAlias] });
      }
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

  const onSetPropertyId = useCallback(
    (propertyId?: string) => {
      if (!propertyId) {
        onChange({ ...query, propertyIds: undefined });
      } else if (query.propertyIds) {
        onChange({ ...query, propertyIds: [...query.propertyIds, propertyId] });
      } else {
        onChange({ ...query, propertyIds: [propertyId] });
      }
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

    setPropertyAliases(getSelectableTemplateVariables());

    Promise.allSettled([
      assetId && cache.getAssetInfo(assetId).then(setAsset),
      assetId &&
        cache.listAssetProperties(assetId).then((assetProperties) => {
          setAssetProperties(
            assetProperties?.map(({ id, name }) => ({
              value: id,
              label: name,
            })) ?? []
          );
        }),
      cache.getAssetPickerOptions().then(setAssets),
    ])
      .catch(console.error)
      .finally(() => setIsLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  let selectedPropertyAliases = propertyAliases.filter(
    (alias) => alias.value && query.propertyAliases?.includes(alias.value)
  );
  if (selectedPropertyAliases.length === 0 && query.propertyAliases?.length) {
    if (isLoading) {
      selectedPropertyAliases = query.propertyAliases.map((alias) => ({
        value: alias,
        label: 'loading...',
      }));
    } else {
      selectedPropertyAliases = query.propertyAliases.map((alias) => ({
        value: alias,
        label: alias,
      }));
    }
  }

  let currentAsset = assets.filter((a) => (a.value ? query.assetIds?.includes(a.value) : false));
  if (currentAsset.length === 0 && query.assetIds?.length) {
    if (isLoading) {
      currentAsset = query.assetIds.map((id) => ({
        value: id,
        label: 'loading...',
      }));
    } else {
      currentAsset = query.assetIds.map((id) => ({
        value: id,
        label: `ID: ${assetId}`,
      }));
    }
  }

  const isAssociatedAssets = isListAssociatedAssetsQuery(query);
  const showProp = Boolean(!isAssociatedAssets && (query.propertyIds || query.assetIds));

  const showQuality = Boolean(
    query.propertyIds ||
      (query.propertyAliases && isAssetPropertyAggregatesQuery(query)) ||
      isAssetPropertyValueHistoryQuery(query) ||
      isAssetPropertyInterpolatedQuery(query)
  );

  const showOptionsRow = shouldShowOptionsRow(query, showProp);

  let currentAssetProperty = assetProperties.filter((property) =>
    property.value ? query.propertyIds?.includes(property.value) : false
  );
  if (currentAssetProperty.length === 0 && query.propertyIds?.length) {
    if (isLoading) {
      currentAssetProperty = query.propertyIds.map((propertyId) => ({
        value: propertyId,
        label: 'loading...',
      }));
    } else {
      currentAssetProperty = query.propertyIds.map((propertyId) => ({
        value: propertyId,
        label: `ID: ${propertyId}`,
      }));
    }
  }

  return (
    <>
      {!isAssociatedAssets && (
        <EditorRow>
          <EditorField label="Property Alias" tooltip={<QueryTooltip />} tooltipInteractive htmlFor="alias" width={80}>
            <Select
              id="alias"
              inputId="alias"
              aria-label="Property alias"
              isMulti
              options={propertyAliases}
              value={selectedPropertyAliases}
              onChange={onAliasChange}
              placeholder="optional alias that identifies the property, such as an OPC-UA server data stream path"
              allowCustomValue
              isClearable
              isSearchable
              onCreateOption={onSetPropertyAlias}
              formatCreateLabel={(txt) => `Property Alias: ${txt}`}
            />
          </EditorField>
        </EditorRow>
      )}

      {(!Boolean(query.propertyAliases?.length) || isAssociatedAssets) && (
        <>
          <EditorRow>
            <EditorFieldGroup>
              <EditorField label="Asset" tooltip={<AssetTooltip />} tooltipInteractive htmlFor="asset" width={30}>
                <Select
                  id="asset"
                  inputId="asset"
                  aria-label="Asset"
                  isMulti
                  key={query.region ?? DEFAULT_REGION}
                  isLoading={isLoading}
                  options={assets}
                  value={currentAsset}
                  onChange={onAssetChange}
                  placeholder="Select an asset"
                  allowCustomValue
                  isClearable
                  isSearchable
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
                    // Disabled multi-selection until a better UX is designed around pairing assets and properties
                    isMulti={false}
                    isLoading={isLoading}
                    options={assetProperties}
                    value={currentAssetProperty}
                    onChange={onPropertyChange}
                    placeholder="Select a property"
                    allowCustomValue
                    isClearable
                    isSearchable
                    onCreateOption={onSetPropertyId}
                    formatCreateLabel={(txt) => `Property ID: ${txt}`}
                    menuPlacement="auto"
                  />
                </EditorField>
              </EditorFieldGroup>

              <EditorFieldGroup>
                {showQuality && isAssetPropertyAggregatesQuery(query) && (
                  <AggregationSettings query={query} onChange={onChange} />
                )}

                {isAssetPropertyInterpolatedQuery(query) && (
                  <InterpolatedResolutionSettings query={query} onChange={onChange} />
                )}
              </EditorFieldGroup>
            </EditorRow>
          )}
        </>
      )}

      {Boolean(query.propertyAliases?.length) && isAssetPropertyAggregatesQuery(query) && (
        <EditorRow>
          <AggregationSettings query={query} onChange={onChange} />
        </EditorRow>
      )}

      {Boolean(query.propertyAliases?.length) && isAssetPropertyInterpolatedQuery(query) && (
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
              showQuality={Boolean(query.propertyIds?.length) || Boolean(query.propertyAliases?.length)}
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
    Set the asset ID. It can be either the actual ID in UUID format, or else &quot;externalId:&quot; followed by the
    external ID, if it has one.
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
