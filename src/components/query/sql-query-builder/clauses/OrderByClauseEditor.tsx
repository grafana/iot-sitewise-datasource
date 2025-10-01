import React from 'react';
import { Select, FieldSet, Stack } from '@grafana/ui';
import { AccessoryButton, EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { OrderByField } from '../types';

interface OrderByClauseEditorProps {
  orderByFields: OrderByField[];
  updateQuery: (update: Partial<{ orderByFields: OrderByField[] }>) => void;
  availableProperties: Array<{ id: string; name: string }>;
}

/**
 * Renders an UI editor  that allows users to define one or more ORDER BY clauses
 * for their query by selecting columns and sort directions (ASC/DESC).
 * Provides UI to add or remove sorting fields.
 */
export const OrderByClauseEditor: React.FC<OrderByClauseEditorProps> = ({
  orderByFields,
  updateQuery,
  availableProperties,
}) => {
  /**
   * Adds a new empty ORDER BY field to the query
   */
  const addOrderByField = () => {
    updateQuery({ orderByFields: [...orderByFields, { column: '', direction: 'ASC' }] });
  };

  /**
   * Removes an ORDER BY field from the list
   * If only one field exists, it resets it instead of removing completely
   *
   * @param index - Index of the field to remove
   */
  const removeOrderByField = (index: number) => {
    const newFields =
      orderByFields.length === 1 ? [{ column: '', direction: 'ASC' }] : orderByFields.filter((_, i) => i !== index);
    updateQuery({
      orderByFields: newFields.map((f) => ({
        column: f.column,
        direction: f.direction as 'ASC' | 'DESC',
      })),
    });
  };

  /**
   * Updates a specific ORDER BY field based on index
   *
   * @param index - Index of the field to update
   * @param field - Partial update containing new column or direction
   */
  const updateOrderByField = (index: number, field: Partial<OrderByField>) => {
    const newFields = [...orderByFields];
    newFields[index] = { ...newFields[index], ...field };
    updateQuery({ orderByFields: newFields });
  };

  return (
    <EditorRow>
      <FieldSet label="Order By">
        <Stack gap={3} direction="column">
          {orderByFields.map((field, index) => (
            <EditorFieldGroup key={index}>
              {/* Column selector dropdown */}
              <EditorField label="Column" htmlFor={`order-column-${index}`} width={30}>
                <Select
                  options={availableProperties.map((prop) => ({
                    label: prop.name,
                    value: prop.id,
                  }))}
                  inputId={`order-column-${index}`}
                  value={field.column ? { label: field.column, value: field.column } : null}
                  onChange={(option) => updateOrderByField(index, { column: option?.value || '' })}
                  placeholder="Select column..."
                />
              </EditorField>

              {/* Direction selector (ASC/DESC) - only shown if column is selected */}
              {field.column && (
                <EditorField label="Direction" htmlFor={`order-direction-${index}`} width={30}>
                  <Select
                    options={[
                      { label: 'Ascending', value: 'ASC' },
                      { label: 'Descending', value: 'DESC' },
                    ]}
                    inputId={`order-direction-${index}`}
                    value={
                      field.direction
                        ? { label: field.direction === 'ASC' ? 'Ascending' : 'Descending', value: field.direction }
                        : null
                    }
                    onChange={(option) =>
                      updateOrderByField(index, { direction: (option?.value as 'ASC' | 'DESC') || 'ASC' })
                    }
                    placeholder="Direction"
                  />
                </EditorField>
              )}

              {/* Add/Remove buttons for each ORDER BY field */}
              <Stack gap={1} alignItems="flex-end">
                {index === orderByFields.length - 1 && (
                  <AccessoryButton
                    aria-label="Add Order By field"
                    icon="plus"
                    variant="secondary"
                    onClick={addOrderByField}
                  />
                )}
                {orderByFields.length > 1 && (
                  <AccessoryButton
                    aria-label="Remove ORDER BY field"
                    icon="times"
                    variant="secondary"
                    onClick={() => removeOrderByField(index)}
                  />
                )}
              </Stack>
            </EditorFieldGroup>
          ))}
        </Stack>
      </FieldSet>
    </EditorRow>
  );
};
