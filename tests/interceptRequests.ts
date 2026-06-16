import { type Page } from '@playwright/test';

/**
 * Intercept query requests and replace responses with mock data.
 *
 * TODO: Reduce duplication of query requests which have different refIds.
 * TODO: Use factory to create responses.
 */
export async function interceptRequests(page: Page) {
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
    } else if (requestBody?.includes('ListAssets')) {
      const responseBody = JSON.stringify({
        results: {
          A: {
            status: 200,
            frames: [
              {
                schema: {
                  refId: 'A',
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
    } else if (requestBody?.includes('ListAssociatedAssets')) {
      const responseBody = JSON.stringify({
        results: {
          A: {
            frames: [
              {
                schema: {
                  refId: 'A',
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
                    ['Demo Turbine Asset'],
                    ['78b941f4-c6a7-40c1-820c-b6b4183e31a6'],
                    ['75ee1c3c-564a-4b72-a8af-f44940b9e815'],
                    ['arn:aws:iotsitewise:us-east-1:526544423884:asset/78b941f4-c6a7-40c1-820c-b6b4183e31a6'],
                    [1606184309000],
                    [1606184309000],
                    ['ACTIVE'],
                    [null],
                    ['[]'],
                  ],
                },
              },
            ],
          },
        },
      });
      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('getModels')) {
      const responseBody = JSON.stringify({
        results: {
          getModels: {
            frames: [
              {
                schema: {
                  refId: 'getModels',
                  meta: {
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
                      name: 'description',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
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
                      name: 'error',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
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
                  ],
                },
                data: {
                  values: [
                    ['Demo Wind Farm Asset Model'],
                    ['This is an asset model used in the IoT SiteWise Demo for representing a wind farm.'],
                    ['cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    ['arn:aws:iotsitewise:us-east-1:135710515793:asset-model/cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    [null],
                    ['ACTIVE'],
                    [1606184309000],
                    [1606184309000],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('ListAssetModels')) {
      const responseBody = JSON.stringify({
        results: {
          A: {
            frames: [
              {
                schema: {
                  refId: 'A',
                  meta: {
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
                      name: 'description',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
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
                      name: 'error',
                      type: 'string',
                      typeInfo: {
                        frame: 'string',
                        nullable: true,
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
                  ],
                },
                data: {
                  values: [
                    ['Demo Wind Farm Asset Model'],
                    ['This is an asset model used in the IoT SiteWise Demo for representing a wind farm.'],
                    ['cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    ['arn:aws:iotsitewise:us-east-1:135710515793:asset-model/cec092ac-b034-4d4b-bbd8-1eca007c5750'],
                    [null],
                    ['ACTIVE'],
                    [1606184309000],
                    [1606184309000],
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
    } else if (requestBody?.includes('PropertyValueHistory')) {
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
                  values: [
                    [1709158500000, 1709158500001, 1709158500002],
                    [15614.641075268504, 14312.2342346346, 16283.893249239],
                    ['GOOD', 'GOOD', 'GOOD'],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('PropertyAggregate')) {
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
                  values: [
                    [1709158500000, 1709158500001, 1709158500002],
                    [15614.641075268504, 14312.2342346346, 16283.893249239],
                    ['GOOD', 'GOOD', 'GOOD'],
                  ],
                },
              },
            ],
          },
        },
      });

      await route.fulfill({ body: responseBody });
    } else if (requestBody?.includes('PropertyInterpolated')) {
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
                  values: [
                    [1709158500000, 1709158500001, 1709158500002],
                    [15614.641075268504, 14312.2342346346, 16283.893249239],
                    ['GOOD', 'GOOD', 'GOOD'],
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
