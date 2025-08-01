import React from 'react';
import { render, screen } from '@testing-library/react';
import { StyledLabel } from './StyledLabel';

describe('StyledLabel', () => {
  it('renders the label text', () => {
    render(<StyledLabel text="FROM" />);
    expect(screen.getByText('FROM')).toBeInTheDocument();
  });

  it('applies inline styles for color and fontWeight when not using Tailwind', () => {
    render(<StyledLabel text="Custom" color="#ff0000" bold />);
    const label = screen.getByText('Custom');
    expect(label).toHaveStyle('color: #ff0000');
    expect(label).toHaveStyle('font-weight: bold');
  });

  it('uses Tailwind class for color when passed as a utility class', () => {
    render(<StyledLabel text="Tailwind" color="text-green-500" />);
    const label = screen.getByText('Tailwind');
    expect(label.className).toMatch(/text-green-500/);
  });

  it('does not apply inline color when Tailwind class is used', () => {
    render(<StyledLabel text="Tailwind" color="text-green-500" />);
    const label = screen.getByText('Tailwind');
    expect(label).not.toHaveStyle('color: text-green-500');
  });

  it('applies custom fontSize when provided', () => {
    render(<StyledLabel text="Sized" fontSize="18px" />);
    const label = screen.getByText('Sized');
    expect(label).toHaveStyle('font-size: 18px');
  });

  it('applies custom className if provided', () => {
    render(<StyledLabel text="Styled" className="custom-css" />);
    const label = screen.getByText('Styled');
    expect(label.className).toMatch(/custom-css/);
  });

  it('renders tooltip info icon when tooltip is enabled', () => {
    render(<StyledLabel text="FROM" tooltip />);
    const icon = screen.getByTestId('info-circle');
    expect(icon).toBeInTheDocument();
  });
});
