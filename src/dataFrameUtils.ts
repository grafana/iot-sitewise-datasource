import { AbsoluteTimeRange, ArrayVector, DataFrame, FieldType } from '@grafana/data';

interface TrimParams {
  dataFrame: DataFrame;
  timeRange: AbsoluteTimeRange;
  lastObservation?: boolean;
};

/**
 * Trim the time series data frame to the specified time range.
 * @param param0 - The parameters for trimming the data frame.
 * @param param0.dataFrame - The data frame to trim.
 * @param param0.timeRange - The time range to trim to.
 * @param param0.lastObservation - Whether to include the last observation in the range.
 * @returns The trimmed data frame.
 */
export function trimTimeSeriesDataFrame({
  dataFrame,
  timeRange: { from, to },
  lastObservation,
}: TrimParams): DataFrame {
  const { fields } = dataFrame;
  if (fields == null || fields.length === 0) {
    return {
      ...dataFrame,
      fields: [],
      length: 0,
    }
  }

  const timeField = fields.find(field => field.name === 'time' && field.type === FieldType.time);
  if (timeField == null) {
    // return the original data frame if a time field cannot be found
    return dataFrame;
  }

  let timeValues = timeField.values.toArray();

  let fromIndex = timeValues.findIndex(time => time > from);  // from is exclusive
  if (fromIndex === -1) {
    // no time value within range; include no data in the slice
    fromIndex = timeValues.length ;
  } else if (lastObservation) {
    // Keeps 1 extra data point before the range
    fromIndex = Math.max(fromIndex - 1, 0);
  }

  let toIndex = timeValues.findIndex(time => time > to);  // to is inclusive
  if (toIndex === -1) {
    // all time values before `to`
    toIndex = timeValues.length;
  }

  const trimmedFields = fields.map(field => ({
    ...field,
    values: new ArrayVector(field.values.toArray().slice(fromIndex, toIndex)),
  }));
  
  return {
    ...dataFrame,
    fields: trimmedFields,
    length: trimmedFields[0].values.length,
  };
}

/**
 * Trim the time series data frame to the specified time range where the time field is in reversed order.
 * @param param0 - The parameters for trimming the data frame.
 * @param param0.dataFrame - The data frame to trim.
 * @param param0.timeRange - The time range to trim to.
 * @param param0.lastObservation - Whether to include the last observation in the range.
 * @returns The trimmed data frame.
 */
export function trimTimeSeriesDataFrameReversedTime({
  dataFrame,
  timeRange: { from, to },
  lastObservation,
}: TrimParams): DataFrame {
  const { fields } = dataFrame;
  if (fields == null || fields.length === 0) {
    return {
      ...dataFrame,
      fields: [],
      length: 0,
    }
  }

  const timeField = fields.find(field => field.name === 'time' && field.type === FieldType.time);
  if (timeField == null) {
    // return the original data frame if a time field cannot be found
    return dataFrame;
  }

  // Copy before reverse in place
  let timeValues = [...timeField.values.toArray()].reverse();
  
  let fromIndex = timeValues.findIndex(time => time > from);  // from is exclusive
  if (fromIndex === -1) {
    // no time value within range; include no data in the slice
    fromIndex = timeValues.length ;
  } else if (lastObservation) {
    // Keeps 1 extra data point before the range
    fromIndex = Math.max(fromIndex - 1, 0);
  }

  let toIndex = timeValues.findIndex(time => time > to);  // to is inclusive
  if (toIndex === -1) {
    // all time values before `to`
    toIndex = timeValues.length;
  }

  const trimmedFields = fields.map(field => {
    const dataValues = [...field.values.toArray()].reverse().slice(fromIndex, toIndex);

    return {
      ...field,
      values: new ArrayVector(dataValues.reverse()),
    };
  });
  
  return {
    ...dataFrame,
    fields: trimmedFields,
    length: trimmedFields[0].values.length,
  };
}
