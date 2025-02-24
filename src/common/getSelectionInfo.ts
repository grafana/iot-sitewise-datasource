import { SelectableValue } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';

export interface SelectionInfo<T> {
  options: Array<SelectableValue<T>>;
  current?: SelectableValue<T>;
}

const handleValueNotFound = <T = any>(value: T, showNotFound = true) => {
  const current = {
    label: `${value} ${showNotFound ? '(not found)' : ''}`,
    value,
  };
  if (current.label!.indexOf('$') >= 0) {
    current.label = `${value}`;
    const escaped = getTemplateSrv().replace(String(value));
    if (escaped !== current.label) {
      current.label += ` (variable)`;
    } else {
      current.label += ` ${escaped}`;
    }
  }
  return current;
};

export function getSelectionInfo<T>(
  v?: T,
  options?: Array<SelectableValue<T>>,
  templateVars?: Array<SelectableValue<T>>,
  allowCustom?: boolean
): SelectionInfo<T> {
  if (v && !options) {
    const current = { label: `${v}`, value: v };
    return { options: [current], current };
  }
  if (!options) {
    options = [];
  }
  let current = options.find((item) => item.value === v);
  if (templateVars) {
    if (!current) {
      current = templateVars.find((item) => item.value === v);
    }
    options = [{ label: 'Use template variable', options: templateVars, icon: 'link-h' }, ...options];
  }

  if (v && !current) {
    current = handleValueNotFound(v, false);
    options.push(current);
  }
  return {
    options,
    current,
  };
}
