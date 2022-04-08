import React, { PureComponent } from 'react';
import { ListAssetsQueryEditor } from './ListAssetsQueryEditor';
import { SitewiseQueryEditorProps } from './types';
import { SelectableValue } from '@grafana/data';
import { ListAssetsQuery } from 'types';

type Props = SitewiseQueryEditorProps<ListAssetsQuery>;

export default class VariableQueryEditor extends PureComponent<Props> {
  onVariableChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, filter: sel.value as 'ALL' | 'TOP_LEVEL' });
    onRunQuery();
  };
  render() {
    let { query } = this.props;
    return <ListAssetsQueryEditor {...this.props} query={query as ListAssetsQuery} onChange={this.onVariableChange} />;
  }
}
