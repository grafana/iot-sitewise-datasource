import React, { useMemo } from 'react';
import { SelectableValue } from '@grafana/data';
import { SiteWiseResolution, AssetPropertyInterpolatedQuery } from 'types';
import { Select } from '@grafana/ui';
import { EditorField, EditorFieldGroup } from '@grafana/plugin-ui';
import { getSelectionInfo } from 'common/getSelectionInfo';
import { getVariableOptions } from 'common/getVariableOptions';

const RESOLUTIONS: Array<SelectableValue<SiteWiseResolution>> = [
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

export const InterpolatedResolutionSettings = ({
  onChange,
  query,
}: {
  query: AssetPropertyInterpolatedQuery;
  onChange: (value: AssetPropertyInterpolatedQuery) => void;
}) => {
  const onResolutionChange = (sel: SelectableValue<string>) => {
    onChange({ ...query, resolution: sel.value as SiteWiseResolution });
  };

  const resolution = useMemo(
    () =>
      getSelectionInfo(
        query.resolution || RESOLUTIONS[0].value,
        RESOLUTIONS,
        getVariableOptions({ keepVarSyntax: true })
      ),
    [query]
  );

  return (
    <EditorFieldGroup>
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
