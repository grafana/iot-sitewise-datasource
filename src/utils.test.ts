import { dataFrameFromJSON, DataFrameJSON } from '@grafana/data';

import { frameToMetricFindValues } from './utils';

describe('Test utils', () => {
  it('convert simple values', () => {
    const df = dataFrameFromJSON(listAssetsResponse as DataFrameJSON);
    expect(frameToMetricFindValues(df)).toMatchInlineSnapshot(`
      Array [
        Object {
          "text": "WaterTankSimulatorAsset1",
          "value": "3091bd85-8371-4842-8c82-8ece1bf992bb",
        },
      ]
    `);
    const dfPropertyValue = dataFrameFromJSON(getPropertyValueResponse as DataFrameJSON);
    expect(frameToMetricFindValues(dfPropertyValue)).toMatchInlineSnapshot(`Array []`);
  });
});

const listAssetsResponse = {
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
      ['WaterTankSimulatorAsset1'],
      ['3091bd85-8371-4842-8c82-8ece1bf992bb'],
      ['a54eec1f-5433-4a7d-863b-0a7a252bbdd8'],
      ['arn:aws:iotsitewise:us-east-1:166800769179:asset/3091bd85-8371-4842-8c82-8ece1bf992bb'],
      [1636079884000],
      [1636079884000],
      ['ACTIVE'],
      [null],
      ['[]'],
    ],
  },
};

const getPropertyValueResponse = {
  schema: {
    name: 'WaterTankSimulatorAsset1',
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
        name: 'min_temp',
        type: 'number',
        typeInfo: {
          frame: 'float64',
        },
        config: {},
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
    values: [[1636069024000], [275.83007612612244], ['GOOD']],
  },
};
