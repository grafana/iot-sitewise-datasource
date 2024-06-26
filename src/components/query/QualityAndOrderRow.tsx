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
  SiteWiseResponseFormat,
  QueryType,
} from 'types';
import { Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { EditorField } from '@grafana/experimental';

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
  { value: SiteWiseResolution.TenSec, label: '10 Seconds', description: '1 point every 10 seconds' },
  { value: SiteWiseResolution.Min, label: 'Minute', description: '1 point every minute' },
  { value: SiteWiseResolution.TenMin, label: '10 Minutes', description: '1 point every 10 minutes' },
  { value: SiteWiseResolution.Hour, label: 'Hour', description: '1 point every hour' },
  { value: SiteWiseResolution.TenHour, label: '10 Hours', description: '1 point every 10 hours' },
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

export const FORMAT_OPTIONS: Array<SelectableValue<SiteWiseResponseFormat>> = [
  { label: 'Table', value: SiteWiseResponseFormat.Table },
  { label: 'Time series', value: SiteWiseResponseFormat.TimeSeries },
];

export class QualityAndOrderRow extends PureComponent<Props> {
  onQualityChange = (sel: SelectableValue<SiteWiseQuality>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, quality: sel.value });
  };

  onResponseFormatChange = (sel: SelectableValue<SiteWiseResponseFormat>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, responseFormat: sel.value });
  };

  onResolutionChange = (sel: SelectableValue<SiteWiseResolution>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, resolution: sel.value });
  };

  onMaxPageAggregations = (event: React.FormEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;

    onChange({ ...query, maxPageAggregations: +event.currentTarget.value });
  };

  timeOrderField = () => {
    const { onChange, query } = this.props;

    // PropertyInterpolated has no time ordering support
    if (query.queryType === QueryType.PropertyInterpolated) {
      return null;
    }

    const onOrderChange = (sel: SelectableValue<SiteWiseTimeOrder>) => {
      onChange({ ...query, timeOrdering: sel.value });
    };

    return (
      <EditorField label="Time" width={10} htmlFor="time">
        <Select
          id="time"
          aria-label="Time"
          options={ordering}
          value={ordering.find((v) => v.value === query.timeOrdering) ?? ordering[0]}
          onChange={onOrderChange}
          isSearchable={true}
          menuPlacement="auto"
        />
      </EditorField>
    );
  };

  render() {
    const { query } = this.props;
    return (
      <>
        <EditorField label="Quality" width={15} htmlFor="quality">
          <Select
            id="quality"
            aria-label="Quality"
            options={qualities}
            value={qualities.find((v) => v.value === query.quality) ?? qualities[0]}
            onChange={this.onQualityChange}
            isSearchable={true}
            menuPlacement="auto"
          />
        </EditorField>
        {this.timeOrderField()}
        <EditorField label="Format" width={10} htmlFor="format">
          <Select
            id="format"
            aria-label="Format"
            value={query.responseFormat || SiteWiseResponseFormat.Table}
            onChange={this.onResponseFormatChange}
            options={FORMAT_OPTIONS}
          />
        </EditorField>
        {isAssetPropertyInterpolatedQuery(query) && (
          <EditorField label="Resolution" width={25} htmlFor="resolution">
            <Select
              id="resolution"
              aria-label="Resolution"
              options={interpolatedResolutions}
              value={interpolatedResolutions.find((v) => v.value === query.resolution) || interpolatedResolutions[0]}
              onChange={this.onResolutionChange}
              menuPlacement="auto"
            />
          </EditorField>
        )}
      </>
    );
  }
}
