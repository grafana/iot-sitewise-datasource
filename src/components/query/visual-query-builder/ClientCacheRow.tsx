import React from 'react';
import { Switch } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';

interface Props {
  clientCache?: boolean;
  onClientCacheChange: (evt: React.SyntheticEvent<HTMLInputElement>) => void;
}

export const ClientCacheRow = ({ clientCache, onClientCacheChange }: Props) => {
  return (
    <EditorRow>
      <EditorFieldGroup>
        <EditorField
          label="Client cache"
          htmlFor="clientCache"
          tooltip="Enable to cache results in the browser that are older than 15 minutes. Note: Dashboard variable query result will not update when client cache is enabled."
        >
          <Switch id="clientCache" value={clientCache} onChange={onClientCacheChange} />
        </EditorField>
      </EditorFieldGroup>
    </EditorRow>
  );
};
