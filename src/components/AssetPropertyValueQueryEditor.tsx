import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { SitewiseQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';
import { AssetPickerRow } from './AssetPickerRow';

type Props = SitewiseQueryEditorProps<SitewiseQuery>;

export class AssetPropertyValueQueryEditor extends PureComponent<Props> {
  onPropertyIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, propertyId: sel.value! });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    const properties: Array<SelectableValue<string>> = [];

    if (query.propertyId) {
      properties.push({
        label: query.propertyId,
        value: query.propertyId,
      });
    }

    return (
      <>
        <AssetPickerRow {...this.props} />
        <div className="gf-form">
          <InlineField label="Property" labelWidth={10} grow={true}>
            <Select
              options={properties}
              value={properties.find(v => v.value === query.propertyId) || undefined}
              onChange={this.onPropertyIdChange}
              placeholder="Select a property"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              formatCreateLabel={txt => `Property: ${txt}`}
            />
          </InlineField>
        </div>
      </>
    );
  }
}
