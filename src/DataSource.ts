import IoTSiteWise, {
  AssetProperty,
  AssetPropertyValueHistory,
  GetAssetPropertyValueHistoryRequest,
} from 'aws-sdk/clients/iotsitewise';
import { merge, Observable, Subscriber } from 'rxjs';

import {
  CircularDataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  FieldType,
} from '@grafana/data';

import { SitewiseDataSourceOptions, SitewiseQuery, SitewiseQueryType } from './types';
import { AWSError } from 'aws-sdk';
import { MutableDataFrame } from '@grafana/data/dataframe/MutableDataFrame';
import { SitewiseDataFrameRow, SitewiseDataProvider } from './data/sitewiseDataProvider';

interface FetchDataParams extends GetAssetPropertyValueHistoryRequest {
  region?: string;
}

interface BackfillDataParams {
  client: IoTSiteWise;
  region: string;
  params: FetchDataParams;
}

type FetchDataCallback = (results: AssetPropertyValueHistory, done: boolean, err: AWSError | any) => void;

const fetchData = (args: BackfillDataParams, callback: FetchDataCallback): void => {
  const { client, region, params } = args;
  // my guess is this wont work for cross-region. will likely need to initialize a client per region
  client.config.region = region;
  let request: GetAssetPropertyValueHistoryRequest = { ...params };
  client
    .getAssetPropertyValueHistory(request)
    .promise()
    .then(value => {
      callback(value.assetPropertyValueHistory, value.nextToken === null, value.$response.error);
      if (value.nextToken) {
        request.nextToken = value.nextToken;
        fetchData(
          {
            client: client,
            region: region,
            params: request,
          },
          callback
        );
      }
    })
    .catch(reason => {
      callback([], true, reason);
    });
};

const notify = (subscriber: Subscriber<DataQueryResponse>, refId: string, frame: MutableDataFrame) => {
  subscriber.next({
    data: [frame],
    key: refId,
  });
};

const queryKey = (request: DataQueryRequest, queryTarget: SitewiseQuery): string => {
  const { asset, property, refId } = queryTarget;
  const { dashboardId, panelId } = request;
  // return `${assetId}/${propertyId}/${refId}`;
  const key = [dashboardId, panelId, asset.assetId, property.id, refId].join('/');
  console.log('KEY: ', key);
  return key;
};

const fieldTypeForProperty = (property: AssetProperty): FieldType => {
  switch (property.dataType) {
    case 'BOOLEAN':
      return FieldType.boolean;
    case 'DOUBLE':
      return FieldType.number;
    case 'INTEGER':
      return FieldType.number;
    case 'STRING':
      return FieldType.string;
    default:
      return FieldType.other;
  }
};

const addAggregationFields = (query: SitewiseQuery, frame: MutableDataFrame) => {
  for (const agg of query.aggregateTypes) {
    switch (agg) {
      case 'AVERAGE':
        frame.addField({ name: 'average', type: FieldType.number });
        continue;
      case 'COUNT':
        frame.addField({ name: 'count', type: FieldType.number });
        continue;
      case 'MAXIMUM':
        frame.addField({ name: 'maximum', type: FieldType.number });
        continue;
      case 'MINIMUM':
        frame.addField({ name: 'minimum', type: FieldType.number });
        continue;
      case 'SUM':
        frame.addField({ name: 'sum', type: FieldType.number });
        continue;
      case 'STANDARD_DEVIATION':
        frame.addField({ name: 'standardDeviation', type: FieldType.number });
    }
  }
};

const newDataFrame = (query: SitewiseQuery): CircularDataFrame => {
  const { asset, property, refId } = query;

  const assetName = asset.assetName;
  const propertyName = property.name;

  let frame = new CircularDataFrame({
    append: 'tail',
    capacity: 1000000, // 100k capacity to account for large 'raw' data streams. Should this be configurable?
  });

  frame.addField({ name: 'time', type: FieldType.time });
  frame.name = `${assetName} (${propertyName})`;
  frame.refId = refId;

  if (SitewiseQueryType.PropertyHistory === query.queryType) {
    frame.addField({ name: 'value', type: fieldTypeForProperty(property) });
  }

  if (SitewiseQueryType.Aggregate === query.queryType) {
    addAggregationFields(query, frame);
  }

  return frame;
};

