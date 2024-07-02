import React, { ChangeEvent, useState } from 'react';
import { SelectableValue } from '@grafana/data';
import { ListTimeSeriesQuery } from 'types';
import { InlineField, Input, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/experimental';
import { firstLabelWith } from './QueryEditor';

interface Props extends SitewiseQueryEditorProps<ListTimeSeriesQuery> {
  newFormStylingEnabled?: boolean;
}

const timeSeriesTypes = [
  {
    label: 'ALL',
    value: 'ALL',
    description: 'All time series data',
  },
  {
    label: 'ASSOCIATED',
    value: 'ASSOCIATED',
    description: "The time series is associated with an asset property.",
  },
  { label: 'DISASSOCIATED', value: 'DISASSOCIATED', description: "The time series isn't associated with any asset property." },
];

export const ListTimeSeriesQueryEditorFunction = (props: Props) => {

  const [lastInput, setLastInput] = useState<string>()

  const onAliasPrefixChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    setLastInput("prefix")
    onChange({ ...query, aliasPrefix: e.target.value });
  };

  const onAssetIdChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    setLastInput("id")
    onChange({ ...query, assetId: e.target.value });
  };

  const onTimeSeriesTypeChange = (sel: SelectableValue<string>) => {
    const { onChange, query } = props;
    if (sel.value === 'ALL' && lastInput === "prefix"){
      onChange({ ...query, timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL' , assetId: undefined});
    }
    else if (sel.value === 'ALL' && lastInput === "id"){
      onChange({ ...query, timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL' , aliasPrefix: undefined});
    } else {
      onChange({ ...query, timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL' });
    }
    
  };

  const { query, newFormStylingEnabled } = props;

  return newFormStylingEnabled ?
    (
      <EditorRow>
        <EditorFieldGroup>
          <EditorField label="timeSeriesType" htmlFor="timeSeriesType" width={20}>
            <Select
              id="timeSeriesType"
              aria-label="timeSeriesType"
              options={timeSeriesTypes}
              value={timeSeriesTypes.find((v: { value: string; }) => v.value === query.timeSeriesType) || timeSeriesTypes[0]}
              onChange={onTimeSeriesTypeChange}
              placeholder="Select a property"
              menuPlacement="auto"
            />
          </EditorField>
          <EditorField label="aliasPrefix" htmlFor="aliasPrefix" width={30}>
            <Input
              id="aliasPrefix"
              aria-label="Alias Prefix"
              value={query.aliasPrefix}
              onChange={onAliasPrefixChange}
              placeholder="Optional: The alias prefix of the time series."
            />
          </EditorField>
          <EditorField label="assetId" htmlFor="assetId" width={30}>
            <Input
              id="assetId"
              aria-label="Asset Id"
              value={query.assetId}
              onChange={onAssetIdChange}
              placeholder="The ID of the asset in which the asset property was created. This can be either the actual ID in UUID format, or else externalId: followed by the external ID, if it has one"
            />
          </EditorField>
        </EditorFieldGroup>
      </EditorRow>
    ) : (
      <>
        <div className="gf-form">
          <InlineField htmlFor="timeSeriesType" label="Time Series Type" labelWidth={firstLabelWith} grow={true}>
            <Select
              inputId="timeSeriesType"
              options={timeSeriesTypes}
              value={timeSeriesTypes.find((v: { value: string; }) => v.value === query.timeSeriesType) || timeSeriesTypes[0]}
              onChange={onTimeSeriesTypeChange}
              placeholder="Select a time series type"
              menuPlacement="bottom"
            />
          </InlineField>
        </div>
        {!Boolean(query.timeSeriesType ==="ASSOCIATED") &&
          <div className="gf-form">
            <InlineField htmlFor="aliasPrefix" label="Alias Prefix" labelWidth={firstLabelWith} grow={true}>
              <Input
                id="aliasPrefix"
                value={query.aliasPrefix}
                onChange={onAliasPrefixChange}
                placeholder="Optional: The alias prefix of the time series."
              />
            </InlineField>
          </div>}
        {!Boolean(query.timeSeriesType === "DISASSOCIATED") &&
          <div className="gf-form">
            <InlineField htmlFor="assetId" label="Asset Id" labelWidth={firstLabelWith} grow={true}>
              <Input
                id="assetId"
                value={query.assetId}
                onChange={onAssetIdChange}
                placeholder="The ID of the asset in which the asset property was created. This can be either the actual ID in UUID format, or else externalId: followed by the external ID, if it has one"
              />
            </InlineField>
          </div>}
      </>
    );
};



