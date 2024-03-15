import React, { PureComponent } from 'react';

import difference from 'lodash/difference';

import { Select } from '@grafana/ui';

import { Registry, SelectableValue } from '@grafana/data';
import { AggregateType, AssetPropertyInfo } from 'types';

interface Props {
  assetPropInfo?: AssetPropertyInfo;
  onChange: (stats: AggregateType[]) => void;
  stats: AggregateType[];
  allowMultiple?: boolean;
  defaultStat?: AggregateType;
  className?: string;
  menuPlacement?: 'auto' | 'bottom' | 'top';
}

//type AggChecker = (p:AssetPropertyInfo) => boolean;
const AnyTypeOK = (p: AssetPropertyInfo) => true;
const OnlyNumbers = (p: AssetPropertyInfo) => p.DataType !== 'STRING';

export const aggReg = new Registry(() => [
  { id: AggregateType.AVERAGE, name: 'Average', isValid: OnlyNumbers },
  { id: AggregateType.COUNT, name: 'Count', isValid: AnyTypeOK },
  { id: AggregateType.MAXIMUM, name: 'Max', isValid: OnlyNumbers },
  { id: AggregateType.MINIMUM, name: 'Min', isValid: OnlyNumbers },
  { id: AggregateType.SUM, name: 'Sum', isValid: OnlyNumbers },
  { id: AggregateType.STANDARD_DEVIATION, name: 'Stddev', description: 'Standard Deviation', isValid: OnlyNumbers },
]);

export class AggregatePicker extends PureComponent<Props> {
  static defaultProps: Partial<Props> = {
    allowMultiple: true,
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
      const found = current.map((v) => v.id);
      const notFound = difference(stats, found);
      console.warn('Unknown stats', notFound, stats);
      onChange(current.map((stat) => stat.id));
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

  onSelectionChange = (item: SelectableValue<AggregateType>) => {
    const { onChange } = this.props;
    if (Array.isArray(item)) {
      onChange(item.map((v) => v.value));
    } else {
      onChange(item && item.value ? [item.value] : []);
    }
  };

  render() {
    const { stats, allowMultiple, defaultStat, className, menuPlacement, assetPropInfo } = this.props;

    const select = aggReg.selectOptions(stats);
    if (assetPropInfo && assetPropInfo.DataType === 'STRING') {
      select.options = aggReg.list().filter((a) => a.isValid(assetPropInfo));
    }
    return (
      <Select
        inputId="aggregate-picker"
        aria-label="Aggregate picker"
        value={select.current}
        className={className}
        isClearable={!defaultStat}
        isMulti={allowMultiple}
        isSearchable={true}
        options={select.options as any}
        onChange={this.onSelectionChange}
        menuPlacement={menuPlacement}
      />
    );
  }
}
