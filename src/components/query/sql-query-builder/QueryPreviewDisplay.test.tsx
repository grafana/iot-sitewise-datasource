import React from 'react';
import { render, screen } from '@testing-library/react';
import { QueryPreviewDisplay } from './QueryPreviewDisplay';
import { ValidationError } from './types';

describe('QueryPreviewDisplay', () => {
  it('renders the preview text when there are no errors', () => {
    render(<QueryPreviewDisplay preview="SELECT * FROM turbine" errors={[]} />);

    expect(screen.getByText('SELECT * FROM turbine')).toBeInTheDocument();

    const box = screen.getByText('SELECT * FROM turbine').closest('div');
    expect(box).toBeInTheDocument();
  });

  it('renders the preview text with error color when errors are present', () => {
    const errors: ValidationError[] = [{ error: 'Syntax error', type: 'column' }];
    render(<QueryPreviewDisplay preview="SELECT * FROM turbine WHERE" errors={errors} />);

    expect(screen.getByText('SELECT * FROM turbine WHERE')).toBeInTheDocument();

    const box = screen.getByText('SELECT * FROM turbine WHERE').closest('div');
    expect(box).toBeInTheDocument();
  });

  it('renders inside a Box with marginTop', () => {
    const { container } = render(<QueryPreviewDisplay preview="SELECT 1" errors={[]} />);
    const boxElement = container.querySelector('div');
    expect(boxElement).toBeInTheDocument();
  });
});
