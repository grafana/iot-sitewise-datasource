import { DataFrame, dataFrameFromJSON, DataFrameJSON, toDataFrame } from '@grafana/data';

import { frameToMetricFindValues } from './utils';

describe('Test utils', () => {
  it('convert simple values', () => {
    const df = dataFrameFromJSON(listAssetsResponse as DataFrameJSON);
    expect(frameToMetricFindValues(df)).toMatchInlineSnapshot(`Array []`); /
  });
});

// TODO real response!!!
const listAssetsResponse = {
  schema: {
    refId: 'A',
    name: 'js_libraries.csv',
    fields: [
      {
        name: 'Library',
        type: 'string',
        typeInfo: {
          frame: 'string',
          nullable: true,
        },
        config: {
          custom: {
            align: 'auto',
            displayMode: 'auto',
            inspect: false,
          },
          color: {
            mode: 'thresholds',
          },
          mappings: [],
          thresholds: {
            mode: 'absolute',
            steps: [
              {
                value: null,
                color: 'green',
              },
              {
                value: 80,
                color: 'red',
              },
            ],
          },
        },
      },
      {
        name: 'Github Stars',
        type: 'number',
        typeInfo: {
          frame: 'int64',
          nullable: true,
        },
        config: {
          custom: {
            align: 'auto',
            displayMode: 'auto',
            inspect: false,
          },
          color: {
            mode: 'thresholds',
          },
          mappings: [],
          thresholds: {
            mode: 'absolute',
            steps: [
              {
                value: null,
                color: 'green',
              },
              {
                value: 80,
                color: 'red',
              },
            ],
          },
        },
      },
      {
        name: 'Forks',
        type: 'number',
        typeInfo: {
          frame: 'int64',
          nullable: true,
        },
        config: {
          custom: {
            align: 'auto',
            displayMode: 'auto',
            inspect: false,
          },
          color: {
            mode: 'thresholds',
          },
          mappings: [],
          thresholds: {
            mode: 'absolute',
            steps: [
              {
                value: null,
                color: 'green',
              },
              {
                value: 80,
                color: 'red',
              },
            ],
          },
        },
      },
      {
        name: 'Watchers',
        type: 'number',
        typeInfo: {
          frame: 'int64',
          nullable: true,
        },
        config: {
          custom: {
            align: 'auto',
            displayMode: 'auto',
            inspect: false,
          },
          color: {
            mode: 'thresholds',
          },
          mappings: [],
          thresholds: {
            mode: 'absolute',
            steps: [
              {
                value: null,
                color: 'green',
              },
              {
                value: 80,
                color: 'red',
              },
            ],
          },
        },
      },
    ],
  },
  data: {
    values: [
      ['React.js', 'Vue', 'Angular', 'JQuery', 'Meteor', 'Aurelia'],
      [169000, 184000, 73400, 54900, 42400, 11600],
      [34000, 29100, 19300, 20000, 5200, 684],
      [6700, 6300, 3200, 3300, 1700, 442],
    ],
  },
};
