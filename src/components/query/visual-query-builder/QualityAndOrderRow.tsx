import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SiteWiseTimeOrder,
  SiteWiseQuality,
  SiteWiseResponseFormat,
  QueryType,
  AssetPropertyValueHistoryQuery,
  AssetPropertyAggregatesQuery,
} from 'types';
import { Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { EditorField } from '@grafana/plugin-ui';

const qualities: Array<SelectableValue<SiteWiseQuality>> = [
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

export class QualityAndOrderRow extends PureComponent<SitewiseQueryEditorProps> {
  onQualityChange = (sel: SelectableValue<SiteWiseQuality>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, quality: sel.value });
  };

  onResponseFormatChange = (sel: SelectableValue<SiteWiseResponseFormat>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, responseFormat: sel.value });
  };

  timeOrderField = () => {
    const { onChange, query } = this.props;

    // PropertyInterpolated has no time ordering support
    if (query.queryType === QueryType.PropertyInterpolated) {
      return null;
    }

    const onOrderChange = (sel: SelectableValue<SiteWiseTimeOrder>) => {
      onChange({ ...query, timeOrdering: sel.value } as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery);
    };

    return (
      <EditorField label="Time" width={10} htmlFor="time">
        <Select
          id="time"
          aria-label="Time"
          options={ordering}
          value={
            ordering.find(
              (v) => v.value === (query as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery).timeOrdering
            ) ?? ordering[0]
          }
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
      </>
    );
  }
}
