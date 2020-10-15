import React, { PureComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { AssetPropertyValueQuery } from '../types';
import { InlineField, Select } from '@grafana/ui';
import { SitewiseQueryEditorProps } from './types';

type Props = SitewiseQueryEditorProps<AssetPropertyValueQuery>;

export class QueryPropertyValueEditor extends PureComponent<Props> {
  onAssetIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, assetId: sel.value! });
    onRunQuery();
  };

  onPropertyIdChange = (sel: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, propertyId: sel.value! });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    const assets: Array<SelectableValue<string>> = [];
    const properties: Array<SelectableValue<string>> = [];

    if (query.assetId) {
      assets.push({
        label: query.assetId,
        value: query.assetId,
      });
    }

    if (query.propertyId) {
      properties.push({
        label: query.propertyId,
        value: query.propertyId,
      });
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Asset" labelWidth={10} grow={true}>
            <Select
              options={assets}
              value={assets.find(v => v.value === query.assetId) || undefined}
              onChange={this.onAssetIdChange}
              placeholder="Select an asset"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              formatCreateLabel={txt => `Asset: ${txt}`}
            />
          </InlineField>
        </div>
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
