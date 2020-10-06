import React, { ChangeEvent } from 'react';
import { SelectableValue } from '@grafana/data';
import { SitewiseAggregateType, SitewiseQueryType, SitewiseResolution } from '../../types';
import { FormField, FormInlineField } from '../layout/Fields';
import { FormFieldMultiSelector, FormFieldSelector } from '../layout/FormFieldSelector';
import { Checkbox, Input, Select } from '@grafana/ui';

interface Props {
  onAggregateChange: (value: Array<SelectableValue<SitewiseAggregateType>>) => void;
  aggregateTypes: SitewiseAggregateType[];

  onResolutionChange: (value: SelectableValue<SitewiseResolution>) => void;
  dataResolution: SitewiseResolution;

  onQueryTypeChange: (value: SelectableValue<SitewiseQueryType>) => void;
  queryType: SitewiseQueryType;

  onStreamingChange: (event: ChangeEvent<HTMLInputElement>) => void;
  streaming: boolean;

  onStreamingIntervalBlur: (event: ChangeEvent<HTMLInputElement>) => void;
  onStreamingIntervalChange: (event: ChangeEvent<HTMLInputElement>) => void;
  streamingInterval: number;
}

const queryTypes: { [queryType: string]: SelectableValue<SitewiseQueryType> } = {
  [SitewiseQueryType.PropertyHistory]: {
    value: SitewiseQueryType.PropertyHistory,
    label: 'Property Value',
  },
  [SitewiseQueryType.Aggregate]: {
    value: SitewiseQueryType.Aggregate,
    label: 'Property Aggregate',
  },
};

const aggregateOptions: { [aggregate: string]: SelectableValue<SitewiseAggregateType> } = {
  ['AVERAGE']: { label: 'Average', value: 'AVERAGE' },
  ['COUNT']: { label: 'Count', value: 'COUNT' },
  ['MAXIMUM']: { label: 'Max', value: 'MAXIMUM' },
  ['MINIMUM']: { label: 'Min', value: 'MINIMUM' },
  ['SUM']: { label: 'Sum', value: 'SUM' },
  ['STANDARD_DEVIATION']: { label: 'Std. Deviation', value: 'STANDARD_DEVIATION' },
};

const resolutionOptions: { [resolution: string]: SelectableValue<SitewiseResolution> } = {
  ['1m']: { label: 'One Minute', value: '1m' },
  ['1h']: { label: 'One Hour', value: '1h' },
  ['1d']: { label: 'One Day', value: '1d' },
};

export const QueryTypePicker: React.FC<Props> = props => {
  const {
    queryType,
    aggregateTypes,
    dataResolution,
    streaming,
    streamingInterval,
    onStreamingIntervalChange,
    onStreamingIntervalBlur,
    onQueryTypeChange,
    onAggregateChange,
    onStreamingChange,
    onResolutionChange,
  } = props;

  return (
    <FormInlineField label="Query Type">
      <Select
        options={Object.values(queryTypes)}
        onChange={onQueryTypeChange}
        value={(queryType && queryTypes[queryType]) || queryTypes[0]}
        width={20}
      />

      {queryType === SitewiseQueryType.Aggregate && (
        <>
          <FormFieldMultiSelector
            label="Aggregate"
            placeholder="Select aggregations"
            options={Object.values(aggregateOptions)}
            onChange={onAggregateChange}
            value={aggregateTypes && aggregateTypes.map(ag => aggregateOptions[ag])}
            width={48}
          />

          <FormFieldSelector<SitewiseResolution>
            label="Resolution"
            placeholder="Select resolution"
            options={Object.values(resolutionOptions)}
            onChange={onResolutionChange}
            value={dataResolution && resolutionOptions[dataResolution]}
            width={24}
          />
        </>
      )}
      {(SitewiseQueryType.PropertyHistory === queryType || SitewiseQueryType.DisableStreaming === queryType) && (
        <FormInlineField label="Streaming">
          <FormField label="Enabled">
            <Checkbox value={streaming} onChange={onStreamingChange} />
          </FormField>
          <FormField label="Interval (ms)">
            <Input
              value={`${streamingInterval}`}
              type="text"
              placeholder="5000"
              onBlur={onStreamingIntervalBlur}
              onChange={onStreamingIntervalChange}
              width={16}
            />
          </FormField>
        </FormInlineField>
      )}
    </FormInlineField>
  );
};
