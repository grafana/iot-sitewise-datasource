import { DateTime, RawTimeRange, TimeRange, dateTime } from '@grafana/data';

/**
 * Checks if the subject TimeRange covers the start time of the object TimeRange.
 *
 * @param subjectRange - The TimeRange to check if it covers the start time of the object.
 * @param objectRange - The TimeRange to check if its start time is covered by the subject.
 * @returns True if the subject TimeRange covers the start time of the object, false otherwise.
 */
export function isTimeRangeCoveringStart(subjectRange: TimeRange, objectRange: TimeRange): boolean {
  const { from: subjectFrom, to: subjectTo } = subjectRange;
  const { from: objectFrom } = objectRange;

  /*
   * True if both time ranges start at the same time.
   *
   * Positive example (same from time):
   *   subject: <from>...
   *   object:  <from>...
   */
  if (objectFrom.isSame(subjectFrom)) {
    return true;
  }

  /*
   * True if subject starts before object starts and overlaps the object start time
   *
   * Positive example (subject from and to wrap around object from):
   *   subject: <from>......<to>
   *   object:  ......<from>....(disregard to)
   *
   * Negative example (subject from and to both before object):
   *   subject: <from>.<to>.......
   *   object:  ...........<from>.(disregard to)
   */
  if (subjectFrom.isBefore(objectFrom) && objectFrom.isBefore(subjectTo)) {
    return true;
  }

  return false;
}

export function minDateTime(firstDateTimes: DateTime, ...dateTimes: DateTime[]) {
  const minValue = Math.min(firstDateTimes.valueOf(), ...dateTimes.map((dateTime) => dateTime.valueOf()));
  return dateTime(minValue);
}

export function isRelativeFromNow(timeRange: RawTimeRange): boolean {
  const { from, to } = timeRange;

  if (typeof from !== 'string' || typeof to !== 'string') {
    return false;
  }

  return from.startsWith('now-') && to === 'now';
}
