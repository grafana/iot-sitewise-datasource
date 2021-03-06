import React, { PureComponent } from 'react';
import { onUpdateDatasourceResetOption, updateDatasourcePluginSecureJsonDataOption } from '@grafana/data';
import { SitewiseOptions, SitewiseSecureJsonData } from '../types';
import { ConnectionConfig, ConnectionConfigProps } from '@grafana/aws-sdk';
import { Alert, Button } from '@grafana/ui';
import { standardRegions } from '../regions';

export type Props = ConnectionConfigProps<SitewiseOptions, SitewiseSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
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
          <ConnectionConfig {...this.props} standardRegions={standardRegions} />

          {defaultRegion === 'Edge' && (
            <div>
              {!endpoint && (
                <Alert title="Edge region requires an explicit endpoint configured above" severity="warning" />
              )}

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
