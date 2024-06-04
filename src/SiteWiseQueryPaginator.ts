import { DataQueryRequest, DataQueryResponse, LoadingState } from '@grafana/data';
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
  // The cached response to set as the initial response.
  cachedResponse?: {
    start?: DataQueryResponse,
    end?: DataQueryResponse,
  },
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
    const { request: { requestId }, cachedResponse } = this.options;
    const subject = new Subject<DataQueryResponse>();

    if (cachedResponse?.start) {
      subject.next({ ...cachedResponse.start, state: LoadingState.Streaming, key: requestId });
    }

    this.paginateQuery(subject);

    return subject;
  }

  /**
   * Performs the pagination logic for the query response.
   * @param subject The subject to emit the query responses to.
   */
  private async paginateQuery(subject: Subject<DataQueryResponse>) {
    const { request: initialRequest, queryFn, cachedResponse } = this.options;
    const { requestId } = initialRequest;
    let paginatingRequest = initialRequest;

    try {
      let retrievedData = cachedResponse?.start?.data;
      let nextQueries: SitewiseNextQuery[] | undefined;
      let errorEncountered = false;  // whether there's a error response
      let count = 1;

      do {
        if (nextQueries != null) {
          paginatingRequest = {
            ...paginatingRequest,
            requestId: `${requestId}.${++count}`,
            targets: nextQueries,
          };
        }

        const response = await queryFn(paginatingRequest);
        if (retrievedData == null) {
          retrievedData = response.data;
        } else {
          retrievedData = appendMatchingFrames(retrievedData, response.data);
        }

        if (response.state === LoadingState.Error) {
          subject.next({ ...response, data: retrievedData, state: LoadingState.Error, key: requestId });
          break;
        }

        nextQueries = getNextQueries(paginatingRequest, response);
        const loadingState = nextQueries || cachedResponse?.end != null ? LoadingState.Streaming : LoadingState.Done;

        subject.next({ ...response, data: retrievedData, state: loadingState, key: requestId });
      } while (nextQueries != null && !subject.closed);

      if (cachedResponse?.end != null && !errorEncountered) {
        retrievedData = appendMatchingFrames(retrievedData, cachedResponse.end.data);
        subject.next({ ...cachedResponse.end, data: retrievedData , state: LoadingState.Done, key: requestId });
      }

      subject.complete();
    } catch (err) {
      subject.error(err);
    }
  }
}
