import React, { InputHTMLAttributes, FunctionComponent } from 'react';
import { InlineFormLabel } from '@grafana/ui';

export interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  children?: React.ReactNode;
  width?: number | 'auto';
}

// stolen from https://github.com/grafana/grafana/blob/dde5b724e860b7899c59629ea8cc7be74d3c5b11/public/app/plugins/datasource/cloudwatch/components/Forms.tsx#L19
export const FormField: FunctionComponent<Partial<Props>> = ({ label, width = 8, children }) => {
  return (
    <>
      <InlineFormLabel width={width} className="query-keyword">
        {label}
      </InlineFormLabel>
      {children}
    </>
  );
};

export const FormInlineField: FunctionComponent<Props> = ({ ...props }) => {
  return (
    <div className={'gf-form-inline'}>
      <div className="gf-form gf-form--grow">
        <FormField {...props} />
        <div className="gf-form-label gf-form-label--grow" />
      </div>
    </div>
  );
};
