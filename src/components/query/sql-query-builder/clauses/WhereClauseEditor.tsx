import React, { useMemo } from 'react';
import { Select, IconButton, Tooltip } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { isFunctionOfType, WhereCondition, whereOperators } from '../types';
import { VariableSuggestInput } from '../VariableSuggestInput';
import { StyledLabel } from '../StyledLabel';

interface WhereClauseEditorProps {
  whereConditions: WhereCondition[];
  updateQuery: (updatedFields: Partial<{ whereConditions: WhereCondition[] }>) => void;
  availableProperties: Array<{ id: string; name: string }>;
}

/**
 * Renders an UI editor for managing `WHERE` clause conditions in a SQL-like query builder.
 * Allows users to select columns, operators, and input values (including template variables),
 * with support for logical operators (AND/OR)
 */
export const WhereClauseEditor: React.FC<WhereClauseEditorProps> = ({
  whereConditions,
  updateQuery,
  availableProperties,
}) => {
  /**
   * Memoized list of available property options for the column selector.
   * Converts `{ id, name }` into `{ label, value }` objects for use with the `Select` component.
   */
  const columnOptions = useMemo(
    () => availableProperties.map((prop) => ({ label: prop.name, value: prop.id })),
    [availableProperties]
  );

  /**
   * Supports partial updates to keys like column, operator, value, etc.
   *
   * @param index - Index of the condition being updated
   */
  const handleUpdate = (index: number) => (key: keyof WhereCondition, value: any) => {
    const updated = [...whereConditions];
    updated[index] = { ...updated[index], [key]: value };
    if (key === 'operator') {
      updated[index].operator === 'BETWEEN' ? (updated[index].operator2 = 'AND') : delete updated[index].operator2;
      updated[index].value = '';
      updated[index].value2 = '';
    }
    updateQuery({ whereConditions: updated });
  };

  /**
   * Adds a new empty condition row to the list.
   * Defaults to `column: '', operator: '=', value: '', logicalOperator: 'AND'`
   */
  const addWhereCondition = () => {
    updateQuery({
      whereConditions: [...whereConditions, { column: '', operator: '=', value: '', logicalOperator: 'AND' }],
    });
  };

  /**
   * Removes a condition at the given index.
   * If only one condition exists, it resets it instead of removing.
   *
   * @param index - Index of the condition to remove
   */
  const removeWhereCondition = (index: number) => {
    const updatedConditions =
      whereConditions.length === 1
        ? [{ column: '', operator: '', value: '' }]
        : whereConditions.filter((_, i) => i !== index);

    updateQuery({ whereConditions: updatedConditions });
  };

  return (
    <>
      {whereConditions.map((condition, index) => (
        <EditorRow key={index}>
          <EditorFieldGroup>
            {/* Show the 'WHERE' label */}
            <StyledLabel text={index === 0 ? 'WHERE' : ''} width={15} tooltip={index === 0} />

            {/* Column selector */}
            <EditorField label="" width={30}>
              <Select
                options={columnOptions}
                value={condition.column ? { label: condition.column, value: condition.column } : null}
                onChange={(o) => handleUpdate(index)('column', o?.value || '')}
                placeholder="Select column..."
              />
            </EditorField>

            {/* Operator selector (e.g., =, !=, BETWEEN) */}
            <EditorField label="" width={15}>
              <Select
                options={whereOperators}
                value={condition.operator ? { label: condition.operator, value: condition.operator } : null}
                onChange={(o) => handleUpdate(index)('operator', o?.value || '')}
              />
            </EditorField>

            {/* Value input for function operators except IS NULL/IS NOT NULL */}
            {!isFunctionOfType(condition.operator, 'val') && (
              <>
                <EditorField label="" width={30}>
                  <VariableSuggestInput value={condition.value} onChange={(val) => handleUpdate(index)('value', val)} />
                </EditorField>
              </>
            )}

            {/* BETWEEN operator: adds extra field and static "AND" operator */}
            {condition.operator === 'BETWEEN' && (
              <>
                <EditorField label="" width={10}>
                  <Select options={[{ label: 'AND', value: 'AND' }]} value="AND" onChange={() => {}} disabled />
                </EditorField>

                <EditorField label="" width={30}>
                  <VariableSuggestInput
                    value={condition.value2 || ''}
                    onChange={(val) => handleUpdate(index)('value2', val)}
                  />
                </EditorField>
              </>
            )}

            {/* Logical operator (AND/OR) shown if not the last condition */}
            {index < whereConditions.length - 1 && (
              <EditorField label={''} width={15}>
                <Select
                  options={[
                    { label: 'AND', value: 'AND' },
                    { label: 'OR', value: 'OR' },
                  ]}
                  value={condition.logicalOperator}
                  onChange={(o) => handleUpdate(index)('logicalOperator', o?.value || 'AND')}
                />
              </EditorField>
            )}

            {/* Action buttons: Add/Remove condition */}
            <EditorField label="" width={10}>
              <div>
                {index === whereConditions.length - 1 && (
                  <Tooltip content="Add condition">
                    <IconButton name="plus" onClick={addWhereCondition} aria-label="Add condition" />
                  </Tooltip>
                )}

                <Tooltip content="Remove condition">
                  <IconButton name="minus" onClick={() => removeWhereCondition(index)} aria-label="Remove condition" />
                </Tooltip>
              </div>
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>
      ))}
    </>
  );
};
