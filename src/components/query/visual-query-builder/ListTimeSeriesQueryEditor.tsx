import React, { ChangeEvent, useState } from 'react';
import { SelectableValue } from '@grafana/data';
import { ListTimeSeriesQuery } from 'types';
import { Input, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';

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
    description: 'The time series is associated with an asset property.',
  },
  {
    label: 'DISASSOCIATED',
    value: 'DISASSOCIATED',
    description: "The time series isn't associated with any asset property.",
  },
];

export const ListTimeSeriesQueryEditorFunction = (props: Props) => {
  const [lastInput, setLastInput] = useState<string>();

  const onAliasPrefixChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    setLastInput('prefix');
    onChange({ ...query, aliasPrefix: e.target.value });
  };

  const onAssetIdChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    setLastInput('id');
    onChange({ ...query, assetId: e.target.value });
  };

  const onTimeSeriesTypeChange = (sel: SelectableValue<string>) => {
    const { onChange, query } = props;
    if (sel.value === 'ALL' && lastInput === 'prefix') {
      onChange({ ...query, timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL', assetId: undefined });
    } else if (sel.value === 'ALL' && lastInput === 'id') {
      onChange({
        ...query,
        timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL',
        aliasPrefix: undefined,
      });
    } else {
      onChange({ ...query, timeSeriesType: sel.value as 'ASSOCIATED' | 'DISASSOCIATED' | 'ALL' });
    }
  };

  const { query } = props;

  return (
    <EditorRow>
      <EditorFieldGroup>
        <EditorField label="Time Series Type" htmlFor="timeSeriesType" width={20}>
          <Select
            inputId="timeSeriesType"
            options={timeSeriesTypes}
            value={
              timeSeriesTypes.find((v: { value: string }) => v.value === query.timeSeriesType) || timeSeriesTypes[0]
            }
            onChange={onTimeSeriesTypeChange}
            placeholder="Select a property"
            menuPlacement="auto"
          />
        </EditorField>
        {query.timeSeriesType !== 'ASSOCIATED' && (
          <EditorField
            label="Alias Prefix"
            htmlFor="aliasPrefix"
            width={30}
            tooltip={'The alias prefix of the time series.'}
          >
            <Input
              id="aliasPrefix"
              value={query.aliasPrefix}
              onChange={onAliasPrefixChange}
              placeholder="Optional: alias prefix"
            />
          </EditorField>
        )}
        {query.timeSeriesType !== 'DISASSOCIATED' && (
          <EditorField
            label="Asset Id"
            htmlFor="assetId"
            width={30}
            tooltip={
              'The ID of the asset in which the asset property was created. This can be either the actual ID in UUID format, or else externalId: followed by the external ID, if it has one'
            }
          >
            {/* eslint-disable-next-line @typescript-eslint/no-deprecated */}
            <Input id="assetId" value={query.assetId} onChange={onAssetIdChange} placeholder="Optional: asset id" />
          </EditorField>
        )}
      </EditorFieldGroup>
    </EditorRow>
  );
};
