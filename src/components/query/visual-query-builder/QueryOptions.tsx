import { SitewiseQuery, shouldShowL4eOptions, shouldShowLastObserved, shouldShowQualityAndOrderComponent } from 'types';
import { CollapsableSection, Switch, useTheme2 } from '@grafana/ui';
import React from 'react';
import { EditorField, EditorFieldGroup } from '@grafana/plugin-ui';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { QualityAndOrderRow } from './QualityAndOrderRow';
import { DataSource } from 'SitewiseDataSource';

interface Props {
  query: SitewiseQuery;
  showProp: boolean;
  showQuality: boolean;
  onLastObservationChange: (e?: React.FormEvent<HTMLInputElement>) => void;
  onFlattenL4eChange: (e?: React.FormEvent<HTMLInputElement>) => void;
  datasource: DataSource;
  onChange: (value: SitewiseQuery) => void;
}

export function QueryOptions({
  query,
  datasource,
  onChange,
  showProp,
  showQuality,
  onLastObservationChange,
  onFlattenL4eChange,
}: Props) {
  const theme = useTheme2();
  const style = getStyles(theme);

  return (
    <div className={style.collapseRow}>
      <CollapsableSection
        className={style.collapse}
        label={
          <span data-testid="collapse-title" className={style.collapseLabel}>
            Query options
          </span>
        }
        isOpen={false}
      >
        <EditorFieldGroup>
          {shouldShowLastObserved(query.queryType) && !query.propertyAliases?.length && showProp && (
            <EditorField
              label="Expand Time Range"
              htmlFor="expand"
              tooltip="Expand query to include last observed value before the current time range, and next observed value after the time range. "
            >
              <Switch value={query.lastObservation} onChange={onLastObservationChange} />
            </EditorField>
          )}
          {shouldShowL4eOptions(query.queryType) && !query.propertyAliases?.length && showProp && (
            <EditorField
              label="Format L4E Anomaly Result"
              htmlFor="l4e"
              tooltip="Format query to parse L4E anomaly result."
            >
              <Switch value={query.flattenL4e} onChange={onFlattenL4eChange} />
            </EditorField>
          )}
          {shouldShowQualityAndOrderComponent(query.queryType) && (showProp || query.propertyAliases?.length) && showQuality && (
            <QualityAndOrderRow query={query} onChange={onChange} datasource={datasource} />
          )}
        </EditorFieldGroup>
      </CollapsableSection>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => ({
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
  collapseLabel: css({
    fontSize: theme.typography.body.fontSize,
  }),
});
