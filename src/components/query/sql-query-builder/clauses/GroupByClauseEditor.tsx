import React, { useMemo, useCallback } from 'react';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { Select, ActionMeta, FieldSet, Stack } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

interface PropertyOption {
  id: string;
  name: string;
}

interface GroupByClauseEditorProps {
  availablePropertiesForGrouping: PropertyOption[];
  groupByTags: string[];
  updateQuery: (fields: Partial<{ groupByTags: string[] }>) => void;
}

/**
 *
 * A React component that renders a multi-select dropdown to choose one or more columns
 * for the SQL `GROUP BY` clause. Used in a sql query builder.
 */
export const GroupByClauseEditor: React.FC<GroupByClauseEditorProps> = ({
  availablePropertiesForGrouping,
  groupByTags,
  updateQuery,
}) => {
  /**
   * Memoized transformation of available columns
   */
  const groupByOptions: Array<SelectableValue<string>> = useMemo(
    () =>
      availablePropertiesForGrouping.map(({ id, name }) => ({
        value: id,
        label: name,
      })),
    [availablePropertiesForGrouping]
  );

  /**
   * Memoized selection of currently selected columns for GROUP BY options
   */
  const selectedGroupByOptions: Array<SelectableValue<string>> = useMemo(
    () =>
      groupByTags.map((tag) => {
        return groupByOptions.find((opt) => opt.value === tag) || { value: tag, label: tag };
      }),
    [groupByTags, groupByOptions]
  );

  /**
   * Handler for when user updates selected GROUP BY columns.
   * Converts selected options to an array of string IDs and
   * calls the updateQuery callback with new state.
   *
   * @param options - selected options from the <Select> dropdown
   */
  const handleGroupByTagsChange = useCallback(
    (options: SelectableValue<string> | Array<SelectableValue<string>>, _meta?: ActionMeta) => {
      const tags: string[] = options.map((opt: any) => opt.value).filter(Boolean) as string[];
      const nextState: Partial<{ groupByTags: string[] }> = {
        groupByTags: tags,
      };
      updateQuery(nextState);
    },
    [updateQuery]
  );

  return (
    <EditorRow>
      <FieldSet label="Group By">
        <Stack direction="row" gap={4} alignItems="center">
          {/* Choose GROUP BY columns */}
          <EditorField label="" width={30}>
            <Select
              options={groupByOptions}
              value={selectedGroupByOptions}
              onChange={handleGroupByTagsChange}
              isMulti
              placeholder="Select column(s)..."
            />
          </EditorField>
        </Stack>
      </FieldSet>
    </EditorRow>
  );
};
