import React, { useEffect, useRef, useState } from 'react';
import { Input } from '@grafana/ui';
import { getSelectableTemplateVariables } from 'variables';

interface Props {
  value: string;
  onChange: (val: string) => void;
  placeholder?: string;
}

// Wrapper style for positioning suggestion list relative to input
const wrapperStyle: React.CSSProperties = {
  position: 'relative',
};

// Styles for the suggestion dropdown list
const suggestionListStyle: React.CSSProperties = {
  position: 'absolute',
  zIndex: 10,
  background: 'var(--page-bg, #121212)',
  marginTop: 4,
  width: '100%',
  maxHeight: 152,
  overflowY: 'auto',
  listStyle: 'none',
  padding: 0,
  border: '1px solid var(--panel-border-color, #444)',
  borderRadius: 4,
};

// Style for each item in the suggestion list
const suggestionItemStyle: React.CSSProperties = {
  padding: '8px 12px',
  cursor: 'pointer',
  fontSize: 13,
};

const hoverStyle: React.CSSProperties = {
  backgroundColor: 'var(--input-hover-bg, #2a2a2a)',
};

export const VariableSuggestInput: React.FC<Props> = ({ value = '', onChange, placeholder }) => {
  const [allVars, setAllVars] = useState<string[]>([]); // Holds all available template variables
  const [suggestions, setSuggestions] = useState<string[]>([]); // Filtered suggestions based on input
  const [showSuggestions, setShowSuggestions] = useState(false); // Toggle for showing suggestion list
  const inputRef = useRef<HTMLInputElement>(null); // Ref to focus back after selecting a suggestion
  const [hoveredItem, setHoveredItem] = useState<string | null>(null); // Track hover for inline style

  /**
   * On component mount, fetch all selectable template variable names
   * and clean them by removing `${}` wrapping.
   */
  useEffect(() => {
    const vars = getSelectableTemplateVariables().map((v) => v.value.replace(/\$\{|\}/g, ''));
    setAllVars(vars);
  }, []);

  /**
   * Handles input changes and updates suggestions list if user types `$`
   *
   * @param e - The input change event
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputVal = e.currentTarget.value;
    onChange(inputVal);

    const matchStart = inputVal.lastIndexOf('$');
    if (matchStart !== -1) {
      const currentWord = inputVal.slice(matchStart + 1);
      const matched = allVars.filter((v) => v.startsWith(currentWord));
      setSuggestions(matched);
      setShowSuggestions(true);
    } else {
      setShowSuggestions(false);
    }
  };

  /**
   * Inserts selected suggestion into the input field at the last `$` position,
   * wraps the variable in `${}` format, and refocuses the input.
   *
   * @param suggestion - The selected variable name
   */
  const insertSuggestion = (suggestion: string) => {
    const matchStart = value.lastIndexOf('$');
    if (matchStart !== -1) {
      const before = value.substring(0, matchStart);
      const replaced = `${before}\${${suggestion}}`;
      onChange(replaced);
      setShowSuggestions(false);
      setTimeout(() => inputRef.current?.focus(), 0);
    }
  };

  return (
    <div style={wrapperStyle}>
      <Input
        ref={inputRef}
        value={value}
        onChange={handleChange}
        placeholder={placeholder || 'Enter value or $variable'}
      />
      {/* Render suggestion list only when needed */}
      {showSuggestions && suggestions.length > 0 && (
        <ul style={suggestionListStyle}>
          {suggestions.map((s) => (
            <li
              key={s}
              style={{
                ...suggestionItemStyle,
                ...(hoveredItem === s ? hoverStyle : {}),
              }}
              onMouseDown={() => insertSuggestion(s)}
              onMouseEnter={() => setHoveredItem(s)}
              onMouseLeave={() => setHoveredItem(null)}
            >
              ${s}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
