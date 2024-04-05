import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { assign } from 'lodash';
import { QueryType, SitewiseQuery } from './types';
import { DataSource } from './DataSource';
import { DataQueryRequest, DataQueryResponse, CustomVariableSupport, DataFrameView } from '@grafana/data';
import { QueryEditor } from './components/query/QueryEditor';
import { AssetModelSummary } from 'queryResponseTypes';

export class SitewiseVariableSupport extends CustomVariableSupport<DataSource, SitewiseQuery, SitewiseQuery> {
  constructor(private readonly datasource: DataSource) {
    super();
    this.datasource = datasource;
    this.query = this.query.bind(this);
  }

  editor = QueryEditor;

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    assign(request.targets, [{ ...request.targets[0], refId: 'A' }]);
    let response = this.datasource.query(request);
    switch (request.targets[0].queryType) {
      // These are the only query types that return non-timeseries results
      case QueryType.ListAssetModels:
      case QueryType.ListAssets:
      case QueryType.ListAssociatedAssets:
        return this.parseOptions(response);
      default:
        return response
    }
  }

  parseOptions(response: Observable<DataQueryResponse> ): Observable<DataQueryResponse> {
    return response.pipe(
      map((res) => {
        let data = []
        if (res.data.length) {
          data = res.data[0]
        }
        return {data: new DataFrameView<AssetModelSummary>(data)};
      }),
      map((res) => {
        let newData = res.data.map((m)=>{
          return {
          value: m.id,
          text: m.name,
        }})
        return {data:newData}
      })
    )
  }

}
