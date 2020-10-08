import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData, SitewiseQuery } from '../types';
import CommonConfig from '../common/CommonConfig';

export type Props = DataSourcePluginOptionsEditorProps<SitewiseOptions, SitewiseSecureJsonData>;

interface State {
  schemaState?: Partial<SitewiseQuery>;
}

export class ConfigEditor extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <>
        <div>
          <CommonConfig {...this.props} />
        </div>
      </>
    );
  }
}
