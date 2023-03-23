import { DataSource } from './DataSource';
import { DataSourceInstanceSettings, PluginMeta, ScopedVars, TypedVariableModel } from '@grafana/data';
import { QueryType, SitewiseOptions, SitewiseQuery } from './types';
import { assetIdVariableArray, assetIdVariableConstant } from './__mocks__/variableMocks';

const testInstanceSettings = (
  overrides?: Partial<DataSourceInstanceSettings<SitewiseOptions>>
): DataSourceInstanceSettings<SitewiseOptions> => ({
  id: 1,
  uid: 'sitewise-test',
  type: 'sitewise',
  name: 'sitewise-test',
  meta: {} as PluginMeta,
  readOnly: false,
  jsonData: {} as SitewiseOptions,
  access: 'direct',
  ...overrides,
});

jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
  getTemplateSrv: () => {
    return {
      getVariableName(variableId: string) {
        return variableId;
      },
      getVariables(): TypedVariableModel[] {
        return [assetIdVariableArray, assetIdVariableConstant];
      },
      replace(str: string) {
        return str;
      },
    };
  },
}));

describe('Sitewise Datasource', () => {
  describe('Variable support', () => {
    it('should correctly replace assetIds in the query if variable is a constant', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetIds: ['assetIdConstant'],
        propertyAlias: '',
        region: 'default',
        propertyId: '',
      };
      expect(datasource.applyTemplateVariables(query, {} as ScopedVars)).toEqual({
        ...query,
        assetIds: ['valueConstant'],
      });
    });
    it('should correctly replace assetIds in the query if variable is an array of values', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetIds: ['assetIdArray'],
        propertyAlias: '',
        region: 'default',
        propertyId: '',
      };
      expect(datasource.applyTemplateVariables(query, {} as ScopedVars)).toEqual({
        ...query,
        assetIds: ['array1', 'array2', 'array3'],
      });
    });
    it('should correctly replace assetIds in the query if variable is a mix of string constant and array values', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetIds: ['assetIdConstant', 'assetIdArray'],
        propertyAlias: '',
        region: 'default',
        propertyId: '',
      };
      expect(datasource.applyTemplateVariables(query, {} as ScopedVars)).toEqual({
        ...query,
        assetIds: ['valueConstant', 'array1', 'array2', 'array3'],
      });
    });
  });
});
