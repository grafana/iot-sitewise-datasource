import { CoreApp, DataQueryRequest, DataSourceInstanceSettings, dateTime } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { QueryType, SitewiseOptions, SitewiseQuery } from 'types';
import { variableFormatter, SitewiseVariableSupport } from 'variables';
import { of } from 'rxjs';

const request: DataQueryRequest<SitewiseQuery> = {
  targets: [],
  range: { from: dateTime(), to: dateTime(), raw: { from: dateTime(), to: dateTime() } },
  interval: '1s',
  intervalMs: 1000,
  scopedVars: {},
  timezone: 'UTC',
  requestId: '1',
  app: CoreApp.Dashboard,
  startTime: 1234567890,
};

const mockedDatasourceQuery = jest.fn(() => of({ data: [] }));

describe('SiteWiseVariableSupport', () => {
  describe('query filtering', () => {
    afterEach(() => {
      jest.clearAllMocks();
    });
    const mockDatasource = new DataSource({} as DataSourceInstanceSettings<SitewiseOptions>);
    mockDatasource.query = mockedDatasourceQuery;
    const variableSupport = new SitewiseVariableSupport(mockDatasource);
    test.each([
      { refId: 'A', queryType: QueryType.PropertyInterpolated },
      { refId: 'A', queryType: QueryType.PropertyInterpolated, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.PropertyAggregate },
      { refId: 'A', queryType: QueryType.PropertyAggregate, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.PropertyValueHistory },
      { refId: 'A', queryType: QueryType.PropertyValueHistory, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.PropertyValue },
      { refId: 'A', queryType: QueryType.PropertyValue, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.ListAssociatedAssets },
      { refId: 'A', queryType: QueryType.ListAssets },
    ])('Filters out queries that are missing any required fields', (query: SitewiseQuery) => {
      variableSupport.query({ ...request, targets: [query] });
      expect(mockedDatasourceQuery).not.toHaveBeenCalled();
    });
    test.each([
      { refId: 'A', queryType: QueryType.PropertyInterpolated, assetIds: ['assetId'], propertyIds: ['propertyId'] },
      { refId: 'A', queryType: QueryType.PropertyAggregate, assetIds: ['assetId'], propertyIds: ['propertyId'] },
      { refId: 'A', queryType: QueryType.PropertyValueHistory, assetIds: ['assetId'], propertyIds: ['propertyId'] },
      { refId: 'A', queryType: QueryType.PropertyValue, assetIds: ['assetId'], propertyIds: ['propertyId'] },
      { refId: 'A', queryType: QueryType.ListAssetModels },
      { refId: 'A', queryType: QueryType.ListAssociatedAssets, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.ListAssets, modelId: 'modelId', filter: 'ALL' },
      { refId: 'A', queryType: QueryType.ListTimeSeries },
      { refId: 'A', queryType: QueryType.ListTimeSeries },
    ])('Does not filter out queries that have all the required data', (query: SitewiseQuery) => {
      variableSupport.query({ ...request, targets: [query] });
      expect(mockedDatasourceQuery).toHaveBeenCalled();
    });
  });
});

describe('variableFormatter', () => {
  it('formats a single value correctly', () => {
    expect(variableFormatter('abc')).toBe("'abc'");
  });

  it('formats a number value correctly', () => {
    expect(variableFormatter(123)).toBe("'123'");
  });

  it('formats an array of strings correctly', () => {
    expect(variableFormatter(['a', 'b', 'c'])).toBe("('a', 'b', 'c')");
  });

  it('formats an array of numbers correctly', () => {
    expect(variableFormatter([1, 2, 3])).toBe("('1', '2', '3')");
  });
});
