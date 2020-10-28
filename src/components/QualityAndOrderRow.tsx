import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import {
  SiteWiseTimeOrder,
  AssetPropertyValueHistoryQuery,
  AssetPropertyAggregatesQuery,
  SiteWiseQuality,
} from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';

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

  render() {
    const { query } = this.props;

    return (
      <>
        <div className="gf-form">
          <InlineField label="Quality" labelWidth={10}>
            <Select
              width={20}
              options={qualities}
              value={qualities.find(v => v.value === query.quality) ?? qualities[0]}
              onChange={this.onQualityChange}
              isSearchable={true}
            />
          </InlineField>
          <InlineField label="Time" labelWidth={8}>
            <Select
              options={ordering}
              value={ordering.find(v => v.value === query.timeOrdering) ?? ordering[0]}
              onChange={this.onOrderChange}
              isSearchable={true}
            />
          </InlineField>
        </div>
      </>
    );
  }
}
