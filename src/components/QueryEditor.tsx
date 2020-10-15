import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../DataSource';
import { SitewiseQuery, SitewiseOptions, QueryType } from '../types';
import { Select } from '@grafana/ui';
import { QueryTypeInfo, siteWisteQueryTypes, changeQueryType } from '../queryInfo';
import { QueryInlineField } from './Forms';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTypeChange = (sel: SelectableValue<QueryType>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange(changeQueryType(query, sel as QueryTypeInfo));
    onRunQuery();
  };

  renderQuery(query: SitewiseQuery) {
    if (!query.queryType) {
      return;
    }
    switch (query.queryType) {
    }
    return <div>Missing UI for query type: {query.queryType}</div>;
  }

  render() {
    const { query } = this.props;

    return (
      <>
        <div className={'gf-form'}>
          <QueryInlineField label="Query type">
            <Select
              options={siteWisteQueryTypes}
              value={siteWisteQueryTypes.find(v => v.value === query.queryType)}
              onChange={this.onQueryTypeChange}
              placeholder="Select query type"
            />
            <div>TODO: region picker</div>
          </QueryInlineField>
        </div>
        {this.renderQuery(query)}
      </>
    );
  }
}
