import React, { PureComponent } from 'react';
import { MetadataInspectorProps, DataFrame } from '@grafana/data';
import { DataSource } from '../SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions, SitewiseCustomMeta } from '../types';
import { Tag } from '@grafana/ui';

export type Props = MetadataInspectorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export class MetaInspector extends PureComponent<Props> {
  state = { index: 0 };

  renderInfo = (frame: DataFrame, idx: number) => {
    const custom = frame.meta?.custom as SitewiseCustomMeta;
    if (!custom) {
      return null;
    }

    return (
      <div key={idx}>
        {custom.resolution && (
          <div>
            <h3>Resolution</h3>
            <Tag name={custom.resolution} colorIndex={1} />
            <br />
            <br />
          </div>
        )}

        {custom.aggregates?.length && (
          <div>
            <h3>Aggregates</h3>
            {custom.aggregates.map((agg) => {
              return (
                <>
                  <Tag name={agg} key={agg} colorIndex={1} /> &nbsp;
                </>
              );
            })}
          </div>
        )}
      </div>
    );
  };

  render() {
    const { data } = this.props;
    if (!data || !data.length) {
      return <div>No Data</div>;
    }
    return (
      <div>
        {data.map((frame, idx) => {
          return this.renderInfo(frame, idx);
        })}
      </div>
    );
  }
}
