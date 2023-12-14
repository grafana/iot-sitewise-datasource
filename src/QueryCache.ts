import {
  closestIdx,
  DataFrame,
  DataQueryRequest,
  dateTime,
  durationToMilliseconds,
  Field,
  incrRoundDn,
  isValidDuration,
  parseDuration,
} from '@grafana/data';
import { SitewiseQuery } from './types';

export type Table = [times: number[], ...values: any[][]];

// prevTable and nextTable are assumed sorted ASC on reference [0] arrays
// nextTable is assumed to be contiguous, only edges are checked for overlap
// ...so prev: [1,2,5] + next: [3,4,6] -> [1,2,3,4,6]
export function amendTable(prevTable: Table, nextTable: Table): Table {
  let [prevTimes] = prevTable;
  let [nextTimes] = nextTable;

  let pLen = prevTimes.length;
  let pStart = prevTimes[0];
  let pEnd = prevTimes[pLen - 1];

  let nLen = nextTimes.length;
  let nStart = nextTimes[0];
  let nEnd = nextTimes[nLen - 1];

  let outTable: Table;

  if (pLen) {
    if (nLen) {
      // append, no overlap
      if (nStart > pEnd) {
        outTable = prevTable.map((_, i) => prevTable[i].concat(nextTable[i])) as Table;
      }
      // prepend, no overlap
      else if (nEnd < pStart) {
        outTable = nextTable.map((_, i) => nextTable[i].concat(prevTable[i])) as Table;
      }
      // full replace
      else if (nStart <= pStart && nEnd >= pEnd) {
        outTable = nextTable;
      }
      // partial replace
      else if (nStart > pStart && nEnd < pEnd) {
      }
      // append, with overlap
      else if (nStart >= pStart) {
        let idx = closestIdx(nStart, prevTimes);
        idx = prevTimes[idx] < nStart ? idx - 1 : idx;
        outTable = prevTable.map((_, i) => prevTable[i].slice(0, idx).concat(nextTable[i])) as Table;
      }
      // prepend, with overlap
      else if (nEnd >= pStart) {
        let idx = closestIdx(nEnd, prevTimes);
        idx = prevTimes[idx] < nEnd ? idx : idx + 1;
        outTable = nextTable.map((_, i) => nextTable[i].concat(prevTable[i].slice(idx))) as Table;
      }
    } else {
      outTable = prevTable;
    }
  } else {
    if (nLen) {
      outTable = nextTable;
    } else {
      outTable = [[]];
    }
  }

  return outTable!;
}

export function trimTable(table: Table, fromTime: number, toTime: number): Table {
  let [times, ...vals] = table;
  let fromIdx: number | undefined;
  let toIdx: number | undefined;

  // trim to bounds
  if (times[0] < fromTime) {
    fromIdx = closestIdx(fromTime, times);

    if (times[fromIdx] < fromTime) {
      fromIdx++;
    }
  }

  if (times[times.length - 1] > toTime) {
    toIdx = closestIdx(toTime, times);

    if (times[toIdx] > toTime) {
      toIdx--;
    }
  }

  if (fromIdx != null || toIdx != null) {
    times = times.slice(fromIdx ?? 0, toIdx);
    vals = vals.map((vals2) => vals2.slice(fromIdx ?? 0, toIdx));
  }

  return [times, ...vals];
}

// dashboardUID + panelId + refId
// (must be stable across query changes, time range changes / interval changes / panel resizes / template variable changes)
type TargetIdent = string;

// query + template variables + interval + raw time range
// used for full target cache busting -> full range re-query
type TargetSig = string;

type TimestampMs = number;

type SupportedQueryTypes = SitewiseQuery;

// TODO: UPDATE THIS
export const defaultQueryOverlapWindow = '10m';

interface TargetCache {
  sig: TargetSig;
  prevFrom: TimestampMs;
  prevTo: TimestampMs;
  frames: DataFrame[];
}

export interface CacheRequestInfo<T extends SupportedQueryTypes> {
  requests: Array<DataQueryRequest<T>>;
  targetSigs: Map<TargetIdent, TargetSig>;
  shouldCache: boolean;
}

/**
 * Get field identity
 * This is the string used to uniquely identify a field within a "target"
 * @param field
 */
export const getFieldIdent = (field: Field) => `${field.type}|${field.name}|${JSON.stringify(field.labels ?? '')}`;

/**
 * NOMENCLATURE
 * Target: The request target (DataQueryRequest), i.e. a specific query reference within a panel
 * Ident: Identity: the string that is not expected to change
 * Sig: Signature: the string that is expected to change, upon which we wipe the cache fields
 */
export class QueryCache<T extends SupportedQueryTypes> {
  private overlapWindowMs: number;
  private getTargetSignature: (request: DataQueryRequest<T>, target: T) => string;

  cache = new Map<TargetIdent, TargetCache>();

  constructor(options: {
    getTargetSignature: (request: DataQueryRequest<T>, target: T) => string;
    overlapString: string;
  }) {
    const unverifiedOverlap = options.overlapString;
    if (isValidDuration(unverifiedOverlap)) {
      const duration = parseDuration(unverifiedOverlap);
      this.overlapWindowMs = durationToMilliseconds(duration);
    } else {
      const duration = parseDuration(defaultQueryOverlapWindow);
      this.overlapWindowMs = durationToMilliseconds(duration);
    }

    this.getTargetSignature = options.getTargetSignature;
  }

