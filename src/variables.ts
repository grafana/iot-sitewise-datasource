import { Observable } from 'rxjs';
import { assign } from 'lodash';
import { SitewiseQuery } from './types';
import { DataSource } from './DataSource';
import { DataQueryRequest, DataQueryResponse, CustomVariableSupport } from '@grafana/data';
import { QueryEditor } from './components/query/QueryEditor';

export class SitewiseVariableSupport extends CustomVariableSupport<DataSource, SitewiseQuery, SitewiseQuery> {
  constructor(private readonly datasource: DataSource) {
    super();
    this.datasource = datasource;
    this.query = this.query.bind(this);
  }

  editor = QueryEditor;

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    assign(request.targets, [{ ...request.targets[0], refId: 'A' }]);
    return this.datasource.query(request);
  }
}
