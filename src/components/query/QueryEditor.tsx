import React, { ChangeEvent } from 'react';
import {
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceJsonDataOptionSelect,
  onUpdateDatasourceResetOption,
  SelectableValue,
  updateDatasourcePluginJsonDataOption,
  updateDatasourcePluginSecureJsonDataOption,
} from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData } from '../types';
import { ConnectionConfig, ConnectionConfigProps, Divider } from '@grafana/aws-sdk';
import { Alert, Button, Field, FieldSet, InlineField, InlineFieldRow, Input, Select } from '@grafana/ui';
import { standardRegions } from '../regions';
import { config } from '@grafana/runtime';
import { ConfigSection } from '@grafana/experimental';

export type Props = ConnectionConfigProps<SitewiseOptions, SitewiseSecureJsonData>;

const edgeAuthMethods: Array<SelectableValue<string>> = [
  { value: 'default', label: 'Standard', description: 'Use the authentication provider configured above' },
  { value: 'linux', label: 'Linux', description: 'Linux-based authentication' },
  { value: 'ldap', label: 'LDAP', description: 'LDAP-based authentication' },
];

export function ConfigEditor(props: Props) {
  const onUserChange = (event: ChangeEvent<HTMLInputElement>) => {
    updateDatasourcePluginJsonDataOption(props, 'edgeAuthUser', event.target.value);
  };
  // @ts-ignore
  const newFormStylingEnabled = config.featureToggles.awsDatasourcesNewFormStyling;

  function onPasswordChange(event: ChangeEvent<HTMLInputElement>) {
    const { options, onOptionsChange } = props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        edgeAuthPass: event.target.value,
      },
    });
  }

  function onResetPassword() {
    const { options, onOptionsChange } = props;
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
  }

  function renderEdgeConfig() {
    const { options } = props;
    const { jsonData } = options;
    const { endpoint } = jsonData;

    const edgeAuthMode = edgeAuthMethods.find((f) => f.value === jsonData.edgeAuthMode) ?? edgeAuthMethods[0];
    const hasEdgeAuth = edgeAuthMode !== edgeAuthMethods[0];
    const labelWidth = 28;
    const regions = standardRegions.map((value) => ({ value, label: value }));

    if (newFormStylingEnabled) {
      return (
        <div className="width-30">
          {hasEdgeAuth && (
            <ConfigSection title="Connection Details" data-testid="connection-config">
              <Field
                label="Endpoint"
                description="Optionally, specify a custom endpoint for the service"
                htmlFor="endpoint"
              >
                <Input
                  id="endpoint"
                  placeholder={endpoint ?? 'https://{service}.{region}.amazonaws.com'}
                  value={endpoint || ''}
                  onChange={onUpdateDatasourceJsonDataOption(props, 'endpoint')}
                />
              </Field>
              <Field label="Default Region">
                <Select
                  value={regions.find((region) => region.value === options.jsonData.defaultRegion)}
                  options={regions}
                  defaultValue={options.jsonData.defaultRegion}
                  allowCustomValue={true}
                  onChange={onUpdateDatasourceJsonDataOptionSelect(props, 'defaultRegion')}
                  formatCreateLabel={(r) => `Use region: ${r}`}
                />
              </Field>
            </ConfigSection>
          )}
          {!hasEdgeAuth && (
            <ConnectionConfig {...props} standardRegions={standardRegions} newFormStylingEnabled={true} />
          )}

<Divider/>
          <ConfigSection title="Edge settings" data-testid="edge-settings">
            {!endpoint && (
              <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" data-testid="endpoint-warning" />
            )}

            <Field label="Authentication Mode" htmlFor="edgeAuthMethods">
              <Select
              id="edgeAuthMethods"
              aria-label="Authentication Mode"
                options={edgeAuthMethods}
                value={edgeAuthMode}
                onChange={(v) => {
                  updateDatasourcePluginJsonDataOption(props, 'edgeAuthMode', v.value);
                }}
              />
            </Field>
            {hasEdgeAuth && (
              <>
                <Field label="Username" description="The username set to local authentication proxy" htmlFor="username">
                  <Input
                  id="username"
                    name="username"
                    value={jsonData.edgeAuthUser}
                    autoComplete="off"
                    className="width-30"
                    onChange={onUserChange}
                    required
                  />
                </Field>
                <Field label="Password" description="The password sent to local authentication proxy" htmlFor="password">
                  <Input
                  id="password"
                    type="password"
                    name="password"
                    autoComplete="off"
                    placeholder={options.secureJsonFields?.edgeAuthPass ? 'configured' : ''}
                    value={options.secureJsonData?.edgeAuthPass ?? ''}
                    onChange={onPasswordChange}
                    onReset={onResetPassword}
                    className="width-30"
                    required
                  />
                </Field>
              </>
            )}
            <Field label="SSL Certificate" description="Certificate for SSL enabled authentication." htmlFor="certificate">
              {options.secureJsonFields?.cert ? (
                <Button variant="secondary" type="reset" onClick={onUpdateDatasourceResetOption(props as any, 'cert')} aria-label="Reset certificate input">
                  Reset
                </Button>
              ) : (
                <textarea
                id = "certificate"
                  rows={7}
                  className="gf-form-input gf-form-textarea width-30"
                  onChange={(event) => {
                    updateDatasourcePluginSecureJsonDataOption(props as any, 'cert', event.target.value);
                  }}
                  placeholder="Begins with -----BEGIN CERTIFICATE------"
                  required
                />
              )}
            </Field>
          </ConfigSection>
        </div>
      );
    } else {
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
                  placeholder={endpoint ?? 'https://{service}.{region}.amazonaws.com'}
                  value={endpoint || ''}
                  onChange={onUpdateDatasourceJsonDataOption(props, 'endpoint')}
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
                  onChange={onUpdateDatasourceJsonDataOptionSelect(props, 'defaultRegion')}
                  formatCreateLabel={(r) => `Use region: ${r}`}
                />
              </InlineField>
            </FieldSet>
          )}
          {!hasEdgeAuth && <ConnectionConfig {...props} standardRegions={standardRegions} />}

          <FieldSet label={'Edge settings'} data-testid="edge-connection">
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
                    updateDatasourcePluginJsonDataOption(props, 'edgeAuthMode', v.value);
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
                      onChange={onUserChange}
                      required
                    />
                  </InlineField>
                </InlineFieldRow>
                <InlineFieldRow>
                  <InlineField
                    label="Password"
                    labelWidth={labelWidth}
                    tooltip="The password sent to local authentication proxy"
                  >
                    <Input
                      type="password"
                      name="password"
                      autoComplete="off"
                      placeholder={options.secureJsonFields?.edgeAuthPass ? 'configured' : ''}
                      value={options.secureJsonData?.edgeAuthPass ?? ''}
                      onChange={onPasswordChange}
                      onReset={onResetPassword}
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
                    onClick={onUpdateDatasourceResetOption(props as any, 'cert')}
                  >
                    Reset
                  </Button>
                ) : (
                  <textarea
                    rows={7}
                    className="gf-form-input gf-form-textarea width-30"
                    onChange={(event) => {
                      updateDatasourcePluginSecureJsonDataOption(props as any, 'cert', event.target.value);
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
  }
  // Simple
  if (props.options.jsonData.defaultRegion === 'Edge') {
    return renderEdgeConfig();
  }
  return (
    <ConnectionConfig {...props} standardRegions={standardRegions} newFormStylingEnabled={newFormStylingEnabled} />
  );
}
