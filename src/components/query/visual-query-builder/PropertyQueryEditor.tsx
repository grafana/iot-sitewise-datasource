import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  AssetInfo,
  isAssetPropertyAggregatesQuery,
  isAssetPropertyValueHistoryQuery,
  ListAssociatedAssetsQuery,
  isListAssociatedAssetsQuery,
  isAssetPropertyInterpolatedQuery,
  shouldShowOptionsRow,
  SitewiseQuery,
} from 'types';
import { LinkButton, Select, Icon } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetBrowser } from '../../browser/AssetBrowser';
import { aggReg } from './AggregationSettings/AggregatePicker';
import { getAssetProperty, getDefaultAggregate } from 'queryInfo';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { css } from '@emotion/css';
import { QueryOptions } from './QueryOptions';
import { AggregationSettings } from './AggregationSettings/AggregationSettings';
import { InterpolatedResolutionSettings } from './InterpolatedResolutionSettings';
import { DEFAULT_REGION } from '../../../regions';
import { getSelectableTemplateVariables } from 'variables';
import { getPropertyIdPickerOptions } from 'sitewiseCache';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89abAB][0-9a-f]{3}-[0-9a-f]{12}$/i;
interface State {
  assetId?: string;
  asset?: AssetInfo;
  assets: Array<SelectableValue<string>>;
  assetProperties: Array<SelectableValue<string>>;
  propertyAliases: Array<SelectableValue<string>>;
  loading: boolean;
  openModal: boolean;
}

const ALL_HIERARCHIES = '*';

export class PropertyQueryEditor extends PureComponent<Props, State> {
  state: State = {
    assetId: this.props.query.assetIds && this.props.query.assetIds[0],
    assets: [],
    assetProperties: [],
    propertyAliases: [],
    loading: true,
    openModal: false,
  };

  async updateInfo() {
    const { onChange, query, datasource } = this.props;
    const update: State = {
      loading: false,
    } as State;

    const cache = datasource.getCache(query.region);

    update.propertyAliases = getSelectableTemplateVariables();

    // TODO: handle user selecting two assets from different asset models
    if (query?.assetIds?.length) {
      try {
        for (let i = 0; i < query?.assetIds.length; i++) {
          update.asset = await cache.getAssetInfo(query.assetIds[i]);
          update.assetProperties = getPropertyIdPickerOptions(update.asset);
          // Update external ids to asset ids in the query
          if (update.asset?.id && query.assetIds[0].startsWith('externalId')) {
            query.assetIds[0] = update.asset.id;
            onChange(query);
          }
        }
      } catch (err) {
        console.warn('error reading asset info', err);
      }
    }

    try {
      update.assets = await cache.getAssetPickerOptions();
    } catch (err) {
      console.warn('error getting options', err);
    }
    this.setState(update);
  }

  async componentDidMount() {
    this.updateInfo();
  }

  async componentDidUpdate(oldProps: Props) {
    const { query } = this.props;

    const assetChanged = query?.assetIds !== oldProps?.query?.assetIds;
    const propChanged = query?.propertyIds !== oldProps?.query?.propertyIds;
    const regionChanged = query?.region !== oldProps?.query?.region;

    if (assetChanged || propChanged || regionChanged) {
      if (!query.assetIds?.length && !regionChanged) {
        this.setState({ assetId: undefined, asset: undefined, loading: false });
      } else {
        this.setState({ loading: true });
        this.updateInfo();
      }
    }
  }

  onAliasChange = (sel: SelectableValue<string> | Array<SelectableValue<string>>) => {
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

    const { onChange, query } = this.props;

    const newQuery =
      Array.isArray(sel) && sel.length === 0
        ? {
            ...query,
            propertyAliases: [],
          }
        : {
            ...query,
            propertyAliases: [...propertyAliases],
          };

    onChange(newQuery);
  };

  onAssetChange(sel: SelectableValue<string> | Array<SelectableValue<string>>) {
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

    const { onChange, query } = this.props;

    const newQuery =
      Array.isArray(sel) && sel.length === 0
        ? {
            ...query,
            assetIds: [],
            propertyIds: undefined,
          }
        : {
            ...query,
            assetIds: [...assetIds],
          };

    onChange(newQuery);
  }

