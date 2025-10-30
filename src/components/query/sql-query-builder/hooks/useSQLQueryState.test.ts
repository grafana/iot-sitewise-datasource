import { renderHook, act, waitFor } from '@testing-library/react';
import { useSQLQueryState } from './useSQLQueryState';
import * as validateModule from '../utils/validateQuery';
import * as generatorModule from '../utils/queryGenerator';
import { SitewiseQueryState, queryReferenceViews, defaultSitewiseQueryState } from '../types';

jest.useFakeTimers();

// Mocks
jest.mock('../utils/validateQuery', () => ({
  validateQuery: jest.fn().mockReturnValue([]),
}));

jest.mock('../utils/queryGenerator', () => ({
  generateQueryPreview: jest.fn().mockResolvedValue('SELECT column FROM asset;'),
}));

const mockPreview = 'SELECT column FROM asset;';
const mockErrors = ['Missing WHERE condition'];

const mockQuery: SitewiseQueryState = {
  ...defaultSitewiseQueryState,
  selectedAssetModel: 'asset',
  selectFields: [{ column: 'asset_name', aggregation: '', alias: 'name' }],
  whereConditions: [{ column: 'asset_id', operator: '=', value: '123' }],
  groupByTags: ['department'],
  orderByFields: [{ column: 'asset_name', direction: 'ASC' }],
  limit: 500,
  rawSQL: '',
};

describe('useSQLQueryState', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (validateModule.validateQuery as jest.Mock).mockReturnValue(mockErrors);
    (generatorModule.generateQueryPreview as jest.Mock).mockReturnValue(mockPreview);
  });

  it('initializes with full query state', async () => {
    const onChange = jest.fn();

    const { result } = renderHook(() => useSQLQueryState({ initialQuery: mockQuery, onChange }));

    await waitFor(() => {
      expect(result.current.queryState).toEqual(mockQuery);
      expect(result.current.preview).toBe(mockPreview);
      expect(result.current.validationErrors).toEqual(mockErrors);
    });
  });

  it('updates query state and triggers validation + preview', async () => {
    const onChange = jest.fn();
    (generatorModule.generateQueryPreview as jest.Mock).mockReturnValue('MOCK_PREVIEW');
    (validateModule.validateQuery as jest.Mock).mockReturnValue(['MOCK_ERROR']);

    const { result } = renderHook(() => useSQLQueryState({ initialQuery: mockQuery, onChange }));

    await act(async () => {
      result.current.setQueryState((prev) => ({
        ...prev,
        limit: 200,
      }));
    });

    await waitFor(() => {
      expect(result.current.preview).toBe('MOCK_PREVIEW');
      expect(result.current.validationErrors).toEqual(['MOCK_ERROR']);
    });

    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({ limit: 200 }));
  });

  it('computes selectedModel and availablePropertiesForGrouping correctly', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useSQLQueryState({ initialQuery: mockQuery, onChange }));

    const selectedModel = queryReferenceViews.find((m) => m.id === 'asset');
    const availableProperties = selectedModel?.properties ?? [];

    const availablePropertiesForGrouping = availableProperties.filter((prop) =>
      mockQuery.selectFields.some((field) => field.column === prop.name)
    );

    expect(result.current.selectedModel).toEqual(selectedModel);
    expect(result.current.availableProperties).toEqual(availableProperties);
    expect(result.current.availablePropertiesForGrouping).toEqual(availablePropertiesForGrouping);
  });

  it('can update deeply nested fields like selectFields', () => {
    const onChange = jest.fn();
    const { result } = renderHook(() => useSQLQueryState({ initialQuery: mockQuery, onChange }));

    act(() => {
      result.current.updateQuery({
        selectFields: [{ column: 'asset_name', aggregation: '', alias: 'name' }],
      });
    });

    expect(result.current.queryState.selectFields).toEqual([{ column: 'asset_name', aggregation: '', alias: 'name' }]);
  });
});
