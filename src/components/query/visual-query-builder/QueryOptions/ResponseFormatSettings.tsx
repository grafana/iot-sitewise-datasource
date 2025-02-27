import { EditorField } from '@grafana/plugin-ui';
import { Select } from '@grafana/ui';
import React from 'react';
import { SitewiseQuery, SiteWiseResponseFormat } from 'types';
import { SelectableValue } from '@grafana/data';

const FORMAT_OPTIONS: Array<SelectableValue<SiteWiseResponseFormat>> = [
  { label: 'Table', value: SiteWiseResponseFormat.Table },
  { label: 'Time series', value: SiteWiseResponseFormat.TimeSeries },
];

export const ResponseFormatSettings = ({
  onChange,
  query,
}: {
  query: SitewiseQuery;
  onChange: (value: SitewiseQuery) => void;
}) => {
  const onResponseFormatChange = (sel: SelectableValue<SiteWiseResponseFormat>) => {
    onChange({ ...query, responseFormat: sel.value });
  };

  return (
    <EditorField label="Format" width={10} htmlFor="format">
      <Select
        id="format"
        inputId="format"
        aria-label="Format"
        value={query.responseFormat || SiteWiseResponseFormat.Table}
        onChange={onResponseFormatChange}
        options={FORMAT_OPTIONS}
      />
    </EditorField>
  );
};
