import { type MetadataInspectorProps } from '@grafana/data';
import { Tag } from '@grafana/ui';
import React, { type CSSProperties } from 'react';
import { type DataSource } from '../SitewiseDataSource';
import type { SitewiseQuery, SitewiseOptions, SitewiseCustomMetadata } from '../types';

const resolutionContainerStyles: CSSProperties = { marginBottom: '16px' };

function isSiteWiseCustomMetadata(u: SitewiseCustomMetadata | unknown): u is SitewiseCustomMetadata {
  const resolutionKey = 'resolution' satisfies keyof SitewiseCustomMetadata;
  const aggregatesKey = 'aggregates' satisfies keyof SitewiseCustomMetadata;

  return typeof u === 'object' && u != null && resolutionKey in u && aggregatesKey in u;
}

export function MetadataInspector({ data }: MetadataInspectorProps<DataSource, SitewiseQuery, SitewiseOptions>) {
  if (!data || !data.length) {
    return <div>No Data</div>;
  }

  return (
    <div>
      {data.map(({ meta: { custom: siteWiseMetadata = {} } = {} }, idx) => {
        if (!isSiteWiseCustomMetadata(siteWiseMetadata)) {
          return null;
        }

        const { resolution, aggregates = [] } = siteWiseMetadata;

        return (
          <div key={idx}>
            {resolution && (
              <div style={resolutionContainerStyles}>
                <h3>Resolution</h3>
                <Tag name={resolution} colorIndex={1} />
              </div>
            )}

            {aggregates.length > 0 && (
              <div>
                <h3>Aggregates</h3>
                {aggregates.map((agg) => (
                  <>
                    <Tag name={agg} key={agg} colorIndex={1} /> &nbsp;
                  </>
                ))}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
