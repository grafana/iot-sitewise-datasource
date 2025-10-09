import React, { useMemo } from 'react';
import { Select, FieldSet, Stack, Text, Alert, Box } from '@grafana/ui';
import { EditorRow, EditorField, EditorFieldGroup, AccessoryButton } from '@grafana/plugin-ui';
import { isFunctionOfType, ValidationError, WhereCondition, whereOperators } from '../types';
import { getSelectableTemplateVariables } from 'variables';

interface WhereClauseEditorProps {
  whereConditions: WhereCondition[];
  validationErrors: ValidationError[];
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
  validationErrors,
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
   * Retrieves all available Grafana template variables using getSelectableTemplateVariables(),
     then maps them into { label, value } objects suitable for Select components.
   */
  const variableOptions = useMemo(() => {
    return getSelectableTemplateVariables().map(({ value }) => {
      return {
        label: value,
        value: value,
      };
    });
  }, []);

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
    <EditorRow>
      <FieldSet label="Where">
        <Stack gap={3} direction="column">
          {whereConditions.map((condition, index) => (
            <EditorFieldGroup key={index}>
              {/* Column selector */}
              <EditorField label="Column" htmlFor={`where-column-${index}`} width={30}>
                <Select
                  inputId={`where-column-${index}`}
                  options={columnOptions}
                  placeholder="Select column..."
                  value={condition.column ? { label: condition.column, value: condition.column } : null}
                  onChange={(o) => handleUpdate(index)('column', o?.value || '')}
                />
              </EditorField>
              {/* Operator selector (e.g., =, !=, BETWEEN) */}
              <EditorField label="Operator" htmlFor={`where-operator-${index}`} width={15}>
                <Select
                  inputId={`where-operator-${index}`}
                  options={whereOperators}
                  value={condition.operator ? { label: condition.operator, value: condition.operator } : null}
                  onChange={(o) => handleUpdate(index)('operator', o?.value || '')}
                />
              </EditorField>
              {/* Value input for function operators except IS NULL/IS NOT NULL */}
              {!isFunctionOfType(condition.operator, 'val') && (
                <EditorField label="Value" htmlFor={`where-value-${index}`} width={30}>
                  <Select
                    inputId={`where-value-${index}`}
                    placeholder="Enter value or $variable"
                    options={variableOptions}
                    value={condition.value ? { label: condition.value, value: condition.value } : null}
                    allowCustomValue
                    onChange={(o) => handleUpdate(index)('value', o?.value || '')}
                    isClearable
                  />
                </EditorField>
              )}
              {/* BETWEEN operator: adds extra field and static "AND" operator */}
              {condition.operator === 'BETWEEN' && (
                <>
                  <EditorField label="AND" htmlFor={`where-between-and-${index}`} width={10}>
                    <Select
                      inputId={`where-between-and-${index}`}
                      options={[{ label: 'AND', value: 'AND' }]}
                      value={{ label: 'AND', value: 'AND' }}
                      onChange={() => {}}
                      disabled
                    />
                  </EditorField>

                  <EditorField label="Value 2" htmlFor={`where-value2-${index}`} width={30}>
                    <Select
                      inputId={`where-value2-${index}`}
                      placeholder="Enter value or $variable"
                      options={variableOptions}
                      value={condition.value2 ? { label: condition.value2, value: condition.value2 } : null}
                      allowCustomValue
                      onChange={(o) => handleUpdate(index)('value2', o?.value || '')}
                      isClearable
                    />
                  </EditorField>
                </>
              )}
              {/* Logical operator (AND/OR) shown if not the last condition */}
              {index < whereConditions.length - 1 && (
                <EditorField label="Logical" htmlFor={`where-logical-${index}`} width={15}>
                  <Select
                    inputId={`where-logical-${index}`}
                    options={[
                      { label: 'AND', value: 'AND' },
                      { label: 'OR', value: 'OR' },
                    ]}
                    value={
                      condition.logicalOperator
                        ? { label: condition.logicalOperator, value: condition.logicalOperator }
                        : { label: 'AND', value: 'AND' }
                    }
                    onChange={(o) => handleUpdate(index)('logicalOperator', o?.value || 'AND')}
                  />
                </EditorField>
              )}
              {/* Action buttons: Add/Remove condition */}
              <Stack gap={1} alignItems="flex-end">
                {index === whereConditions.length - 1 && (
                  <AccessoryButton
                    aria-label="Add condition"
                    icon="plus"
                    variant="secondary"
                    onClick={addWhereCondition}
                  />
                )}
                <AccessoryButton
                  aria-label="Remove condition"
                  icon="times"
                  variant="secondary"
                  onClick={() => removeWhereCondition(index)}
                />
              </Stack>
            </EditorFieldGroup>
          ))}
        </Stack>
        {validationErrors?.length > 0 &&
          validationErrors.map(
            (err, idx) =>
              err.type === 'where' && (
                <Box marginTop={1} key={idx}>
                  <Alert title="" severity="error">
                    <Text variant="code" color="error">
                      {err.error}
                    </Text>
                  </Alert>
                </Box>
              )
          )}
      </FieldSet>
    </EditorRow>
  );
};
