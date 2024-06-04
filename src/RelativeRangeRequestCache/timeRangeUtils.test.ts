import { dateTime } from '@grafana/data';
import { getRefreshRequestRange, isCacheableTimeRange } from './timeRangeUtils';

describe('isCacheableTimeRange()', () => {
  it('returns true for TimeRange with relative time from before refresh minutes to now', () => {
    expect(isCacheableTimeRange({
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T00:30:00Z'),
      raw: {
        from: 'now-30m',
        to: 'now'
      ,}
    })).toBe(true);
  });

  it('returns false for undefined TimeRange', () => {
    expect(isCacheableTimeRange(undefined)).toBe(false);
  });

  it('returns false for TimeRange with absolute time', () => {
    expect(isCacheableTimeRange({
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T00:15:00Z'),
      raw: {
        from: '2024-05-28T00:00:00Z',
        to: '2024-05-28T00:15:00Z'
      ,}
    })).toBe(false);
  });

  it('returns false for TimeRange with relative time not greater than refresh minutes', () => {
    expect(isCacheableTimeRange({
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T00:15:00Z'),
      raw: {
        from: 'now-15m',
        to: 'now'
      ,}
    })).toBe(false);
  });
});

describe('getRefreshRequestRange()', () => {
  it('returns time range with cache ending time as start', () => {
    const requestRange = {
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T02:00:00Z'),
      raw: {
        from: 'now-2h',
        to: 'now'
      ,}
    };

    const cacheRange = {
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T01:00:00Z'),
      raw: {
        from: 'now-2h',
        to: 'now'
      ,}
    };

    const resultTimeRange = getRefreshRequestRange(requestRange, cacheRange);
    
    expect(resultTimeRange.from.isSame(dateTime('2024-05-28T01:00:00Z'))).toBe(true);
    expect(resultTimeRange.to.isSame(dateTime('2024-05-28T02:00:00Z'))).toBe(true);
  });

  it('returns time range with refresh minutes when time ranges are the same', () => {
    const requestRange = {
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T02:00:00Z'),
      raw: {
        from: 'now-2h',
        to: 'now'
      ,}
    };

    const cacheRange = {
      from: dateTime('2024-05-28T00:00:00Z'),
      to: dateTime('2024-05-28T02:00:00Z'),
      raw: {
        from: 'now-2h',
        to: 'now'
      ,}
    };

    const resultTimeRange = getRefreshRequestRange(requestRange, cacheRange);
    
    expect(resultTimeRange.from.isSame(dateTime('2024-05-28T01:45:00Z'))).toBe(true);
    expect(resultTimeRange.to.isSame(dateTime('2024-05-28T02:00:00Z'))).toBe(true);
  });
});
