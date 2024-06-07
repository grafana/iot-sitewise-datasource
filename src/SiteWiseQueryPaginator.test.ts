import { DataQueryRequest, DataQueryResponse, LoadingState, dateTime } from '@grafana/data';

import { SitewiseQueryPaginator } from 'SiteWiseQueryPaginator';
import { QueryType, SitewiseNextQuery, SitewiseQuery } from 'types';
import { first, last } from 'rxjs/operators';

// Request for SiteWise data
const dataQueryRequest: DataQueryRequest<SitewiseQuery> = {
  app: 'panel-editor',
  requestId: 'Q112',
  timezone: 'browser',
  panelId: 2,
  dashboardUID: 'OPixSZySk',
  range: {
    from: dateTime('2024-05-28T20:59:49.659Z'),
    to: dateTime('2024-05-28T21:29:49.659Z'),
    raw: {
      from: 'now-30m',
      to: 'now'
    }
  },
  interval: '2s',
  intervalMs: 2000,
  targets: [
    {
      datasource: {
        type: 'grafana-iot-sitewise-datasource',
        uid: 's0PWceLIz'
      },
      assetIds: [
        '0af4b18d-8c44-4944-8f59-9001b8824362'
      ],
      flattenL4e: true,
      maxPageAggregations: 1,
      propertyId: '3e44a93c-eb71-4dfd-8aec-bb825cfdf7dd',
      queryType: QueryType.PropertyValue,
      refId: 'A'
    }
  ],
  maxDataPoints: 711,
  scopedVars: {
    __interval: {
      text: '2s',
      value: '2s'
    },
    __interval_ms: {
      text: '2000',
      value: 2000
    }
  },
  startTime: 1716931789659,
  rangeRaw: {
    from: 'now-30m',
    to: 'now'
  }
};

// Response with SiteWise data
const dataQueryResponse: DataQueryResponse = {
  data: [
    {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: 'time',
          typeInfo: {
            frame: 'time.Time'
          },
          config: {},
          values: [
            1716931550000
          ],
          entities: {}
        },
        {
          name: 'RotationsPerSecond',
          type: 'number',
          typeInfo: {
            frame: 'float64'
          },
          config: {
            unit: 'RPS'
          },
          values: [
            0.45253960150485795
          ],
          entities: {}
        },
        {
          name: 'quality',
          type: 'string',
          typeInfo: {
            frame: 'string'
          },
          config: {},
          values: [
            'GOOD'
          ],
          entities: {}
        }
      ],
      length: 1
    }
  ],
  state: LoadingState.Done,
};

// Request with SiteWise next token
const dataQueryRequestPaginating: DataQueryRequest<SitewiseNextQuery> = {
  app: 'panel-editor',
  requestId: 'Q112.2',
  timezone: 'browser',
  panelId: 2,
  dashboardUID: 'OPixSZySk',
  range: {
    from: dateTime('2024-05-28T20:59:49.659Z'),
    to: dateTime('2024-05-28T21:29:49.659Z'),
    raw: {
      from: 'now-30m',
      to: 'now'
    }
  },
  interval: '2s',
  intervalMs: 2000,
  targets: [
    {
      datasource: {
        type: 'grafana-iot-sitewise-datasource',
        uid: 's0PWceLIz'
      },
      assetIds: [
        '0af4b18d-8c44-4944-8f59-9001b8824362'
      ],
      flattenL4e: true,
      maxPageAggregations: 1,
      propertyId: '3e44a93c-eb71-4dfd-8aec-bb825cfdf7dd',
      queryType: QueryType.PropertyValue,
      refId: 'A',
      nextToken: 'mock-next-token-value',
      nextTokens: {},
    }
  ],
  maxDataPoints: 711,
  scopedVars: {
    __interval: {
      text: '2s',
      value: '2s'
    },
    __interval_ms: {
      text: '2000',
      value: 2000
    }
  },
  startTime: 1716931789659,
  rangeRaw: {
    from: 'now-30m',
    to: 'now'
  }
};

