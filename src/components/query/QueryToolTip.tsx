import React from 'react';
import { Icon, LinkButton } from '@grafana/ui';
import { QueryTypeInfo } from 'queryInfo';

export const QueryToolTip = ({ description, helpURL }: QueryTypeInfo) => {
  return (
    <div>
      {description} <br />
      <LinkButton href={helpURL} target="_blank">
        API Docs <Icon name="external-link-alt" />
      </LinkButton>
    </div>
  );
};
