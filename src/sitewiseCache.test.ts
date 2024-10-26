import { DataSource } from 'SitewiseDataSource';
import { SitewiseCache } from './sitewiseCache';
import { DataFrameView, DataQueryRequest, DataSourceInstanceSettings } from '@grafana/data';
import { of } from 'rxjs';
import { AssetModelSummary, AssetSummary } from './queryResponseTypes';
import { SitewiseOptions, SitewiseQuery } from './types';
import * as mocked from 'dataRequestResponses';

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

jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
  getTemplateSrv: () => ({
    getVariables: () => [],
    replace: (v: string) => v,
  }),
}));

describe('SitewiseCache', () => {
  let ds: DataSource;
  let cache: SitewiseCache;

  beforeEach(() => {
    ds = new DataSource(instanceSettings);
    cache = new SitewiseCache(ds, 'us-west-2');
  });

  describe('getAssetInfo', () => {
    it('should return cached asset info if available', async () => {
      cache['assetsById'].set('1', mocked.mockedAssetInfo);

      const result = await cache.getAssetInfo('1');
      expect(result).toEqual(mocked.mockedAssetInfo);
    });

    it('should fetch and cache asset info if not available', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedAssetInfoResponse));

      const result = await cache.getAssetInfo('1');
      console.log(result);
      expect(result).toEqual(mocked.mockedAssetInfo);
      expect(cache['assetsById'].get('1')).toEqual(result);
    });
  });

  describe('getAssetInfoSync', () => {
    it('should return cached asset info if available', () => {
      cache['assetsById'].set('1', mocked.mockedAssetInfo);

      const result = cache.getAssetInfoSync('1');
      expect(result).toEqual(mocked.mockedAssetInfo);
    });

    it('should fetch and cache asset info if not available', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedAssetInfoResponse));

      const result = cache.getAssetInfoSync('1');
      expect(result).toEqual(mocked.mockedAssetInfo);
      expect(cache['assetsById'].get('1')).toEqual(result);
    });
  });

  describe('listAssetProperties', () => {
    it('should return cached asset properties if available', async () => {
      const assetProperties = new DataFrameView<{ id: string; name: string }>(mocked.mockedAssetProperties);
      cache['assetPropertiesByAssetId'].set('1', assetProperties);

      const result = await cache.listAssetProperties('1');
      expect(result).toEqual(assetProperties);
    });

    it('should fetch and cache asset properties if not available', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedListAssetPropertiesResponse));

      const assetProperties = new DataFrameView<{ id: string; name: string }>(mocked.mockedAssetProperties);
      const result = await cache.listAssetProperties('1');

      expect(result).toEqual(assetProperties);
      expect(cache['assetPropertiesByAssetId'].get('1')).toEqual(assetProperties);
    });
  });

  describe('getModels', () => {
    it('should return cached models if available', async () => {
      const models = new DataFrameView<AssetModelSummary>(mocked.mockedAssetModelSummary);
      cache['models'] = models;

      const result = await cache.getModels();
      expect(result).toEqual(models);
    });

    it('should fetch and cache models if not available', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedAssetModelSummaryResponse));

      const models = new DataFrameView<AssetModelSummary>(mocked.mockedAssetModelSummary);
      const result = await cache.getModels();

      expect(result).toEqual(models);
      expect(cache['models']).toEqual(models);
    });
  });

  describe('getAssetsOfType', () => {
    it('should fetch assets of a specific type', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedAssetSummaryResponse));

      const assets = new DataFrameView<AssetSummary>(mocked.mockedAssetSummary);

      const result = await cache.getAssetsOfType('modelId');
      expect(result).toEqual(assets);
    });
  });

  describe('getAssociatedAssets', () => {
    it('should fetch associated assets', async () => {
      const assets = new DataFrameView<AssetSummary>(mocked.mockedAssetSummary);

      const result = await cache.getAssociatedAssets('assetId');
      expect(result).toEqual(assets);
    });
  });

  describe('getTopLevelAssets', () => {
    it('should return cached top-level assets if available', async () => {
      const assets = new DataFrameView<AssetSummary>(mocked.mockedAssetSummary);
      cache['topLevelAssets'] = assets;

      const result = await cache.getTopLevelAssets();
      expect(result).toEqual(assets);
    });

    it('should fetch and cache top-level assets if not available', async () => {
      jest
        .spyOn(DataSource.prototype, 'query')
        .mockImplementation((request: DataQueryRequest<SitewiseQuery>) => of(mocked.mockedAssetSummaryResponse));

      const assets = new DataFrameView<AssetSummary>(mocked.mockedAssetSummary);
      const result = await cache.getTopLevelAssets();

      expect(result).toEqual(assets);
      expect(cache['topLevelAssets']).toEqual(assets);
    });
  });

  describe('getAssetPickerOptions', () => {
    it('should return asset picker options', async () => {
      const assets = new DataFrameView<AssetSummary>(mocked.mockedAssetSummary);
      cache['topLevelAssets'] = assets;

      const options = await cache.getAssetPickerOptions();
      expect(options).toEqual([
        {
          description: 'arn:aws:iot:us-west-2:123456789012:thing/Asset1',
          label: 'Asset 1',
          value: '1',
        },
      ]);
    });
  });
});
