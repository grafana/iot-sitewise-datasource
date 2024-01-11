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
} from 'types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { firstLabelWith } from './QueryEditor';
import { EditorField } from '@grafana/experimental';

interface Props
  extends SitewiseQueryEditorProps<
    AssetPropertyValueHistoryQuery | AssetPropertyAggregatesQuery | AssetPropertyInterpolatedQuery
  > {
  newFormStylingEnabled?: boolean;
}

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

  onOrderChange = (sel: SelectableValue<SiteWiseTimeOrder>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, timeOrdering: sel.value });
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

  render() {
    const { query } = this.props;
    return this.props.newFormStylingEnabled ? (
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
        <EditorField label="Time" width={10} htmlFor="time">
          <Select
            id="time"
            aria-label="Time"
            options={ordering}
            value={ordering.find((v) => v.value === query.timeOrdering) ?? ordering[0]}
            onChange={this.onOrderChange}
            isSearchable={true}
            menuPlacement="auto"
          />
        </EditorField>
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
    ) : (
      <>
        <div className="gf-form">
          <InlineField htmlFor="quality" label="Quality" labelWidth={firstLabelWith}>
            <Select
              inputId="quality"
              width={20}
              options={qualities}
              value={qualities.find((v) => v.value === query.quality) ?? qualities[0]}
              onChange={this.onQualityChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField htmlFor="time" label="Time" labelWidth={8}>
            <Select
              inputId="time"
              options={ordering}
              value={ordering.find((v) => v.value === query.timeOrdering) ?? ordering[0]}
              onChange={this.onOrderChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>

          <InlineField htmlFor="format" label="Format" labelWidth={8}>
            <Select
              inputId="format"
              value={query.responseFormat || SiteWiseResponseFormat.Table}
              onChange={this.onResponseFormatChange}
              options={FORMAT_OPTIONS}
            />
          </InlineField>

          {isAssetPropertyInterpolatedQuery(query) && (
            <InlineField htmlFor="resolution" label="Resolution" labelWidth={10}>
              <Select
                inputId="resolution"
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
