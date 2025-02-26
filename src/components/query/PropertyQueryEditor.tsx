import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SitewiseQuery,
  AssetInfo,
  AssetPropertyAggregatesQuery,
  isAssetPropertyAggregatesQuery,
  isAssetPropertyValueHistoryQuery,
  AssetPropertyInfo,
  ListAssociatedAssetsQuery,
  isListAssociatedAssetsQuery,
  isAssetPropertyInterpolatedQuery,
  shouldShowOptionsRow,
} from 'types';
import { LinkButton, Select, Input, Icon } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetBrowser } from '../browser/AssetBrowser';
import { aggReg } from './AggregationSettings/AggregatePicker';
import { getAssetProperty, getDefaultAggregate } from 'queryInfo';
import { QualityAndOrderRow } from './QualityAndOrderRow';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { css } from '@emotion/css';
import { QueryOptions } from './QueryOptions';
import { AggregationSettings } from './AggregationSettings/AggregationSettings';

type Props = SitewiseQueryEditorProps<SitewiseQuery | AssetPropertyAggregatesQuery | ListAssociatedAssetsQuery>;

const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89abAB][0-9a-f]{3}-[0-9a-f]{12}$/i;
interface State {
  assetId?: string;
  asset?: AssetInfo;
  property?: AssetPropertyInfo;
  assets: Array<SelectableValue<string>>;
  assetProperties: Array<{ id: string; name: string }>;
  loading: boolean;
  openModal: boolean;
}

const ALL_HIERARCHIES = '*';

export class PropertyQueryEditor extends PureComponent<Props, State> {
  state: State = {
    assetId: this.props.query.assetIds && this.props.query.assetIds[0],
    assets: [],
    assetProperties: [],
    loading: true,
    openModal: false,
  };

  async updateInfo() {
    const { onChange, query, datasource } = this.props;
    const update: State = {
      loading: false,
    } as State;

    const cache = datasource.getCache(query.region);
    if (query?.assetIds?.length) {
      try {
        update.asset = await cache.getAssetInfo(query.assetIds[0]);
        const ps = await cache.listAssetProperties(query.assetIds[0]);
        update.assetProperties = ps?.map(({ id, name }) => ({ id, name })) || [];
        // Update external ids to asset ids in the query
        if (update.asset?.id && query.assetIds[0].startsWith('externalId')) {
          query.assetIds[0] = update.asset.id;
          onChange(query);
        }
      } catch (err) {
        console.warn('error reading asset info', err);
        update.property = undefined;
      }
    }
    update.property = getAssetProperty(update.asset, query.propertyId);

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
    const propChanged = query?.propertyId !== oldProps?.query?.propertyId;
    const regionChanged = query?.region !== oldProps?.query?.region;

    if (assetChanged || propChanged || regionChanged) {
      if (!query.assetIds?.length && !regionChanged) {
        this.setState({ assetId: undefined, asset: undefined, property: undefined, loading: false });
      } else {
        this.setState({ loading: true });
        this.updateInfo();
      }
    }
  }

  onAliasChange = (evt: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, propertyAlias: evt.currentTarget.value });
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
            propertyId: undefined,
          }
        : {
            ...query,
            assetIds: [...assetIds],
          };

    onChange(newQuery);
  }

  onPropertyChange = (sel: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    const update = { ...query, propertyId: sel.value! };

    // Make sure the selected aggregates are actually supported
    if (isAssetPropertyAggregatesQuery(update)) {
      if (update.propertyId) {
        const info = getAssetProperty(this.state.asset, update.propertyId);

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
  };

  onSetAssetId = (assetId?: string) => {
    const { onChange, query } = this.props;
    if (!assetId) {
      this.setState({ assetId: undefined });
      onChange({ ...query, assetIds: undefined });
    } else {
      const validId = uuidRegex.test(assetId) || assetId.startsWith('externalId:') || assetId.startsWith('$');
      const assetIds = validId ? [assetId] : [`externalId:${assetId}`];
      this.setState({ assetId: assetIds[0] });
      onChange({ ...query, assetIds });
    }
  };

  onSetPropertyId = (propertyId?: string) => {
    const { onChange, query } = this.props;
    onChange({ ...query, propertyId });
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
    const { query, datasource } = this.props;
    const { loading, assets, assetProperties } = this.state;

    let current = assets.filter((a) => (a.value ? query.assetIds?.includes(a.value) : false));
    if (current.length === 0 && query.assetIds?.length) {
      if (loading) {
        current = query.assetIds.map((assetId) => ({ label: 'loading...', value: assetId }));
      } else {
        current = query.assetIds.map((assetId) => ({ label: `ID: ${this.state.assetId}`, value: assetId }));
      }
    }

    const isAssociatedAssets = isListAssociatedAssetsQuery(query);
    const showProp = !!(!isAssociatedAssets && (query.propertyId || query.assetIds));
    const assetPropertyOptions = showProp
      ? assetProperties.map(({ id, name }) => ({ id, name, value: id, label: name }))
      : [];

    const showQuality = !!(
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
              <Input
                id="alias"
                aria-label="Property alias"
                value={query.propertyAlias}
                onChange={this.onAliasChange}
                placeholder="optional alias that identifies the property, such as an OPC-UA server data stream path"
              />
            </EditorField>
          </EditorRow>
        )}

        {(!Boolean(query.propertyAlias) || isAssociatedAssets) && (
          <>
            <EditorRow>
              <EditorFieldGroup>
                <EditorField label="Asset" tooltip={assetTooltip} tooltipInteractive htmlFor="asset" width={30}>
                  <Select
                    id="asset"
                    inputId="asset"
                    aria-label="Asset"
                    isMulti={true}
                    key={query.region ? query.region : 'default'}
                    isLoading={loading}
                    options={assets}
                    value={current}
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
                      isLoading={loading}
                      options={assetPropertyOptions}
                      value={currentAssetPropertyOption ?? null}
                      onChange={this.onPropertyChange}
                      placeholder="Select a property"
                      allowCustomValue={true}
                      isSearchable={true}
                      onCreateOption={this.onSetPropertyId}
                      formatCreateLabel={(txt) => `Property ID: ${txt}`}
                      menuPlacement="auto"
                    />
                  </EditorField>
                </EditorFieldGroup>

                <EditorFieldGroup>
                  {showQuality && isAssetPropertyAggregatesQuery(query) && (
                    <AggregationSettings query={query} property={this.state.property} onChange={this.props.onChange} />
                  )}
                </EditorFieldGroup>
              </EditorRow>
            )}
          </>
        )}

        {query.propertyAlias && isAssetPropertyAggregatesQuery(query) && (
          <EditorRow>
            <AggregationSettings query={query} property={this.state.property} onChange={this.props.onChange} />
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
                showProp={showProp}
                showQuality={!!(query.propertyId || query.propertyAlias)}
                onLastObservationChange={this.onLastObservationChange}
                onFlattenL4eChange={this.onFlattenL4eChange}
                qualityAndOrderComponent={<QualityAndOrderRow {...(this.props as any)} />}
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
