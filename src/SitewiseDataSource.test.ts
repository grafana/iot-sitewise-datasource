import { DataSource } from './SitewiseDataSource';
import { DataSourceInstanceSettings, PluginMeta, ScopedVar, ScopedVars } from '@grafana/data';
import { QueryType, SitewiseOptions, SitewiseQuery, SiteWiseResolution } from './types';

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
    // ref: https://github.com/grafana/grafana/blob/main/public/app/features/variables/utils.ts#L17
    const variableRegex = /\$(\w+)|\[\[(\w+?)(?::(\w+))?\]\]|\${(\w+)(?:\.([^:^\}]+))?(?::([^\}]+))?}/g;
    const globalVars: ScopedVars = {
      assetIdConstant: { text: 'valueConstant', value: 'valueConstant' },
      assetIdArray: { text: ['array1', 'array2', 'array3'], value: ['array1', 'array2', 'array3'] },
    };
    return {
      // Approximate mock of replace function, with 'csv' format
      // ref: https://github.com/grafana/grafana/blob/main/public/app/features/templating/template_srv.mock.ts#L30
      replace(str: string, scopedVars?: ScopedVars, format?: string | Function) {
        return str.replace(variableRegex, (match, var1, var2, fmt2, var3, fieldPath, fmt3) => {
          const variableName = var1 || var2 || var3;

          let varMatch: ScopedVar | undefined;
          if (!!scopedVars) {
            varMatch = scopedVars[variableName];
          }
          varMatch = varMatch ?? globalVars[variableName];
          if (Array.isArray(varMatch?.value)) {
            return varMatch?.value.join(',');
          }
          return varMatch?.value ?? '';
        });
      },
    };
  },
}));

describe('Sitewise Datasource', () => {
  describe('Variable support', () => {
    it('should correctly replace resolution in the query if variable is a constant', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.PropertyAggregate,
        assetId: '',
        assetIds: [],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        resolution: '${resolution}' as SiteWiseResolution,
        region: '',
        propertyId: '',
        rawSQL: '',
      };

      expect(
        datasource.applyTemplateVariables(query, {
          resolution: { text: '15m', value: '15m' },
        })
      ).toEqual({
        ...query,
        resolution: '15m',
      });
    });

    it('should correctly replace assetIds in the query if variable is a constant', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetId: '',
        assetIds: ['${assetIdConstant}'],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: '',
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
        assetId: '',
        assetIds: ['${assetIdArray}'],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: '',
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
        assetIds: ['${assetIdConstant}', '${assetIdArray}'],
        assetId: '',
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: '',
      };

      expect(datasource.applyTemplateVariables(query, {} as ScopedVars)).toEqual({
        ...query,
        assetIds: ['valueConstant', 'array1', 'array2', 'array3'],
      });
    });

    it('should correctly prioritize scopedVars over globalVars', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetId: '',
        assetIds: ['${assetIdConstant}'],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: '',
      };

      expect(
        datasource.applyTemplateVariables(query, {
          assetIdConstant: { text: 'scopedValueConstant', value: 'scopedValueConstant' },
        })
      ).toEqual({
        ...query,
        assetIds: ['scopedValueConstant'],
      });
    });

    it('should correctly prioritize scopedVars over globalVars and handle a mix of array and non array vars', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ListAssociatedAssets,
        assetId: '',
        assetIds: ['${assetIdConstant}', 'noVar', '${assetIdArray}'],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: '',
      };

      expect(
        datasource.applyTemplateVariables(query, {
          assetIdConstant: { text: 'scopedValueConstant', value: 'scopedValueConstant' },
        })
      ).toEqual({
        ...query,
        assetIds: ['scopedValueConstant', 'noVar', 'array1', 'array2', 'array3'],
      });
    });

    it('should replace single-select variable in rawSQL using assetIdConstant', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ExecuteQuery,
        assetId: '',
        assetIds: [],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: "SELECT * FROM table WHERE assetId = '${assetIdConstant}'",
      };

      expect(
        datasource.applyTemplateVariables(query, {
          assetIdConstant: { text: 'singleAsset', value: 'singleAsset' },
        })
      ).toEqual({
        ...query,
        rawSQL: "SELECT * FROM table WHERE assetId = 'singleAsset'",
      });
    });

    it('should replace multi-select variable in rawSQL', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ExecuteQuery,
        assetId: '',
        assetIds: [],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: 'SELECT * FROM table WHERE assetId IN (${assetIdArray})',
      };

      expect(datasource.applyTemplateVariables(query, {} as ScopedVars)).toEqual({
        ...query,
        rawSQL: 'SELECT * FROM table WHERE assetId IN (array1,array2,array3)',
      });
    });

    it('should correctly replace multiple variables in rawSQL', async () => {
      const datasource = new DataSource(testInstanceSettings());
      const query: SitewiseQuery = {
        refId: 'RefA',
        queryType: QueryType.ExecuteQuery,
        assetId: '',
        assetIds: [],
        propertyIds: [],
        propertyAlias: '',
        propertyAliases: [],
        region: '',
        propertyId: '',
        rawSQL: 'SELECT * FROM table WHERE id = ${assetIdConstant} AND region = ${region}',
      };

      expect(
        datasource.applyTemplateVariables(query, {
          assetIdConstant: { text: 'sqlValue', value: 'sqlValue' },
          region: { text: 'us-east-1', value: 'us-east-1' },
        })
      ).toEqual({
        ...query,
        rawSQL: 'SELECT * FROM table WHERE id = sqlValue AND region = us-east-1',
      });
    });
  });
});
