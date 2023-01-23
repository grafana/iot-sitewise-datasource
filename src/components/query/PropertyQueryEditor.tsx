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
  isAssetPropertyInterpolatedQuery,
  shouldShowLastObserved,
} from 'types';
import { InlineField, LinkButton, Select, Input, Icon, InlineSwitch } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetBrowser } from '../browser/AssetBrowser';
import { AggregatePicker, aggReg } from '../AggregatePicker';
import { getAssetProperty, getDefaultAggregate } from 'queryInfo';
import { QualityAndOrderRow } from './QualityAndOrderRow';
import { firstLabelWith } from './QueryEditor';

type Props = SitewiseQueryEditorProps<SitewiseQuery | AssetPropertyAggregatesQuery | ListAssociatedAssetsQuery>;

const resolutions: Array<SelectableValue<SiteWiseResolution>> = [
  {
    value: SiteWiseResolution.Auto,
    label: 'Auto',
    description:
      'Picks a resolution based on the time window. ' +
      'Will switch to raw data if higher than 1m resolution is needed',
  },
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
    if ((query?.assetIds?.length ?? []) > 0) {
      try {
        update.asset = await cache.getAssetInfo(query.assetIds![0]);
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
        this.setState({ asset: undefined, property: undefined, loading: false });
      } else {
        this.setState({ loading: true });
        this.updateInfo();
      }
    }
  }

  onAliasChange = (evt: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, propertyAlias: evt.currentTarget.value });
    onRunQuery();
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
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetIds: [...assetIds] }); // TODO append????
    onRunQuery();
  }

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
          update.aggregates = update.aggregates.filter((a) => aggReg.get(a).isValid(info));
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
    onChange({ ...query, assetIds: assetId ? [assetId] : undefined });
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

  onLastObservationChange = () => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, lastObservation: !query.lastObservation });
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
            value={resolutions.find((v) => v.value === query.resolution) || resolutions[0]}
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

    let current = hierarchies.find((v) => v.value === query.hierarchyId);
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
            formatCreateLabel={(txt) => `Hierarchy Id: ${txt}`}
            menuPlacement="bottom"
          />
        </InlineField>
      </div>
    );
  }

  render() {
    const { query, datasource } = this.props;
    const { loading, asset, assets } = this.state;

    let current = assets.filter((a) => (a.value ? query.assetIds?.includes(a.value) : false));
    if (current.length === 0 && query.assetIds?.length) {
      if (loading) {
        current = query.assetIds.map((assetId) => ({ label: 'loading...', value: assetId }));
      } else {
        current = query.assetIds.map((assetId) => ({ label: `ID: ${assetId}`, value: assetId }));
      }
    }

    const isAssociatedAssets = isListAssociatedAssetsQuery(query);
    const showProp = !isAssociatedAssets && (query.propertyId || query.assetIds);
    const properties = showProp ? (asset ? asset.properties : []) : [];
    const showQuality =
      (query.propertyId && isAssetPropertyAggregatesQuery(query)) ||
      isAssetPropertyValueHistoryQuery(query) ||
      isAssetPropertyInterpolatedQuery(query);

    let currentProperty = properties.find((p) => p.Id === query.propertyId);
    if (!currentProperty && query.propertyId) {
      currentProperty = {
        value: query.propertyId,
        label: 'ID: ' + query.propertyId,
      } as AssetPropertyInfo;
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

    return (
      <>
        <div className="gf-form">
          <InlineField label="Property Alias" labelWidth={firstLabelWith} grow={true} tooltip={queryTooltip}>
            <Input
              value={query.propertyAlias}
              onChange={this.onAliasChange}
              placeholder="optional alias that identifies the property, such as an OPC-UA server data stream path"
            />
          </InlineField>
        </div>
        {!Boolean(query.propertyAlias) && (
          <>
            <div className="gf-form">
              <InlineField label="Asset" labelWidth={firstLabelWith} grow={true}>
                <Select
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
                  menuPlacement="bottom"
                />
              </InlineField>
              <AssetBrowser
                datasource={datasource}
                region={query.region}
                assetId={query.assetIds?.[0]}
                onAssetChanged={this.onSetAssetId}
              />
            </div>
            {showProp && (
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
                    formatCreateLabel={(txt) => `Property ID: ${txt}`}
                    menuPlacement="bottom"
                  />
                </InlineField>
                {shouldShowLastObserved(query.queryType) && (
                  <InlineField
                    label="Expand Time Range"
                    tooltip="Expand query to include last observed value before the current time range, and next observed value after the time range. "
                  >
                    <InlineSwitch value={query.lastObservation} onChange={this.onLastObservationChange} />
                  </InlineField>
                )}
              </div>
            )}
          </>
        )}
        {(showProp || query.propertyAlias) && showQuality && (
          <>
            {isAssetPropertyAggregatesQuery(query) && this.renderAggregateRow(query)}
            <QualityAndOrderRow {...(this.props as any)} />
          </>
        )}

        {isAssociatedAssets && this.renderAssociatedAsset(query as ListAssociatedAssetsQuery)}
      </>
    );
  }
}
