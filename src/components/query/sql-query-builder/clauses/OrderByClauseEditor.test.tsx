import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { OrderByClauseEditor } from './OrderByClauseEditor';
import { OrderByField } from '../types';

const availableProperties = [
  { id: 'asset_id', name: 'Asset ID' },
  { id: 'asset_name', name: 'Asset Name' },
];

describe('OrderByClauseEditor', () => {
  const mockUpdateQuery = jest.fn();

  const setup = (orderByFields: OrderByField[] = [{ column: 'asset_id', direction: 'ASC' }]) => {
    mockUpdateQuery.mockReset();
    const user = userEvent.setup();
    render(
      <OrderByClauseEditor
        orderByFields={orderByFields}
        updateQuery={mockUpdateQuery}
        availableProperties={availableProperties}
      />
    );
    return { user };
  };

  it('renders with initial ORDER BY field', () => {
    setup();

    expect(screen.getByText('Order By')).toBeInTheDocument();
    expect(screen.getByText('asset_id')).toBeInTheDocument();
    expect(screen.getByText('Ascending')).toBeInTheDocument();
  });

  it('adds a new ORDER BY field when plus button is clicked', () => {
    setup();
    const addButton = screen.getByRole('button', { name: /Add Order By field/i });
    fireEvent.click(addButton);

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      orderByFields: [
        { column: 'asset_id', direction: 'ASC' },
        { column: '', direction: 'ASC' },
      ],
    });
  });

  it('removes an ORDER BY field when minus button is clicked', () => {
    const twoFields: OrderByField[] = [
      { column: 'asset_id', direction: 'ASC' },
      { column: 'asset_name', direction: 'DESC' },
    ];
    setup(twoFields);
    const removeButtons = screen.getAllByRole('button', { name: /remove order by field/i });
    fireEvent.click(removeButtons[0]);

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      orderByFields: [{ column: 'asset_name', direction: 'DESC' }],
    });
  });

  it('updates the column when a different column is selected', async () => {
    const { user } = setup();
    const columnSelect = screen.getByText('asset_id');
    await user.click(columnSelect);
    await user.click(screen.getByText('Asset Name'));

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      orderByFields: [{ column: 'asset_name', direction: 'ASC' }],
    });
  });

  it('updates the sort direction when selected', async () => {
    const { user } = setup();
    const directionSelect = screen.getByText('Ascending');
    await user.click(directionSelect);
    await user.click(await screen.findByText('Descending'));

    expect(mockUpdateQuery).toHaveBeenCalledWith({
      orderByFields: [{ column: 'asset_id', direction: 'DESC' }],
    });
  });

  it('falls back to empty string when column option value is undefined', () => {
    setup([{ column: undefined as any, direction: 'ASC' }]);
    const comboboxes = screen.getAllByRole('combobox');

    expect(comboboxes[0]).toHaveValue('');
  });
});
