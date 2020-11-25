import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SitewiseQuery,
  AssetInfo,
  AssetPropertyAggregatesQuery,
  AggregateType,
  SiteWiseResolution,
  isAssetPropertyAggregatesQuery,
  isAssetPropertyValueHistoryQuery,
  AssetPropertyInfo,
  ListAssociatedAssetsQuery,
  isListAssociatedAssetsQuery,
} from 'types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetBrowser } from '../browser/AssetBrowser';
import { AggregatePicker, aggReg } from '../AggregatePicker';
import { getAssetProperty, getDefaultAggregate } from 'queryInfo';
import { QualityAndOrderRow } from './QualityAndOrderRow';
import { firstLabelWith } from './QueryEditor';

type Props = SitewiseQueryEditorProps<SitewiseQuery | AssetPropertyAggregatesQuery | ListAssociatedAssetsQuery>;

const resolutions: Array<SelectableValue<SiteWiseResolution>> = [
  { value: SiteWiseResolution.Auto, label: 'Auto', description: 'Pick a resolution based on the time window' },
  { value: SiteWiseResolution.Min, label: 'Minute', description: '1 point every minute' },
  { value: SiteWiseResolution.Hour, label: 'Hour', description: '1 point every hour' },
  { value: SiteWiseResolution.Day, label: 'Day', description: '1 point every day' },
];

interface State {
  asset?: AssetInfo;
  property?: AssetPropertyInfo;
  assets: Array<SelectableValue<string>>;
  loading: boolean;
  openModal: boolean;
}

export class PropertyQueryEditor extends PureComponent<Props, State> {
  state: State = {
    assets: [],
    loading: true,
    openModal: false,
  };

