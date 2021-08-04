import React, { PureComponent, ChangeEvent } from 'react';
import {
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceJsonDataOptionSelect,
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

  renderEdgeConfig() {
    const { options } = this.props;
    const { jsonData } = options;
    const { endpoint } = jsonData;

    const edgeAuthMode = edgeAuthMethods.find((f) => f.value === jsonData.edgeAuthMode) ?? edgeAuthMethods[0];
    const hasEdgeAuth = edgeAuthMode !== edgeAuthMethods[0];
    const labelWidth = 28;
    const regions = standardRegions.map((value) => ({ value, label: value }));

    return (
      <>
        {hasEdgeAuth && (
          <FieldSet label={'Connection Details'} data-testid="connection-config">
            <InlineField
              label="Endpoint"
              labelWidth={28}
              tooltip="Optionally, specify a custom endpoint for the service"
            >
              <Input
                className="width-30"
                placeholder={this.props.defaultEndpoint ?? 'https://{service}.{region}.amazonaws.com'}
                value={options.jsonData.endpoint || ''}
                onChange={onUpdateDatasourceJsonDataOption(this.props, 'endpoint')}
              />
            </InlineField>
            <InlineField
              label="Default Region"
              labelWidth={28}
              tooltip="Specify the region, such as for US West (Oregon) use ` us-west-2 ` as the region."
            >
              <Select
                className="width-30"
                value={regions.find((region) => region.value === options.jsonData.defaultRegion)}
                options={regions}
                defaultValue={options.jsonData.defaultRegion}
                allowCustomValue={true}
                onChange={onUpdateDatasourceJsonDataOptionSelect(this.props, 'defaultRegion')}
                formatCreateLabel={(r) => `Use region: ${r}`}
              />
            </InlineField>
          </FieldSet>
        )}
        {!hasEdgeAuth && <ConnectionConfig {...this.props} standardRegions={standardRegions} />}

        <FieldSet label={'Edge settings'} data-testid="edge-connection">
          {!endpoint && <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" />}
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
                <InlineField
                  label="Username"
                  labelWidth={labelWidth}
                  tooltip="The username set to local authentication proxy"
                >
                  <Input
                    name="username"
                    value={jsonData.edgeAuthUser}
                    autoComplete="off"
                    className="width-30"
                    onChange={this.onUserChange}
                    required
                  />
                </InlineField>
              </InlineFieldRow>
              <InlineFieldRow>
                <InlineField
                  label="Password"
                  labelWidth={labelWidth}
                  tooltip="The password sent to local authenticaion proxy"
                >
                  <Input
                    type="password"
                    name="password"
                    autoComplete="off"
                    placeholder={options.secureJsonFields?.edgeAuthPass ? 'configured' : ''}
                    value={options.secureJsonData?.edgeAuthPass ?? ''}
                    onChange={this.onPasswordChange}
                    onReset={this.onResetPassword}
                    className="width-30"
                    required
                  />
                </InlineField>
              </InlineFieldRow>
            </>
          )}
          <InlineFieldRow>
            <InlineField
              label="SSL Certificate"
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
      </>
    );
  }

  // Simple
  render() {
    const { options } = this.props;
    if (options.jsonData.defaultRegion === 'Edge') {
      return this.renderEdgeConfig();
    }
    return <ConnectionConfig {...this.props} standardRegions={standardRegions} />;
  }
}