  onPropertyChange = (sel: SelectableValue<string> | Array<SelectableValue<string>>) => {
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

    const { onChange, query } = this.props;

    const newQuery =
      Array.isArray(sel) && sel.length === 0
        ? {
            ...query,
            propertyIds: [],
          }
        : {
            ...query,
            propertyIds: [...propertyIds],
          };

    // Make sure the selected aggregates are actually supported
    if (isAssetPropertyAggregatesQuery(newQuery)) {
      newQuery.aggregates = newQuery.aggregates || [];
      newQuery.propertyIds.forEach((propertyId) => {
        const info = getAssetProperty(this.state.asset, propertyId);
        if (info) {
          newQuery.aggregates = newQuery.aggregates.filter((a) => aggReg.get(a).isValid(info));
        }
      });

      if (!newQuery.aggregates.length) {
        newQuery.aggregates = [getDefaultAggregate()];
      }
    }

    onChange(newQuery);
  };

  onSetPropertyAlias = (propertyAlias?: string) => {
    const { onChange, query } = this.props;
    if (!propertyAlias) {
      onChange({ ...query, propertyAliases: undefined });
    } else if (query.propertyAliases) {
      onChange({ ...query, propertyAliases: [...query.propertyAliases, propertyAlias] });
    } else {
      onChange({ ...query, propertyAliases: [propertyAlias] });
    }
  };

  onSetAssetId = (assetId?: string) => {
    const { onChange, query } = this.props;
    if (!assetId) {
      this.setState({ assetId: undefined });
      onChange({ ...query, assetIds: undefined });
    } else {
      // TODO: handle entering multiple externalIds with other assetIds
      const validId = uuidRegex.test(assetId) || assetId.startsWith('externalId:') || assetId.startsWith('$');
      const assetIds = validId ? [assetId] : [`externalId:${assetId}`];
      this.setState({ assetId: assetIds[0] });
      if (query.assetIds) {
        onChange({ ...query, assetIds: [...query.assetIds, ...assetIds] });
      } else {
        onChange({ ...query, assetIds });
      }
    }
  };

  onSetPropertyId = (propertyId?: string) => {
    const { onChange, query } = this.props;
    if (!propertyId) {
      onChange({ ...query, propertyIds: undefined });
    } else if (query.propertyIds) {
      onChange({ ...query, propertyIds: [...query.propertyIds, propertyId] });
    } else {
      onChange({ ...query, propertyIds: [propertyId] });
    }
  };

  onSetHierarchyId = (hierarchyId?: string) => {
    const { onChange, query } = this.props;
    onChange({ ...(query as any), hierarchyId });
  };

  onHierarchyIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query } = this.props;
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
  };

  onLastObservationChange = () => {
    const { onChange, query } = this.props;
    onChange({ ...query, lastObservation: !query.lastObservation });
  };

  onFlattenL4eChange = () => {
    const { onChange, query } = this.props;
    onChange({ ...query, flattenL4e: !query.flattenL4e });
  };

  renderAssociatedAsset(query: ListAssociatedAssetsQuery) {
    const { asset, loading } = this.state;
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
          isLoading={loading}
          options={hierarchies}
          value={current}
          onChange={this.onHierarchyIdChange}
          placeholder="Select..."
          allowCustomValue={true}
          backspaceRemovesValue={true}
          isClearable={true}
          isSearchable={true}
          onCreateOption={this.onSetHierarchyId}
          formatCreateLabel={(txt) => `Hierarchy Id: ${txt}`}
          menuPlacement="auto"
        />
      </EditorField>
    );
  }

  render() {
    const { query, datasource, onChange } = this.props;
    const { loading, assets, assetProperties, propertyAliases } = this.state;

    let currentPropertyAlias = propertyAliases.filter((alias) =>
      alias.value ? query.propertyAliases?.includes(alias.value) : false
    );
    if (currentPropertyAlias.length === 0 && query.propertyAliases?.length) {
      if (loading) {
        currentPropertyAlias = query.propertyAliases.map((alias) => ({ label: 'loading...', value: alias }));
      } else {
        currentPropertyAlias = query.propertyAliases.map((alias) => ({ label: alias, value: alias }));
      }
    }

    let currentAsset = assets.filter((asset) => (asset.value ? query.assetIds?.includes(asset.value) : false));
    if (currentAsset.length === 0 && query.assetIds?.length) {
      if (loading) {
        currentAsset = query.assetIds.map((assetId) => ({ label: 'loading...', value: assetId }));
      } else {
        currentAsset = query.assetIds.map((assetId) => ({ label: `ID: ${this.state.assetId}`, value: assetId }));
      }
    }

    const isAssociatedAssets = isListAssociatedAssetsQuery(query);
    const showProp = !!(!isAssociatedAssets && (query.propertyIds || query.assetIds));

    const showQuality = !!(
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
      if (loading) {
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

    const queryTooltip = (
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

    const assetTooltip = (
      <div>
        Set the asset ID. It can be either the actual ID in UUID format, or else "externalId:" followed by the external
        ID, if it has one.
        <LinkButton
          href="https://docs.aws.amazon.com/iot-sitewise/latest/userguide/object-ids.html#external-ids"
          target="_blank"
        >
          API Docs <Icon name="external-link-alt" />
        </LinkButton>
      </div>
    );

    return (
      <>
        {!isAssociatedAssets && (
          <EditorRow>
            <EditorField label="Property Alias" tooltip={queryTooltip} tooltipInteractive htmlFor="alias" width={80}>
              <Select
                id="alias"
                inputId="alias"
                aria-label="Property alias"
                isMulti={true}
                options={propertyAliases}
                value={currentPropertyAlias}
                onChange={this.onAliasChange}
                placeholder="optional alias that identifies the property, such as an OPC-UA server data stream path"
                allowCustomValue={true}
                isClearable={true}
                isSearchable={true}
                onCreateOption={this.onSetPropertyAlias}
                formatCreateLabel={(txt) => `Property Alias: ${txt}`}
                menuPlacement="auto"
              />
            </EditorField>
          </EditorRow>
        )}

        {(!query.propertyAliases?.length || isAssociatedAssets) && (
          <>
            <EditorRow>
              <EditorFieldGroup>
                <EditorField label="Asset" tooltip={assetTooltip} tooltipInteractive htmlFor="asset" width={30}>
                  <Select
                    id="asset"
                    inputId="asset"
                    aria-label="Asset"
                    isMulti={true}
                    key={query.region ?? DEFAULT_REGION}
                    isLoading={loading}
                    options={assets}
                    value={currentAsset}
                    onChange={(sel) => this.onAssetChange(sel)}
                    placeholder="Select an asset"
                    allowCustomValue={true}
                    isClearable={true}
                    isSearchable={true}
                    onCreateOption={this.onSetAssetId}
                    formatCreateLabel={(txt) => `Asset ID: ${txt}`}
                    menuPlacement="auto"
                  />
                </EditorField>

                <div className={styles.exploreContainer}>
                  <AssetBrowser
                    datasource={datasource}
                    region={query.region}
                    assetId={query.assetIds?.[0]}
                    onAssetChanged={this.onSetAssetId}
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
                      isMulti={true}
                      isLoading={loading}
                      options={assetProperties}
                      value={currentAssetProperty}
                      onChange={this.onPropertyChange}
                      placeholder="Select a property"
                      allowCustomValue={true}
                      isClearable={true}
                      isSearchable={true}
                      onCreateOption={this.onSetPropertyId}
                      formatCreateLabel={(txt) => `Property ID: ${txt}`}
                      menuPlacement="auto"
                    />
                  </EditorField>
                </EditorFieldGroup>

                <EditorFieldGroup>
                  {showQuality && isAssetPropertyAggregatesQuery(query) && (
                    <AggregationSettings query={query} onChange={this.props.onChange} />
                  )}

                  {isAssetPropertyInterpolatedQuery(query) && (
                    <InterpolatedResolutionSettings query={query} onChange={onChange} />
                  )}
                </EditorFieldGroup>
              </EditorRow>
            )}
          </>
        )}

        {!!query.propertyAliases?.length && isAssetPropertyAggregatesQuery(query) && (
          <EditorRow>
            <AggregationSettings query={query} onChange={this.props.onChange} />
          </EditorRow>
        )}

        {!!query.propertyAliases?.length && isAssetPropertyInterpolatedQuery(query) && (
          <EditorRow>
            <InterpolatedResolutionSettings query={query} onChange={onChange} />
          </EditorRow>
        )}

        {isAssociatedAssets && (
          <EditorRow>
            <EditorFieldGroup>{this.renderAssociatedAsset(query as ListAssociatedAssetsQuery)}</EditorFieldGroup>
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
                showQuality={!!(query.propertyIds?.length || query.propertyAliases?.length)}
                onLastObservationChange={this.onLastObservationChange}
                onFlattenL4eChange={this.onFlattenL4eChange}
              />
            </EditorFieldGroup>
          </EditorRow>
        )}
      </>
    );
  }
}

const styles = {
  exploreContainer: css({ display: 'flex', alignItems: 'flex-end' }),
};
