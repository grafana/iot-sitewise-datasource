import React from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ConfigEditor } from './ConfigEditor';
import { render, screen } from '@testing-library/react';
import { SitewiseOptions, SitewiseSecureJsonData } from 'types';
import { config } from '@grafana/runtime';
const datasourceOptions = {
  id: 1,
  uid: 'sitewise',
  orgId: 1,
  name: 'sitewise-name',
  typeLogoUrl: 'http://',
  type: 'type',
  typeName: 'typeName',
  access: 'proxy',
  url: 'https://',
  user: 'user',
  database: 'database',
  basicAuth: true,
  basicAuthUser: 'bAUser',
  isDefault: true,
  jsonData: { defaultRegion: 'us-east-1' },
  secureJsonFields: {},
  readOnly: true,
  withCredentials: false,
};
const defaultProps: DataSourcePluginOptionsEditorProps<SitewiseOptions, SitewiseSecureJsonData> = {
  options: datasourceOptions,
  onOptionsChange: jest.fn(),
};
const originalToggleValue = config.featureToggles.awsDatasourcesNewFormStyling

describe('ConfigEditor', () => {
  beforeEach(() => {
    config.featureToggles.awsDatasourcesNewFormStyling = true;
  })
 afterEach(() => {
     config.featureToggles.awsDatasourcesNewFormStyling = originalToggleValue;
 })
  describe('edge configuration', () => {
    it('should show correct fields if Standard authentication', () => {
      render(
        <ConfigEditor {...defaultProps} options={{ ...datasourceOptions, jsonData: { defaultRegion: 'Edge' } }} />
      );
      expect(screen.getByText('Edge settings')).toBeInTheDocument();
      expect(screen.getByText('Authentication Provider')).toBeInTheDocument();
    });
    it('should show correct fields if linux authentication', () => {
      render(
        <ConfigEditor
          {...defaultProps}
          options={{ ...datasourceOptions, jsonData: { defaultRegion: 'Edge', edgeAuthMode: 'linux' } }}
        />
      );
      expect(screen.getByText('Edge settings')).toBeInTheDocument();
      expect(screen.getByText('Username')).toBeInTheDocument();
      expect(screen.getByText('Password')).toBeInTheDocument();
    });
    it('should display warning if region is Edge but no endpoint is specified', () => {
      render(
        <ConfigEditor
          {...defaultProps}
          options={{ ...datasourceOptions, jsonData: { defaultRegion: 'Edge', endpoint: '' } }}
        />
      );
      expect(screen.getByText('Edge settings')).toBeInTheDocument();
      expect(screen.getByTestId('endpoint-warning')).toBeInTheDocument();
    });
  });
  describe('non-edge configuration', () => {
    it('should show correct fields if region is not edge', () => {
      render(
        <ConfigEditor {...defaultProps} options={{ ...datasourceOptions, jsonData: { defaultRegion: 'us-east-2' } }} />
      );
      expect(screen.queryByText('Edge settings')).not.toBeInTheDocument();
      expect(screen.getByText('Authentication Provider')).toBeInTheDocument();
    });
  });
});
