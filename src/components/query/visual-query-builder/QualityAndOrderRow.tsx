import { type SelectableValue } from '@grafana/data';
import { EditorField } from '@grafana/plugin-ui';
import { Select } from '@grafana/ui';
import React, { useCallback } from 'react';
import {
  SiteWiseTimeOrder,
  SiteWiseQuality,
  SiteWiseResponseFormat,
  QueryType,
  type AssetPropertyValueHistoryQuery,
  type AssetPropertyAggregatesQuery,
} from 'types';
import type { SitewiseQueryEditorProps } from './types';

const QUALITY_OPTIONS = [
  { value: SiteWiseQuality.GOOD, label: 'GOOD' },
  { value: SiteWiseQuality.BAD, label: 'BAD' },
  { value: SiteWiseQuality.UNCERTAIN, label: 'UNCERTAIN' },
] satisfies Array<SelectableValue<SiteWiseQuality>>;

const ORDERING_OPTIONS = [
  { value: SiteWiseTimeOrder.ASCENDING, label: 'ASCENDING' },
  { value: SiteWiseTimeOrder.DESCENDING, label: 'DESCENDING' },
] satisfies Array<SelectableValue<SiteWiseTimeOrder>>;

export const FORMAT_OPTIONS = [
  { label: 'Table', value: SiteWiseResponseFormat.Table },
  { label: 'Time series', value: SiteWiseResponseFormat.TimeSeries },
] satisfies Array<SelectableValue<SiteWiseResponseFormat>>;

export const QualityAndOrderRow = ({ onChange, query }: SitewiseQueryEditorProps) => {
  const onQualityChange = useCallback(
    (sel: SelectableValue<SiteWiseQuality>) => {
      onChange({ ...query, quality: sel.value });
    },
    [onChange, query]
  );

  const onResponseFormatChange = useCallback(
    (sel: SelectableValue<SiteWiseResponseFormat>) => {
      onChange({ ...query, responseFormat: sel.value });
    },
    [onChange, query]
  );

  const onOrderChange = useCallback(
    (sel: SelectableValue<SiteWiseTimeOrder>) => {
      onChange({ ...query, timeOrdering: sel.value } as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery);
    },
    [onChange, query]
  );

  return (
    <>
      <EditorField label="Quality" width={15} htmlFor="quality">
        <Select
          id="quality"
          aria-label="Quality"
          options={QUALITY_OPTIONS}
          value={QUALITY_OPTIONS.find((v) => v.value === query.quality) ?? QUALITY_OPTIONS[0]}
          onChange={onQualityChange}
          isSearchable
          menuPlacement="auto"
        />
      </EditorField>

      {query.queryType !== QueryType.PropertyInterpolated && (
        <EditorField label="Time" width={10} htmlFor="time">
          <Select
            id="time"
            aria-label="Time"
            options={ORDERING_OPTIONS}
            value={
              ORDERING_OPTIONS.find(
                (v) => v.value === (query as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery).timeOrdering
              ) ?? ORDERING_OPTIONS[0]
            }
            onChange={onOrderChange}
            isSearchable={true}
            menuPlacement="auto"
          />
        </EditorField>
      )}

      <EditorField label="Format" width={10} htmlFor="format">
        <Select
          id="format"
          aria-label="Format"
          value={query.responseFormat || SiteWiseResponseFormat.Table}
          onChange={onResponseFormatChange}
          options={FORMAT_OPTIONS}
        />
      </EditorField>
    </>
  );
};
