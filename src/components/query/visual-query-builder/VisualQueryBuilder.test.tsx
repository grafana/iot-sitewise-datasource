import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { DataQueryRequest, DataSourceInstanceSettings, QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { QueryType, SitewiseOptions, SitewiseQuery } from 'types';
import { VisualQueryBuilder } from './VisualQueryBuilder';
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
    <VisualQueryBuilder
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
const mockDatasource = new DataSource(instanceSettings);
mockDatasource.runQuery = jest.fn().mockReturnValue(of({ data: [] }));

const defaultProps: QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions> = {
  datasource: mockDatasource,
  query: { refId: 'A', queryType: QueryType.DescribeAsset },
  onRunQuery: jest.fn(),
  onChange: jest.fn(),
};

describe('VisualQueryBuilder', () => {
  it('should display correct fields for query type PropertyAggregate', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyIds: ['prop'],
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
      // GOOD as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyAggregate and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyAliases: ['propAlias'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Aggregate')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Time')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
      // GOOD as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type Interpolated Property', async () => {
    await setup({
      queryType: QueryType.PropertyInterpolated,
      propertyIds: ['prop'],
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
      expect(screen.getByText('Resolution')).toBeInTheDocument();
      // GOOD as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type  Interpolated Property and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyInterpolated,
      propertyAliases: ['propAlias'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Quality')).toBeInTheDocument();
      expect(screen.getByText('Resolution')).toBeInTheDocument();
      expect(screen.getByText('Format')).toBeInTheDocument();
      // GOOD as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValueHistory', async () => {
    await setup({
      queryType: QueryType.PropertyValueHistory,
      propertyIds: ['prop'],
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
      // GOOD as the quality default
      expect(screen.getByText('GOOD')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValueHistory and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyAggregate,
      propertyAliases: ['propAlias'],
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
      propertyIds: ['prop'],
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.getByText('Property')).toBeInTheDocument();
      expect(screen.getByText('Format L4E Anomaly Result')).toBeInTheDocument();
      expect(screen.getByText('Client cache')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type PropertyValue and using Property alias', async () => {
    await setup({
      queryType: QueryType.PropertyValue,
      propertyAliases: ['propAlias'],
    });
    await waitFor(() => {
      expect(screen.getByText('Property Alias')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type ListAssets', async () => {
    await setup({
      queryType: QueryType.ListAssets,
      propertyIds: ['prop'],
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Model ID')).toBeInTheDocument();
      expect(screen.getByText('Filter')).toBeInTheDocument();
    });
  });

  it('should display correct fields for query type ListAssociatedAssets', async () => {
    await setup({
      queryType: QueryType.ListAssociatedAssets,
      assetIds: ['asset'],
    });
    await waitFor(() => {
      expect(screen.getByText('Asset Hierarchy')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.queryByText('Property Alias')).toBeNull();
    });
  });

  it('should display correct fields for query type ListAssociatedAssets if property Alias is defined', async () => {
    await setup({
      queryType: QueryType.ListAssociatedAssets,
      propertyAliases: ['prop'],
    });
    await waitFor(() => {
      expect(screen.getByText('Asset Hierarchy')).toBeInTheDocument();
      expect(screen.getByText('Asset')).toBeInTheDocument();
      expect(screen.queryByText('Property Alias')).toBeNull();
    });
  });

  it('should clear property when the only asset is deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyIds: ['prop'],
        assetIds: ['asset'],
      },
      {
        ...defaultProps,
        onChange,
      }
    );

    const clearButton = (await screen.findAllByRole('button', { name: 'Clear value' }))[1];
    expect(clearButton).toBeInTheDocument();
    await userEvent.click(clearButton);

    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        assetIds: [],
        propertyIds: [],
      })
    );
  });

  it('should clear property when all assets are deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyIds: ['prop'],
        assetIds: ['asset1', 'asset2', 'asset3'],
      },
      {
        ...defaultProps,
        onChange,
      }
    );

    const clearButton = (await screen.findAllByRole('button', { name: 'Clear value' }))[1];
    expect(clearButton).toBeInTheDocument();
    await userEvent.click(clearButton);

    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({
        assetIds: [],
        propertyIds: [],
      })
    );
  });

  it('should not clear property when only one of multiple assets is deselected', async () => {
    const onChange = jest.fn();

    await setup(
      {
        queryType: QueryType.PropertyValue,
        propertyIds: ['prop'],
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
        propertyIds: ['prop'],
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
