import React from 'react';
import { Select, IconButton, Tooltip, Input } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { HavingCondition } from '../types';
import { StyledLabel } from '../StyledLabel';

interface HavingClauseEditorProps {
  havingConditions: HavingCondition[];
  updateQuery: (updatedFields: Partial<{ havingConditions: HavingCondition[] }>) => void;
  availableProperties: Array<{ id: string; name: string }>;
}

const aggregationOptions = ['COUNT', 'SUM', 'AVG', 'MAX', 'MIN'].map((val) => ({ label: val, value: val }));

/**
 * UI editor for building SQL-style `HAVING` clauses.
 * Supports aggregation selection, column targeting, logical operators, and variable-based value inputs.
 */
export const HavingClauseEditor: React.FC<HavingClauseEditorProps> = ({
  havingConditions,
  updateQuery,
  availableProperties,
}) => {
  const columnOptions = availableProperties.map((prop) => ({ label: prop.name, value: prop.id }));
  const operatorOptions = ['=', '!=', '>', '<', '>=', '<='].map((op) => ({ label: op, value: op }));

  /**
   * Updates a specific field in one of the `havingConditions`.
   *
   * @param index - Index of the condition being edited
   * @param key - Key in the condition object to update
   * @param value - New value for the key
   */
  const updateCondition = (index: number, key: keyof HavingCondition, value: any) => {
    const updated = [...havingConditions];
    updated[index] = { ...updated[index], [key]: value };
    updateQuery({ havingConditions: updated });
  };

  /**
   * Adds a new HAVING condition row with default values.
   * Appends to the current `havingConditions` array.
   */
  const addCondition = () => {
    updateQuery({
      havingConditions: [
        ...havingConditions,
        { aggregation: 'SUM', column: '', operator: '>', value: '', logicalOperator: 'AND' },
      ],
    });
  };

  /**
   * Removes a condition at the specified index.
   * If there's only one condition left, it resets it instead of removing.
   *
   * @param index - Index of the condition to remove
   */
  const removeCondition = (index: number) => {
    const updatedConditions =
      havingConditions.length === 1
        ? [{ aggregation: 'COUNT', column: '', operator: '=', value: '' }]
        : havingConditions.filter((_, i) => i !== index);
    updateQuery({
      havingConditions: updatedConditions as HavingCondition[],
    });
  };

  return (
    <>
      {havingConditions.map((cond, index) => (
        <EditorRow key={index}>
          <EditorFieldGroup>
            {/* Show the 'HAVING' label */}
            <StyledLabel text={index === 0 ? 'HAVING' : ''} width={15} tooltip={index === 0} />

            {/* Aggregation function dropdown */}
            <EditorField label="" width={10}>
              <Select
                options={aggregationOptions}
                value={{ label: cond.aggregation, value: cond.aggregation }}
                onChange={(o) => updateCondition(index, 'aggregation', o?.value)}
              />
            </EditorField>

            {/* Column selection dropdown */}
            <EditorField label="" width={25}>
              <Select
                options={columnOptions}
                value={cond.column ? { label: cond.column, value: cond.column } : null}
                onChange={(o) => updateCondition(index, 'column', o?.value || '')}
                placeholder="Select column..."
              />
            </EditorField>

            {/* Operator dropdown (e.g., =, !=, >, <) */}
            <EditorField label="" width={5}>
              <Select
                options={operatorOptions}
                value={{ label: cond.operator, value: cond.operator }}
                onChange={(o) => updateCondition(index, 'operator', o?.value)}
              />
            </EditorField>

            {/* Value input */}
            <EditorField label="" width={25}>
              <Input
                value={cond.value}
                onChange={(e) => updateCondition(index, 'value', e.currentTarget.value)}
                placeholder="Enter value"
              />
            </EditorField>

            {/* Logical operator (AND/OR) dropdown shown for all but last condition */}
            {index < havingConditions.length - 1 && (
              <EditorField label="" width={10}>
                <Select
                  options={[
                    { label: 'AND', value: 'AND' },
                    { label: 'OR', value: 'OR' },
                  ]}
                  value={{ label: cond.logicalOperator || 'AND', value: cond.logicalOperator || 'AND' }}
                  onChange={(o) => updateCondition(index, 'logicalOperator', o?.value)}
                />
              </EditorField>
            )}

            {/* Action buttons to add or remove condition */}
            <EditorField label="" width={10}>
              <div>
                {index === havingConditions.length - 1 && (
                  <Tooltip content="Add condition">
                    <IconButton name="plus" onClick={addCondition} aria-label="Add condition" />
                  </Tooltip>
                )}
                <Tooltip content="Remove condition">
                  <IconButton name="minus" onClick={() => removeCondition(index)} aria-label="Remove condition" />
                </Tooltip>
              </div>
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>
      ))}
    </>
  );
};
