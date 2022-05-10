import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SiteWiseTimeOrder,
  AssetPropertyValueHistoryQuery,
  AssetPropertyAggregatesQuery,
  AssetPropertyInterpolatedQuery,
  SiteWiseQuality,
  SiteWiseResolution,
  isAssetPropertyInterpolatedQuery,
} from 'types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { firstLabelWith } from './QueryEditor';

type Props = SitewiseQueryEditorProps<
  AssetPropertyValueHistoryQuery | AssetPropertyAggregatesQuery | AssetPropertyInterpolatedQuery
>;

const interpolatedResolutions: Array<SelectableValue<SiteWiseResolution>> = [
  {
    value: SiteWiseResolution.Auto,
    label: 'Auto',
    description:
      'Picks a resolution based on the time window. ' +
      'Will switch to raw data if higher than 1m resolution is needed',
  },
  { value: SiteWiseResolution.Sec, label: 'Second', description: '1 point every second' },
  { value: SiteWiseResolution.Sec, label: '10 Seconds', description: '1 point every 10 seconds' },
  { value: SiteWiseResolution.Min, label: 'Minute', description: '1 point every minute' },
  { value: SiteWiseResolution.Sec, label: '10 Minutes', description: '1 point every 10 minutes' },
  { value: SiteWiseResolution.Hour, label: 'Hour', description: '1 point every hour' },
  { value: SiteWiseResolution.Day, label: 'Day', description: '1 point every day' },
];

const qualities: Array<SelectableValue<SiteWiseQuality>> = [
  { value: SiteWiseQuality.ANY, label: 'ANY' },
  { value: SiteWiseQuality.GOOD, label: 'GOOD' },
  { value: SiteWiseQuality.BAD, label: 'BAD' },
  { value: SiteWiseQuality.UNCERTAIN, label: 'UNCERTAIN' },
];

const ordering: Array<SelectableValue<SiteWiseTimeOrder>> = [
  { value: SiteWiseTimeOrder.ASCENDING, label: 'ASCENDING' },
  { value: SiteWiseTimeOrder.DESCENDING, label: 'DESCENDING' },
];

export class QualityAndOrderRow extends PureComponent<Props> {
  onQualityChange = (sel: SelectableValue<SiteWiseQuality>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, quality: sel.value });
    onRunQuery();
  };

  onOrderChange = (sel: SelectableValue<SiteWiseTimeOrder>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, timeOrdering: sel.value });
    onRunQuery();
  };

  onResolutionChange = (sel: SelectableValue<SiteWiseResolution>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, resolution: sel.value } as any);
    onRunQuery();
  };

  onMaxPageAggregations = (event: React.FormEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;

    onChange({ ...query, maxPageAggregations: +event.currentTarget.value });
    onRunQuery();
  };

  render() {
    const { query } = this.props;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Quality" labelWidth={firstLabelWith}>
            <Select
              width={20}
              options={qualities}
              value={qualities.find((v) => v.value === query.quality) ?? qualities[0]}
              onChange={this.onQualityChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Time" labelWidth={8}>
            <Select
              options={ordering}
              value={ordering.find((v) => v.value === query.timeOrdering) ?? ordering[0]}
              onChange={this.onOrderChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>
          {isAssetPropertyInterpolatedQuery(query) && (
            <InlineField label="Resolution" labelWidth={10}>
              <Select
                width={18}
                options={interpolatedResolutions}
                value={interpolatedResolutions.find((v) => v.value === query.resolution) || interpolatedResolutions[0]}
                onChange={this.onResolutionChange}
                menuPlacement="bottom"
              />
            </InlineField>
          )}
          {/*<InlineField label="Pages per Query" labelWidth={8}>*/}
          {/*  <Input*/}
          {/*    type="number"*/}
          {/*    min="0"*/}
          {/*    value={query.maxPageAggregations ?? 1}*/}
          {/*    placeholder="enter a number"*/}
          {/*    onChange={this.onMaxPageAggregations}*/}
          {/*    width={8}*/}
          {/*    css=""*/}
          {/*  />*/}
          {/*</InlineField>*/}
        </div>
      </>
    );
  }
}
