import { DataFrame, FieldType, dateTime } from '@grafana/data';
import { QueryType, SiteWiseTimeOrder } from 'types';
import { trimCachedQueryDataFramesAtStart, trimCachedQueryDataFramesEnding } from './dataFrameUtils';

describe('trimCachedQueryDataFrames', () => {
  const absolutionRange = {
    from: dateTime('2024-05-28T00:00:00Z').valueOf(),
    to: dateTime('2024-05-28T00:15:00Z').valueOf(),
  };

  const dataFrame: DataFrame = {
    name: 'Demo Turbine Asset 1',
    refId: 'A',
    fields: [
      {
        name: 'time',
        type: FieldType.time,
        config: {},
        values: [
          1716854400000,  // 2024-05-28T00:00:00Z
          1716854400001,  // 2024-05-28T00:00:00Z + 1ms
          1716855300000,  // 2024-05-28T00:15:00Z
          1716855300001,  // 2024-05-28T00:15:00Z + 1ms
        ],
      },
      {
        name: 'RotationsPerSecond',
        type: FieldType.number,
        config: {
          unit: 'RPS'
        },
        values: [
          1,
          2,
          3,
          4,
        ],
      },
    ],
    length: 4
  };

  const dataFrameDescending: DataFrame = {
    name: 'Demo Turbine Asset 1',
    refId: 'A',
    fields: [
      {
        name: 'time',
        type: FieldType.time,
        config: {},
        values: [
          1716855300001,  // 2024-05-28T00:15:00Z + 1ms
          1716855300000,  // 2024-05-28T00:15:00Z
          1716854400001,  // 2024-05-28T00:00:00Z + 1ms
          1716854400000,  // 2024-05-28T00:00:00Z
        ],
      },
      {
        name: 'RotationsPerSecond',
        type: FieldType.number,
        config: {
          unit: 'RPS'
        },
        values: [4,3,2,1],
      },
    ],
    length: 4
  };

  it('excludes data of PropertyValue query', () => {
    const cachedQueryInfo = {
      query: {
        queryType: QueryType.PropertyValue,
        refId: 'A'
      },
      dataFrame,
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual({
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [],
      length: 0,
    });
  });

  it.each([
    QueryType.ListAssetModels,
    QueryType.ListAssets,
    QueryType.ListAssociatedAssets,
    QueryType.ListAssetProperties,
    QueryType.DescribeAsset,
  ])('does not modify data of non-time-series type - %s', (queryType: QueryType) => {
    const cachedQueryInfo = {
      query: {
        queryType,
        refId: 'A'
      },
      dataFrame,
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(dataFrame);
  });

  it.each([
    QueryType.PropertyAggregate,
    QueryType.PropertyInterpolated,
    QueryType.PropertyValueHistory,
  ])('trims time series data of time-series type - "%s"', (queryType: QueryType) => {
    const cachedQueryInfo = {
      query: {
        queryType,
        refId: 'A'
      },
      dataFrame,
    };
    const expectedDataFrame: DataFrame = {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: FieldType.time,
          config: {},
          values: [
            1716854400001,  // +1ms
            1716855300000,  // 2024-05-28T00:15:00Z
          ],
        },
        {
          name: 'RotationsPerSecond',
          type: FieldType.number,
          config: {
            unit: 'RPS'
          },
          values: [
            2,
            3,
          ],
        },
      ],
      length: 2
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(expectedDataFrame);
  });

  it.each([
    QueryType.PropertyAggregate,
    QueryType.PropertyValueHistory,
  ])('trims descending time series data of time-series type - "%s"', (queryType: QueryType) => {
    const cachedQueryInfo = {
      query: {
        queryType,
        refId: 'A',
        timeOrdering: SiteWiseTimeOrder.DESCENDING,
      },
      dataFrame: dataFrameDescending,
    };
    const expectedDataFrame: DataFrame = {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: FieldType.time,
          config: {},
          values: [
            1716855300000,  // 2024-05-28T00:15:00Z
            1716854400001,  // 2024-05-28T00:00:00Z+1ms
          ],
        },
        {
          name: 'RotationsPerSecond',
          type: FieldType.number,
          config: {
            unit: 'RPS'
          },
          values: [3,2],
        },
      ],
      length: 2
    };
    const dataFrames = trimCachedQueryDataFramesEnding([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(expectedDataFrame);
  });

  it('keeps all data when all time values within range', () => {
    const cachedQueryInfo = {
      query: {
        queryType: QueryType.PropertyValueHistory,
        refId: 'A'
      },
      dataFrame: {
        name: 'Demo Turbine Asset 1',
        refId: 'A',
        fields: [
          {
            name: 'time',
            type: FieldType.time,
            config: {},
            values: [
              1716854400001,  // 2024-05-28T00:00:00Z+1ms
              1716855300000,  // 2024-05-28T00:15:00Z
            ],
          },
          {
            name: 'RotationsPerSecond',
            type: FieldType.number,
            config: {
              unit: 'RPS'
            },
            values: [
              1,
              2,
            ],
          },
        ],
        length: 2
      },
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(cachedQueryInfo.dataFrame);
  });

  it('includes no time series data when all time values are before start time', () => {
    const cachedQueryInfo = {
      query: {
        queryType: QueryType.PropertyValueHistory,
        refId: 'A'
      },
      dataFrame: {
        name: 'Demo Turbine Asset 1',
        refId: 'A',
        fields: [
          {
            name: 'time',
            type: FieldType.time,
            config: {},
            values: [
              1716854399999,
              1716854400000,  // 2024-05-28T00:00:00Z
            ],
          },
          {
            name: 'RotationsPerSecond',
            type: FieldType.number,
            config: {
              unit: 'RPS'
            },
            values: [
              1,
              2,
            ],
          },
        ],
        length: 2
      },
    };
    const expectedDataFrame: DataFrame = {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: FieldType.time,
          config: {},
          values: [],
        },
        {
          name: 'RotationsPerSecond',
          type: FieldType.number,
          config: {
            unit: 'RPS'
          },
          values: [],
        },
      ],
      length: 0
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(expectedDataFrame);
  });

  it('includes no time series data when all time values are after end time', () => {
    const cachedQueryInfo = {
      query: {
        queryType: QueryType.PropertyValueHistory,
        refId: 'A'
      },
      dataFrame: {
        name: 'Demo Turbine Asset 1',
        refId: 'A',
        fields: [
          {
            name: 'time',
            type: FieldType.time,
            config: {},
            values: [
              1716855300001,  // 2024-05-28T00:15:00Z +1ms
              1716855300002,
            ],
          },
          {
            name: 'RotationsPerSecond',
            type: FieldType.number,
            config: {
              unit: 'RPS'
            },
            values: [
              1,
              2,
            ],
          },
        ],
        length: 2
      },
    };
    const expectedDataFrame: DataFrame = {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: FieldType.time,
          config: {},
          values: [],
        },
        {
          name: 'RotationsPerSecond',
          type: FieldType.number,
          config: {
            unit: 'RPS'
          },
          values: [],
        },
      ],
      length: 0
    };
    const dataFrames = trimCachedQueryDataFramesAtStart([cachedQueryInfo], absolutionRange);

    expect(dataFrames).toHaveLength(1);
    expect(dataFrames).toContainEqual(expectedDataFrame);
  });
});
