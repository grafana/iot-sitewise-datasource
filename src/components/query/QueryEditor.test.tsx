import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { DataQueryRequest, DataSourceInstanceSettings, QueryEditorProps } from '@grafana/data';
import { DataSource } from 'DataSource';
import { QueryType, SitewiseOptions, SitewiseQuery } from 'types';
import { QueryEditor } from './QueryEditor';
import { of } from 'rxjs';
import userEvent from '@testing-library/user-event';

const instanceSettings: DataSourceInstanceSettings<SitewiseOptions> = {
  id: 0,
  uid: 'test',
  name: 'sitewise',
  type: 'datasource',
  access: 'direct',
  url: 'http://localhost',
  database: '',
  basicAuth: '',
  isDefault: false,
  jsonData: {},
  readOnly: false,
  withCredentials: false,
  meta: {} as any,
};

const setup = async (query: Partial<SitewiseQuery>, props = defaultProps) => {
  render(
    <QueryEditor
      {...props}
      query={{
        ...defaultProps.query,
        ...query,
      }}
    />
  );

  await openOptionsCollapse();
};
jest.spyOn(DataSource.prototype, 'query').mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of());
jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
  getTemplateSrv: () => ({
    getVariables: () => [],
    replace: (v: string) => v,
  }),
}));
const defaultProps: QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions> = {
  datasource: new DataSource(instanceSettings),
  query: { refId: 'A', queryType: QueryType.DescribeAsset, region: 'default' },
  onRunQuery: jest.fn(),
  onChange: jest.fn(),
};

describe('QueryEditor', () => {
  it('should display correct fields for query type PropertyAggregate', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyId: 'prop',
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Aggregate')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Resolution')).toBeInTheDocument();
      expect(screen.getByText('Expand Time Range')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyAggregate and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyAlias: 'propAlias',
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Aggregate')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type Interpolated Property', async () => {
    await setup({
      queryType: QueryType.PropertyInterpolated,
      propertyId: 'prop',
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
      expect(screen.getByText('Resolution')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type  Interpolated Property and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyInterpolated,
      propertyAlias: 'propAlias',
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Resolution')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();

      // Interpolated Property queries should not have ANY as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
      expect(screen.queryByText('ANY')).not.toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValueHistory', async () => {
    await setup({
      queryType: QueryType.PropertyValueHistory,
      propertyId: 'prop',
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Expand Time Range')).toBeInTheDocument();
      expect(screen.getByText('Format L4E Anomaly Result')).toBeInTheDocument();
      expect(screen.getByText('Client cache')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValueHistory and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyAlias: 'propAlias',
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValue', async () => {
    await setup({
      queryType: QueryType.PropertyValue,
      propertyId: 'prop',
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Format L4E Anomaly Result')).toBeInTheDocument();
      expect(screen.getByText('Client cache')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValue and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyValue,
      propertyAlias: 'propAlias',
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      // temporary condition - in the old form version, the following fields are not displayed, but they should be in the new one
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type ListAssets', async () => {
    await setup({
      queryType: QueryType.ListAssets,
      propertyId: 'prop',
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Model ID')).toBeInTheDocument();
      expect(screen.getByText('Filter')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type ListAssociatedAssets if assetId is defined', async () => {
    await setup({
      queryType: QueryType.ListAssociatedAssets,
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Show')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type ListAssociatedAssets if property Alias is defined', async () => {
    await setup({
      queryType: QueryType.ListAssociatedAssets,
      propertyAlias: 'prop',
    });
    await waitFor(() => {
      expect(screen.getByText('Show')).toBeInTheDocument();
    });
  });

  it('should clear property when the only asset is deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyId: 'prop',
        assetIds: ['asset'],
      },
      {
        ...defaultProps,
        onChange,
      }
    );

    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: 'select-clear-value' })[1]).toBeInTheDocument();
    });

    await userEvent.click(screen.getAllByRole('button', { name: 'select-clear-value' })[1]);

    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        assetIds: [],
        propertyId: undefined,
      })
    );
  });

  it('should clear property when all assets are deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyId: 'prop',
        assetIds: ['asset1', 'asset2', 'asset3'],
      },
      {
        ...defaultProps,
        onChange,
      }
    );

    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: 'select-clear-value' })[1]).toBeInTheDocument();
    });

    await userEvent.click(screen.getAllByRole('button', { name: 'select-clear-value' })[1]);

    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        assetIds: [],
        propertyId: undefined,
      })
    );
  });

  it('should not clear property when only one of multiple assets is deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyId: 'prop',
        assetIds: ['asset1', 'asset2', 'asset3'],
      },
      {
        ...defaultProps,
        onChange,
      }
    );

    await waitFor(() => {
      expect(screen.getAllByRole('button', { name: 'Remove' })[1]).toBeInTheDocument();
    });

    await userEvent.click(screen.getAllByRole('button', { name: 'Remove' })[1]);

    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        assetIds: ['asset1', 'asset3'],
        propertyId: 'prop',
      })
    );
  });
});

async function openOptionsCollapse() {
  const collapseLabel = screen.queryByTestId('collapse-title');
  if (collapseLabel) {
    return userEvent.click(collapseLabel);
  }
}
