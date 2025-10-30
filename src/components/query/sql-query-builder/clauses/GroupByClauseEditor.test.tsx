import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { GroupByClauseEditor } from './GroupByClauseEditor';
import userEvent from '@testing-library/user-event';

const mockProperties = [
  { id: 'asset_id', name: 'assetId' },
  { id: 'asset_name', name: 'assetName' },
  { id: 'property_id', name: 'propertyId' },
];

const setup = (groupByTags: string[] = [], availableProperties = mockProperties) => {
  const updateQuery = jest.fn();
  render(
    <GroupByClauseEditor
      availablePropertiesForGrouping={availableProperties}
      groupByTags={groupByTags}
      updateQuery={updateQuery}
    />
  );
  return { updateQuery };
};

describe('GroupByClauseEditor', () => {
  it('renders the Group By label', () => {
    setup();
    expect(screen.getByText('Group By')).toBeInTheDocument();
  });
  it('renders the Group By placeholder', () => {
    setup();
    expect(screen.getByText('Select column(s)...')).toBeInTheDocument();
  });

  it('renders selected groupByTags if present', () => {
    setup(['asset_name']);
    expect(screen.getByText('assetName')).toBeInTheDocument();
  });

  it('displays unknown groupByTags that are not in available options', () => {
    setup(['unknown_property']);
    expect(screen.getByText('unknown_property')).toBeInTheDocument();
  });

  it('applies expected CSS classes to the dropdown container', () => {
    setup();
    const container = screen.getByText('Select column(s)...').closest('div');
    expect(container?.className).toMatch(/css-/);
  });

  it('displays all options when opened', async () => {
    setup();
    const dropdown = screen.getByText('Select column(s)...');
    fireEvent.mouseDown(dropdown);

    expect(await screen.findByText('assetId')).toBeInTheDocument();
    expect(screen.getByText('assetName')).toBeInTheDocument();
    expect(screen.getByText('propertyId')).toBeInTheDocument();
  });

  it('renders empty dropdown when no availableProperties are provided', async () => {
    setup([], []);

    const dropdown = screen.getByText('Select column(s)...');
    fireEvent.mouseDown(dropdown);

    expect(screen.queryByText('assetId')).not.toBeInTheDocument();
    expect(screen.queryByText('assetName')).not.toBeInTheDocument();
    expect(screen.queryByText('propertyId')).not.toBeInTheDocument();
  });

  it('calls updateQuery with selected tags when option selected', async () => {
    const { updateQuery } = setup();

    const dropdown = screen.getByText('Select column(s)...');
    fireEvent.mouseDown(dropdown);

    const option = await screen.findByText('assetName');
    fireEvent.click(option);

    expect(updateQuery).toHaveBeenCalledWith({
      groupByTags: ['asset_name'],
    });
  });

  it('calls updateQuery with multiple selected tags', async () => {
    const user = userEvent.setup();
    const { updateQuery } = setup();

    const dropdown = screen.getByText('Select column(s)...');
    await user.click(dropdown);

    const tempOption = await screen.findByText('assetId');
    await user.click(tempOption);

    expect(updateQuery).toHaveBeenCalledWith({
      groupByTags: ['asset_id'],
    });

    const input = screen.getByRole('combobox');
    await user.click(input);

    const pressureOption = await screen.findByText('propertyId');
    await user.click(pressureOption);

    expect(updateQuery).toHaveBeenLastCalledWith({
      groupByTags: ['property_id'],
    });
  });

  it('handles custom groupByTags not in options', () => {
    setup(['custom_column']);
    expect(screen.getByText('custom_column')).toBeInTheDocument();
  });

  it('calls updateQuery with empty array when cleared', async () => {
    const { updateQuery } = setup(['asset_id']);
    const tag = screen.getByText('assetId');

    const clearBtn = tag.parentElement?.querySelector('[aria-label="Remove"]');
    if (clearBtn) {
      fireEvent.click(clearBtn);
      expect(updateQuery).toHaveBeenCalledWith({
        groupByTags: [],
      });
    }
  });
});