  // can be used to change full range request to partial, split into multiple requests
  requestInfo(request: DataQueryRequest<T>): CacheRequestInfo<T> {
    // TODO: align from/to to interval to increase probability of hitting backend cache

    const newFrom = request.range.from.valueOf();
    const newTo = request.range.to.valueOf();

    //TODO: REVISIT THIS
    // only cache 'now'-relative queries (that can benefit from a backfill cache)
    //const shouldCache = request.rangeRaw?.to?.toString() === 'now';
    const shouldCache = true;

    // all targets are queried together, so we check for any that causes group cache invalidation & full re-query
    let doPartialQuery = shouldCache;
    let prevTo: TimestampMs | undefined = undefined;
    let prevFrom: TimestampMs | undefined = undefined;

    // pre-compute reqTargetSigs
    const reqTargetSigs = new Map<TargetIdent, TargetSig>();
    request.targets.forEach((target) => {
      let targetIdent = `${request.dashboardUID}|${request.panelId}|${target.refId}`;
      let targetSig = this.getTargetSignature(request, target); // ${request.maxDataPoints} ?

      reqTargetSigs.set(targetIdent, targetSig);
    });

    // figure out if new query range or new target props trigger full cache invalidation & re-query
    for (const [targetIdent, targetSig] of reqTargetSigs) {
      let cached = this.cache.get(targetIdent);
      let cachedSig = cached?.sig;

      if (cachedSig !== targetSig) {
        doPartialQuery = false;
      } else {
        // only do partial queries when new request range follows prior request range (possibly with overlap)
        // e.g. now-6h with refresh <= 6h
        prevTo = cached?.prevTo ?? Infinity;
        prevFrom = cached?.prevFrom ?? Infinity;

        doPartialQuery = newTo > prevTo && newFrom <= prevTo && newFrom >= prevFrom;
      }

      if (!doPartialQuery) {
        break;
      }
    }

    if (doPartialQuery && prevTo) {
      // clamp to make sure we don't re-query previous 10m when newFrom is ahead of it (e.g. 5min range, 30s refresh)
      let newFromPartial = Math.max(prevTo - this.overlapWindowMs, newFrom);

      const newToDate = dateTime(newTo);
      const newFromPartialDate = dateTime(incrRoundDn(newFromPartial, request.intervalMs));

      // modify to partial query
      request = {
        ...request,
        range: {
          ...request.range,
          from: newFromPartialDate,
          to: newToDate,
        },
      };
    } else {
      reqTargetSigs.forEach((targetSig, targetIdent) => {
        this.cache.delete(targetIdent);
      });
    }

    return {
      requests: [request],
      targetSigs: reqTargetSigs,
      shouldCache,
    };
  }

  // should amend existing cache with new frames and return full response
  procFrames(
    request: DataQueryRequest<T>,
    requestInfo: CacheRequestInfo<T> | undefined,
    respFrames: DataFrame[]
  ): DataFrame[] {
    if (requestInfo?.shouldCache) {
      const newFrom = request.range.from.valueOf();
      const newTo = request.range.to.valueOf();

      // group frames by targets
      const respByTarget = new Map<TargetIdent, DataFrame[]>();

      respFrames.forEach((frame: DataFrame) => {
        let targetIdent = `${request.dashboardUID}|${request.panelId}|${frame.refId}`;

        let frames = respByTarget.get(targetIdent);

        if (!frames) {
          frames = [];
          respByTarget.set(targetIdent, frames);
        }

        frames.push(frame);
      });

      let outFrames: DataFrame[] = [];

      respByTarget.forEach((respFrames, targetIdent) => {
        let cachedFrames = (targetIdent ? this.cache.get(targetIdent)?.frames : null) ?? [];

        respFrames.forEach((respFrame: DataFrame) => {
          // skip empty frames
          if (respFrame.length === 0 || respFrame.fields.length === 0) {
            return;
          }

          // frames are identified by their second (non-time) field's name + labels
          // TODO: maybe also frame.meta.type?
          let respFrameIdent = getFieldIdent(respFrame.fields[1]);

          let cachedFrame = cachedFrames.find((cached) => getFieldIdent(cached.fields[1]) === respFrameIdent);

          if (!cachedFrame) {
            // append new unknown frames
            cachedFrames.push(respFrame);
          } else {
            // we assume that fields cannot appear/disappear and will all exist in same order

            // amend & re-cache
            // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
            let prevTable: Table = cachedFrame.fields.map((field) => field.values) as Table;
            // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
            let nextTable: Table = respFrame.fields.map((field) => field.values) as Table;

            let amendedTable = amendTable(prevTable, nextTable);
            if (amendedTable) {
              for (let i = 0; i < amendedTable.length; i++) {
                cachedFrame.fields[i].values = amendedTable[i];
              }
              cachedFrame.length = cachedFrame.fields[0].values.length;
            }
          }
        });

        // trim all frames to in-view range, evict those that end up with 0 length
        let nonEmptyCachedFrames: DataFrame[] = [];

        cachedFrames.forEach((frame) => {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          let table: Table = frame.fields.map((field) => field.values) as Table;

          let trimmed = trimTable(table, newFrom, newTo);

          if (trimmed[0].length > 0) {
            for (let i = 0; i < trimmed.length; i++) {
              frame.fields[i].values = trimmed[i];
            }
            nonEmptyCachedFrames.push(frame);
          }
        });

        this.cache.set(targetIdent, {
          sig: requestInfo.targetSigs.get(targetIdent)!,
          frames: nonEmptyCachedFrames,
          prevFrom: newFrom,
          prevTo: newTo,
        });

        outFrames.push(...nonEmptyCachedFrames);
      });

      respFrames = outFrames.map((frame) => ({
        ...frame,
        fields: frame.fields.map((field) => ({
          ...field,
          config: {
            ...field.config,
          },
          values: Array.from(field.values).slice(),
        })),
      }));
    }

    return respFrames;
  }
}
