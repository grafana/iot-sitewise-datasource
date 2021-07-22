import React, { PureComponent } from 'react';
import {
  onUpdateDatasourceResetOption,
  SelectableValue,
  updateDatasourcePluginJsonDataOption,
  updateDatasourcePluginSecureJsonDataOption,
} from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData } from '../types';
import { ConnectionConfig, ConnectionConfigProps } from '@grafana/aws-sdk';
import { Alert, Button, InlineField, InlineFieldRow, Select } from '@grafana/ui';
import { standardRegions } from '../regions';

export type Props = ConnectionConfigProps<SitewiseOptions, SitewiseSecureJsonData>;

const edgeAuthMethods: Array<SelectableValue<string>> = [
  { value: 'default', label: 'Default', description: 'Default AWS authentication methods.' },
  { value: 'linux', label: 'Linux', description: 'Linux-based authentication' },
  { value: 'ldap', label: 'LDAP', description: 'LDAP-based authentication' },
];
export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  render() {
    const { options } = this.props;
    const jsonData = options.jsonData;
    const { defaultRegion, endpoint } = jsonData;
    const edgeAuthMode = edgeAuthMethods.find((f) => f.value === jsonData.edgeAuthMode) ?? edgeAuthMethods[0];

    return (
      <>
        <div>
          <ConnectionConfig {...this.props} standardRegions={standardRegions} />

          {defaultRegion === 'Edge' && (
            <div>
              {!endpoint && (
                <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" />
              )}
              <InlineFieldRow>
                <InlineField label="Auth Mode"
                labelWidth={28}
                tooltip="Specify which authentication method to use.">
                  <Select
                    width={30}
                    options={edgeAuthMethods}
                    value={edgeAuthMode}
                    onChange={(v) => {
                      updateDatasourcePluginJsonDataOption(this.props as any, 'edgeAuthMode' as never, v.value);
                    }}
                  />
                </InlineField>
              </InlineFieldRow>
              <InlineFieldRow>
                <InlineField label="Certification"
                labelWidth={28}
                tooltip="Certificate for SSL enabled authentication.">
                  {options.secureJsonFields?.cert ? (
                    <Button
                      variant="secondary"
                      type="reset"
                      onClick={onUpdateDatasourceResetOption(this.props as any, 'cert')}
                    >
                      Reset
                    </Button>
                  ) : (
                    <textarea
                      rows={7}
                      className="gf-form-input gf-form-textarea width-30"
                      onChange={(event) => {
                        updateDatasourcePluginSecureJsonDataOption(this.props as any, 'cert', event.target.value);
                      }}
                      placeholder="Begins with -----BEGIN CERTIFICATE------"
                      required
                    />
                  )}
                </InlineField>
              </InlineFieldRow>
            </div>
          )}
        </div>
      </>
    );
  }
}
