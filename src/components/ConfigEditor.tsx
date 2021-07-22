import React, { PureComponent, ChangeEvent } from 'react';
import {
  onUpdateDatasourceResetOption,
  SelectableValue,
  updateDatasourcePluginJsonDataOption,
  updateDatasourcePluginSecureJsonDataOption,
} from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData } from '../types';
import { ConnectionConfig, ConnectionConfigProps } from '@grafana/aws-sdk';
import { Alert, Button, FieldSet, InlineField, InlineFieldRow, Input, Select } from '@grafana/ui';
import { standardRegions } from '../regions';

export type Props = ConnectionConfigProps<SitewiseOptions, SitewiseSecureJsonData>;

const edgeAuthMethods: Array<SelectableValue<string>> = [
  { value: 'default', label: 'Standard', description: 'Use the authentication provider configured above' },
  { value: 'linux', label: 'Linux', description: 'Linux-based authentication' },
  { value: 'ldap', label: 'LDAP', description: 'LDAP-based authentication' },
];
export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  onUserChange = (event: ChangeEvent<HTMLInputElement>) => {
    updateDatasourcePluginJsonDataOption(this.props, 'edgeAuthUser', event.target.value);
  };

  onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { options, onOptionsChange } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        edgeAuthPass: event.target.value,
      },
    });
  };

  onResetPassword = () => {
    const { options, onOptionsChange } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        edgeAuthPass: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const jsonData = options.jsonData;
    const { defaultRegion, endpoint } = jsonData;
    const edgeAuthMode = edgeAuthMethods.find((f) => f.value === jsonData.edgeAuthMode) ?? edgeAuthMethods[0];
    const hasEdgeAuth = edgeAuthMode !== edgeAuthMethods[0];
    const labelWidth = 28;

    return (
      <>
        <div>
          <ConnectionConfig {...this.props} standardRegions={standardRegions} />

          {defaultRegion === 'Edge' && (
            <FieldSet label={'Edge settings'} data-testid="connection-config">
              {!endpoint && (
                <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" />
              )}
              <InlineFieldRow>
                <InlineField
                  label="Authentication Mode"
                  labelWidth={labelWidth}
                  tooltip="Specify which authentication method to use."
                >
                  <Select
                    className="width-30"
                    options={edgeAuthMethods}
                    value={edgeAuthMode}
                    onChange={(v) => {
                      updateDatasourcePluginJsonDataOption(this.props, 'edgeAuthMode', v.value);
                    }}
                  />
                </InlineField>
              </InlineFieldRow>
              {hasEdgeAuth && (
                <>
                  <InlineFieldRow>
                    <InlineField label="Username" labelWidth={labelWidth} tooltip="Specify the username to use.">
                      <Input
                        name="username"
                        value={jsonData.edgeAuthUser}
                        autoComplete="off"
                        className="width-30"
                        onChange={this.onUserChange}
                      />
                    </InlineField>
                  </InlineFieldRow>
                  <InlineFieldRow>
                    <InlineField label="Password" labelWidth={labelWidth} tooltip="Specify the password to use.">
                      <Input
                        type="password"
                        name="password"
                        autoComplete="off"
                        placeholder={options.secureJsonFields?.edgeAuthPass ? 'configured' : ''}
                        value={options.secureJsonData?.edgeAuthPass ?? ''}
                        onChange={this.onPasswordChange}
                        onReset={this.onResetPassword}
                        className="width-30"
                      />
                    </InlineField>
                  </InlineFieldRow>
                </>
              )}
              <InlineFieldRow>
                <InlineField
                  label="Certification"
                  labelWidth={labelWidth}
                  tooltip="Certificate for SSL enabled authentication."
                >
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
            </FieldSet>
          )}
        </div>
      </>
    );
  }
}
