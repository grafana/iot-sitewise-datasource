import { CustomVariableSupport, DataFrameView, DataQueryRequest, MetricFindValue } from '@grafana/data';
import { QueryEditor } from 'components/query/QueryEditor';
import { HierarchyInfo } from './queryResponseTypes';

import { DataSource } from 'DataSource';
import { QueryType, SitewiseQuery } from 'types';
import { mergeMap } from 'rxjs/operators';
import { of, from } from 'rxjs';

export class SitewiseVariableSupport extends CustomVariableSupport<DataSource, SitewiseQuery> {
  constructor(private readonly datasource: DataSource) {
    super();
  }

  editor = QueryEditor;

  query = (request: DataQueryRequest<SitewiseQuery>) => {
    const { targets } = request;
    if (targets && targets.length === 1) {
      const query = targets[0];
      if (query.queryType === QueryType.DescribeAsset) {
        if (!query.assetId) {
          return of({ data: [] });
        }

        return from(this.datasource.getCache(query.region).getAssetInfo(query.assetId)).pipe(
          mergeMap((rsp) => {
            console.log('CHECK', rsp);
            const data: MetricFindValue[] = [];
            for (const p of rsp.properties) {
              data.push({
                text: p.Name,
                value: p.Id,
              });
            }
            return of({ data });
          })
        );
      }

      if (query.queryType === QueryType.ListAssets || query.queryType === QueryType.ListAssetModels) {
        return this.datasource.query(request).pipe(
          mergeMap((rsp) => {
            if (!rsp.data?.length) {
              return of({ data: [] });
            }

            const assets = new DataFrameView<HierarchyInfo>(rsp.data[0]);
            const data: MetricFindValue[] = assets.map((a) => ({
              text: a.name,
              value: a.id,
            }));
            return of({ data });
          })
        );
      }
    }

    return this.datasource.query(request);
  };
}
