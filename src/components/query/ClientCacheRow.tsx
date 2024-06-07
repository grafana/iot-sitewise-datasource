import React from 'react';
import { InlineField, InlineSwitch, Switch } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/experimental';

interface Props {
  clientCache?: boolean;
  newFormStylingEnabled?: boolean;
  onClientCacheChange: (evt: React.SyntheticEvent<HTMLInputElement>) => void;
}

export const ClientCacheRow = ({clientCache, newFormStylingEnabled, onClientCacheChange}: Props) => {
  if (newFormStylingEnabled) {
    return (
      <EditorRow>
        <EditorFieldGroup>
          <EditorField
            label="Client cache"
            htmlFor="clientCache"
            tooltip="Enable to cache results in the browser that are older than 15 minutes. This will improve performance for repeated queries with relative time range."
          >
            <Switch id="clientCache" value={clientCache} onChange={onClientCacheChange} />
          </EditorField>
        </EditorFieldGroup>
      </EditorRow>
    );
  }

  return (
    <div className="gf-form">
      <InlineField
          label="Client cache"
          htmlFor="clientCache"
          tooltip="Enable to cache results from the query. This will improve performance for repeated queries with relative time range."
      >
        <InlineSwitch value={clientCache} onChange={onClientCacheChange} />
      </InlineField>
    </div>
  );
};