// Response with SiteWise next token
const dataQueryResponsePaginating: DataQueryResponse = {
  data: [
    {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: 'time',
          typeInfo: {
            frame: 'time.Time'
          },
          config: {},
          values: [
            1716931549000
          ],
          entities: {}
        },
        {
          name: 'RotationsPerSecond',
          type: 'number',
          typeInfo: {
            frame: 'float64'
          },
          config: {
            unit: 'RPS'
          },
          values: [
            1
          ],
          entities: {}
        },
        {
          name: 'quality',
          type: 'string',
          typeInfo: {
            frame: 'string'
          },
          config: {},
          values: [
            'GOOD'
          ],
          entities: {}
        }
      ],
      length: 1,
      meta: {
        custom: {
          nextToken: 'mock-next-token-value',
          resolution: 'RAW'
        }
      },
    }
  ],
  state: LoadingState.Done,
};

// Response with data combined from `dataQueryResponse` and `dataQueryResponsePaginating`
const dataQueryResponseCombined: DataQueryResponse = {
  data: [
    {
      name: 'Demo Turbine Asset 1',
      refId: 'A',
      fields: [
        {
          name: 'time',
          type: 'time',
          typeInfo: {
            frame: 'time.Time'
          },
          config: {},
          values: [
            1716931549000,
            1716931550000
          ],
          entities: {}
        },
        {
          name: 'RotationsPerSecond',
          type: 'number',
          typeInfo: {
            frame: 'float64'
          },
          config: {
            unit: 'RPS'
          },
          values: [
            1,
            0.45253960150485795
          ],
          entities: {}
        },
        {
          name: 'quality',
          type: 'string',
          typeInfo: {
            frame: 'string'
          },
          config: {},
          values: [
            'GOOD',
            'GOOD'
          ],
          entities: {}
        }
      ],
      length: 2
    }
  ],
  state: LoadingState.Done,
};

