import React from 'react';
import { Select } from '@grafana/ui';
import { EditorField, EditorFieldGroup } from '@grafana/plugin-ui';
import { StyledLabel } from '../StyledLabel';

interface FromClauseEditorProps {
  assetModels: Array<{ id: string; name: string }>;
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

/**
 * Renders the "FROM" clause UI for the query builder.
 * It allows users to select an asset model (data source view) from a dropdown list.
 * Once model changed from the dropdown it resets the previous query fields (SELECT, WHERE, GROUP BY, ORDER BY)
 * to their default/initial state.
 */
export const FromClauseEditor: React.FC<FromClauseEditorProps> = ({ assetModels, selectedModelId, updateQuery }) => {
  return (
    <EditorFieldGroup>
      {/* Section label with tooltip */}
      <StyledLabel text="FROM" width={15} tooltip />

      {/* Dropdown to select a model */}
      <EditorField label="" width={30}>
        <Select
          options={assetModels.map((model) => ({
            label: model.name,
            value: model.id,
          }))}
          value={
            selectedModelId
              ? {
                  label: assetModels.find((m) => m.id === selectedModelId)?.name || '',
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
    </EditorFieldGroup>
  );
};
