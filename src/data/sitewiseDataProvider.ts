import { DataQuery, DataQueryRequest } from '@grafana/data';
import { SitewiseAggregateValue, SitewiseQuery, SitewiseQueryType } from '../types';
import IoTSiteWise, {
  AssetProperty,
  AssetPropertyValue,
  GetAssetPropertyAggregatesRequest,
  GetAssetPropertyAggregatesResponse,
  GetAssetPropertyValueHistoryRequest,
  GetAssetPropertyValueHistoryResponse,
  GetAssetPropertyValueRequest,
  GetAssetPropertyValueResponse,
} from 'aws-sdk/clients/iotsitewise';

export interface SitewiseDataFrameRow {
  time: number;
  value?: number | string | boolean;
  minimum?: number;
  maximum?: number;
  standardDeviation?: number;
  sum?: number;
  count?: number;
  average?: number;
}

export type SitewiseDataFrameValues = SitewiseDataFrameRow[];

export interface SitewiseDataPageResponse {
  values: SitewiseDataFrameValues;
  done?: boolean;
}

export type GetDataCallback<V> = (err: Error | null, results?: V) => void;

export interface GetDataRequest<T extends DataQuery> {
  request: DataQueryRequest<T>;
  query: T;
}

export interface DataProvider<T extends DataQuery, V> {
  provide(request: GetDataRequest<T>, callback: GetDataCallback<V>): void;
}

const getAssetPropetyVariantValue = (
  property: AssetProperty,
  value: AssetPropertyValue
): undefined | string | boolean | number => {
  switch (property.dataType) {
    case 'BOOLEAN':
      return value.value.booleanValue;
    case 'DOUBLE':
      return value.value.doubleValue;
    case 'INTEGER':
      return value.value.integerValue;
    case 'STRING':
      return value.value.stringValue;
    default:
      return undefined;
  }
};

const propertyValueToDataFrame = (property: AssetProperty, value: AssetPropertyValue): SitewiseDataFrameRow => {
  let timestamp = value.timestamp.timeInSeconds * 1000;
  if (value.timestamp.offsetInNanos) {
    timestamp = timestamp + value.timestamp.offsetInNanos / 1000;
  }
  return { time: timestamp, value: getAssetPropetyVariantValue(property, value) };
};

const aggregatedValueToDataFrame = (value: SitewiseAggregateValue): SitewiseDataFrameRow => {
  return {
    time: value.timestamp.getTime(),
    ...value.value,
  };
};

export class SitewiseDataProvider implements DataProvider<SitewiseQuery, SitewiseDataPageResponse> {
  client: IoTSiteWise;

  constructor({ client }: { client: IoTSiteWise }) {
    this.client = client;
  }

  provide(args: GetDataRequest<SitewiseQuery>, callback: GetDataCallback<SitewiseDataPageResponse>): void {
    const { query, request } = args;
    const { range } = request;
    const { asset, property, aggregateTypes, dataResolution } = query;
    const from = range!.from.valueOf();
    const to = range!.to.valueOf();

    switch (query.queryType) {
      case SitewiseQueryType.PropertyHistory:
        this.getAssetPropertyHistory(
          {
            assetId: asset.assetId,
            propertyId: property.id,
            startDate: new Date(from),
            endDate: new Date(to),
            maxResults: 250,
          },
          (err, results) => {
            if (!err && results?.assetPropertyValueHistory) {
              const data = results?.assetPropertyValueHistory.map(value => propertyValueToDataFrame(property, value));
              callback(null, { values: data, done: results.nextToken === undefined });
            } else {
              callback(err);
            }
          }
        );
        return;
      case SitewiseQueryType.Aggregate:
        this.getAssetPropertyAggregate(
          {
            aggregateTypes: aggregateTypes,
            assetId: asset.assetId,
            maxResults: 250,
            propertyId: property.id,
            resolution: dataResolution,
            startDate: new Date(from),
            endDate: new Date(to),
          },
          (err, results) => {
            if (!err && results?.aggregatedValues) {
              const data = results.aggregatedValues.map(value => aggregatedValueToDataFrame(value));
              callback(null, { values: data, done: results.nextToken === undefined });
            } else {
              callback(err);
            }
          }
        );
        return;
      default:
        callback(new Error(`Unable to determine query action for: ${query.queryType}`));
    }
  }

  private getAssetPropertyHistory(
    request: GetAssetPropertyValueHistoryRequest,
    callback: GetDataCallback<GetAssetPropertyValueHistoryResponse>
  ) {
    let req: GetAssetPropertyValueHistoryRequest = { ...request };

    this.client.getAssetPropertyValueHistory(req, (err, resp) => {
      if (err) {
        callback(err);
      } else {
        callback(null, resp);
        if (resp.nextToken) {
          req.nextToken = resp.nextToken;
          this.getAssetPropertyHistory(req, callback);
        }
      }
    });
  }

  private getAssetPropertyAggregate(
    request: GetAssetPropertyAggregatesRequest,
    callback: GetDataCallback<GetAssetPropertyAggregatesResponse>
  ) {
    let req: GetAssetPropertyAggregatesRequest = { ...request };

    this.client.getAssetPropertyAggregates(req, (err, resp) => {
      if (err) {
        callback(err);
      } else {
        callback(null, resp);
        if (resp.nextToken) {
          req.nextToken = resp.nextToken;
          this.getAssetPropertyAggregate(req, callback);
        }
      }
    });
  }

  private getAssetPropertyValueInternal(
    request: GetAssetPropertyValueRequest,
    callback: GetDataCallback<GetAssetPropertyValueResponse>
  ) {
    this.client.getAssetPropertyValue(request, (err, resp) => {
      callback(err, resp);
    });
  }

  getAssetPropertyValue(args: GetDataRequest<SitewiseQuery>, callback: GetDataCallback<SitewiseDataFrameRow>) {
    const { asset, property } = args.query;

    this.getAssetPropertyValueInternal({ assetId: asset.assetId, propertyId: property.id }, (err, results) => {
      if (!err && results?.propertyValue) {
        callback(null, propertyValueToDataFrame(property, results.propertyValue));
      } else {
        callback(err);
      }
    });
  }
}
