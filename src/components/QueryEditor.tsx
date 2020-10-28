import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../DataSource';
import {
  SitewiseQuery,
  SitewiseOptions,
  QueryType,
  ListAssetsQuery,
  AssetPropertyValueHistoryQuery,
  AssetPropertyAggregatesQuery,
} from '../types';
import { InlineField, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWisteQueryTypes, changeQueryType } from '../queryInfo';
import { standardRegions } from 'common/types';
import { PropertyValueHistoryEditor } from './PropertyValueHistoryEditor';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';
import { AssetPropPickerRows } from './AssetPropPickerRows';
import { PropertyAggregatesEditor } from './PropertyAggregatesEditor';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTypeChange = (sel: SelectableValue<QueryType>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange(changeQueryType(query, sel as QueryTypeInfo));
    onRunQuery();
  };

  onRegionChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, region: sel.value });
    onRunQuery();
  };

  renderQuery(query: SitewiseQuery) {
    if (!query.queryType) {
      return;
    }
    switch (query.queryType) {
      case QueryType.ListAssetModels:
        return null; // nothing required
      case QueryType.ListAssets:
        return <ListAssetsQueryEditor {...this.props} query={query as ListAssetsQuery} />;
      case QueryType.PropertyValue:
        return <AssetPropPickerRows {...this.props} />;
      case QueryType.PropertyValueHistory:
        return <PropertyValueHistoryEditor {...this.props} query={query as AssetPropertyValueHistoryQuery} />;
      case QueryType.PropertyAggregate:
        return <PropertyAggregatesEditor {...this.props} query={query as AssetPropertyAggregatesQuery} />;
    }
    return <div>Missing UI for query type: {query.queryType}</div>;
  }

  render() {
    const { query, datasource } = this.props;

    const defaultRegion = { label: `Default`, desctiption: datasource.options?.defaultRegion, value: undefined };
    const regions = query.region ? [defaultRegion, ...standardRegions] : standardRegions;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={10} grow={true}>
            <Select
              options={siteWisteQueryTypes}
              value={siteWisteQueryTypes.find(v => v.value === query.queryType)}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
            />
          </InlineField>
          <InlineField label="Region" labelWidth={10}>
            <Select
              width={18}
              options={regions}
              value={standardRegions.find(v => v.value === query.region) || defaultRegion}
              onChange={this.onRegionChange}
              backspaceRemovesValue={true}
              allowCustomValue={true}
              isClearable={true}
            />
          </InlineField>
        </div>
        {this.renderQuery(query)}
      </>
    );
  }
}
