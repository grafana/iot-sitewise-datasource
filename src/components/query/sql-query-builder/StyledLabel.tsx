import React from 'react';
import { EditorField } from '@grafana/plugin-ui';
import { InlineLabel } from '@grafana/ui';
import { tooltipMessages } from './types';

interface StyledLabelProps {
  text: string;
  width?: number;
  color?: string; // HEX, RGB, or CSS color string
  tooltip?: boolean;
  bold?: boolean;
  fontSize?: string;
  className?: string; // allows external CSS if needed
}

/**
 * StyledLabel
 *
 * A reusable label component for Grafana plugin editors that supports:
 * - Custom font size and color
 * - Optional bold styling
 * - Tooltip integration (from tooltipMessages)
 * - External CSS classes via `className`
 */
export const StyledLabel: React.FC<StyledLabelProps> = ({
  text,
  width = 10,
  color = '#6e9fff',
  tooltip = false,
  bold = true,
  fontSize,
  className,
}) => {
  const style: React.CSSProperties = {
    color,
    fontWeight: bold ? 'bold' : undefined,
    fontSize,
  };

  const tooltipText = tooltip ? tooltipMessages[text] : undefined;

  return (
    <EditorField label="" width={width}>
      <InlineLabel width="auto" style={style} tooltip={tooltipText} className={className}>
        {text}
      </InlineLabel>
    </EditorField>
  );
};
