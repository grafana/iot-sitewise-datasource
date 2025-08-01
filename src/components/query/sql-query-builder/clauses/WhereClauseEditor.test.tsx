import React from 'react';
import { render, screen } from '@testing-library/react';
import { WhereCondition } from '../types';
import { WhereClauseEditor } from './WhereClauseEditor';
import userEvent from '@testing-library/user-event';

jest.mock('@grafana/runtime', () => {
  const actual = jest.requireActual('@grafana/runtime');
  return {
    ...actual,
    getTemplateSrv: () => ({
      getVariables: () => [{ name: 'var1' }, { name: 'var2' }],
    }),
  };
});

const availableProperties = [
  { id: 'asset_id', name: 'assetId' },
  { id: 'asset_name', name: 'assetName' },
];
const setup = (
  whereConditions: WhereCondition[] = [{ column: '', operator: '=', value: '', logicalOperator: 'AND' }],
  updateQuery = jest.fn()
) => {
  render(
    <WhereClauseEditor
      whereConditions={whereConditions}
      updateQuery={updateQuery}
      availableProperties={availableProperties}
    />
  );
  return { updateQuery };
};

describe('WhereClauseEditor', () => {
  it('renders default condition row', () => {
    setup();
    expect(screen.getByText('WHERE')).toBeInTheDocument();
    expect(screen.getByText('Select column...')).toBeInTheDocument();
    expect(screen.getByText('=')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter value or $variable')).toBeInTheDocument();
    expect(screen.getByLabelText('Add condition')).toBeInTheDocument();
    expect(screen.getByLabelText('Remove condition')).toBeInTheDocument();
  });

  it('changes column value', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);

    const columnDropdown = screen.getByText('Select column...');
    await user.click(columnDropdown);
    const option = screen.getByText('assetId');
    await user.click(option);

    expect(updateQuery).toHaveBeenCalledWith(
      expect.objectContaining({
        whereConditions: [expect.objectContaining({ column: 'asset_id' })],
      })
    );
  });

  it('changes operator and resets value and operator2', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    const condition: WhereCondition = {
      column: 'asset_id',
      operator: '=',
      value: '123',
      value2: '456',
      logicalOperator: 'AND',
    };
    setup([condition], updateQuery);

    const operatorDropdown = screen.getByText('=');
    await user.click(operatorDropdown);
    const betweenOption = screen.getByText('BETWEEN');
    await user.click(betweenOption);

    expect(updateQuery).toHaveBeenCalledWith(
      expect.objectContaining({
        whereConditions: [
          expect.objectContaining({
            operator: 'BETWEEN',
            operator2: 'AND',
            value: '',
            value2: '',
          }),
        ],
      })
    );
  });

  it('renders BETWEEN operator input fields', () => {
    const updateQuery = jest.fn();
    const condition: WhereCondition = {
      column: 'asset_id',
      operator: 'BETWEEN',
      value: '10',
      value2: '20',
      logicalOperator: 'AND',
    };

    setup([condition], updateQuery);
    expect(screen.getByDisplayValue('10')).toBeInTheDocument();
    expect(screen.getByDisplayValue('20')).toBeInTheDocument();
    expect(screen.getByText('AND')).toBeInTheDocument(); // static AND dropdown
  });

  it('adds new where condition row', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    setup(undefined, updateQuery);

    const addButton = screen.getByLabelText('Add condition');
    await user.click(addButton);

    expect(updateQuery).toHaveBeenCalledWith(
      expect.objectContaining({
        whereConditions: expect.arrayContaining([
          expect.objectContaining({ column: '', operator: '=', value: '', logicalOperator: 'AND' }),
        ]),
      })
    );
  });

  it('removes a condition row', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    const whereConditions = [
      { column: 'asset_id', operator: '=', value: '123' },
      { column: 'asset_name', operator: '=', value: 'model' },
    ];
    setup(whereConditions, updateQuery);

    const removeButtons = screen.getAllByLabelText('Remove condition');
    await user.click(removeButtons[1]);

    expect(updateQuery).toHaveBeenCalledWith({
      whereConditions: [{ column: 'asset_id', operator: '=', value: '123' }],
    });
  });

  it('removes last condition but keeps a blank default condition', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    const whereConditions = [{ column: 'asset_id', operator: '=', value: '123' }];
    setup(whereConditions, updateQuery);

    const removeButton = screen.getByLabelText('Remove condition');
    await user.click(removeButton);

    expect(updateQuery).toHaveBeenCalledWith({
      whereConditions: [{ column: '', operator: '', value: '' }],
    });
  });

  it('updates logical operator between multiple conditions', async () => {
    const user = userEvent.setup();
    const updateQuery = jest.fn();
    const whereConditions: WhereCondition[] = [
      { column: 'asset_id', operator: '=', value: '123', logicalOperator: 'AND' },
      { column: 'asset_name', operator: '=', value: 'model', logicalOperator: 'AND' },
    ];
    setup(whereConditions, updateQuery);

    const andDropdown = screen.getByText('AND');
    await user.click(andDropdown);
    const orOption = screen.getByText('OR');
    await user.click(orOption);

    expect(updateQuery).toHaveBeenCalledWith({
      whereConditions: [
        { column: 'asset_id', operator: '=', value: '123', logicalOperator: 'OR' },
        { column: 'asset_name', operator: '=', value: 'model', logicalOperator: 'AND' },
      ],
    });
  });

  it('does not show VariableSuggestInput for IS NULL operator', () => {
    setup([{ column: 'asset_id', operator: 'IS NULL', value: '123' }]);

    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
  });

  it('shows VariableSuggestInput for value-based operators (like =)', () => {
    setup([{ column: 'asset_id', operator: '=', value: '123' }]);

    const input = screen.getByRole('textbox');
    expect(input).toBeInTheDocument();
  });

  it('does not show logicalOperator Select if only one condition', () => {
    const updateQuery = jest.fn();
    const whereConditions = [{ column: 'asset_id', operator: '=', value: '123' }];
    setup(whereConditions, updateQuery);

    expect(screen.queryByText('AND')).not.toBeInTheDocument();
    expect(screen.queryByText('OR')).not.toBeInTheDocument();
  });

  it('shows logicalOperator Select when there are multiple conditions', () => {
    const updateQuery = jest.fn();
    const whereConditions: WhereCondition[] = [
      { column: 'asset_id', operator: '=', value: '123', logicalOperator: 'AND' },
      { column: 'asset_name', operator: '=', value: 'model', logicalOperator: 'OR' },
    ];
    setup(whereConditions, updateQuery);

    expect(screen.getByText('AND')).toBeInTheDocument();
    const selectDropdown = screen.getAllByRole('combobox');
    expect(selectDropdown.length).toBeGreaterThan(0);
  });
});
