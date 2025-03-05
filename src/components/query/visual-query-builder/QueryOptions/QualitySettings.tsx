import React from 'react';
import { SelectableValue } from '@grafana/data';
import { SiteWiseQuality, SitewiseQuery } from 'types';
import { Select } from '@grafana/ui';
import { EditorField } from '@grafana/plugin-ui';

const qualities: Array<SelectableValue<SiteWiseQuality>> = [
  { value: SiteWiseQuality.GOOD, label: 'GOOD' },
  { value: SiteWiseQuality.BAD, label: 'BAD' },
  { value: SiteWiseQuality.UNCERTAIN, label: 'UNCERTAIN' },
];

export const QualitySettings = ({
  onChange,
  query,
}: {
  query: SitewiseQuery;
  onChange: (value: SitewiseQuery) => void;
}) => {
  const onQualityChange = (sel: SelectableValue<SiteWiseQuality>) => {
    onChange({ ...query, quality: sel.value });
  };

  return (
    <EditorField label="Quality" width={15} htmlFor="quality">
      <Select
        id="quality"
        inputId="quality"
        aria-label="Quality"
        options={qualities}
        value={qualities.find((v) => v.value === query.quality) ?? qualities[0]}
        onChange={onQualityChange}
        isSearchable={true}
        menuPlacement="auto"
      />
    </EditorField>
  );
};
