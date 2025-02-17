import { MacroType } from '@grafana/plugin-ui';

const COLUMN = 'column',
  TIMESTAMP_FORMAT = "'yyyy-MM-dd HH:mm:ss'";

export const TABLE_MACRO = '$__table';

export const MACROS = [
  {
    id: '$__selectAll',
    name: '$__selectAll',
    text: '$__selectAll',
    args: [],
    type: MacroType.Column,
    description: 'Will be replaced by all the fields of the current table',
  },
  {
    id: '$__rawTimeFrom',
    name: '$__rawTimeFrom',
    text: '$__rawTimeFrom',
    args: [],
    type: MacroType.Filter,
    description: 'Will return the current starting time of the time range',
  },
  {
    id: '$__rawTimeTo',
    name: '$__rawTimeTo',
    text: '$__rawTimeTo',
    args: [],
    type: MacroType.Filter,
    description: 'Will return the current ending time of the time range',
  },
  {
    id: '$__unixEpochFilter()',
    name: '$__unixEpochFilter()',
    text: '$__unixEpochFilter',
    args: [COLUMN],
    type: MacroType.Filter,
    description:
      'Will be replaced by a time range filter using the specified column name with times represented as Unix timestamp. For example, column >= 1624406400 AND column <= 1624410000',
  },
  {
    id: '$__resolution()',
    name: '$__resolution()',
    text: '$__resolution',
    args: [],
    type: MacroType.Value,
    description:
      'Will be replaced by an appropriate resolution (1m, 15m, 1h, 1d) based on the panel interval to be used on precomputed_aggregates queries.',
  },
  {
    id: '$__column',
    name: '$__column',
    text: '$__column',
    args: [],
    type: MacroType.Column,
    description: 'Will be replaced by the query column.',
  },
  {
    id: TABLE_MACRO,
    name: TABLE_MACRO,
    text: TABLE_MACRO,
    args: [],
    type: MacroType.Table,
    description: 'Will be replaced by the query table.',
  },
  {
    id: '$__parseTime',
    name: '$__parseTime',
    text: '$__parseTime',
    args: [COLUMN, TIMESTAMP_FORMAT],
    type: MacroType.Value,
    description: 'Will cast a varchar as a timestamp using the provided format.',
  },
];
