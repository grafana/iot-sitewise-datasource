import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../DataSource';
import { SitewiseQuery, SitewiseOptions, QueryType, AssetPropertyValueQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWisteQueryTypes, changeQueryType } from '../queryInfo';
import { standardRegions } from 'common/types';
import { QueryPropertyValueEditor } from './QueryPropertyValueEditor';

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
        return <div>Maybe add a search value?</div>;
      case QueryType.ListAssets:
        return <div>TODO: pick asset model</div>;
      case QueryType.PropertyValue:
        return <QueryPropertyValueEditor {...this.props} query={query as AssetPropertyValueQuery} />;
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
