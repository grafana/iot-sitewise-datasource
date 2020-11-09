import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SiteWiseTimeOrder,
  AssetPropertyValueHistoryQuery,
  AssetPropertyAggregatesQuery,
  SiteWiseQuality,
} from 'types';
import { InlineField, Input, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { firstLabelWith } from './QueryEditor';

type Props = SitewiseQueryEditorProps<AssetPropertyValueHistoryQuery | AssetPropertyAggregatesQuery>;

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
              value={qualities.find(v => v.value === query.quality) ?? qualities[0]}
              onChange={this.onQualityChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Time" labelWidth={8}>
            <Select
              options={ordering}
              value={ordering.find(v => v.value === query.timeOrdering) ?? ordering[0]}
              onChange={this.onOrderChange}
              isSearchable={true}
              menuPlacement="bottom"
            />
          </InlineField>
          <InlineField label="Pages per Query" labelWidth={8}>
            <Input
              type="number"
              min="0"
              value={query.maxPageAggregations ?? 1}
              placeholder="enter a number"
              onChange={this.onMaxPageAggregations}
              width={8}
              css=""
            />
          </InlineField>
        </div>
      </>
    );
  }
}
