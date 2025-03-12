import React from 'react';
import { Switch } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';

interface Props {
  clientCache?: boolean;
  onClientCacheChange: (evt: React.SyntheticEvent<HTMLInputElement>) => void;
  queryRefId: string;
}

export const ClientCacheRow = ({ clientCache, onClientCacheChange, queryRefId }: Props) => {
  const cacheSwitchId = `client-cache-switch-${queryRefId}`;

  return (
    <EditorRow>
      <EditorFieldGroup>
        <EditorField
          label="Client cache"
          htmlFor={cacheSwitchId}
          tooltip="Enable to cache results in the browser that are older than 15 minutes. Note: Dashboard variable query result will not update when client cache is enabled."
        >
          <Switch id={cacheSwitchId} value={clientCache} onChange={onClientCacheChange} />
        </EditorField>
      </EditorFieldGroup>
    </EditorRow>
  );
};
