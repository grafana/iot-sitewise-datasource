import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../DataSource';
import { SitewiseQuery, SitewiseOptions } from '../types';

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;
interface State {
  schemaState?: Partial<SitewiseQuery>;
}

export class QueryEditor extends PureComponent<Props, State> {
  state: State = {};

  render() {
    return (
      <>
        <div>TODO -- query editor here</div>
      </>
    );
  }
}
