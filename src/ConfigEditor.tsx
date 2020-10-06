import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SitewiseDataSourceOptions } from './types';

const { FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<SitewiseDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.jsonData.accessKeyId = event.target.value;

    onOptionsChange({
      ...options,
    });
  };

  onSecretAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.jsonData.secretAccessKey = event.target.value;

    onOptionsChange({
      ...options,
    });
  };

  onSessionTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.jsonData.sessionToken = event.target.value;

    onOptionsChange({
      ...options,
    });
  };

  onDefaultRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.jsonData.defaultRegion = event.target.value;

    onOptionsChange({
      ...options,
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData } = options;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Default AWS Region"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onDefaultRegionChange}
            value={jsonData.defaultRegion || ''}
            placeholder={'EX: us-west-2'}
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Access Key ID"
            labelWidth={10}
            inputWidth={30}
            onChange={this.onAccessKeyIdChange}
            value={jsonData.accessKeyId || ''}
            placeholder="Access Key ID"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <FormField
              value={jsonData.secretAccessKey || ''}
              label="Secret Access Key"
              placeholder="Secret Access Key"
              labelWidth={10}
              inputWidth={30}
              onChange={this.onSecretAccessKeyIdChange}
            />
          </div>
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <FormField
              value={jsonData.sessionToken || ''}
              label="Session Token"
              placeholder="Session Token"
              labelWidth={10}
              inputWidth={30}
              onChange={this.onSessionTokenChange}
            />
          </div>
        </div>
      </div>
    );
  }
}
