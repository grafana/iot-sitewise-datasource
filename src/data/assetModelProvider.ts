import IoTSiteWise from 'aws-sdk/clients/iotsitewise';
import { SitewiseModelSummary } from '../types';

export class AssetModelProvider {
  private client: IoTSiteWise;
  //TODO: cache results in a better manner
  private static cache: { [region: string]: SitewiseModelSummary[] } = {};

  constructor(client: IoTSiteWise) {
    this.client = client;
  }

  async provide(): Promise<SitewiseModelSummary[] | undefined> {
    const region = this.client.config.region;

    if (region && AssetModelProvider.cache[region]) {
      return AssetModelProvider.cache[region];
    }

    const models: SitewiseModelSummary[] = [];
    let token: string | undefined = undefined;
    do {
      // @ts-ignore
      const { assetModelSummaries, nextToken } = await this.client.listAssetModels({ nextToken: token }).promise();
      token = nextToken;
      models.push(...assetModelSummaries);
    } while (token);

    if (region) {
      AssetModelProvider.cache[region] = models;
    }

    return models;
  }
}
