import { type Page } from '@playwright/test';
import { test, expect } from '@grafana/plugin-e2e';
import { SiteWiseOptions, SitewiseSecureJsonData } from '../src/types';

async function interceptRequests(page: Page) {
  await page.route('', async (route) => {
    const requestBody = route.request().postData();

    if (requestBody?.includes('topLevelAssets')) {
      const responseBody = JSON.stringify({
        results: {
          topLevelAssets: {
            status: 200,
            frames: [
              {
                schema: {
                  refId: 'topLevelAssets',
                  meta: {
                    typeVersion: [0, 0],
                    custom: {},
                  },
                  fields: [
                    {
                      name: 'name',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'id',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'model_id',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'arn',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'creation_date',
                      type: 'time',
                      typeInfo: {
                        frame: 'time.Time',
                      },
                    },
                    {
                      name: 'last_update',
                      type: 'time',
                      typeInfo: {
                        frame: 'time.Time',
                      },
                    },
                    {
                      name: 'state',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'error',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
                      },
                    },
                    {
                      name: 'hierarchies',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                  ],
                },
                data: {
                  values: [
                    ['Demo Wind Farm Asset'],
                    ['6edf67ad-e647-45bd-b609-4974a86729ce'],
                    ['cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    ['arn:aws:iotsitewise:us-east-1:526544423884:asset/6edf67ad-e647-45bd-b609-4974a86729ce'],
                    [1606184309000],
                    [1606184309000],
                    ['ACTIVE'],
                    [null],
                    ['[{"Id":"883165ce-ea4d-4bac-a223-783e79c5b271","Name":"Turbine Asset Model"}]'],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('getAssetInfo')) {
      const responseBody = JSON.stringify({
        results: {
          getAssetInfo: {
            status: 200,
            frames: [
              {
                schema: {
                  refId: 'getAssetInfo',
                  fields: [
                    {
                      name: 'name',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'id',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'arn',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'model_id',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'state',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'error',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
                      },
                    },
                    {
                      name: 'creation_date',
                      type: 'time',
                      typeInfo: {
                        frame: 'time.Time',
                      },
                    },
                    {
                      name: 'last_update',
                      type: 'time',
                      typeInfo: {
                        frame: 'time.Time',
                      },
                    },
                    {
                      name: 'hierarchies',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'properties',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                  ],
                },
                data: {
                  values: [
                    ['Demo Wind Farm Asset'],
                    ['6edf67ad-e647-45bd-b609-4974a86729ce'],
                    [''],
                    ['cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    ['ACTIVE'],
                    [null],
                    [1606184309000],
                    [1606184309000],
                    ['[{"Id":"883165ce-ea4d-4bac-a223-783e79c5b271","Name":"Turbine Asset Model"}]'],
                    [
                      '[{"Alias":null,"DataType":"STRING","DataTypeSpec":null,"Id":"3016a465-b862-47bc-8e24-0cc7b619347e","Name":"Reliability Manager","Notification":{"State":"DISABLED","Topic":"$aws/sitewise/asset-models/cec092ac-b034-4d4b-bbd8-1eca007c5750/assets/6edf67ad-e647-45bd-b609-4974a86729ce/properties/3016a465-b862-47bc-8e24-0cc7b619347e"},"Unit":null},{"Alias":null,"DataType":"INTEGER","DataTypeSpec":null,"Id":"3a0025fa-5a2a-4837-8023-4421eff2bf20","Name":"Code","Notification":{"State":"DISABLED","Topic":"$aws/sitewise/asset-models/cec092ac-b034-4d4b-bbd8-1eca007c5750/assets/6edf67ad-e647-45bd-b609-4974a86729ce/properties/3a0025fa-5a2a-4837-8023-4421eff2bf20"},"Unit":null},{"Alias":null,"DataType":"STRING","DataTypeSpec":null,"Id":"8ab8b7b2-118b-4bd8-93f6-4125e1a7bd8e","Name":"Location","Notification":{"State":"DISABLED","Topic":"$aws/sitewise/asset-models/cec092ac-b034-4d4b-bbd8-1eca007c5750/assets/6edf67ad-e647-45bd-b609-4974a86729ce/properties/8ab8b7b2-118b-4bd8-93f6-4125e1a7bd8e"},"Unit":null},{"Alias":null,"DataType":"DOUBLE","DataTypeSpec":null,"Id":"cd66c574-350a-4031-9c18-bedb8d84fa90","Name":"Total Average Power","Notification":{"State":"DISABLED","Topic":"$aws/sitewise/asset-models/cec092ac-b034-4d4b-bbd8-1eca007c5750/assets/6edf67ad-e647-45bd-b609-4974a86729ce/properties/cd66c574-350a-4031-9c18-bedb8d84fa90"},"Unit":"Watts"},{"Alias":null,"DataType":"DOUBLE","DataTypeSpec":null,"Id":"23d44fd0-3a3c-45ac-a385-0e29ce4b8652","Name":"Total Overdrive State Time","Notification":{"State":"DISABLED","Topic":"$aws/sitewise/asset-models/cec092ac-b034-4d4b-bbd8-1eca007c5750/assets/6edf67ad-e647-45bd-b609-4974a86729ce/properties/23d44fd0-3a3c-45ac-a385-0e29ce4b8652"},"Unit":"seconds"}]',
                    ],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('listAssetProperties')) {
      const responseBody = JSON.stringify({
        results: {
          listAssetProperties: {
            status: 200,
            frames: [
              {
                schema: {
                  refId: 'listAssetProperties',
                  meta: {
                    custom: {},
                  },
                  fields: [
                    {
                      name: 'id',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                    {
                      name: 'name',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                  ],
                },
                data: {
                  values: [
                    [
                      '17913b18-8d82-4a72-910d-b1fa6ef9f44a',
                      'baca4874-edf4-45f8-9d74-fb53ae1d2362',
                      '3a53e29a-f032-40ac-8115-0e86a9d16b69',
                      '599f5c30-c631-4d78-aea6-63d395f562f0',
                      '13cea911-f7ae-4cb8-951d-70c57809627f',
                    ],
                    ['Reliability Manager', 'Code', 'Location', 'Total Average Power', 'Total Overdrive State Time'],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('PropertyValue')) {
      const responseBody = JSON.stringify({
        results: {
          A: {
            frames: [
              {
                schema: {
                  name: 'Demo Wind Farm Asset',
                  refId: 'A',
                  fields: [
                    {
                      name: 'time',
                      type: 'time',
                      typeInfo: {
                        frame: 'time.Time',
                      },
                    },
                    {
                      name: 'Total Average Power',
                      type: 'number',
                      typeInfo: {
                        frame: 'float64',
                      },
                      config: {
                        unit: 'watt',
                      },
                    },
                    {
                      name: 'quality',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                      },
                    },
                  ],
                },
                data: {
                  values: [[1709158500000], [15614.641075268504], ['GOOD']],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else {
      // Pass through undefined requests
      await route.continue();
    }
  });
}

test.describe('Query Editor', () => {
  test.describe('Queries', () => {
    test('Get property value', async ({ page, panelEditPage, readProvisionedDataSource }) => {
      await interceptRequests(page);

      /* Configure data source */

      const ds = await readProvisionedDataSource<SiteWiseOptions, SitewiseSecureJsonData>({
        fileName: 'mock-iot-sitewise.e2e.yaml',
      });
      await panelEditPage.datasource.set(ds.name);
      await panelEditPage.setVisualization('Table');

      /* Select query type */

      const queryTypeSelect = panelEditPage.getByTestIdOrAriaLabel('Query type');
      const propertyAliasInput = panelEditPage.getByTestIdOrAriaLabel('Property alias');
      const assetSelect = panelEditPage.getByTestIdOrAriaLabel('Asset');

      await expect(queryTypeSelect).toBeVisible();
      await expect(propertyAliasInput).not.toBeVisible();
      await expect(assetSelect).not.toBeVisible();

      await queryTypeSelect.click();
      await page.getByText('Get property value', { exact: true }).click();

      await expect(propertyAliasInput).toBeVisible();
      await expect(assetSelect).toBeVisible();

      /* Select asset */

      const propertySelect = panelEditPage.getByTestIdOrAriaLabel('Property');
      const demoAssetOption = page.getByText('Demo Wind Farm Asset', { exact: true });

      await expect(propertySelect).not.toBeVisible();

      await assetSelect.click();
      await demoAssetOption.waitFor();
      await demoAssetOption.click();

      await expect(propertySelect).toBeVisible();

      /* Select asset property */

      const qualitySelect = panelEditPage.getByTestIdOrAriaLabel('Quality');
      const timeSelect = panelEditPage.getByTestIdOrAriaLabel('Time');
      const formatSelect = panelEditPage.getByTestIdOrAriaLabel('Format');
      const demoPropertyOption = page.getByText('Total Average Power', { exact: true });

      await expect(qualitySelect).not.toBeVisible();
      await expect(timeSelect).not.toBeVisible();
      await expect(formatSelect).not.toBeVisible();
      await expect(page.getByText('No data')).toBeVisible();

      await propertySelect.click();
      await demoPropertyOption.waitFor();
      await demoPropertyOption.click();

      await expect(qualitySelect).toBeVisible();
      await expect(timeSelect).toBeVisible();
      await expect(formatSelect).toBeVisible();
      await expect(page.getByText('No data')).not.toBeVisible();
      await expect(panelEditPage.panel.data).toContainText(['15.6 kW', 'GOOD']);
    });
  });
});
