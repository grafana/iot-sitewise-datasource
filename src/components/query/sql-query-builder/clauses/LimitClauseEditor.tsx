import React from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorFieldGroup } from '@grafana/plugin-ui';
import { StyledLabel } from '../StyledLabel';

interface LimitClauseEditorProps {
  limit?: number;
  updateQuery: (newState: { limit?: number }) => void;
}

/**
 * A numeric input field for setting a LIMIT clause in a query editor.
 * Automatically updates the parent query state when changed.
 */
export const LimitClauseEditor: React.FC<LimitClauseEditorProps> = ({ limit, updateQuery }) => {
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
    <EditorFieldGroup>
      {/* Show the 'HAVING' label */}
      <StyledLabel text="LIMIT" width={15} tooltip />

      {/* Input field for numeric limit value */}
      <EditorField label="" width={30}>
        <Input type="number" min={1} placeholder="Defaults to 100" value={limit ?? ''} onChange={handleChange} />
      </EditorField>
    </EditorFieldGroup>
  );
};
