import { SelectableValue } from '@grafana/data';
import { EditorField } from '@grafana/plugin-ui';
import { Select } from '@grafana/ui';
import React from 'react';
import { SiteWiseTimeOrder, AssetPropertyAggregatesQuery, AssetPropertyValueHistoryQuery, SitewiseQuery } from 'types';

const ORDERING: Array<SelectableValue<SiteWiseTimeOrder>> = [
  { value: SiteWiseTimeOrder.ASCENDING, label: 'ASCENDING' },
  { value: SiteWiseTimeOrder.DESCENDING, label: 'DESCENDING' },
];

export const TimeOrderSettings = ({
  onChange,
  query,
}: {
  query: SitewiseQuery;
  onChange: (value: SitewiseQuery) => void;
}) => {
  const onOrderChange = (sel: SelectableValue<SiteWiseTimeOrder>) => {
    onChange({ ...query, timeOrdering: sel.value } as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery);
  };

  return (
    <EditorField label="Time" width={10} htmlFor="time">
      <Select
        id="time"
        aria-label="Time"
        options={ORDERING}
        value={
          ORDERING.find(
            (v) => v.value === (query as AssetPropertyAggregatesQuery | AssetPropertyValueHistoryQuery).timeOrdering
          ) ?? ORDERING[0]
        }
        onChange={onOrderChange}
        isSearchable={true}
        menuPlacement="auto"
      />
    </EditorField>
  );
};