export class SitewiseDatasource extends DataSourceApi<SitewiseQuery, SitewiseDataSourceOptions> {
  client: IoTSiteWise;
  dataProvider: SitewiseDataProvider;
  settings: SitewiseDataSourceOptions;
  streamIntervalHandlers: Map<string, number | NodeJS.Timeout>; // compiler keeps complaining... so number | Timeout

  constructor(instanceSettings: DataSourceInstanceSettings<SitewiseDataSourceOptions>) {
    super(instanceSettings);
    this.settings = { ...instanceSettings.jsonData };
    this.streamIntervalHandlers = new Map();

    this.client = new IoTSiteWise({
      region: this.settings.defaultRegion,
      credentials: {
        accessKeyId: this.settings.accessKeyId,
        secretAccessKey: this.settings.secretAccessKey,
        sessionToken: this.settings.sessionToken,
      },
    });

    this.dataProvider = new SitewiseDataProvider({ client: this.client });
  }

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    const data = this.handlePreProcessing(request)
      .filter(q => q.asset && q.property)
      .map(query => {
        const { refId, streaming, streamingInterval } = query;

        return new Observable<DataQueryResponse>(subscriber => {
          let frame = newDataFrame(query);

          this.dataProvider.provide({ request: request, query: query }, (err, results) => {
            if (err) {
              subscriber.error(err);
            }

            results?.values.forEach(value => frame.add(value));
            notify(subscriber, refId, frame);

            if (results?.done && streaming && streamingInterval) {
              const intervalId = this.startLiveDataStream(request, query, (err, value) => {
                if (err) {
                  subscriber.error(err);
                }
                frame.add(value);
                notify(subscriber, query.refId, frame);
              });
              this.addStreamIntervalHandler(request, query, intervalId);
            }
          });
        });
      });

    return merge(...data);
  }

  private handlePreProcessing(request: DataQueryRequest<SitewiseQuery>): SitewiseQuery[] {
    return request.targets.filter(q => {
      console.log('processing:', q.queryType);
      switch (q.queryType) {
        case SitewiseQueryType.DisableStreaming:
          // todo: this should probably use the handle to 'subscriber'
          this.removeStreamIntervalHandler(request, q);
          return false;
        default:
          return true;
      }
    });
  }

  private startLiveDataStream(
    request: DataQueryRequest<SitewiseQuery>,
    queryTarget: SitewiseQuery,
    callback: (err: Error | null, value: SitewiseDataFrameRow | undefined) => void
  ): number | NodeJS.Timeout {
    const { streamingInterval } = queryTarget;

    return setInterval(async () => {
      this.dataProvider.getAssetPropertyValue(
        {
          request: request,
          query: queryTarget,
        },
        (err, results) => {
          callback(err, results);
        }
      );
    }, streamingInterval);
  }

  private addStreamIntervalHandler(
    request: DataQueryRequest,
    query: SitewiseQuery,
    intervalId: number | NodeJS.Timeout
  ) {
    if (this.getStreamIntervalHandler(request, query)) {
      this.removeStreamIntervalHandler(request, query);
    }
    this.streamIntervalHandlers.set(queryKey(request, query), intervalId);
  }

  private getStreamIntervalHandler(
    request: DataQueryRequest,
    query: SitewiseQuery
  ): number | NodeJS.Timeout | undefined {
    return this.streamIntervalHandlers.get(queryKey(request, query));
  }

  private removeStreamIntervalHandler(request: DataQueryRequest, query: SitewiseQuery) {
    console.log(this.streamIntervalHandlers);
    const intervalHandler = this.getStreamIntervalHandler(request, query);
    if (intervalHandler) {
      clearInterval(Number(intervalHandler));
      this.streamIntervalHandlers.delete(queryKey(request, query));
    }
  }

  async testDatasource() {
    // Implement a health check for your data source.
    const response = await this.client
      .listAssets({
        filter: 'TOP_LEVEL',
        maxResults: 1,
      })
      .promise();

    if (response.$response.error) {
      return {
        status: 'failure',
        message: response.$response.error.message,
      };
    }

    return {
      status: 'success',
      message: 'Success',
    };
  }
}
