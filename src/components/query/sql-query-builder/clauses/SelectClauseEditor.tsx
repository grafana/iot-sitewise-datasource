import React, { useEffect, useMemo, useRef } from 'react';
import { Select, Input, Cascader, FieldSet, Stack, Text, Alert, Box } from '@grafana/ui';
import { AccessoryButton, EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { allFunctions, FUNCTION_ARGS, isFunctionOfType, SelectField, ValidationError } from '../types';

interface SelectClauseEditorProps {
  selectFields: SelectField[];
  validationErrors: ValidationError[];
  updateQuery: (updatedFields: Partial<{ selectFields: SelectField[] }>) => void;
  availableProperties: Array<{ id: string; name: string }>;
}

/**
 * Renders a UI editor for managing `SELECT` clause fields in a SQL-like query builder.
 * Provides controls for choosing columns, aggregation functions, function arguments,
 * and optional aliases. Supports adding/removing multiple select fields dynamically.
 */
export const SelectClauseEditor: React.FC<SelectClauseEditorProps> = ({
  selectFields,
  validationErrors,
  updateQuery,
  availableProperties,
}) => {
  // Tracks version of the select field list to ensure React re-renders components correctly
  const listVersionRef = useRef(0);
  const prevLengthRef = useRef(selectFields.length);
  useEffect(() => {
    if (selectFields.length !== prevLengthRef.current) {
      listVersionRef.current += 1;
      prevLengthRef.current = selectFields.length;
    }
  }, [selectFields.length]);

  /**
   * Memoized list of available properties for the "Column" dropdown.
   * Converts `{ id, name }` into `{ label, value }` for use with the `Select` component.
   */
  const columnOptions = useMemo(
    () => availableProperties.map((prop) => ({ label: prop.name, value: prop.id })),
    [availableProperties]
  );

  /**
   * Adds a new empty select field at the end of the list.
   * Defaults: column: '', aggregation: 'Raw Values', alias: ''.
   */
  const addSelectField = () => {
    const newFields = [...selectFields, { column: '', aggregation: 'Raw Values', alias: '' }];
    updateQuery({ selectFields: newFields });
  };

  /**
   * Removes a select field by index.
   * If only one field remains, it prevents removal.
   *
   * @param index - Index of the select field to remove
   */
  const removeSelectField = (index: number) => {
    if (selectFields.length > 1) {
      const newFields = selectFields.filter((_, i) => i !== index).map((field) => ({ ...field }));
      updateQuery({ selectFields: newFields });
    }
  };

  /**
   * Updates a specific field in the selectFields array by index.
   * Supports partial updates (e.g., only updating `alias` or `aggregation`).
   *
   * @param index - Index of the select field to update
   * @param field - Partial update object (column, aggregation, alias, etc.)
   */
  const updateSelectField = (index: number, field: Partial<SelectField>) => {
    const newFields = [...selectFields];
    newFields[index] = { ...newFields[index], ...field };
    updateQuery({ selectFields: newFields });
  };

  const shouldShowInput1 = (agg: string) => isFunctionOfType(agg, 'date', 'math', 'str', 'coalesce');

  const shouldShowInput2 = (agg: string) => isFunctionOfType(agg, 'str');

  /**
   * Returns a list of valid function arguments depending on aggregation type.
   * - DATE functions → DATE args
   * - CAST functions → CAST types
   * - CONCAT functions → Column options
   * - Default → empty list
   */
  const getFunctionArgs = (agg: string) => {
    if (isFunctionOfType(agg, 'date')) {
      return FUNCTION_ARGS.DATE.map((arg) => ({ label: arg, value: arg }));
    }
    if (isFunctionOfType(agg, 'cast')) {
      return FUNCTION_ARGS.CAST.map((arg) => ({ label: arg, value: arg }));
    }
    if (isFunctionOfType(agg, 'concat')) {
      return columnOptions;
    }
    return [];
  };

  return (
    <EditorRow>
      <FieldSet label="Select">
        <Stack gap={3} direction="column">
          {selectFields.map((field, index) => {
            const functionArgs = getFunctionArgs(field.aggregation || '');
            const showInput1 = shouldShowInput1(field.aggregation || '');
            const showInput2 = shouldShowInput2(field.aggregation || '');
            const uniqueKey = `${listVersionRef.current}-${index}`;

            return (
              <EditorFieldGroup key={uniqueKey}>
                {/* Column selector */}
                <EditorField label="Column" htmlFor={`column-${index}`} width={30}>
                  <Select
                    options={columnOptions}
                    inputId={`column-${index}`}
                    value={field.column ? { label: field.column, value: field.column } : null}
                    onChange={(option) => updateSelectField(index, { column: option?.value || '' })}
                    placeholder="Select column..."
                  />
                </EditorField>

                {/* Aggregation function selector */}
                <EditorField label="Aggregation" htmlFor={`aggregation-${index}`} width={30}>
                  <Cascader
                    key={uniqueKey}
                    options={allFunctions}
                    id={`aggregation-${index}`}
                    initialValue={field.aggregation || 'Raw Values'}
                    onSelect={(val: string) => {
                      updateSelectField(index, {
                        aggregation: val,
                        functionArg: '',
                        functionArgValue: '',
                        functionArgValue2: '',
                      });
                    }}
                  />
                </EditorField>

                {/* Function arguments (arg type or column for concat) */}
                {functionArgs.length > 0 && (
                  <EditorField
                    label={isFunctionOfType(field.aggregation, 'concat') ? 'Select column' : 'Arg type'}
                    htmlFor={`function-arg-${index}`}
                    width={20}
                  >
                    <Select
                      inputId={`function-arg-${index}`}
                      options={functionArgs}
                      value={field.functionArg ? { label: field.functionArg, value: field.functionArg } : null}
                      onChange={(v) =>
                        updateSelectField(index, {
                          functionArg: (v as any)?.value || '',
                        })
                      }
                    ></Select>
                  </EditorField>
                )}

                {/* Arg value input 1 */}
                {showInput1 && (
                  <EditorField label="Arg value 1" htmlFor={`function-arg-value-1-${index}`} width={20}>
                    <Input
                      id={`function-arg-value-1-${index}`}
                      value={field.functionArgValue || ''}
                      onChange={(e) => updateSelectField(index, { functionArgValue: e.currentTarget.value })}
                      placeholder="Enter value"
                    />
                  </EditorField>
                )}

                {/* Arg value input 2 */}
                {showInput2 && (
                  <EditorField label="Arg value 2" htmlFor={`function-arg-value-2-${index}`} width={20}>
                    <Input
                      id={`function-arg-value-2-${index}`}
                      value={field.functionArgValue2 || ''}
                      onChange={(e) => updateSelectField(index, { functionArgValue2: e.currentTarget.value })}
                      placeholder="Enter value"
                    />
                  </EditorField>
                )}

                {/* Alias input */}
                <EditorField label="Optional alias" htmlFor={`alias-${index}`} width={30}>
                  <Input
                    id={`alias-${index}`}
                    value={field.alias}
                    onChange={(e) => updateSelectField(index, { alias: e.currentTarget.value })}
                    placeholder="Optional alias"
                  />
                </EditorField>

                {/* Action buttons: add/remove select field */}
                <Stack gap={1} alignItems="flex-end">
                  {index === selectFields.length - 1 && (
                    <AccessoryButton aria-label="Add field" icon="plus" variant="secondary" onClick={addSelectField} />
                  )}
                  {selectFields.length > 1 && (
                    <AccessoryButton
                      aria-label="Remove field"
                      icon="times"
                      variant="secondary"
                      onClick={() => removeSelectField(index)}
                    />
                  )}
                </Stack>
              </EditorFieldGroup>
            );
          })}
        </Stack>
        {validationErrors?.length > 0 &&
          validationErrors.map(
            (err, idx) =>
              err.type === 'select' && (
                <Box marginTop={1} key={idx}>
                  <Alert title="" severity="error">
                    <Text variant="code">{err.error}</Text>
                  </Alert>
                </Box>
              )
          )}
      </FieldSet>
    </EditorRow>
  );
};
