import { SitewiseQueryState, ValidationError } from '../types';

export const validateQuery = (state: SitewiseQueryState): ValidationError[] => {
  const errors: ValidationError[] = [];

  // SELECT clause validation
  if (!state.selectFields?.length || !state.selectFields.some((field) => field.column)) {
    errors.push({ type: 'select', error: 'At least one column must be selected in the SELECT clause.' });
  }

  // FROM clause validation
  if (!state.selectedAssetModel) {
    errors.push({ type: 'from', error: 'A source (e.g., asset model or table) must be specified in the FROM clause.' });
  }

  // WHERE clause validation
  if (
    state.whereConditions?.some(
      (cond) => cond.column && (!cond.operator || cond.value === undefined || cond.value === null)
    )
  ) {
    errors.push({
      type: 'where',
      error: 'Each WHERE condition must include both an operator and a value when a column is selected.',
    });
  }

  // LIMIT clause validation
  if (state.limit !== undefined) {
    if (isNaN(state.limit)) {
      errors.push({ type: 'limit', error: 'Limit must be a valid number.' });
    } else if (state.limit <= 0) {
      errors.push({ type: 'limit', error: 'Limit must be greater than 0.' });
    } else if (state.limit > 100000) {
      errors.push({ type: 'limit', error: 'Limit must not exceed 100,000 rows.' });
    }
  }

  // GROUP BY clause validation (optional)
  if (state.groupByTags?.some((tag) => !tag || tag.trim() === '')) {
    errors.push({ type: 'group', error: 'Group by tags must not contain empty values.' });
  }

  return errors;
};
