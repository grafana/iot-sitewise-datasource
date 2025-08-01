import React from 'react';
import { Select, IconButton, Tooltip } from '@grafana/ui';
import { EditorField, EditorFieldGroup, EditorRow } from '@grafana/plugin-ui';
import { StyledLabel } from '../StyledLabel';
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
    <>
      {orderByFields.map((field, index) => (
        <EditorRow key={index}>
          <EditorFieldGroup>
            {/* Show 'ORDER BY' label */}
            <StyledLabel text={index === 0 ? 'ORDER BY' : ''} width={15} tooltip={index === 0} />

            {/* Column selector dropdown */}
            <EditorField label="" width={30}>
              <Select
                options={availableProperties.map((prop) => ({
                  label: prop.name,
                  value: prop.id,
                }))}
                value={field.column ? { label: field.column, value: field.column } : null}
                onChange={(option) => updateOrderByField(index, { column: option?.value || '' })}
                placeholder="Select column..."
              />
            </EditorField>

            {/* Direction selector (ASC/DESC) - only shown if column is selected */}
            {field.column ? (
              <EditorField label="" width={30}>
                <Select
                  options={[
                    { label: 'Ascending', value: 'ASC' },
                    { label: 'Descending', value: 'DESC' },
                  ]}
                  value={field.direction}
                  onChange={(option) =>
                    updateOrderByField(index, { direction: (option?.value as 'ASC' | 'DESC') || 'ASC' })
                  }
                  placeholder="Direction"
                />
              </EditorField>
            ) : null}

            {/* Add/Remove buttons for each ORDER BY field */}
            <EditorField label="" width={15}>
              <div>
                <Tooltip content="Add ORDER BY field">
                  <IconButton name="plus" onClick={addOrderByField} aria-label="Add order by field" />
                </Tooltip>
                <Tooltip content="Remove ORDER BY field">
                  <IconButton
                    name="minus"
                    onClick={() => removeOrderByField(index)}
                    aria-label="Remove order by field"
                  />
                </Tooltip>
              </div>
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>
      ))}
    </>
  );
};
