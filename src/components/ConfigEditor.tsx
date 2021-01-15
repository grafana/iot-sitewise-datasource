import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData, SitewiseQuery } from '../types';
import ConnectionConfig from '../common/ConnectionConfig';
import { Alert, TLSAuthSettings } from '@grafana/ui';

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
    const { options } = this.props;
    const jsonData = options.jsonData;
    const { defaultRegion, endpoint } = jsonData;

    return (
      <>
        <div>
          <ConnectionConfig {...this.props} />

          {defaultRegion === 'Edge' && (
            <div>
              {!endpoint && (
                <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" />
              )}
              TODO... show TLS config editor
            </div>
          )}
        </div>
      </>
    );
  }
}
