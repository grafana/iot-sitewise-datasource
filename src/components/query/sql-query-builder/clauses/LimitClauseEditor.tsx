import React from 'react';
import { FieldSet, Input, Stack } from '@grafana/ui';
import { EditorField } from '@grafana/plugin-ui';
import { ValidationError } from '../types';

interface LimitClauseEditorProps {
  limit?: number;
  validationErrors: ValidationError[];
  updateQuery: (newState: { limit?: number }) => void;
}

/**
 * A numeric input field for setting a LIMIT clause in a query editor.
 * Automatically updates the parent query state when changed.
 */
export const LimitClauseEditor: React.FC<LimitClauseEditorProps> = ({ limit, validationErrors, updateQuery }) => {
  /**
   * Handles input value changes in the limit field and updates to the query state.
   *
   * @param e - Input change event
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value.trim();
    const parsed = parseInt(value, 10);

    if (value === '') {
      updateQuery({ limit: undefined });
    } else if (!isNaN(parsed)) {
      updateQuery({ limit: parsed });
    }
  };

  return (
    <FieldSet label="Limit" style={{ marginBottom: 0 }}>
      <Stack direction="row" gap={4} alignItems="center">
        <EditorField label="Max rows" htmlFor="limit-input" width={30}>
          {/* Input field for numeric limit value */}
          <Input
            id="limit-input"
            type="number"
            min={1}
            placeholder="Defaults to 100"
            value={limit ?? ''}
            onChange={handleChange}
          />
        </EditorField>
      </Stack>
      {validationErrors?.length > 0 &&
        validationErrors.map(
          (err, idx) =>
            err.type === 'limit' && (
              <div key={idx} className="text-error text-sm">
                <div>{err.error}</div>
              </div>
            )
        )}
    </FieldSet>
  );
};
