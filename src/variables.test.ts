import { CoreApp, DataQueryRequest, dateTime } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { testInstanceSettings } from 'SitewiseDataSource.test';
import { QueryType, SitewiseQuery } from 'types';
import { SitewiseVariableSupport } from 'variables';
import { of } from 'rxjs';

const mockedDatasourceQuery = jest.fn().mockReturnValue(of({ data: [] }));

describe('template variable support', () => {
  describe('query filtering', () => {
    const mockDatasource = new DataSource(testInstanceSettings());
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
      { refId: 'A', queryType: QueryType.ListAssetModels },
      { refId: 'A', queryType: QueryType.ListAssociatedAssets },
      { refId: 'A', queryType: QueryType.ListAssets },
    ])('Filters out queries that are missing any required fields', (query: SitewiseQuery) => {
      const request: DataQueryRequest<SitewiseQuery> = {
        targets: [query],
        range: { from: dateTime(), to: dateTime(), raw: { from: dateTime(), to: dateTime() } },
        interval: '1s',
        intervalMs: 1000,
        scopedVars: {},
        timezone: 'UTC',
        requestId: '1',
        app: CoreApp.Dashboard,
        startTime: 1234567890,
      };
      variableSupport.query(request);
      expect(mockedDatasourceQuery).not.toHaveBeenCalled();
    });
    test.each([
      { refId: 'A', queryType: QueryType.PropertyInterpolated, assetIds: ['assetId'], propertyId: 'propertyId' },
      { refId: 'A', queryType: QueryType.PropertyAggregate, assetIds: ['assetId'], propertyId: 'propertyId' },
      { refId: 'A', queryType: QueryType.PropertyValueHistory, assetIds: ['assetId'], propertyId: 'propertyId' },
      { refId: 'A', queryType: QueryType.PropertyValue, assetIds: ['assetId'], propertyId: 'propertyId' },
      { refId: 'A', queryType: QueryType.ListAssetModels, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.ListAssociatedAssets, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.ListAssets, assetIds: ['assetId'] },
      { refId: 'A', queryType: QueryType.ListTimeSeries },
      { refId: 'A', queryType: QueryType.ListTimeSeries },
    ])('Does not filter out queries that have all the required data', (query: SitewiseQuery) => {
      jest.clearAllMocks();
      const request: DataQueryRequest<SitewiseQuery> = {
        targets: [query],
        range: { from: dateTime(), to: dateTime(), raw: { from: dateTime(), to: dateTime() } },
        interval: '1s',
        intervalMs: 1000,
        scopedVars: {},
        timezone: 'UTC',
        requestId: '1',
        app: CoreApp.Dashboard,
        startTime: 1234567890,
      };
      variableSupport.query(request);
      expect(mockedDatasourceQuery).toHaveBeenCalled();
    });
  });
});
