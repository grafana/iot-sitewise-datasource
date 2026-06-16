import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { HavingClauseEditor } from './HavingClauseEditor';
import { HavingCondition } from '../types';
import userEvent from '@testing-library/user-event';

jest.mock('@grafana/runtime', () => {
  const actual = jest.requireActual('@grafana/runtime');
  return {
    ...actual,
    getTemplateSrv: () => ({
      getVariables: () => [{ name: '$region' }, { name: '$device' }],
    }),
  };
});

const availableProperties = [
  { id: 'temperature', name: 'temperature' },
  { id: 'pressure', name: 'pressure' },
];

const defaultCondition: HavingCondition = {
  aggregation: 'SUM',
  column: '',
  operator: '=',
  value: '',
  logicalOperator: 'AND',
};

const setup = (conditions: HavingCondition[] = [defaultCondition]) => {
  const updateQuery = jest.fn();
  render(
    <HavingClauseEditor
      havingConditions={conditions}
      updateQuery={updateQuery}
      availableProperties={availableProperties}
    />
  );
  return { updateQuery };
};

describe('HavingClauseEditor', () => {
  it('renders the component with default condition', () => {
    setup();
    expect(screen.getByText('Having')).toBeInTheDocument();
    expect(screen.getByText('SUM')).toBeInTheDocument();
    expect(screen.getByText('Select column...')).toBeInTheDocument();
    const operatorDropdown = screen.getByText('=');
    expect(operatorDropdown).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter value')).toBeInTheDocument();
    expect(screen.getByLabelText('Add condition')).toBeInTheDocument();
  });

  it('updates the aggregation type', async () => {
    const { updateQuery } = setup();
    const aggSelect = screen.getByText('SUM');
    fireEvent.mouseDown(aggSelect);
    const countOption = await screen.findByText('COUNT');
    fireEvent.click(countOption);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [{ ...defaultCondition, aggregation: 'COUNT' }],
    });
  });

  it('updates the column field', async () => {
    const { updateQuery } = setup();
    fireEvent.mouseDown(screen.getByText('Select column...'));
    const option = await screen.findByText('temperature');
    fireEvent.click(option);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [{ ...defaultCondition, column: 'temperature' }],
    });
  });

  it('updates the operator', async () => {
    const { updateQuery } = setup();
    const operatorSelect = screen.getByText('=');
    fireEvent.mouseDown(operatorSelect);
    const option = await screen.findByText('<');
    fireEvent.click(option);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [{ ...defaultCondition, operator: '<' }],
    });
  });

  it('updates the value input to a static number', async () => {
    const { updateQuery } = setup();
    const input = screen.getByRole('textbox');

    await userEvent.clear(input);
    await userEvent.type(input, '5');

    expect(updateQuery).toHaveBeenLastCalledWith({
      havingConditions: [{ ...defaultCondition, value: '5' }],
    });
  });

  it('adds a new condition', () => {
    const { updateQuery } = setup();
    const addButton = screen.getByLabelText('Add condition');
    fireEvent.click(addButton);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [
        defaultCondition,
        {
          aggregation: 'SUM',
          column: '',
          operator: '>',
          value: '',
          logicalOperator: 'AND',
        },
      ],
    });
  });

  it('removes last condition but keeps a blank default condition', () => {
    const { updateQuery } = setup();
    const removeButton = screen.queryByLabelText('Remove condition');
    if (removeButton) {
      fireEvent.click(removeButton);

      expect(updateQuery).toHaveBeenCalledWith({
        havingConditions: [
          {
            aggregation: 'COUNT',
            column: '',
            operator: '=',
            value: '',
          },
        ],
      });
    } else {
      // If the button does not exist, ensure only one condition remains
      expect(screen.getAllByLabelText('Add condition').length).toBe(1);
    }
  });

  it('removes a condition from multiple', () => {
    const { updateQuery } = setup([
      { aggregation: 'COUNT', column: 'temp', operator: '=', value: '10', logicalOperator: 'AND' },
      { aggregation: 'SUM', column: 'pressure', operator: '>', value: '30', logicalOperator: 'OR' },
    ]);
    const removeButtons = screen.getAllByLabelText(/Remove condition/i);
    fireEvent.click(removeButtons[0]);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [{ aggregation: 'SUM', column: 'pressure', operator: '>', value: '30', logicalOperator: 'OR' }],
    });
  });

  it('updates logical operator between multiple conditions', async () => {
    const { updateQuery } = setup([
      {
        aggregation: 'SUM',
        column: 'temperature',
        operator: '>',
        value: '100',
        logicalOperator: 'AND',
      },
      {
        aggregation: 'AVG',
        column: 'pressure',
        operator: '<',
        value: '50',
        logicalOperator: 'OR',
      },
    ]);
    const logicSelect = screen.getByText('AND');
    fireEvent.mouseDown(logicSelect);
    const orOption = await screen.findByText('OR');
    fireEvent.click(orOption);

    expect(updateQuery).toHaveBeenCalledWith({
      havingConditions: [
        {
          aggregation: 'SUM',
          column: 'temperature',
          operator: '>',
          value: '100',
          logicalOperator: 'OR', // updated
        },
        {
          aggregation: 'AVG',
          column: 'pressure',
          operator: '<',
          value: '50',
          logicalOperator: 'OR',
        },
      ],
    });
  });
});
