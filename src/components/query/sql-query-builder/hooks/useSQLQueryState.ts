import { useState, useEffect, useRef } from 'react';
import { isEqual } from 'lodash';
import { SitewiseQueryState, AssetProperty, mockAssetModels } from '../types';
import { validateQuery } from '../utils/validateQuery';
import { generateQueryPreview } from '../utils/queryGenerator';

interface UseSQLQueryStateOptions {
  initialQuery: SitewiseQueryState;
  onChange: (query: SitewiseQueryState) => void;
}

interface UseSQLQueryStateResult {
  queryState: SitewiseQueryState;
  setQueryState: React.Dispatch<React.SetStateAction<SitewiseQueryState>>;
  preview: string;
  validationErrors: string[];
  updateQuery: (newState: Partial<SitewiseQueryState>) => Promise<void>;
  selectedModel: any | undefined;
  availableProperties: AssetProperty[];
  availablePropertiesForGrouping: AssetProperty[];
}

/**
 * A custom React hook to manage the state of a SQL-like query builder
 *
 * Responsibilities:
 * - Manage current query state (`SitewiseQueryState`)
 * - Generate SQL preview from query state with  validation and collect errors
 * - Notify parent component when query changes
 *
 * @param initialQuery Initial state of the query
 * @param onChange Callback when query state is updated
 * @returns An object containing query state, helpers, and derived values
 */
export const useSQLQueryState = ({ initialQuery, onChange }: UseSQLQueryStateOptions): UseSQLQueryStateResult => {
  const [queryState, setQueryState] = useState<SitewiseQueryState>(initialQuery);
  const [preview, setPreview] = useState('');
  const [validationErrors, setValidationErrors] = useState<string[]>([]);
  const queryStateRef = useRef(queryState);

  /**
   * Sync query state changes with parent via onChange callback.
   * Uses lodash `isEqual` to avoid unnecessary updates.
   */
  useEffect(() => {
    if (!isEqual(queryStateRef.current, queryState)) {
      queryStateRef.current = queryState;
      onChange(queryState);
    }
  }, [queryState, onChange]);

  /**
   * Auto-validate and generate SQL preview whenever query state changes.
   * This effect runs asynchronously and guards against updates after unmount.
   */
  useEffect(() => {
    let isMounted = true;
    const validateAndGenerate = async () => {
      const errors = validateQuery(queryState);
      const preview = await generateQueryPreview(queryState);

      if (isMounted) {
        setValidationErrors(errors);
        setPreview(preview);
      }
    };
    validateAndGenerate();
    return () => {
      isMounted = false;
    };
  }, [queryState]);

  /**
   * Allows updates to the query state. Regenerates SQL preview string
   * and updates the query state including `rawSQL`.
   *
   * @param newState updates to merge into existing query state
   */
  const updateQuery = async (newState: Partial<SitewiseQueryState>) => {
    const updatedStateBeforeSQL = { ...queryStateRef.current, ...newState };
    const rawSQL = await generateQueryPreview(updatedStateBeforeSQL);
    setQueryState({ ...updatedStateBeforeSQL, rawSQL });
  };

  const selectedModel = mockAssetModels.find((model) => model.id === queryState.selectedAssetModel);
  const availableProperties = selectedModel?.properties || [];
  const availablePropertiesForGrouping = availableProperties.filter((prop) =>
    queryState.selectFields.some((field) => field.column === prop.name)
  );

  // Return all state and helper values to the calling component
  return {
    queryState,
    setQueryState,
    preview,
    validationErrors,
    updateQuery,
    selectedModel,
    availableProperties,
    availablePropertiesForGrouping,
  };
};