describe('SitewiseQueryPaginator', () => {
  describe('toObservable()', () => {
    it('handles single page request', async () => {
      const request = dataQueryRequest;
      const queryFn = jest.fn().mockResolvedValue(dataQueryResponse);

      const queryObservable = new SitewiseQueryPaginator({
        request,
        queryFn,
      }).toObservable();

      const firstResponse = queryObservable.pipe(first()).toPromise();
      expect(firstResponse).resolves.toMatchObject(dataQueryResponse);

      const lastResponse = queryObservable.pipe(last()).toPromise();
      expect(lastResponse).resolves.toMatchObject(dataQueryResponse);

      await lastResponse;
      expect(queryFn).toHaveBeenCalledTimes(1);
      expect(queryFn).toHaveBeenCalledWith(request);
    });

    it('handles single page request with cached data', async () => {
      const cachedStartResponse: DataQueryResponse = {
        data: [
          {
            name: 'Demo Turbine Asset 1',
            refId: 'A',
            fields: [
              {
                name: 'time',
                type: 'time',
                typeInfo: {
                  frame: 'time.Time'
                },
                config: {},
                values: [
                  1716931540000
                ],
                entities: {}
              },
              {
                name: 'RotationsPerSecond',
                type: 'number',
                typeInfo: {
                  frame: 'float64'
                },
                config: {
                  unit: 'RPS'
                },
                values: [
                  1
                ],
                entities: {}
              },
              {
                name: 'quality',
                type: 'string',
                typeInfo: {
                  frame: 'string'
                },
                config: {},
                values: [
                  'GOOD'
                ],
                entities: {}
              }
            ],
            length: 1
          }
        ],
        state: LoadingState.Done,
      };

      const cachedEndResponse: DataQueryResponse = {
        data: [
          {
            name: 'Demo Turbine Asset 1',
            refId: 'A',
            fields: [
              {
                name: 'time',
                type: 'time',
                typeInfo: {
                  frame: 'time.Time'
                },
                config: {},
                values: [
                  1716931560000
                ],
                entities: {}
              },
              {
                name: 'RotationsPerSecond',
                type: 'number',
                typeInfo: {
                  frame: 'float64'
                },
                config: {
                  unit: 'RPS'
                },
                values: [
                  3
                ],
                entities: {}
              },
              {
                name: 'quality',
                type: 'string',
                typeInfo: {
                  frame: 'string'
                },
                config: {},
                values: [
                  'GOOD'
                ],
                entities: {}
              }
            ],
            length: 1
          }
        ],
        state: LoadingState.Done,
      };

      const expectedResponse: DataQueryResponse = {
        data: [
          {
            name: 'Demo Turbine Asset 1',
            refId: 'A',
            fields: [
              {
                name: 'time',
                type: 'time',
                typeInfo: {
                  frame: 'time.Time'
                },
                config: {},
                values: [
                  1716931540000,
                  1716931550000,
                  1716931560000
                ],
                entities: {}
              },
              {
                name: 'RotationsPerSecond',
                type: 'number',
                typeInfo: {
                  frame: 'float64'
                },
                config: {
                  unit: 'RPS'
                },
                values: [
                  1,
                  0.45253960150485795,
                  3,
                ],
                entities: {}
              },
              {
                name: 'quality',
                type: 'string',
                typeInfo: {
                  frame: 'string'
                },
                config: {},
                values: [
                  'GOOD',
                  'GOOD',
                  'GOOD',
                ],
                entities: {}
              }
            ],
            length: 3
          }
        ],
        state: LoadingState.Done,
      };

      const request = dataQueryRequest;
      const queryFn = jest.fn().mockResolvedValue(dataQueryResponse);

      const queryObservable = new SitewiseQueryPaginator({
        cachedResponse: {
          start: cachedStartResponse,
          end: cachedEndResponse,
        },
        request,
        queryFn,
      }).toObservable();

      const lastResponse = queryObservable.pipe(last()).toPromise();
      expect(lastResponse).resolves.toMatchObject(expectedResponse);

      await lastResponse;
      expect(queryFn).toHaveBeenCalledTimes(1);
      expect(queryFn).toHaveBeenCalledWith(request);
    });

    it('handles more than 1 page request', async () => {
      const request = dataQueryRequest;
      const queryFn = jest.fn()
        .mockResolvedValueOnce(dataQueryResponsePaginating)
        .mockResolvedValueOnce(dataQueryResponse);

      const queryObservable = new SitewiseQueryPaginator({
        request,
        queryFn,
      }).toObservable();

      const firstResponse = queryObservable.pipe(first()).toPromise();
      expect(firstResponse).resolves.toMatchObject({
        ...dataQueryResponsePaginating,
        state: LoadingState.Streaming,
      });

      const lastResponse = queryObservable.pipe(last()).toPromise();
      expect(lastResponse).resolves.toMatchObject(dataQueryResponseCombined);

      await lastResponse;
      expect(queryFn).toHaveBeenCalledTimes(2);
      expect(queryFn).toHaveBeenCalledWith(request);
      expect(queryFn).toHaveBeenCalledWith(dataQueryRequestPaginating);
    });

    it('handles error state response and terminate pagination', async () => {
      const request = dataQueryRequest;
      const queryFn = jest.fn()
        .mockResolvedValueOnce({
          ...dataQueryResponsePaginating,
          state: LoadingState.Error,
        })
        .mockResolvedValueOnce(dataQueryResponse);

      const queryObservable = new SitewiseQueryPaginator({
        request,
        queryFn,
      }).toObservable();

      const firstResponse = queryObservable.pipe(first()).toPromise();
      expect(firstResponse).resolves.toMatchObject({
        ...dataQueryResponsePaginating,
        state: LoadingState.Error,
      });

      const lastResponse = queryObservable.pipe(last()).toPromise();
      expect(lastResponse).resolves.toMatchObject({
        ...dataQueryResponsePaginating,
        state: LoadingState.Error,
      });

      await lastResponse;
      expect(queryFn).toHaveBeenCalledTimes(1);
      expect(queryFn).toHaveBeenCalledWith(request);
    });
  });
});
