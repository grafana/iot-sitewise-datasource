import React, { PureComponent } from 'react';
import { onUpdateDatasourceResetOption, SelectableValue, updateDatasourcePluginJsonDataOption, updateDatasourcePluginOption, updateDatasourcePluginSecureJsonDataOption } from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData } from '../types';
import { ConnectionConfig, ConnectionConfigProps } from '@grafana/aws-sdk';
import { Alert, Button, InlineField, InlineFieldRow, Input, Select } from '@grafana/ui';
import { standardRegions } from '../regions';

export type Props = ConnectionConfigProps<SitewiseOptions, SitewiseSecureJsonData>;

const edgeAuthMethods: Array<SelectableValue<string>> = [
  {value: 'default', label: 'Default', description: 'default aws auth'},
  {value: 'linux', label: 'Linux', description: 'default aws auth'},
  {value: 'ldap', label: 'LDAP', description: 'default aws auth'}
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
    const edgeAuthMode = edgeAuthMethods.find(f=>f.value === jsonData.edgeAuthMode) ?? edgeAuthMethods[0];

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
                <InlineField label="Auth Mode">
                  <Select options={edgeAuthMethods} value={edgeAuthMode} onChange={(v) => {
                        updateDatasourcePluginJsonDataOption(this.props as any, 'edgeAuthType', v.value);
                      }}/>
                </InlineField>
              </InlineFieldRow>
              <div className="gf-form-inline">
                <div className="gf-form gf-form--v-stretch">
                  <label className="gf-form-label width-14">Certification</label>
                </div>

                {options.secureJsonFields?.cert ? (
                  <div className="gf-form">
                    <div className="max-width-30 gf-form-inline">
                      <Button
                        variant="secondary"
                        type="button"
                        onClick={onUpdateDatasourceResetOption(this.props as any, 'cert')}
                      >
                        Reset
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="gf-form gf-form--grow">
                    <textarea
                      rows={7}
                      className="gf-form-input gf-form-textarea width-30"
                      onChange={(event) => {
                        updateDatasourcePluginSecureJsonDataOption(this.props as any, 'cert', event.target.value);
                      }}
                      placeholder="Begins with -----BEGIN CERTIFICATE------"
                      required
                    />
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </>
    );
  }
}
