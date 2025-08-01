import React from 'react';
import { EditorRows, EditorRow } from '@grafana/plugin-ui';
import { SqlQueryBuilderProps, mockAssetModels } from './types';
import { FromClauseEditor } from './clauses/FromClauseEditor';
import { SelectClauseEditor } from './clauses/SelectClauseEditor';
import { WhereClauseEditor } from './clauses/WhereClauseEditor';
import { GroupByClauseEditor } from './clauses/GroupByClauseEditor';
import { LimitClauseEditor } from './clauses/LimitClauseEditor';
import { OrderByClauseEditor } from './clauses/OrderByClauseEditor';
import { QueryPreviewDisplay } from './QueryPreviewDisplay';
import { useSQLQueryState } from './hooks/useSQLQueryState';
import { HavingClauseEditor } from './clauses/HavingClauseEditor';

/**
 * SqlQueryBuilder
 *
 * A SQL query builder component used in the Grafana plugin for constructing queries interactively.
 * It provides editor sections for common SQL clauses such as FROM, SELECT, WHERE, GROUP BY, HAVING,ORDER BY etc.
 *
 * - @param builderState - Initial query state passed from parent.
 * - @param onChange - Callback to notify parent of query changes.
 */
export function SqlQueryBuilder({ builderState, onChange }: SqlQueryBuilderProps) {
  const { queryState, preview, validationErrors, updateQuery, availableProperties, availablePropertiesForGrouping } =
    useSQLQueryState({
      initialQuery: builderState,
      onChange: onChange,
    });

  // HAVING clause is only shown when the query includes a GROUP BY
  const isHavingVisible = queryState.groupByTags.length > 0;

  return (
    <div className="gf-form-group">
      <EditorRows>
        <EditorRow>
          {/* FROM Clause Editor
              - Allows user to select an asset model (table equivalent)
              - Also handles LIMIT clause configuration */}
          <FromClauseEditor
            assetModels={mockAssetModels}
            selectedModelId={queryState.selectedAssetModel || ''}
            updateQuery={updateQuery}
          />

          {/* LIMIT Clause Editor
              - Sets maximum number of rows to return */}
          <LimitClauseEditor limit={queryState.limit} updateQuery={updateQuery} />
        </EditorRow>

        {/* SELECT Clause Editor
            - Allows selecting fields, aggregations and alias */}
        <SelectClauseEditor
          selectFields={queryState.selectFields}
          updateQuery={updateQuery}
          availableProperties={availableProperties}
        />

        {/* WHERE Clause Editor
            - Defines filters and conditions on data */}
        <WhereClauseEditor
          whereConditions={queryState.whereConditions}
          updateQuery={updateQuery}
          availableProperties={availableProperties}
        />

        {/* GROUP BY Clause Editor
            - Enables grouping the query result by specified fields */}
        <GroupByClauseEditor
          availablePropertiesForGrouping={availablePropertiesForGrouping}
          groupByTags={queryState.groupByTags}
          updateQuery={updateQuery}
        />

        {/* HAVING Clause Editor
            - Conditional logic after GROUP BY (only visible if groupBy is used) */}
        {isHavingVisible && (
          <HavingClauseEditor
            havingConditions={queryState.havingConditions}
            updateQuery={updateQuery}
            availableProperties={availablePropertiesForGrouping}
          />
        )}

        {/* ORDER BY Clause Editor
            - Specifies sorting of the query results */}
        <OrderByClauseEditor
          orderByFields={queryState.orderByFields}
          updateQuery={updateQuery}
          availableProperties={availableProperties}
        />
      </EditorRows>

      {/* Query Preview Display
          - Shows the generated SQL query text and any validation errors */}
      <QueryPreviewDisplay preview={preview} errors={validationErrors} />
    </div>
  );
}
