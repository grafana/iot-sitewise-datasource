import React from 'react';
import { render, screen } from '@testing-library/react';
import { QueryPreviewDisplay } from './QueryPreviewDisplay';

describe('QueryPreviewDisplay', () => {
  it('renders only preview when there are no errors', () => {
    render(<QueryPreviewDisplay preview="SELECT * FROM table" errors={[]} />);
    expect(screen.getByText('Query Preview')).toBeInTheDocument();
    expect(screen.getByText('SELECT * FROM table')).toBeInTheDocument();
  });
});
