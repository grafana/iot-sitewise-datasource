import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../DataSource';
import { SitewiseQuery, SitewiseOptions, QueryType, AssetPropertyValueQuery, ListAssetsQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWisteQueryTypes, changeQueryType } from '../queryInfo';
import { standardRegions } from 'common/types';
import { AssetPropertyValueQueryEditor } from './AssetPropertyValueQueryEditor';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';

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
        return <AssetPropertyValueQueryEditor {...this.props} query={query as AssetPropertyValueQuery} />;
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
          <InlineField label="Region">
            <Select
              width={20}
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
