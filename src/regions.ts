import { type SelectableValue } from '@grafana/data';

// see https://docs.aws.amazon.com/general/latest/gr/iot-sitewise.html#iot-sitewise_region-sdk
// order based on order in documentation link
export const supportedRegions = [
  'us-east-2',
  'us-east-1',
  'us-west-2',
  'ap-south-1',
  'ap-northeast-2',
  'ap-southeast-1',
  'ap-southeast-2',
  'ap-northeast-1',
  'ca-central-1',
  'eu-central-1',
  'eu-west-1',
  'us-gov-west-1',
  'cn-north-1',
  'Edge',
] as const;

// backend is configured to use the user's configured default region when /query
// is called with an empty string for the region
export const DEFAULT_REGION = '';

export type DefaultRegion = typeof DEFAULT_REGION;

export type Region = (typeof supportedRegions)[number] | DefaultRegion;

export const regionOptions = supportedRegions.map((v) => ({
  value: v,
  label: v,
})) satisfies Array<SelectableValue<Region>>;

export const isSupportedRegion = (region: Region | string | unknown): region is Region =>
  Boolean(supportedRegions.find((supportedRegion) => supportedRegion === region));
