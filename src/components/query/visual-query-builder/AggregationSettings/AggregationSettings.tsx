import React from 'react';
import { SelectableValue } from '@grafana/data';
import { AggregateType, SiteWiseResolution, AssetPropertyAggregatesQuery } from 'types';
import { Select } from '@grafana/ui';
import { EditorField, EditorFieldGroup } from '@grafana/plugin-ui';
import { getDefaultAggregate } from 'queryInfo';
import { useOptionsWithVariables } from 'common/useOptionsWithVariables';
import { AggregatePicker } from './AggregatePicker';

const RESOLUTIONS: Array<SelectableValue<string>> = [
  {
    value: SiteWiseResolution.Auto as string,
    label: 'Auto',
    description:
      'Picks a resolution based on the time window. ' +
      'Will switch to raw data if higher than 1m resolution is needed',
  },
  { value: SiteWiseResolution.Min as string, label: 'Minute', description: '1 point every minute' },
  { value: SiteWiseResolution.FifteenMin as string, label: '15 Minutes', description: '1 point every 15 minutes' },
  { value: SiteWiseResolution.Hour as string, label: 'Hour', description: '1 point every hour' },
  { value: SiteWiseResolution.Day as string, label: 'Day', description: '1 point every day' },
];

export const AggregationSettings = ({
  onChange,
  query,
}: {
  query: AssetPropertyAggregatesQuery;
  onChange: (value: AssetPropertyAggregatesQuery) => void;
}) => {
  const resolution = useOptionsWithVariables({ current: query.resolution, options: RESOLUTIONS });

  const onAggregateChange = (aggregates: AggregateType[]) => {
    onChange({ ...query, aggregates });
  };

  const onResolutionChange = (sel: SelectableValue<string>) => {
    onChange({ ...query, resolution: sel.value as SiteWiseResolution });
  };

  return (
    <EditorFieldGroup>
      <EditorField label="Aggregate" htmlFor="aggregate-picker" width={40}>
        <AggregatePicker
          stats={query.aggregates ?? []}
          onChange={onAggregateChange}
          defaultStat={getDefaultAggregate()}
          menuPlacement="auto"
        />
      </EditorField>
      <EditorField label="Resolution" htmlFor="resolution" width={25}>
        <Select
          inputId="resolution"
          aria-label="resolution"
          options={resolution.options}
          value={resolution.current}
          onChange={onResolutionChange}
          allowCustomValue
          menuPlacement="auto"
        />
      </EditorField>
    </EditorFieldGroup>
  );
};
