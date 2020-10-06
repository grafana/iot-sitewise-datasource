import React, { useState } from 'react';
import { FormField, FormInlineField } from '../layout/Fields';
import { Input, Select } from '@grafana/ui';
import { SitewiseAsset, SitewiseModelSummary } from '../../types';
import { SelectableValue } from '@grafana/data';
import IoTSiteWise, { AssetSummaries, NextToken } from 'aws-sdk/clients/iotsitewise';

interface Props {
  onAssetIdChange?: (event: React.FormEvent<HTMLInputElement>) => void;
  onAssetIdBlur: (event: React.FormEvent<HTMLInputElement>) => void;
  onAssetNameChange: (value: SelectableValue<string>) => void;
  assetId?: string;
  asset: SitewiseAsset;

  // no need for client as props
  // TODO: mode the client to be dynamically created by region as a singleton
  client: IoTSiteWise;

  // onModelChange: (value: SelectableValue<SitewiseModelSummary>) => void;
  models: Array<SelectableValue<SitewiseModelSummary>>;
  model: SitewiseModelSummary;
}

const assetToSelectableValue = (asset?: SitewiseAsset): SelectableValue<string> => {
  if (asset) {
    return { value: asset.assetId, label: asset.assetName, description: asset.assetId };
  }
  return {};
};

const fetchAssetOptions = async (value: SelectableValue<SitewiseModelSummary>, client: IoTSiteWise) => {
  console.log(value);
  let data = [];
  if (value.value) {
    let token: NextToken | undefined = undefined;
    do {
      // @ts-ignore
      const { assetSummaries, nextToken } = await client
        .listAssets({
          assetModelId: value.value.id,
          nextToken: token,
        })
        .promise();
      token = nextToken;
      data.push(...assetSummaries);
    } while (token);
  }

  return data;
};

export const AssetPicker: React.FC<Props> = p => {
  const { asset, assetId, model, models, onAssetIdChange, onAssetIdBlur, onAssetNameChange, client } = p;

  const [assetIdValue, setAssetIdValue] = useState<string>();

  const [assetOptionsValue, setAssetOptionsValue] = useState<AssetSummaries>();

  const onAssetIdInputChange = (event: React.FormEvent<HTMLInputElement>) => {
    setAssetIdValue(event.currentTarget.value);
    onAssetIdChange && onAssetIdChange(event);
  };

  const onAssetIdInputBlur = (event: React.FormEvent<HTMLInputElement>) => {
    setAssetIdValue(event.currentTarget.value);
    onAssetIdBlur(event);
  };

  return (
    <>
      <FormInlineField label="Model">
        <FormField label="Name" width={3}></FormField>
        <Select
          width={48}
          onChange={async value => {
            const summaries = await fetchAssetOptions(value, client);
            setAssetOptionsValue(summaries);
          }}
          options={models}
          value={model && { value: model, label: model.name, description: model.description }}
        />
      </FormInlineField>

      <FormInlineField label="Asset">
        <FormField label="Name" width={3}>
          <Select
            width={48}
            value={asset ? assetToSelectableValue(asset) : {}}
            options={
              assetOptionsValue
                ? assetOptionsValue.map(v => {
                    return { label: v.name, value: v.id, description: v.id };
                  })
                : []
            }
            onChange={onAssetNameChange}
            placeholder="ex: Wind Turbine 1"
          />
        </FormField>

        <FormField label="ID" width={2}>
          <Input
            width={48}
            name="AssetId"
            value={assetIdValue || assetId}
            onChange={onAssetIdInputChange}
            onBlur={onAssetIdInputBlur}
            placeholder="ex: a999c93c-a2ed-4a31-b832-533a922440f9"
          />
        </FormField>
      </FormInlineField>
    </>
  );
};
