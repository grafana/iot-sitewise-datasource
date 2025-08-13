import React from 'react';
import { render, screen } from '@testing-library/react';
import { QueryPreviewDisplay } from './QueryPreviewDisplay';

describe('QueryPreviewDisplay', () => {
  it('renders only preview when there are no errors', () => {
    render(<QueryPreviewDisplay preview="SELECT * FROM table" errors={[]} />);
    expect(screen.getByText('Query Preview')).toBeInTheDocument();
    expect(screen.getByText('SELECT * FROM table')).toBeInTheDocument();
  });

  it('renders errors and preview when errors exist', () => {
    const errors = ['Missing WHERE clause', 'Invalid function'];
    render(<QueryPreviewDisplay preview="SELECT name" errors={errors} />);
    expect(screen.getByText('Query Errors & Preview')).toBeInTheDocument();
    errors.forEach((err) => {
      expect(screen.getByText(err)).toBeInTheDocument();
    });
    expect(screen.getByText('SELECT name')).toBeInTheDocument();
  });
});
