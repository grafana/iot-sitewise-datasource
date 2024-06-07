import { DataFrame, DataQueryRequest, DataQueryResponse, LoadingState, TimeRange } from '@grafana/data';
import { isTimeRangeCoveringStart } from 'timeRangeUtils';
import { SitewiseQuery } from 'types';
import { RequestCacheId, generateSiteWiseRequestCacheId } from './cacheIdUtils';
import { CachedQueryInfo, TIME_SERIES_QUERY_TYPES } from './types';
import { trimCachedQueryDataFramesAtStart, trimCachedQueryDataFramesEnding } from './dataFrameUtils';
import { getRefreshRequestRange, isCacheableTimeRange } from './timeRangeUtils';

interface DataFrameCacheInfo {
  queries: CachedQueryInfo[],
  range: TimeRange;
}

export interface RelativeRangeCacheInfo {
  cachedResponse: {
    start: DataQueryResponse;
    end: DataQueryResponse;
  };
  refreshingRequest: DataQueryRequest<SitewiseQuery>;
}

/**
 * Cache for relative range queries.
 * It caches the start and end of the range for each query.
 *
 * @internal
 * `RelativeRangeCache` uses an private map `responseDataMap` to track a history of the relative range requests and their data frame responses.
 */
export class RelativeRangeCache {
  constructor(private responseDataMap: Map<RequestCacheId, DataFrameCacheInfo> = new Map<RequestCacheId, DataFrameCacheInfo>()) {}

  /**
   * Set the cache for the given query and response.
   * @param request The query used to get the response
   * @param response The response to set the cache for
   * 
   * @internal
   * This method sets a request/response pair to `responseDataMap`, it performs the followings:
   * 1. validates the `request` is a relative time range and has data 15 minute ago (data within 15 minute are always refreshed)
   * 2. creates a unique cache id for the request
   * 3. creates a `DataFrameCacheInfo` object with the `queries` and `range` from the `response`
   * 4. sets the `DataFrameCacheInfo` object to the `responseDataMap` using the cache id as the key
   */
  set(request: DataQueryRequest<SitewiseQuery>, response: DataQueryResponse) {
    const {
      targets,
      range,
    } = request;

    if (!isCacheableTimeRange(range)) {
      return;
    }

    const requestCacheId = generateSiteWiseRequestCacheId(request);
    
    const queryIdMap = new Map(targets.map(q => [q.refId, q]));

    try {
      const queries = response.data.map((dataFrame: DataFrame) => {
        if (dataFrame.refId == null) {
          console.error('Response data frame without a refId, dataFrame: ', dataFrame);
          throw new Error('Response data frame without a refId!');
        }

        const query = queryIdMap.get(dataFrame.refId);
        if (query == null){
          console.error('Response data frame without a corresponding request target, dataFrame: ', dataFrame);
          throw new Error('Response data frame without a corresponding request target!');
        }

        return {
          query,
          dataFrame: dataFrame,
        };
      });

      this.responseDataMap.set(requestCacheId, {
        queries,
        range,
      });
    } catch (error) {
      // NOOP
    }
  }

  /**
   * Get the cached response for the given request.
   * @param request The request to get the cached response for
   * @returns The cached response if found, undefined otherwise
   * 
   * @internal
   * This method gets the cached response for the given request:
   * 1. validates the `request` is a relative time range and has data 15 minute ago (data within 15 minute are always refreshed)
   * 2. looks up the cached data for the request in `responseDataMap` using the cache id
   * 3. if the cached data is found and covers the request range, it trims the data point till the refreshing request, and then returns the cached data and the refreshing request
   * 4. otherwise, it returns undefined
   */
  get(request: DataQueryRequest<SitewiseQuery>): RelativeRangeCacheInfo | undefined {
    const { range: requestRange } = request;

    if (!isCacheableTimeRange(request.range)) {
      return undefined;
    }

    const cachedDataInfo = this.lookupCachedData(request);
    
    if (cachedDataInfo == null || !isTimeRangeCoveringStart(cachedDataInfo.range, requestRange)) {
      return undefined;
    }

    return RelativeRangeCache.parseCacheInfo(cachedDataInfo, request);
  }

  /**
   * Lookup cached data for the given request.
   * @param request DataQueryRequest<SitewiseQuery> request to lookup cached data for
   * @returns Cached data info if found, undefined otherwise
   */
  private lookupCachedData(request: DataQueryRequest<SitewiseQuery>) {
    const requestCacheId = generateSiteWiseRequestCacheId(request);
    const cachedDataInfo = this.responseDataMap.get(requestCacheId);
    
    return cachedDataInfo;
  }

  private static parseCacheInfo(cachedDataInfo: DataFrameCacheInfo, request: DataQueryRequest<SitewiseQuery>) {
    const { range: requestRange, requestId } = request;

    const refreshingRequestRange = getRefreshRequestRange(requestRange, cachedDataInfo.range);
    const refreshingRequest = RelativeRangeCache.getRefreshingRequest(request, refreshingRequestRange);

    const cacheRange = {
      from: requestRange.from.valueOf(),
      to: refreshingRequestRange.from.valueOf(),
    };
    const cachedDataFrames = trimCachedQueryDataFramesAtStart(cachedDataInfo.queries, cacheRange);
    const cachedDataFramesEnding = trimCachedQueryDataFramesEnding(cachedDataInfo.queries, cacheRange);

    return {
      cachedResponse: {
        start: {
          data: cachedDataFrames,
          key: requestId,
          state: LoadingState.Streaming,
        },
        end: {
          data: cachedDataFramesEnding,
          key: requestId,
          state: LoadingState.Streaming,
        },
      },
      refreshingRequest,
    };
  }

  private static getRefreshingRequest(request: DataQueryRequest<SitewiseQuery>, range: TimeRange) {
    const {
      targets,
    } = request;
    
    return {
      ...request,
      range,
      targets: targets.filter(({ queryType }) => TIME_SERIES_QUERY_TYPES.has(queryType)),
    };
  }
}
