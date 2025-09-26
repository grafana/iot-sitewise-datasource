import React, { useMemo } from 'react';
import { Select, Input, Cascader, FieldSet, Stack } from '@grafana/ui';
import { AccessoryButton, EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { allFunctions, FUNCTION_ARGS, isFunctionOfType, SelectField } from '../types';

interface SelectClauseEditorProps {
  selectFields: SelectField[];
  updateQuery: (updatedFields: Partial<{ selectFields: SelectField[] }>) => void;
  availableProperties: Array<{ id: string; name: string }>;
}

export const SelectClauseEditor: React.FC<SelectClauseEditorProps> = ({
  selectFields,
  updateQuery,
  availableProperties,
}) => {
  const columnOptions = useMemo(
    () => availableProperties.map((prop) => ({ label: prop.name, value: prop.id })),
    [availableProperties]
  );

  const addSelectField = () => {
    const newFields = [...selectFields, { column: '', aggregation: '', alias: '' }];
    updateQuery({ selectFields: newFields });
  };

  const removeSelectField = (index: number) => {
    if (selectFields.length > 1) {
      const newFields = selectFields.filter((_, i) => i !== index);
      updateQuery({ selectFields: newFields });
    }
  };

  const updateSelectField = (index: number, field: Partial<SelectField>) => {
    const newFields = [...selectFields];
    newFields[index] = { ...newFields[index], ...field };
    updateQuery({ selectFields: newFields });
  };

  const shouldShowInput1 = (agg: string) => isFunctionOfType(agg, 'date', 'math', 'str', 'coalesce');

  const shouldShowInput2 = (agg: string) => isFunctionOfType(agg, 'str');

  const getFunctionArgs = (agg: string): Array<{ label: string; value: string }> => {
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

            return (
              <EditorFieldGroup key={index}>
                <EditorField label="Column" htmlFor={`column-${index}`} width={30}>
                  <Select
                    options={columnOptions}
                    inputId={`column-${index}`}
                    value={field.column ? { label: field.column, value: field.column } : null}
                    onChange={(option) => updateSelectField(index, { column: option?.value || '' })}
                    placeholder="Select column..."
                  />
                </EditorField>
                <EditorField label="Aggregation" width={30}>
                  <Cascader
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
                    placeholder="No function"
                  />
                </EditorField>
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
                      placeholder={isFunctionOfType(field.aggregation, 'concat') ? 'Select column' : 'Arg type'}
                    />
                  </EditorField>
                )}
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
                <EditorField label="Alias" htmlFor={`alias-${index}`} width={30}>
                  <Input
                    id={`alias-${index}`}
                    value={field.alias}
                    onChange={(e) => updateSelectField(index, { alias: e.currentTarget.value })}
                    placeholder="Alias"
                  />
                </EditorField>
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
      </FieldSet>
    </EditorRow>
  );
};
