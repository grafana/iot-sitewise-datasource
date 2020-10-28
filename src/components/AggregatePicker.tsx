import React, { PureComponent } from 'react';

import difference from 'lodash/difference';

import { Select } from '@grafana/ui';

import { Registry, SelectableValue } from '@grafana/data';
import { AggregateTypes } from 'types';

interface Props {
  placeholder?: string;
  onChange: (stats: AggregateTypes[]) => void;
  stats: AggregateTypes[];
  allowMultiple?: boolean;
  defaultStat?: AggregateTypes;
  className?: string;
  menuPlacement?: 'auto' | 'bottom' | 'top';
}

const aggReg = new Registry(() => [
  { id: AggregateTypes.AVERAGE, name: 'Average' },
  { id: AggregateTypes.COUNT, name: 'Count' },
  { id: AggregateTypes.MAXIMUM, name: 'Max' },
  { id: AggregateTypes.MINIMUM, name: 'Min' },
  { id: AggregateTypes.SUM, name: 'Sum' },
  { id: AggregateTypes.STANDARD_DEVIATION, name: 'Stddev' },
]);

export class AggregatePicker extends PureComponent<Props> {
  static defaultProps: Partial<Props> = {
    allowMultiple: true,
    defaultStat: AggregateTypes.AVERAGE,
  };

  componentDidMount() {
    this.checkInput();
  }

  componentDidUpdate(prevProps: Props) {
    this.checkInput();
  }

  checkInput = () => {
    const { stats, allowMultiple, defaultStat, onChange } = this.props;

    const current = aggReg.list(stats);
    if (current.length !== stats.length) {
      const found = current.map(v => v.id);
      const notFound = difference(stats, found);
      console.warn('Unknown stats', notFound, stats);
      onChange(current.map(stat => stat.id));
    }

    // Make sure there is only one
    if (!allowMultiple && stats.length > 1) {
      console.warn('Removing extra stat', stats);
      onChange([stats[0]]);
    }

    // Set the reducer from callback
    if (defaultStat && stats.length < 1) {
      onChange([defaultStat]);
    }
  };

  onSelectionChange = (item: SelectableValue<AggregateTypes>) => {
    const { onChange } = this.props;
    if (Array.isArray(item)) {
      onChange(item.map(v => v.value));
    } else {
      onChange(item && item.value ? [item.value] : []);
    }
  };

  render() {
    const { stats, allowMultiple, defaultStat, placeholder, className, menuPlacement } = this.props;

    const select = aggReg.selectOptions(stats);
    return (
      <Select
        value={select.current}
        className={className}
        isClearable={!defaultStat}
        isMulti={allowMultiple}
        isSearchable={true}
        options={select.options as any}
        placeholder={placeholder}
        onChange={this.onSelectionChange}
        menuPlacement={menuPlacement}
      />
    );
  }
}
