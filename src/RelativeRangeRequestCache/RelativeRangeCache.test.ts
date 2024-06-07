import { FieldType, LoadingState, dateTime } from '@grafana/data';
import { RelativeRangeCache } from 'RelativeRangeRequestCache/RelativeRangeCache';
import { QueryType, SiteWiseTimeOrder } from 'types';
import { generateSiteWiseRequestCacheId } from './cacheIdUtils';

describe('RelativeRangeCache', () => {
  const requestId = 'mock-request-id';
  const range = {
    from: dateTime('2024-05-28T00:00:00Z'),
    to: dateTime('2024-05-28T01:00:00Z'),
    raw: {
      from: 'now-1h',
      to: 'now'
      ,
    }
  };

  const request = {
    requestId,
    interval: '5s',
    intervalMs: 5000,
    range,
    scopedVars: {},
    targets: [
      {
        refId: 'A',
        queryType: QueryType.PropertyValueHistory,
      },
      {
        refId: 'B',
        queryType: QueryType.DescribeAsset,
      },
    ],
    timezone: 'browser',
    app: 'dashboard',
    startTime: 1716858000000,
  };

  const requestDisabledCache = {
    ...request,
    targets: [
      {
        ...request.targets[0],
        clientCache: false,
      },
    ],
  };

  describe('get()', () => {
    it('returns undefined when any query with client cache disabled', () => {
      const cachedQueryInfo = [
        {
          query: {
            queryType: QueryType.PropertyValueHistory,
            refId: 'A',
          },
          dataFrame: {
            name: 'Demo Turbine Asset 1',
            refId: 'A',
            fields: [
              {
                name: 'time',
                type: FieldType.time,
                config: {},
                values: [],
              },
            ],
            length: 0
          },
        },
      ];
      const cacheData = {
        [generateSiteWiseRequestCacheId(requestDisabledCache)]: {
          queries: cachedQueryInfo,
          range,
        },
      };
      const cache = new RelativeRangeCache(new Map(Object.entries(cacheData)));

      expect(cache.get(requestDisabledCache)).toBeUndefined();
    });

    it('returns undefined when there is no cached response', () => {
      const cache = new RelativeRangeCache();

      expect(cache.get(request)).toBeUndefined();
    });

    it('returns starting cached data and time series query to fetch', () => {
      const cachedQueryInfo = [
        {
          query: {
            queryType: QueryType.PropertyValueHistory,
            refId: 'A',
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
                  1716854400000,  // 2024-05-28T00:00:00Z
                  1716854400001,  // 2024-05-28T00:15:00Z + 1ms
                  1716855300000,  // 2024-05-28T00:15:00Z
                  1716855300001,  // 2024-05-28T00:15:00Z + 1ms
                  1716857100000,  // 2024-05-28T00:45:00Z
                  1716857100001,  // 2024-05-28T00:45:00Z + 1ms
                ],
              },
              {
                name: 'RotationsPerSecond',
                type: FieldType.number,
                config: {
                  unit: 'RPS'
                },
                values: [0,1,2,3,4,5],
              },
            ],
            length: 6
          },
        },
        {
          query: {
            queryType: QueryType.ListAssociatedAssets,
            refId: 'B',
          },
          dataFrame: {
            name: 'child',
            refId: 'B',
            fields: [
              {
                  name: 'name',
                  type: FieldType.string,
                  config: {},
                  values: [
                    'child'
                  ],
              },
            ],
            length: 1
          },
        }
      ];
      const expectedDataFrames = [
        {
          name: 'Demo Turbine Asset 1',
          refId: 'A',
          fields: [
            {
              name: 'time',
              type: FieldType.time,
              config: {},
              values: [
                1716854400001,  // 2024-05-28T00:15:00Z + 1ms
                1716855300000,  // 2024-05-28T00:15:00Z
                1716855300001,  // 2024-05-28T00:15:00Z + 1ms
                1716857100000,  // 2024-05-28T00:45:00Z
              ],
            },
            {
              name: 'RotationsPerSecond',
              type: FieldType.number,
              config: {
                unit: 'RPS'
              },
              values: [1,2,3,4],
            },
          ],
          length: 4
        },
        {
          name: 'child',
          refId: 'B',
          fields: [
            {
                name: 'name',
                type: FieldType.string,
                config: {},
                values: [
                  'child'
                ],
            },
          ],
          length: 1
        },
      ];

      const cacheData = {
        [generateSiteWiseRequestCacheId(request)]: {
          queries: cachedQueryInfo,
          range,
        },
      };
      const cache = new RelativeRangeCache(new Map(Object.entries(cacheData)));

      const cacheResult = cache.get(request);

      expect(cacheResult).toBeDefined()
      expect(cacheResult?.cachedResponse).toEqual({
        start: {
          data: expectedDataFrames,
          key: requestId,
          state: LoadingState.Streaming,
        },
        end: {
          data: [],
          key: requestId,
          state: LoadingState.Streaming,
        },
      });
      expect(cacheResult?.refreshingRequest).toEqual({
        ...request,
        targets: [
          {
            ...request.targets[0],  // only time series query
          },
        ],
        range: {
          from: dateTime(1716857100000),  // '2024-05-28T01:45:00Z'
          to: dateTime('2024-05-28T01:00:00Z'),
          raw: {
            from: 'now-1h',
            to: 'now'
            ,
          }
        }
      })
    });

    it('returns ending cached data and time series query to fetch', () => {
      const requestDescending = {
        requestId,
        interval: '5s',
        intervalMs: 5000,
        range,
        scopedVars: {},
        targets: [
          {
            refId: 'A',
            queryType: QueryType.PropertyValueHistory,
            timeOrdering: SiteWiseTimeOrder.DESCENDING,
          },
          {
            refId: 'B',
            queryType: QueryType.DescribeAsset,
          },
        ],
        timezone: 'browser',
        app: 'dashboard',
        startTime: 1716858000000,
      };
      const cachedQueryInfoDescending = [
        {
          query: {
            queryType: QueryType.PropertyValueHistory,
            refId: 'A',
            timeOrdering: SiteWiseTimeOrder.DESCENDING,
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
                  1716857100001,  // 2024-05-28T00:45:00Z + 1ms
                  1716857100000,  // 2024-05-28T00:45:00Z
                  1716855300001,  // 2024-05-28T00:15:00Z + 1ms
                  1716855300000,  // 2024-05-28T00:15:00Z
                  1716854400001,  // 2024-05-28T00:15:00Z + 1ms
                  1716854400000,  // 2024-05-28T00:00:00Z
                ],
              },
              {
                name: 'RotationsPerSecond',
                type: FieldType.number,
                config: {
                  unit: 'RPS'
                },
                values: [5,4,3,2,1,0],
              },
            ],
            length: 6
          },
        },
        {
          query: {
            queryType: QueryType.ListAssociatedAssets,
            refId: 'B',
          },
          dataFrame: {
            name: 'child',
            refId: 'B',
            fields: [
              {
                  name: 'name',
                  type: FieldType.string,
                  config: {},
                  values: [
                    'child'
                  ],
              },
            ],
            length: 1
          },
        }
      ];
      const expectedStartDataFrames = [
        {
          name: 'Demo Turbine Asset 1',
          refId: 'A',
          fields: [],
          length: 0,
        },
        {
          name: 'child',
          refId: 'B',
          fields: [
            {
                name: 'name',
                type: FieldType.string,
                config: {},
                values: [
                  'child'
                ],
            },
          ],
          length: 1
        }
      ];
      const expectedEndingDataFrames = [
        {
          name: 'Demo Turbine Asset 1',
          refId: 'A',
          fields: [
            {
              name: 'time',
              type: FieldType.time,
              config: {},
              values: [
                1716857100000,  // 2024-05-28T00:45:00Z
                1716855300001,  // 2024-05-28T00:15:00Z + 1ms
                1716855300000,  // 2024-05-28T00:15:00Z
                1716854400001,  // 2024-05-28T00:15:00Z + 1ms
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
        }
      ];
      
      const cacheData = {
        [generateSiteWiseRequestCacheId(requestDescending)]: {
          queries: cachedQueryInfoDescending,
          range,
        },
      };
      const cache = new RelativeRangeCache(new Map(Object.entries(cacheData)));

      const cacheResult = cache.get(requestDescending);

      expect(cacheResult).toBeDefined()
      expect(cacheResult?.cachedResponse).toEqual({
        start: {
          data: expectedStartDataFrames,
          key: requestId,
          state: LoadingState.Streaming,
        },
        end: {
          data: expectedEndingDataFrames,
          key: requestId,
          state: LoadingState.Streaming,
        },
      });
      expect(cacheResult?.refreshingRequest).toEqual({
        ...requestDescending,
        targets: [
          {
            ...requestDescending.targets[0],  // only time series query
          },
        ],
        range: {
          from: dateTime(1716857100000),  // '2024-05-28T01:45:00Z'
          to: dateTime('2024-05-28T01:00:00Z'),
          raw: {
            from: 'now-1h',
            to: 'now'
            ,
          }
        }
      })
    });
  });

  describe('set()', () => {
    const cachedQueryInfo = [
      {
        query: {
          queryType: QueryType.PropertyValueHistory,
          refId: 'A',
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
                1716854400001,  // 2024-05-28T00:15:00Z + 1ms
                1716855300000,  // 2024-05-28T00:15:00Z
                1716855300001,  // 2024-05-28T00:15:00Z + 1ms
                1716857100000,  // 2024-05-28T00:45:00Z
              ],
            },
            {
              name: 'RotationsPerSecond',
              type: FieldType.number,
              config: {
                unit: 'RPS'
              },
              values: [1,2,3,4],
            },
          ],
          length: 4
        },
      },
      {
        query: {
          queryType: QueryType.DescribeAsset,
          refId: 'B',
        },
        dataFrame: {
          name: 'child',
          refId: 'B',
          fields: [
            {
                name: 'name',
                type: FieldType.string,
                config: {},
                values: [
                  'child'
                ],
            },
          ],
          length: 1
        },
      }
    ];  
    const expectedDataFrames = [
      {
        name: 'Demo Turbine Asset 1',
        refId: 'A',
        fields: [
          {
            name: 'time',
            type: FieldType.time,
            config: {},
            values: [
              1716854400001,  // 2024-05-28T00:15:00Z + 1ms
              1716855300000,  // 2024-05-28T00:15:00Z
              1716855300001,  // 2024-05-28T00:15:00Z + 1ms
              1716857100000,  // 2024-05-28T00:45:00Z
            ],
          },
          {
            name: 'RotationsPerSecond',
            type: FieldType.number,
            config: {
              unit: 'RPS'
            },
            values: [1,2,3,4],
          },
        ],
        length: 4
      },
      {
        name: 'child',
        refId: 'B',
        fields: [
          {
              name: 'name',
              type: FieldType.string,
              config: {},
              values: [
                'child'
              ],
          },
        ],
        length: 1
      },
    ];

    it('does nothing when any query with client cache disabled', () => {
      const cacheMap = new Map();
      const cache = new RelativeRangeCache(cacheMap);
  
      cache.set(requestDisabledCache, {
        data: expectedDataFrames,
      });
  
      expect(cacheMap.size).toBe(0);
    });

    it('set request/response pair', () => {
      const cacheData = {
        [generateSiteWiseRequestCacheId(request)]: {
          queries: cachedQueryInfo,
          range,
        },
      };
      const expectedCacheMap = new Map(Object.entries(cacheData));
  
      const cacheMap = new Map();
      const cache = new RelativeRangeCache(cacheMap);
  
      cache.set(request, {
        data: expectedDataFrames,
      });
  
      expect(cacheMap).toEqual(expectedCacheMap);
    });
  });
});
