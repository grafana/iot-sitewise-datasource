import React, { PureComponent } from 'react';
import { AggregateType, AssetPropertyAggregatesQuery, SiteWiseResolution } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetPropPickerRows } from './AssetPropPickerRows';
import { AggregatePicker } from './AggregatePicker';
import { SelectableValue } from '@grafana/data';

type Props = SitewiseQueryEditorProps<AssetPropertyAggregatesQuery>;

const resolutions: Array<SelectableValue<SiteWiseResolution>> = [
  { value: SiteWiseResolution.Auto, label: 'Auto', description: 'Pick a resolution based on the time window' },
  { value: SiteWiseResolution.Min, label: 'Minute', description: '1 point every minute' },
  { value: SiteWiseResolution.Hour, label: 'Hour', description: '1 point every hour' },
  { value: SiteWiseResolution.Day, label: 'Day', description: '1 point every day' },
];

export class PropertyAggregatesEditor extends PureComponent<Props> {
  onAggregateChange = (aggregates: AggregateType[]) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, aggregates });
    onRunQuery();
  };

  onResolutionChange = (sel: SelectableValue<SiteWiseResolution>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, resolution: sel.value });
    onRunQuery();
  };

  render() {
    const { query } = this.props;

    return (
      <>
        <AssetPropPickerRows {...(this.props as any)} />
        <div className="gf-form">
          <InlineField label="Aggregate" labelWidth={10} grow={true}>
            <AggregatePicker stats={query.aggregates ?? []} onChange={this.onAggregateChange} />
          </InlineField>
          <InlineField label="Resolution" labelWidth={10}>
            <Select
              width={18}
              options={resolutions}
              value={resolutions.find(v => v.value === query.resolution) || resolutions[0]}
              onChange={this.onResolutionChange}
            />
          </InlineField>
        </div>
      </>
    );
  }
}
