import defaults from 'lodash/defaults';

import React, { PropsWithChildren, ReactElement } from 'react';
import { FormField } from './Fields';
import { MultiSelectCommonProps, SelectCommonProps } from '@grafana/ui/components/Select/types';
import { MultiSelect, Select } from '@grafana/ui';

export interface FormFieldSelectorProps<T> extends SelectCommonProps<T> {
  label?: string;
}

export interface FormFieldMultiSelectorProps<T> extends MultiSelectCommonProps<T> {
  label?: string;
}

const defaultProps: Partial<FormFieldSelectorProps<any>> = {
  width: 16,
};

export const FormFieldSelector = <T extends any>(
  p: PropsWithChildren<FormFieldSelectorProps<T>>
): ReactElement | null => {
  const props = defaults(p, defaultProps);

  return (
    <FormField label={props.label}>
      <Select {...props} />
    </FormField>
  );
};

export const FormFieldMultiSelector = <T extends any>(
  p: PropsWithChildren<FormFieldMultiSelectorProps<T>>
): ReactElement | null => {
  const props = defaults(p, defaultProps);

  return (
    <FormField label={props.label}>
      <MultiSelect {...props} />
    </FormField>
  );
};
