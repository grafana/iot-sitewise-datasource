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
import { lastValueFrom, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { frameToMetricFindValues } from 'utils';
import { applyVariableForList, SitewiseVariableSupport, variableFormatter } from 'variables';
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
    this.defaultQuery = 'select $__selectAll from raw_time_series where $__timeFilter(event_timestamp)';
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
      res = await lastValueFrom(this.query(request));
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
    /* eslint-disable @typescript-eslint/no-deprecated */
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

    // Migrate propertyId to propertyIds (v2.1)
    if (query.propertyId) {
      const ids = new Set<string>();
      ids.add(query.propertyId);
      if (query.propertyIds) {
        for (const id of query.propertyIds) {
          ids.add(id);
        }
      }
      query.propertyIds = Array.from(ids);
      delete query.propertyId;
    }

    // Migrate propertyAlias to propertyAliases (v2.1)
    if (query.propertyAlias) {
      const aliases = new Set<string>();
      aliases.add(query.propertyAlias);
      if (query.propertyAliases) {
        for (const alias of query.propertyAliases) {
          aliases.add(alias);
        }
      }
      query.propertyAliases = Array.from(aliases);
      delete query.propertyAlias;
    }

    if (isPropertyQueryType(query.queryType)) {
      return Boolean((query.assetIds?.length && query.propertyIds?.length) || query.propertyAliases?.length);
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

      if (query.propertyIds?.length && info.properties) {
        const properties = info.properties.filter((v) => query.propertyIds?.includes(v.Id));
        if (properties.length > 0) {
          txt += ' / ' + properties.map((p) => p.Name).join('/');
        } else {
          txt += ' / ' + query.propertyIds.join('/');
        }
      }
    } else if (query.propertyAliases?.length) {
      txt += ' / ' + query.propertyAliases.join(' / ');
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
      propertyAliases: applyVariableForList(templateSrv, scopedVars, query.propertyAliases),
      propertyId: templateSrv.replace(query.propertyId || '', scopedVars),
      propertyIds: applyVariableForList(templateSrv, scopedVars, query.propertyIds),
      assetId: templateSrv.replace(query.assetId || '', scopedVars),
      assetIds: applyVariableForList(templateSrv, scopedVars, query.assetIds),
      resolution: query.resolution
        ? (templateSrv.replace(query.resolution, scopedVars) as SiteWiseResolution)
        : undefined,
      rawSQL: templateSrv.replace(query.rawSQL, scopedVars, variableFormatter),
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
    const interpolatedRequest = this.constructVariableInterpolatedRequest(request);
    const cachedInfo = request.range != null ? this.relativeRangeCache.get(interpolatedRequest) : undefined;

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
                this.relativeRangeCache.set(interpolatedRequest, response);
              }
            }
          },
        })
      );
  }

  private constructVariableInterpolatedRequest(
    request: DataQueryRequest<SitewiseQuery>
  ): DataQueryRequest<SitewiseQuery> {
    return {
      ...request,
      targets: this.interpolateVariablesInQueries(request.targets, request.scopedVars, request.filters),
    };
  }
}

let counter = 1000;
