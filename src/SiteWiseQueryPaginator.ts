import { DataQueryRequest, DataQueryResponse, DataQueryResponseData, LoadingState } from '@grafana/data';
import { appendMatchingFrames } from 'appendFrames';
import { getNextQueries } from 'getNextQueries';
import { Subject } from 'rxjs';
import { SitewiseNextQuery, SitewiseQuery } from 'types';

/**
 * Options for the SitewiseQueryPaginator class.
 */
export interface SitewiseQueryPaginatorOptions {
  // The initial query request.
  request: DataQueryRequest<SitewiseQuery>,
  // The function to call to execute the query.
  queryFn: (request: DataQueryRequest<SitewiseQuery>) => Promise<DataQueryResponse>;
}

/**
 * This class is responsible for paginating through the query response
 * and handling any errors that may occur during the pagination process.
 */
export class SitewiseQueryPaginator {
  private options: SitewiseQueryPaginatorOptions;

  /**
   * Constructor for the SitewiseQueryPaginator class.
   * @param options The options for the paginator.
   */
  constructor(options: SitewiseQueryPaginatorOptions) {
    this.options = options;
  }

  /**
   * Returns an observable that can be subscribed to for the paginated query responses.
   * @returns An observable that emits the paginated query responses.
   */
  toObservable() {
    const subject = new Subject<DataQueryResponse>();

    this.paginateQuery(subject);

    return subject;
  }

  /**
   * Performs the pagination logic for the query response.
   * @param subject The subject to emit the query responses to.
   */
  private async paginateQuery(subject: Subject<DataQueryResponse>) {
    const { request: initialRequest, queryFn } = this.options;
    const { requestId } = initialRequest;

    try {
      let retrievedData: DataQueryResponseData[] | undefined;
      let nextQueries: SitewiseNextQuery[] | undefined;
      let count = 1;

      do {
        let request = initialRequest;
        if (nextQueries != null) {
          request = {
            ...request,
            requestId: `${requestId}.${++count}`,
            targets: nextQueries,
          };
        }

        const response = await queryFn(request);
        if (retrievedData == null) {
          retrievedData = response.data;
        } else {
          retrievedData = appendMatchingFrames(retrievedData, response.data);
        }

        if (response.state === LoadingState.Error) {
          subject.next({ ...response, data: retrievedData, state: LoadingState.Error, key: requestId });
          break;
        }

        nextQueries = getNextQueries(request, response);
        const loadingState = nextQueries ? LoadingState.Streaming : LoadingState.Done;

        subject.next({ ...response, data: retrievedData, state: loadingState, key: requestId });
      } while (nextQueries != null && !subject.closed);

      subject.complete();
    } catch (err) {
      subject.error(err);
    }
  }
}
