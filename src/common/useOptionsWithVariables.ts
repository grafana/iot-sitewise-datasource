import { useMemo } from 'react';
import { getSelectionInfo } from './getSelectionInfo';
import { getVariableOptions } from './getVariableOptions';
import { SelectableValue } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';

export const useOptionsWithVariables = ({
  current,
  options,
}: {
  current?: string;
  options: SelectableValue<string>[];
}) => {
  const variableOptions = getVariableOptions({ keepVarSyntax: true });
  const variables = getTemplateSrv().getVariables();
  return useMemo(() => getSelectionInfo(current, options, variableOptions), [current, variables, options]);
};
