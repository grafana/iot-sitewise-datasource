import defaults from 'lodash/defaults';
import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from 'DataSource';
import { SitewiseQuery, SitewiseOptions, QueryType, ListAssetsQuery } from 'types';
import { Icon, InlineField, LinkButton, Select } from '@grafana/ui';
import { QueryTypeInfo, siteWisteQueryTypes, changeQueryType } from 'queryInfo';
import { standardRegions } from 'common/regions';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';
import { PropertyQueryEditor } from './PropertyQueryEditor';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

const queryDefaults: Partial<SitewiseQuery> = {
  maxPageAggregations: 1,
};

export const firstLabelWith = 14;

export class QueryEditor extends PureComponent<Props> {
  onQueryTypeChange = (sel: SelectableValue<QueryType>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange(changeQueryType(query, sel as QueryTypeInfo));
    onRunQuery();
  };

  onRegionChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId: undefined, propertyId: undefined, region: sel.value });
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
      case QueryType.ListAssociatedAssets:
      case QueryType.PropertyValue:
      case QueryType.PropertyAggregate:
      case QueryType.PropertyValueHistory:
        return <PropertyQueryEditor {...this.props} />;
    }
    return <div>Missing UI for query type: {query.queryType}</div>;
  }

  render() {
    const { datasource } = this.props;
    const query = defaults(this.props.query, queryDefaults);

    const defaultRegion = { label: `Default`, desctiption: datasource.options?.defaultRegion, value: undefined };
    const regions = query.region ? [defaultRegion, ...standardRegions] : standardRegions;
    const currentQueryType = siteWisteQueryTypes.find((v) => v.value === query.queryType);
    const queryTooltip = currentQueryType ? (
      <div>
        {currentQueryType.description} <br />
        <LinkButton href={currentQueryType.helpURL} target="_blank">
          API Docs <Icon name="external-link-alt" />
        </LinkButton>
      </div>
    ) : undefined;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Query type" labelWidth={14} grow={true} tooltip={queryTooltip}>
            <Select
              options={siteWisteQueryTypes}
              value={currentQueryType}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Region" labelWidth={14}>
            <Select
              width={18}
              options={regions}
              value={standardRegions.find((v) => v.value === query.region) || defaultRegion}
              onChange={this.onRegionChange}
              backspaceRemovesValue={true}
              allowCustomValue={true}
              isClearable={true}
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
        {this.renderQuery(query)}
      </>
    );
  }
}
