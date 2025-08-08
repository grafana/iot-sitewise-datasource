import React from 'react';
import { EditorField } from '@grafana/plugin-ui';
import { InlineLabel } from '@grafana/ui';
import clsx from 'clsx';
import { tooltipMessages } from './types';

interface StyledLabelProps {
  text: string;
  width?: number;
  color?: string;
  tooltip?: boolean;
  bold?: boolean;
  fontSize?: string;
  className?: string;
}

/**
 * StyledLabel
 *
 * A reusable label component for Grafana plugin editors that supports:
 * - Custom font size and color
 * - Optional bold styling
 * - Tailwind text color classes
 * - Optional tooltip integration
 *
 * Wrapped inside `EditorField` to align with Grafana’s plugin editor layout.
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
  const isTailwindColor = color.startsWith('text-');

  const style: React.CSSProperties = {
    color: isTailwindColor ? undefined : color,
    fontWeight: bold ? 'bold' : undefined,
    fontSize,
  };

  const combinedClassName = clsx(isTailwindColor && color, bold && 'font-bold', className);

  const tooltipText = tooltip ? tooltipMessages[text] : undefined;

  return (
    <EditorField label="" width={width}>
      <InlineLabel width="auto" style={style} tooltip={tooltipText} className={combinedClassName}>
        {text}
      </InlineLabel>
    </EditorField>
  );
};
