import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { LimitClauseEditor } from './LimitClauseEditor';
import { ValidationError } from '../types';

const setup = (limit: number | undefined = undefined, validationErrors: ValidationError[] = []) => {
  const mockUpdateQuery = jest.fn();
  render(<LimitClauseEditor limit={limit} updateQuery={mockUpdateQuery} validationErrors={validationErrors} />);
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

  it('renders "Limit must be a valid number." error when limit is NaN', () => {
    const validationErrors: ValidationError[] = [
      { type: 'limit', error: 'Limit must be a valid number.' },
      { type: 'select', error: 'This should not render here' },
    ];
    setup(undefined, validationErrors);

    expect(screen.getByText('Limit must be a valid number.')).toBeInTheDocument();
    expect(screen.queryByText('This should not render here')).not.toBeInTheDocument();
  });

  it('renders "Limit must be greater than 0." error when limit is zero or negative', () => {
    const validationErrors: ValidationError[] = [
      { type: 'limit', error: 'Limit must be greater than 0.' },
      { type: 'select', error: 'This should not render here' },
    ];
    setup(0, validationErrors);

    expect(screen.getByText('Limit must be greater than 0.')).toBeInTheDocument();
    expect(screen.queryByText('This should not render here')).not.toBeInTheDocument();
  });

  it('renders "Limit must not exceed 100,000 rows." error when limit is too large', () => {
    const validationErrors: ValidationError[] = [
      { type: 'limit', error: 'Limit must not exceed 100,000 rows.' },
      { type: 'select', error: 'This should not render here' },
    ];
    setup(100001, validationErrors);

    expect(screen.getByText('Limit must not exceed 100,000 rows.')).toBeInTheDocument();
    expect(screen.queryByText('This should not render here')).not.toBeInTheDocument();
  });
});