  async updateInfo() {
    const { query, datasource } = this.props;
    const update: State = {
      loading: false,
    } as State;

    const cache = datasource.getCache(query.region);
    if (query?.assetId) {
      try {
        update.asset = await cache.getAssetInfo(query.assetId);
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
    const assetChanged = query?.assetId !== oldProps?.query?.assetId;
    const propChanged = query?.propertyId !== oldProps?.query?.propertyId;
    if (assetChanged || propChanged) {
      if (!query.assetId) {
        this.setState({ asset: undefined, property: undefined, loading: false });
      } else {
        this.setState({ loading: true });
        this.updateInfo();
      }
    }
  }

  onAssetChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId: sel.value! });
    onRunQuery();
  };

  onPropertyChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    const update = { ...query, propertyId: sel.value! };
    // Make sure the selected aggregates are actually supported
    if (isAssetPropertyAggregatesQuery(update)) {
      if (update.propertyId) {
        const info = getAssetProperty(this.state.asset, update.propertyId);
        if (!update.aggregates) {
          update.aggregates = [];
        }
        if (info) {
          update.aggregates = update.aggregates.filter(a => aggReg.get(a).isValid(info));
        }
        if (!update.aggregates.length) {
          update.aggregates = [getDefaultAggregate(info)];
        }
      }
    }
    onChange(update);
    onRunQuery();
  };

  onSetAssetId = (assetId?: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId });
    onRunQuery();
  };

  onSetPropertyId = (propertyId?: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, propertyId });
    onRunQuery();
  };

  onSetHierarchyId = (hierarchyId?: string) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...(query as any), hierarchyId });
    onRunQuery();
  };

  onHierarchyIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    const update = { ...query };
    if (isListAssociatedAssetsQuery(update)) {
      if (sel.value && sel.value.length) {
        update.hierarchyId = sel.value;
      } else {
        delete update.hierarchyId;
      }
    }
    onChange(update);
    onRunQuery();
  };

  //--------------------------------------------------------------------------------
  //
  //--------------------------------------------------------------------------------

  onAggregateChange = (aggregates: AggregateType[]) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, aggregates } as any);
    onRunQuery();
  };

  onResolutionChange = (sel: SelectableValue<SiteWiseResolution>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, resolution: sel.value } as any);
    onRunQuery();
  };

  renderAggregateRow(query: AssetPropertyAggregatesQuery) {
    const { property } = this.state;
    return (
      <div className="gf-form">
        <InlineField label="Aggregate" labelWidth={firstLabelWith} grow={true}>
          <AggregatePicker
            stats={query.aggregates ?? []}
            onChange={this.onAggregateChange}
            defaultStat={getDefaultAggregate(property)}
            menuPlacement="bottom"
          />
        </InlineField>
        <InlineField label="Resolution" labelWidth={10}>
          <Select
            width={18}
            options={resolutions}
            value={resolutions.find(v => v.value === query.resolution) || resolutions[0]}
            onChange={this.onResolutionChange}
            menuPlacement="bottom"
          />
        </InlineField>
      </div>
    );
  }

  renderAssociatedAsset(query: ListAssociatedAssetsQuery) {
    const { asset, loading } = this.state;
    const hierarchies: Array<SelectableValue<string>> = [{ value: '', label: '** Parent **' }];
    if (asset) {
      hierarchies.push(...asset.hierarchy);
    }

    let current = hierarchies.find(v => v.value === query.hierarchyId);
    if (!current) {
      if (query.hierarchyId) {
        current = { value: query.hierarchyId, label: 'ID: ' + query.hierarchyId };
        hierarchies.push(current);
      } else {
        current = hierarchies[0]; // parent
      }
    }

    return (
      <div className="gf-form">
        <InlineField label="Show" labelWidth={firstLabelWith} grow={true}>
          <Select
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
            formatCreateLabel={txt => `Hierarchy Id: ${txt}`}
            menuPlacement="bottom"
          />
        </InlineField>
      </div>
    );
  }

  render() {
    const { query, datasource } = this.props;
    const { loading, asset, assets } = this.state;

    let current = query.assetId ? assets.find(v => v.value === query.assetId) : undefined;
    if (!current && query.assetId) {
      if (loading) {
        current = { label: 'loading...', value: query.assetId };
      } else if (asset) {
        current = { label: asset.name, description: query.assetId, value: query.assetId };
      } else {
        current = { label: `ID: ${query.assetId}`, value: query.assetId };
      }
    }

    const isAssociatedAssets = isListAssociatedAssetsQuery(query);
    const showProp = !isAssociatedAssets && (query.propertyId || query.assetId);
    const properties = showProp ? (asset ? asset.properties : []) : [];
    const showQuality =
      (query.propertyId && isAssetPropertyAggregatesQuery(query)) || isAssetPropertyValueHistoryQuery(query);

    let currentProperty = properties.find(p => p.Id === query.propertyId);
    if (!currentProperty && query.propertyId) {
      currentProperty = {
        value: query.propertyId,
        label: 'ID: ' + query.propertyId,
      } as AssetPropertyInfo;
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Asset" labelWidth={firstLabelWith} grow={true}>
            <Select
              key={query.region ? query.region : 'default'}
              isLoading={loading}
              options={assets}
              value={current}
              onChange={this.onAssetChange}
              placeholder="Select an asset"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              onCreateOption={this.onSetAssetId}
              formatCreateLabel={txt => `Asset ID: ${txt}`}
              menuPlacement="bottom"
            />
          </InlineField>
          <AssetBrowser
            datasource={datasource}
            region={query.region}
            assetId={query.assetId}
            onAssetChanged={this.onSetAssetId}
          />
        </div>
        {showProp && (
          <>
            <div className="gf-form">
              <InlineField label="Property" labelWidth={firstLabelWith} grow={true}>
                <Select
                  isLoading={loading}
                  options={properties}
                  value={currentProperty}
                  onChange={this.onPropertyChange}
                  placeholder="Select a property"
                  allowCustomValue={true}
                  isSearchable={true}
                  onCreateOption={this.onSetPropertyId}
                  formatCreateLabel={txt => `Property ID: ${txt}`}
                  menuPlacement="bottom"
                />
              </InlineField>
            </div>
            {showQuality && (
              <>
                {isAssetPropertyAggregatesQuery(query) && this.renderAggregateRow(query)}
                <QualityAndOrderRow {...(this.props as any)} />
              </>
            )}
          </>
        )}
        {isAssociatedAssets && this.renderAssociatedAsset(query as ListAssociatedAssetsQuery)}
      </>
    );
  }
}
