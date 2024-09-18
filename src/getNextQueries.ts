import { DataFrame, DataQueryRequest, DataQueryResponse } from '@grafana/data';
import { SitewiseCustomMeta, SitewiseNextQuery, SitewiseQuery } from 'types';

export function getNextQueries(request: DataQueryRequest<SitewiseQuery>, rsp?: DataQueryResponse) {
  if (rsp?.data?.length) {
    const next: SitewiseNextQuery[] = [];
    for (const frame of rsp.data as DataFrame[]) {
      const meta = frame.meta?.custom as SitewiseCustomMeta;
      if (meta && meta.nextToken) {
        const query = request.targets.find((t) => t.refId === frame.refId);
        if (query) {
          const existingNextQuery = next.find((v) => v.refId === frame.refId);
          if (existingNextQuery) {
            if (meta.entryId && meta.nextToken) {
              if (!existingNextQuery.nextTokens) {
                existingNextQuery.nextTokens = {};
              }
              existingNextQuery.nextTokens[meta.entryId] = meta.nextToken;
            }
          } else {
            next.push({
              ...query,
              nextToken: meta.nextToken,
              nextTokens: { ...(meta.entryId && meta.nextToken ? { [meta.entryId]: meta.nextToken } : {}) },
            });
          }
        }
      }
    }
    if (next.length) {
      return next;
    }
  }
  return undefined;
}
