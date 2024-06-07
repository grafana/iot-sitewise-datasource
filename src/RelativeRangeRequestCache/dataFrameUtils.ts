import { AbsoluteTimeRange, DataFrame } from '@grafana/data';
import { CachedQueryInfo, TIME_SERIES_QUERY_TYPES } from './types';
import { QueryType, SiteWiseTimeOrder } from 'types';
import { trimTimeSeriesDataFrame, trimTimeSeriesDataFrameReversedTime } from 'dataFrameUtils';

/**
 * Trim cached query data frames based on the query type and time ordering for appending to the start of the data frame.
 * 
 * @remarks
 * This function is used to trim the cached data frames based on the query type and time ordering
 * to ensure that the data frames are properly formatted for rendering.
 * For descending ordered data frames, it will return an empty data frame.
 * For property value queries, it will return an empty data frame.
 * For all other queries, it will return the trimmed data frame.
 *
 * @param cachedQueryInfos - Cached query infos to trim
 * @param cacheRange - Cache range to include
 * @returns Trimmed data frames
 */
export function trimCachedQueryDataFramesAtStart(cachedQueryInfos: CachedQueryInfo[], cacheRange: AbsoluteTimeRange): DataFrame[] {
  return cachedQueryInfos
    .map((cachedQueryInfo) => {
      const { query: { queryType, timeOrdering }, dataFrame } = cachedQueryInfo;
      if (timeOrdering === SiteWiseTimeOrder.DESCENDING) {
        // Descending ordering data frame are added at the end of the request to respect the ordering
        // See related function - trimCachedQueryDataFramesEnding()
        return {
          ...dataFrame,
          fields: [],
          length: 0,
        };
      }

      // Always refresh PropertyValue
      if (queryType === QueryType.PropertyValue) {
        return {
          ...dataFrame,
          fields: [],
          length: 0,
        };
      }

      if (TIME_SERIES_QUERY_TYPES.has(queryType)) {
        return trimTimeSeriesDataFrame({
          dataFrame: cachedQueryInfo.dataFrame,
          timeRange: cacheRange,
          lastObservation: cachedQueryInfo.query.lastObservation,
        });
      }

      // No trimming needed
      return dataFrame;
    });
}

/**
 * Trim cached query data frames based on the time ordering for appending to the end of the data frame.
 *
 * @remarks
 * This function is used to trim the cached data frames based on the time ordering
 * to ensure that the data frames are properly formatted for rendering.
 * For descending ordered data frames, it will return the trimmed data frame.
 * For all other queries, it will return an empty data frame.
 *
 * @param cachedQueryInfos - Cached query infos to trim
 * @param cacheRange - Cache range to include
 * @returns Trimmed data frames
 */
export function trimCachedQueryDataFramesEnding(cachedQueryInfos: CachedQueryInfo[], cacheRange: AbsoluteTimeRange): DataFrame[] {
  return cachedQueryInfos
    .filter((cachedQueryInfo) => (cachedQueryInfo.query.timeOrdering === SiteWiseTimeOrder.DESCENDING))
    .map((cachedQueryInfo) => {
      return trimTimeSeriesDataFrameReversedTime({
        dataFrame: cachedQueryInfo.dataFrame,
        lastObservation: cachedQueryInfo.query.lastObservation,
        timeRange: cacheRange,
      });
    });
}
