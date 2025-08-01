import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { LimitClauseEditor } from './LimitClauseEditor';

const setup = (limit: number | undefined = undefined) => {
  const mockUpdateQuery = jest.fn();
  render(<LimitClauseEditor limit={limit} updateQuery={mockUpdateQuery} />);
  const input = screen.getByRole('spinbutton');
  return { mockUpdateQuery, input };
};

describe('LimitClauseEditor', () => {
  it('renders with empty input if limit is undefined', () => {
    const { input } = setup();
    expect(input).toHaveValue(null);
  });

  it('renders with provided limit', () => {
    const { input } = setup(50);
    expect(input).toHaveValue(50);
  });

  it('calls updateQuery with number on valid input', () => {
    const { input, mockUpdateQuery } = setup();
    fireEvent.change(input, { target: { value: '25' } });
    expect(mockUpdateQuery).toHaveBeenCalledWith({ limit: 25 });
  });

  it('calls updateQuery with undefined on empty input', () => {
    const { input, mockUpdateQuery } = setup(100);
    fireEvent.change(input, { target: { value: '' } });
    expect(mockUpdateQuery).toHaveBeenCalledWith({ limit: undefined });
  });

  it('does not call updateQuery on non-numeric input', () => {
    const { input, mockUpdateQuery } = setup();
    fireEvent.change(input, { target: { value: 'abc' } });
    expect(mockUpdateQuery).not.toHaveBeenCalled();
  });
});
