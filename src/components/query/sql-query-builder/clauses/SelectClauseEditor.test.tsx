import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SelectClauseEditor } from './SelectClauseEditor';
import { SelectField } from '../types';

const availableProperties = [
  { id: 'asset_id', name: 'assetId' },
  { id: 'asset_name', name: 'assetName' },
];

const setup = (selectFields: SelectField[] = [{ column: '', aggregation: '', alias: '' }], updateQuery = jest.fn()) => {
  render(
    <SelectClauseEditor
      selectFields={selectFields}
      updateQuery={updateQuery}
      availableProperties={availableProperties}
    />
  );
  return { updateQuery };
};

describe('SelectClauseEditor', () => {
  it('renders default select field', () => {
    setup();
    expect(screen.getByText('Select column...')).toBeInTheDocument();
    expect(screen.getByText('Raw Values')).toBeInTheDocument();
    expect(screen.getAllByPlaceholderText('Optional alias')[0]).toBeInTheDocument();
  });

  it('adds a new select field when plus button is clicked', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);
    const addButton = screen.getByLabelText('Add field');
    await user.click(addButton);
    expect(updateQuery).toHaveBeenCalled();
  });

  it('removes a select field when minus button is clicked', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    const selectFields = [
      { column: 'asset_id', aggregation: '', alias: '' },
      { column: 'asset_name', aggregation: '', alias: '' },
    ];
    setup(selectFields, updateQuery);
    const removeButton = screen.getAllByLabelText('Remove field')[0];
    await user.click(removeButton);
    expect(updateQuery).toHaveBeenCalledWith({
      selectFields: [{ column: 'asset_name', aggregation: '', alias: '' }],
    });
  });

  it('updates alias input', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);
    const aliasInput = screen.getByPlaceholderText('Optional alias');
    await user.type(aliasInput, 'temp_alias');
    expect(updateQuery).toHaveBeenCalled();
  });

  it('updates column select', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);
    const columnDropdown = screen.getByText('Select column...');
    await user.click(columnDropdown);
    const option = screen.getByText('assetId');
    await user.click(option);
    expect(updateQuery).toHaveBeenCalled();
  });

  it('updates aggregation function', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);
    const funcDropdown = screen.getByText('Raw Values');
    await user.click(funcDropdown);
    const strReplaceOption = screen.getByText('String: STR_REPLACE');
    await user.click(strReplaceOption);
    expect(updateQuery).toHaveBeenCalled();
  });

  it('renders additional inputs for STR_REPLACE', async () => {
    const selectFields = [
      {
        column: 'asset_id',
        aggregation: 'STR_REPLACE',
        alias: '',
        functionArgValue: 'hot',
        functionArgValue2: 'warm',
      },
    ];
    setup(selectFields);
    expect(screen.getAllByPlaceholderText('Enter value')[0]).toBeInTheDocument();
    expect(screen.getAllByPlaceholderText('Enter value')).toHaveLength(2);
  });

  it('renders CAST function and selects BOOLEAN as arg type', async () => {
    const updateQuery = jest.fn();

    const selectFields = [
      {
        column: 'asset_id',
        aggregation: 'CAST',
        alias: '',
        functionArg: 'BOOLEAN',
      },
    ];

    setup(selectFields, updateQuery);
    expect(screen.getByText('asset_id')).toBeInTheDocument();
    expect(screen.getByText('DateTime: CAST')).toBeInTheDocument();
    expect(screen.getByText('BOOLEAN')).toBeInTheDocument();
  });

  it('handles CONCAT with single selection (edge case)', async () => {
    const selectFields = [
      {
        column: 'asset_id',
        aggregation: 'CONCAT',
        alias: '',
        functionArg: 'asset_name',
      },
    ];
    setup(selectFields);
    const allComboboxes = screen.getAllByRole('combobox');
    const argTypeSelect = allComboboxes[3]; // Adjust index based on actual DOM
    expect(screen.getByText('String: CONCAT')).toBeInTheDocument();
    await userEvent.click(argTypeSelect);
    const option = screen.getByText('asset_name');
    await userEvent.click(option);
    expect(screen.getByText('asset_id')).toBeInTheDocument();
  });

  it('does not remove last field', async () => {
    const updateQuery = jest.fn();
    setup([{ column: '', aggregation: '', alias: '' }], updateQuery);
    const removeButtons = screen.queryAllByLabelText('Remove field');
    expect(removeButtons.length).toBe(0);
  });

  it('handles empty aggregation gracefully', () => {
    const selectFields = [{ column: 'asset_id', aggregation: '', alias: '' }];
    setup(selectFields);
    expect(screen.getByText('asset_id')).toBeInTheDocument();
  });
});
