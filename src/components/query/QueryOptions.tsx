import { SitewiseQuery, shouldShowLastObserved } from 'types';
import { CollapsableSection, Switch, Text } from '@grafana/ui';
import React from 'react';
import { EditorField, EditorFieldGroup } from '@grafana/experimental';
import { css } from '@emotion/css';

interface Props {
  qualityAndOrderComponent: JSX.Element;
  query: SitewiseQuery;
  showProp: boolean;
  showQuality: boolean;
  onLastObservationChange: (e?: React.FormEvent<HTMLInputElement>) => void;
}

export function QueryOptions({
  query,
  showProp,
  showQuality,
  qualityAndOrderComponent,
  onLastObservationChange,
}: Props) {
  return (
    <div className={styles.collapseRow}>
      <CollapsableSection
        className={styles.collapse}
        label={
          <Text variant="body" data-testid="collapse-title">
            Query options
          </Text>
        }
        isOpen={false}
      >
        <EditorFieldGroup>
          {shouldShowLastObserved(query.queryType) && !Boolean(query.propertyAlias) && showProp && (
            <EditorField
              label="Expand Time Range"
              htmlFor="expand"
              tooltip="Expand query to include last observed value before the current time range, and next observed value after the time range. "
            >
              <Switch value={query.lastObservation} onChange={onLastObservationChange} />
            </EditorField>
          )}
          {(showProp || query.propertyAlias) && showQuality && qualityAndOrderComponent}
        </EditorFieldGroup>
      </CollapsableSection>
    </div>
  );
}

const styles = {
  collapse: css({
    alignItems: 'flex-start',
    paddingTop: 0,
  }),
  collapseRow: css({
    display: 'flex',
    flexDirection: 'column',
    '>div': {
      alignItems: 'baseline',
      justifyContent: 'flex-end',
    },
    '*[id^="collapse-content-"]': {
      padding: 'unset',
    },
  }),
};
