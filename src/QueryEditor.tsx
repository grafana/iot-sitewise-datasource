import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent, useState } from 'react';
import { Input, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { SitewiseDatasource } from './DataSource';
import {
  Setter,
  SitewiseAggregateType,
  SitewiseDataSourceOptions,
  SitewiseQuery,
  SitewiseQueryType,
  SitewiseResolution,
} from './types';
import {} from '@emotion/core';
import { AssetProperty } from 'aws-sdk/clients/iotsitewise';
import { FormField, FormInlineField } from './components/layout/Fields';
import { AssetExplorerModal } from './components/assets/AssetExplorerModal';
import { QueryTypePicker } from './components/query/QueryTypePicker';
import { AssetPicker } from './components/assets/AssetPicker';
import { AssetModelProvider } from './data/assetModelProvider';

/**
 * Sample resources:
 * - https://github.com/grafana/grafana/blob/master/public/app/plugins/datasource/prometheus/components/PromQueryField.tsx
 * - https://github.com/grafana/grafana/blob/master/public/app/plugins/datasource/prometheus/components/PromQueryEditor.tsx
 * - https://github.com/grafana/grafana/blob/master/public/app/plugins/datasource/cloudwatch/components/MetricsQueryFieldsEditor.tsx#
 * - https://github.com/grafana/grafana/blob/master/public/app/plugins/datasource/cloudwatch/components/MetricsQueryEditor.tsx
 * - https://github.com/grafana/grafana/blob/master/public/app/plugins/datasource/elasticsearch/components/ElasticsearchQueryField.tsx#L72
 *
 * SCSS classes:
 * - https://github.com/grafana/grafana/tree/master/public/sass
 *
 */

export interface SitewiseQueryEditorProps
  extends QueryEditorProps<SitewiseDatasource, SitewiseQuery, SitewiseDataSourceOptions> {}

type SelectableValueSetter<T> = Setter<SelectableValue<T> | undefined>;
type StringValueSetter = Setter<string | undefined>;

const defaultQuery: Partial<SitewiseQuery> = {
  streaming: false,
  streamingInterval: 60000,
  queryType: SitewiseQueryType.PropertyHistory,
};

interface IdentityIdPickerProps extends SitewiseQueryEditorProps {
  onInputChange: (
    setIdentityIdValue: StringValueSetter,
    setIdentityNameValue: SelectableValueSetter<AssetProperty>
  ) => (event: ChangeEvent<HTMLInputElement>) => void;
  onSelectChange: (
    setIdentityIdValue: StringValueSetter,
    setIdentityNameValue: SelectableValueSetter<AssetProperty>
  ) => (event: SelectableValue<AssetProperty>) => void;
  propertyId?: string;
  initialSelectedValue?: SelectableValue<AssetProperty>;
}

// need to move the identity id picker to its own component... TODO: learn react
const IdentityIdInputPicker = (props: IdentityIdPickerProps) => {
  const { query, propertyId, initialSelectedValue, onInputChange, onSelectChange } = defaults(props, defaultQuery);
  const { asset } = query;
  const [propertyNameValue, setPropertyNameValue] = useState<SelectableValue<AssetProperty> | undefined>();
  const [propertyIdValue, setPropertyIdValue] = useState<string>();

  return (
    <FormInlineField label="Property">
      <FormField label="Name" width={3}>
        <Select
          width={48}
          onChange={onSelectChange(setPropertyIdValue, setPropertyNameValue)}
          options={
            asset
              ? asset.assetProperties.map<SelectableValue<AssetProperty>>(p => ({
                  value: p,
                  label: p.name,
                  description: `(${p.unit}) ${p.dataType.toLowerCase()}`,
                }))
              : []
          }
          value={propertyNameValue || initialSelectedValue}
          placeholder="Select property name from drop down"
        />
      </FormField>
      <FormField label="ID" width={2}>
        <Input
          width={48}
          name="PropertyId"
          value={propertyIdValue || propertyId}
          onChange={onInputChange(setPropertyIdValue, setPropertyNameValue)}
          placeholder="ex: 75874949-5292-4eff-8674-bdb32d3797a1"
        />
      </FormField>
    </FormInlineField>
  );
};

export class QueryEditor extends PureComponent<SitewiseQueryEditorProps> {
  private modelProvider: AssetModelProvider;

  constructor(p: SitewiseQueryEditorProps) {
    super(p);

    this.modelProvider = new AssetModelProvider(p.datasource.client);
    this.modelProvider.provide().then(resp => {
      if (resp) {
        const models = resp.map(model => {
          return { value: model, label: model.name, description: model.description };
        });
        this.props.onChange({ ...this.props.query, models: models });
      }
    });
  }

  onRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, region: event.target.value });
    onRunQuery();
  };

  onStreamingChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    const { queryType } = query;

    let state = {
      ...query,
      streaming: event.target.checked,
      queryType: !event.target.checked ? SitewiseQueryType.DisableStreaming : queryType,
    };

    onChange(state);

    onRunQuery();

    // reset the querytype state, or else query ends up disabled
    if (SitewiseQueryType.DisableStreaming === state.queryType) {
      onChange({ ...query, queryType: queryType });
    }
  };

  onStreamingIntervalBlur = (event: ChangeEvent<HTMLInputElement>) => {
    const { query, onRunQuery } = this.props;

    this.onStreamingIntervalChange(event);

    if (query.streaming && query.streamingInterval && query.streamingInterval > 0) {
      onRunQuery();
    }
  };

  onStreamingIntervalChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;

    const intervalVal = Number(event.target.value);
    if (!isNaN(intervalVal)) {
      onChange({ ...query, streamingInterval: intervalVal });
    }
  };

  onAssetIdChange = (setAssetIdValue: StringValueSetter) => (event: ChangeEvent<HTMLInputElement>): void => {
    setAssetIdValue(event.target.value);
  };

  onAssetIdBlur = (setAssetIdValue: StringValueSetter) => async (event: ChangeEvent<HTMLInputElement>) => {
    const { query, onChange, onRunQuery, datasource } = this.props;
    const assetId = event.target.value;

    if (assetId) {
      const response = await datasource.client.describeAsset({ assetId: assetId }).promise();
      onChange({ ...query, asset: response });
      setAssetIdValue(response.assetId);
      // executes the query
      if (query.property) {
        onRunQuery();
      }
    }
  };

  fetchAsset = async (assetId: string): Promise<void> => {
    const { query, onChange, onRunQuery, datasource } = this.props;

    const response = await datasource.client.describeAsset({ assetId: assetId }).promise();
    onChange({ ...query, asset: response });
    // executes the query
    if (query.property) {
      onRunQuery();
    }
  };

  onAssetIDBlur = async (event: React.FormEvent<HTMLInputElement>) => {
    const assetId = event.currentTarget.value;

    if (assetId) {
      await this.fetchAsset(assetId);
    }
  };

  onPropertyIdChange = (
    setIdentityIdValue: StringValueSetter,
    setIdentityNameValue: SelectableValueSetter<AssetProperty>
  ) => (event: ChangeEvent<HTMLInputElement>): void => {
    const { query, onChange, onRunQuery } = this.props;
    const { asset } = query;

    const propertyID = event.target.value;

    setIdentityIdValue(event.target.value);

    if (asset) {
      const property = asset.assetProperties.find(p => p.id === propertyID);

      if (property) {
        setIdentityNameValue({
          value: property,
          label: property.name,
        });
        onChange({ ...query, property: property });
        onRunQuery();
      } else {
        // clear the drop down
        setIdentityNameValue(undefined);
      }
    }
  };

  onPropertySelectChange = (
    setIdentityIdValue: StringValueSetter,
    setIdentityNameValue: SelectableValueSetter<AssetProperty>
  ) => (event: SelectableValue<AssetProperty>): void => {
    const { onChange, query, onRunQuery } = this.props;
    if (event.value) {
      onChange({ ...query, property: event.value });
      setIdentityIdValue(event.value.id);
    }
    setIdentityNameValue(event);
    onRunQuery();
  };

  onAggregateChange = (selectableValues: Array<SelectableValue<SitewiseAggregateType>>) => {
    const { onChange, query, onRunQuery } = this.props;

    const values = selectableValues.map(v => v.value).filter((v): v is string => !!v);

    onChange({ ...query, aggregateTypes: values });

    if (SitewiseQueryType.Aggregate === query.queryType) {
      onRunQuery();
    }
  };

  onResolutionChange = (value: SelectableValue<SitewiseResolution>) => {
    const { onChange, query, onRunQuery } = this.props;

    if (value.value) {
      onChange({ ...query, dataResolution: value.value });
    }

    if (SitewiseQueryType.Aggregate === query.queryType) {
      onRunQuery();
    }
  };

  onQueryTypeChange = (value: SelectableValue<SitewiseQueryType>) => {
    const { onChange, query } = this.props;
    if (value.value && value.value !== query.queryType) {
      onChange({ ...query, queryType: value.value });
    }
  };

  onAssetNameChange = async (value: SelectableValue<string>) => {
    const assetId = value.value;
    if (assetId) {
      await this.fetchAsset(assetId);
    }
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { asset, property, region } = query;
    const { defaultRegion } = this.props.datasource.settings;
    const assetId = asset && asset.assetId;
    const propertyId = property && property.id;

    return (
      <div>
        <FormInlineField label="Search">
          <AssetExplorerModal datasource={this.props.datasource} isOpen={false} />
        </FormInlineField>

        <QueryTypePicker {...this} {...query} />

        <AssetPicker
          onAssetIdBlur={this.onAssetIDBlur}
          onAssetNameChange={this.onAssetNameChange}
          assetId={assetId}
          client={this.props.datasource.client}
          {...query}
        />

        <IdentityIdInputPicker
          propertyId={propertyId}
          initialSelectedValue={
            property && { value: property, label: property.name, description: property.dataType.toLowerCase() }
          }
          onInputChange={this.onPropertyIdChange}
          onSelectChange={this.onPropertySelectChange}
          {...this.props}
        />

        <FormInlineField label="AWS Region">
          <Input
            name="AWS Region"
            value={region || defaultRegion}
            placeholder="ex: us-east-1"
            onChange={this.onRegionChange}
            width={16}
          />
        </FormInlineField>
      </div>
    );
  }
}
