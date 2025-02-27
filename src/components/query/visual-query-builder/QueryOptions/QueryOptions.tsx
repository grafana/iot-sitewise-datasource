import {
  SitewiseQuery,
  isAssetPropertyAggregatesQuery,
  isAssetPropertyValueHistoryQuery,
  shouldShowL4eOptions,
  shouldShowLastObserved,
} from 'types';
import { CollapsableSection, Switch, useTheme2 } from '@grafana/ui';
import React from 'react';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { ResponseFormatSettings } from './ResponseFormatSettings';
import { QualitySettings } from './QualitySettings';
import { TimeOrderSettings } from './TimeOrderSettings';

interface Props {
  query: SitewiseQuery;
  onChange: (value: SitewiseQuery) => void;
}

export function QueryOptions({ query, onChange }: Props) {
  const theme = useTheme2();
  const style = getStyles(theme);

  const onLastObservationChange = () => {
    onChange({ ...query, lastObservation: !query.lastObservation });
  };

  const onFlattenL4eChange = () => {
    onChange({ ...query, flattenL4e: !query.flattenL4e });
  };

  return (
    <EditorRow>
      <EditorFieldGroup>
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
              {shouldShowLastObserved(query.queryType) && (
                <EditorField
                  label="Expand Time Range"
                  htmlFor="expand"
                  tooltip="Expand query to include last observed value before the current time range, and next observed value after the time range. "
                >
                  <Switch value={query.lastObservation} onChange={onLastObservationChange} />
                </EditorField>
              )}
              {shouldShowL4eOptions(query.queryType) && (
                <EditorField
                  label="Format L4E Anomaly Result"
                  htmlFor="l4e"
                  tooltip="Format query to parse L4E anomaly result."
                >
                  <Switch value={query.flattenL4e} onChange={onFlattenL4eChange} />
                </EditorField>
              )}

              <QualitySettings query={query} onChange={onChange} />
              {(isAssetPropertyAggregatesQuery(query) || isAssetPropertyValueHistoryQuery(query)) && (
                <TimeOrderSettings query={query} onChange={onChange} />
              )}
              <ResponseFormatSettings query={query} onChange={onChange} />
            </EditorFieldGroup>
          </CollapsableSection>
        </div>
      </EditorFieldGroup>
    </EditorRow>
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
