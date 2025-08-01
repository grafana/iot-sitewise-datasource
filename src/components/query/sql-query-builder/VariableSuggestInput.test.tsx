import React from 'react';
import { render, fireEvent, screen, waitFor } from '@testing-library/react';
import { VariableSuggestInput } from './VariableSuggestInput';
import * as variableModule from 'variables';

jest.mock('variables', () => ({
  getSelectableTemplateVariables: jest.fn(),
}));

const mockVars = [{ value: '${env}' }, { value: '${region}' }, { value: '${device}' }];

const setup = (value = '', onChange = jest.fn(), placeholder = '') => {
  (variableModule.getSelectableTemplateVariables as jest.Mock).mockReturnValue(mockVars);

  render(<VariableSuggestInput value={value} onChange={onChange} placeholder={placeholder} />);

  const input = screen.getByPlaceholderText(placeholder || 'Enter value or $variable');
  return { input, onChange };
};

describe('VariableSuggestInput', () => {
  it('renders input with default placeholder', () => {
    const { input } = setup();
    expect(input).toBeInTheDocument();
  });

  it('renders input with custom placeholder', () => {
    const { input } = setup('', jest.fn(), 'Custom Placeholder');
    expect(input).toHaveAttribute('placeholder', 'Custom Placeholder');
  });

  it('allows manual value entry without showing suggestions', () => {
    const onChange = jest.fn();
    const { input } = setup('', onChange);
    fireEvent.change(input, { target: { value: 'custom-value' } });
    expect(onChange).toHaveBeenCalledWith('custom-value');
    expect(screen.queryByRole('list')).not.toBeInTheDocument();
    expect(screen.queryByText('$region')).not.toBeInTheDocument();
  });

  it('calls onChange on input change and shows suggestions', () => {
    const { input, onChange } = setup();
    fireEvent.change(input, { target: { value: 'Hello $' } });

    expect(onChange).toHaveBeenCalledWith('Hello $');
    expect(screen.getByText('$env')).toBeInTheDocument();
    expect(screen.getByText('$region')).toBeInTheDocument();
    expect(screen.getByText('$device')).toBeInTheDocument();
  });

  it('filters suggestions based on input', () => {
    const { input } = setup();
    fireEvent.change(input, { target: { value: 'Hello $re' } });

    expect(screen.getByText('$region')).toBeInTheDocument();
    expect(screen.queryByText('$env')).not.toBeInTheDocument();
  });

  it('shows variable suggestions when typing a matching $variable', async () => {
    const { input } = setup('', jest.fn());

    await waitFor(() => {
      expect(variableModule.getSelectableTemplateVariables).toHaveBeenCalled();
    });
    fireEvent.change(input, { target: { value: '$re' } });
    await waitFor(() => {
      expect(screen.getByText('$region')).toBeInTheDocument();
    });
  });

  it('does not show suggestions when match fails', () => {
    const { input } = setup('$xyz');
    fireEvent.change(input, { target: { value: '$xyz' } });

    expect(screen.queryByText('$env')).not.toBeInTheDocument();
    expect(screen.queryByRole('list')).not.toBeInTheDocument();
  });

  it('hides suggestions if input does not contain $', () => {
    const { input } = setup();
    fireEvent.change(input, { target: { value: 'No variables here' } });

    expect(screen.queryByText('$env')).not.toBeInTheDocument();
  });
});
