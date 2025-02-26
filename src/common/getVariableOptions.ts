import { getTemplateSrv } from '@grafana/runtime';
import { SelectableValue } from '@grafana/data';

interface VariableOptions {
  hideValue?: boolean;
  keepVarSyntax?: boolean;
}

export const getVariableOptions = (opts?: VariableOptions) => {
  const templateSrv = getTemplateSrv();
  return templateSrv.getVariables().map((t) => {
    const label = '${' + t.name + '}';
    const info: SelectableValue<string> = {
      label: opts?.hideValue ? `${label}` : `${label} = ${templateSrv.replace(label)}`,
      value: opts?.keepVarSyntax ? label : t.name,
      icon: 'arrow-right', // not sure what makes the most sense
    };
    return info;
  });
};
