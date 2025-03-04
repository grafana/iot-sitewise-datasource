import {
  CoreApp,
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  LoadingState,
  MetricFindValue,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { SitewiseCache } from 'sitewiseCache';
import { isListAssetsQuery, isPropertyQueryType, SitewiseOptions, SitewiseQuery, SiteWiseResolution } from './types';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { frameToMetricFindValues } from 'utils';
import { SitewiseVariableSupport } from 'variables';
import { SitewiseQueryPaginator } from 'SiteWiseQueryPaginator';
import { RelativeRangeCache } from 'RelativeRangeRequestCache/RelativeRangeCache';
import { DEFAULT_REGION, isSupportedRegion, type Region } from './regions';

export class DataSource extends DataSourceWithBackend<SitewiseQuery, SitewiseOptions> {
  // Easy access for QueryEditor
  readonly options: SitewiseOptions;
  readonly defaultQuery: string;
  private cache = new Map<string, SitewiseCache>();
  private relativeRangeCache = new RelativeRangeCache();

  constructor(instanceSettings: DataSourceInstanceSettings<SitewiseOptions>) {
    super(instanceSettings);
    this.options = instanceSettings.jsonData;
    this.defaultQuery = 'select $__selectAll from raw_time_series where $__unixEpochFilter(event_timestamp)';
    this.variables = new SitewiseVariableSupport(this);
  }

  /**
   * Get a region scoped cache
   */
  getCache(
    region: Region = isSupportedRegion(this.options.defaultRegion) ? this.options.defaultRegion : DEFAULT_REGION
  ): SitewiseCache {
    let v = this.cache.get(region);
    if (!v) {
      v = new SitewiseCache(this, region);
      this.cache.set(region, v);
    }
    return v;
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  getDefaultQuery(_: CoreApp): Partial<SitewiseQuery> {
    return {
      region: isSupportedRegion(this.options.defaultRegion) ? this.options.defaultRegion : DEFAULT_REGION,
      rawSQL: this.defaultQuery,
    };
  }

  async metricFindQuery(query: SitewiseQuery, options: any): Promise<MetricFindValue[]> {
    const request = {
      targets: [
        {
          ...query,
          refId: 'metricFindQuery',
        },
      ],
      range: options.range,
      rangeRaw: options.rangeRaw,
    } as DataQueryRequest<SitewiseQuery>;

    let res: DataQueryResponse | undefined;

    try {
      res = await this.query(request).toPromise();
    } catch (err) {
      return Promise.reject(err);
    }

    if (!res || !res.data || res.data.length <= 0) {
      return [];
    }
    return frameToMetricFindValues(res.data[0] as DataFrame);
  }

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: SitewiseQuery): boolean {
    if (!query.queryType) {
      return false; // skip the query
    }
    // Migrate assetId to assetIDs (v1.6)
    if (query.assetId) {
      const ids = new Set<string>();
      ids.add(query.assetId);
      if (query.assetIds) {
        for (const id of query.assetIds) {
          ids.add(id);
        }
      }
      query.assetIds = Array.from(ids);
      delete query.assetId;
    }

    if (isPropertyQueryType(query.queryType)) {
      return Boolean((query.assetIds?.length && query.propertyId) || query.propertyAlias);
    }
    return true; // keep the query
  }

  // returns string that will be shown in the panel header when the panel is collapsed
  getQueryDisplayText(query: SitewiseQuery): string {
    const cache = this.getCache(query.region);
    let txt: string = query.queryType;
    if (query.assetIds?.length) {
      const info = cache.getAssetInfoSync(query.assetIds[0]);
      if (!info) {
        return txt + ' / ' + query.assetIds.join('/');
      }
      txt += ' / ' + info.name;

      if (query.propertyId && info.properties) {
        const p = info.properties.find((v) => v.Id === query.propertyId);
        if (p) {
          txt += ' / ' + p.Name;
        } else {
          txt += ' / ' + query.propertyId;
        }
      }
    } else if (query.propertyAlias) {
      txt += ' / ' + query.propertyAlias;
    }
    return txt;
  }

  /**
   * Supports template variables for region, asset and property
   */
  applyTemplateVariables(query: SitewiseQuery, scopedVars: ScopedVars): SitewiseQuery {
    const templateSrv = getTemplateSrv();
    const interpolatedQuery = {
      ...query,
      propertyAlias: templateSrv.replace(query.propertyAlias, scopedVars),
      region: templateSrv.replace(query.region ?? DEFAULT_REGION, scopedVars) as Region | undefined,
      propertyId: templateSrv.replace(query.propertyId || '', scopedVars),
      assetId: templateSrv.replace(query.assetId || '', scopedVars),
      assetIds: query.assetIds?.flatMap((assetId) => templateSrv.replace(assetId, scopedVars, 'csv').split(',')) ?? [],
      resolution: query.resolution
        ? (templateSrv.replace(query.resolution, scopedVars) as SiteWiseResolution)
        : undefined,
    };
    if (isListAssetsQuery(interpolatedQuery)) {
      interpolatedQuery.modelId = templateSrv.replace(interpolatedQuery.modelId, scopedVars);
    }
    return interpolatedQuery;
  }

  runQuery(query: SitewiseQuery, maxDataPoints?: number): Observable<DataQueryResponse> {
    // @ts-ignore
    return this.query({ targets: [query], requestId: `iot.${counter++}`, maxDataPoints });
  }

  query(request: DataQueryRequest<SitewiseQuery>): Observable<DataQueryResponse> {
    const cachedInfo = request.range != null ? this.relativeRangeCache.get(request) : undefined;

    return new SitewiseQueryPaginator({
      request: cachedInfo?.refreshingRequest || request,
      queryFn: (request: DataQueryRequest<SitewiseQuery>) => {
        return super.query(request).toPromise();
      },
      cachedResponse: cachedInfo?.cachedResponse,
    })
      .toObservable()
      .pipe(
        // Cache the last (done) response
        tap({
          next: (response) => {
            if (response.state === LoadingState.Done) {
              if (response.data.length > 0) {
                this.relativeRangeCache.set(request, response);
              }
            }
          },
        })
      );
  }
}

let counter = 1000;
