import { TimeRange, dateTime } from '@grafana/data';
import { isRelativeFromNow, minDateTime } from 'timeRangeUtils';
import { DEFAULT_TIME_SERIES_REFRESH_MINUTES } from './constants';

/**
 * Check if the given TimeRange is cacheable. A TimeRange is cacheable if it is relative and has data 15 minutes ago.
 * @param TimeRange to check
 * @returns true if the TimeRange is cacheable, false otherwise
 */
export function isCacheableTimeRange(timeRange?: TimeRange): boolean {
  if (!timeRange) {
    return false;
  }

  const { from, to, raw } = timeRange;

  if (!isRelativeFromNow(raw)) {
    return false;
  }

  const defaultRefreshAgo = dateTime(to).subtract(DEFAULT_TIME_SERIES_REFRESH_MINUTES, 'minutes');
  if (!from.isBefore(defaultRefreshAgo)) {
    return false;
  }

  return true;
}

/**
 * Get the refresh TimeRange for a given TimeRange. The refresh TimeRange is the TimeRange that will be used to refresh the cache.
 *
 * @remarks
 * The refresh TimeRange is the TimeRange that will be used to refresh the cache. The TimeRange is usually 15 minutes ago until the end of the request.
 * Unless the cache data ends earlier than the 15 minutes refresh range, then the TimeRange starts from the end cache data until the end of the request.
 * 
 * @param TimeRange to get the refresh TimeRange for
 * @param TimeRange cacheRange the TimeRange that will be used to refresh the cache
 * @returns TimeRange the refresh TimeRange
 */
export function getRefreshRequestRange(requestRange: TimeRange, cacheRange: TimeRange): TimeRange {
  const defaultRefreshAgo = dateTime(requestRange.to).subtract(DEFAULT_TIME_SERIES_REFRESH_MINUTES, 'minutes');
  const from = minDateTime(cacheRange.to, defaultRefreshAgo);

  return {
    from,
    to: requestRange.to,
    raw: requestRange.raw,
  };
}
