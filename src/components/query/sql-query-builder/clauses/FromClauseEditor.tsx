import React from 'react';
import { Select, FieldSet, Stack } from '@grafana/ui';
import { ValidationError } from '../types';
import { EditorField } from '@grafana/plugin-ui';
interface FromClauseEditorProps {
  queryReferenceViews: Array<{ id: string; name: string }>;
  selectedModelId: string;
  validationErrors: ValidationError[];
  updateQuery: (
    updatedFields: Partial<{
      selectedAssetModel: string;
      selectFields: Array<{ column: string; aggregation: string; alias: string }>;
      whereConditions: Array<{ column: string; operator: string; value: string; logicalOperator: 'AND' | 'OR' }>;
      groupByFields?: Array<{ column: string }>;
      orderByFields: Array<{ column: string; direction: 'ASC' | 'DESC' }>;
    }>
  ) => void;
}

/**
 * Renders the "FROM" clause UI for the query builder.
 * It allows users to select an asset model (data source view) from a dropdown list.
 * Once model changed from the dropdown it resets the previous query fields (SELECT, WHERE, GROUP BY, ORDER BY)
 * to their default/initial state.
 */
export const FromClauseEditor: React.FC<FromClauseEditorProps> = ({
  queryReferenceViews,
  selectedModelId,
  validationErrors,
  updateQuery,
}) => {
  return (
    <FieldSet label="From" style={{ marginBottom: 0 }}>
      <Stack direction="row" gap={4} alignItems="center">
        <EditorField label="View" width={40}>
          {/* Dropdown to select a model */}
          <Select
            options={queryReferenceViews.map((model) => ({
              label: model.name,
              value: model.id,
            }))}
            value={
              selectedModelId
                ? {
                    label: queryReferenceViews.find((m) => m.id === selectedModelId)?.name || '',
                    value: selectedModelId,
                  }
                : null
            }
            onChange={(option) =>
              updateQuery({
                selectedAssetModel: option?.value || '',
                selectFields: [{ column: '', aggregation: '', alias: '' }],
                whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
                groupByFields: [{ column: '' }],
                orderByFields: [{ column: '', direction: 'ASC' }],
              })
            }
            placeholder="Select view..."
          />
        </EditorField>
      </Stack>
      {validationErrors?.length > 0 &&
        validationErrors.map(
          (err, idx) =>
            err.type === 'from' && (
              <div key={idx} className="text-error text-sm">
                <div>{err.error}</div>
              </div>
            )
        )}
    </FieldSet>
  );
};
