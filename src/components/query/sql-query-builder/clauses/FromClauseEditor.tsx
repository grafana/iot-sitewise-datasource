import React from 'react';
import { Select, FieldSet, Stack } from '@grafana/ui';
import { EditorField } from '@grafana/plugin-ui';
import { css } from '@emotion/css';
interface FromClauseEditorProps {
  queryReferenceViews: Array<{ id: string; name: string }>;
  selectedModelId: string;
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

const noMarginBottom = css`
  margin-bottom: 0;
`;

/**
 * Renders the "FROM" clause UI for the query builder.
 * It allows users to select an asset model (data source view) from a dropdown list.
 * Once model changed from the dropdown it resets the previous query fields (SELECT, WHERE, GROUP BY, ORDER BY)
 * to their default/initial state.
 */
export const FromClauseEditor: React.FC<FromClauseEditorProps> = ({
  queryReferenceViews,
  selectedModelId,
  updateQuery,
}) => {
  return (
    <FieldSet label="From" className={noMarginBottom}>
      <Stack direction="row" gap={4} alignItems="center">
        <EditorField label="View" htmlFor="view" width={40}>
          {/* Dropdown to select a model */}
          <Select
            inputId="view"
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
                selectFields: [{ column: 'all', aggregation: '', alias: '' }],
                whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
                groupByFields: [{ column: '' }],
                orderByFields: [{ column: '', direction: 'ASC' }],
              })
            }
            placeholder="Select view..."
          />
        </EditorField>
      </Stack>
    </FieldSet>
  );
};
