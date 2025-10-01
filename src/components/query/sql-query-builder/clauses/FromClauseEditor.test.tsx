import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { FromClauseEditor } from './FromClauseEditor';
import { ValidationError } from '../types';

const queryReferenceViews = [
  { id: 'model-1', name: 'Asset Model 1' },
  { id: 'model-2', name: 'Asset Model 2' },
];

const setup = (selectedModelId = '', customModels = queryReferenceViews, validationErrors: ValidationError[] = []) => {
  const mockUpdateQuery = jest.fn();
  render(
    <FromClauseEditor
      queryReferenceViews={customModels}
      selectedModelId={selectedModelId}
      validationErrors={validationErrors}
      updateQuery={mockUpdateQuery}
    />
  );
  return { mockUpdateQuery };
};

describe('FromClauseEditor', () => {
  it('renders the dropdown with model options', () => {
    setup();
    expect(screen.getByText('Select view...')).toBeInTheDocument();
  });

  it('calls updateQuery when a model is selected', async () => {
    const { mockUpdateQuery } = setup();
    const dropdown = screen.getByText('Select view...');
    fireEvent.mouseDown(dropdown);
    const option = screen.queryByText('Asset Model 1');
    if (option) {
      fireEvent.click(option);
    }

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      selectedAssetModel: 'model-1',
      selectFields: [{ column: '', aggregation: '', alias: '' }],
      whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
      groupByFields: [{ column: '' }],
      orderByFields: [{ column: '', direction: 'ASC' }],
    });
  });

  it('calls updateQuery when a different model is selected', async () => {
    const { mockUpdateQuery } = setup('model-1');
    const dropdown = screen.getByText('Asset Model 1');
    await userEvent.click(dropdown);
    const option = screen.getByText('Asset Model 2');
    await userEvent.click(option);

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      selectedAssetModel: 'model-2',
      selectFields: [{ column: '', aggregation: '', alias: '' }],
      whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
      groupByFields: [{ column: '' }],
      orderByFields: [{ column: '', direction: 'ASC' }],
    });
  });

  it('shows the selected model name if selectedModelId is given', () => {
    setup('model-2');
    expect(screen.getByText('Asset Model 2')).toBeInTheDocument();
  });

  it('dropdown has all model options', async () => {
    setup();
    const dropdown = screen.getByText('Select view...');
    await userEvent.click(dropdown);

    for (const model of queryReferenceViews) {
      expect(await screen.findByText(model.name)).toBeInTheDocument();
    }
  });

  it('calls updateQuery when selecting using mouse click', async () => {
    const { mockUpdateQuery } = setup();
    const dropdown = screen.getByText('Select view...');
    fireEvent.mouseDown(dropdown);

    const option = await screen.findByText('Asset Model 1');
    fireEvent.click(option);

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      selectedAssetModel: 'model-1',
      selectFields: [{ column: '', aggregation: '', alias: '' }],
      whereConditions: [{ column: '', operator: '', value: '', logicalOperator: 'AND' }],
      groupByFields: [{ column: '' }],
      orderByFields: [{ column: '', direction: 'ASC' }],
    });
  });

  it('shows no options when queryReferenceViews is empty', async () => {
    setup('', []);
    const dropdown = screen.getByText('Select view...');
    await userEvent.click(dropdown);

    expect(screen.queryByRole('option')).not.toBeInTheDocument();
  });

  it('renders validation errors below the dropdown (only for "from" type)', () => {
    const errors: ValidationError[] = [
      { type: 'from', error: 'You must select a model' },
      { type: 'select', error: 'This should not render here' },
    ];
    setup('', queryReferenceViews, errors);
    expect(screen.getByText('You must select a model')).toBeInTheDocument();
    expect(screen.queryByText('This should not render here')).not.toBeInTheDocument();
  });
});
