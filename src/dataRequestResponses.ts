import { AssetInfo } from 'types';
import { FieldType } from '@grafana/data';

export const mockedAssetInfo: AssetInfo = {
  id: '1',
  name: 'Asset 1',
  properties: [],
  hierarchy: [],
  arn: 'arn:aws:iot:us-west-2:123456789012:thing/Asset1',
  model_id: 'modelId',
};

export const mockedAssetInfoResponse = {
  data: [
    {
      name: 'assetInfo',
      fields: [
        { name: 'name', type: 'string', values: [mockedAssetInfo.name] },
        { name: 'id', type: 'string', values: [mockedAssetInfo.id] },
        { name: 'arn', type: 'string', values: [mockedAssetInfo.arn] },
        { name: 'model_id', type: 'string', values: [mockedAssetInfo.model_id] },
        { name: 'properties', type: 'string', values: [JSON.stringify(mockedAssetInfo.properties)] },
        { name: 'hierarchies', type: 'string', values: [JSON.stringify(mockedAssetInfo.hierarchy)] },
      ],
      length: 1,
    },
  ],
};

export const mockedAssetProperties = {
  name: 'assetProperties',
  fields: [
    { name: 'id', type: FieldType.string, values: ['1'], config: {} },
    { name: 'name', type: FieldType.string, values: ['Property 1'], config: {} },
  ],
  length: 1,
};

export const mockedListAssetPropertiesResponse = {
  data: [mockedAssetProperties],
};

export const mockedAssetModelSummary = {
  name: 'assetModelSummary',
  fields: [
    { name: 'id', type: FieldType.string, values: ['modelId'], config: {} },
    { name: 'name', type: FieldType.string, values: ['Model 1'], config: {} },
    { name: 'arn', type: FieldType.string, values: ['arn:aws:iot:us-west-2:123456789012:thing/Model1'], config: {} },
    { name: 'properties', type: FieldType.string, values: [JSON.stringify([])], config: {} },
    { name: 'hierarchies', type: FieldType.string, values: [JSON.stringify([])], config: {} },
  ],
  length: 1,
};

export const mockedAssetModelSummaryResponse = {
  data: [mockedAssetModelSummary],
};

export const mockedAssetSummary = {
  name: 'assetSummary',
  fields: [
    { name: 'id', type: FieldType.string, values: ['1'], config: {} },
    { name: 'name', type: FieldType.string, values: ['Asset 1'], config: {} },
    { name: 'model_id', type: FieldType.string, values: ['modelId'], config: {} },
    { name: 'arn', type: FieldType.string, values: ['arn:aws:iot:us-west-2:123456789012:thing/Asset1'], config: {} },
    { name: 'creation_date', type: FieldType.time, values: [0], config: {} },
    { name: 'last_update', type: FieldType.time, values: [0], config: {} },
    { name: 'state', type: FieldType.string, values: ['ACTIVE'], config: {} },
    { name: 'hierarchies', type: FieldType.string, values: [JSON.stringify([])], config: {} },
  ],
  length: 1,
};

export const mockedAssetSummaryResponse = {
  data: [mockedAssetSummary],
};
