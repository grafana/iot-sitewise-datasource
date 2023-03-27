import {
  BaseVariableModel,
  ConstantVariableModel,
  CustomVariableModel,
  LoadingState,
  VariableHide,
} from '@grafana/data';
const baseModel: BaseVariableModel = {
  name: '',
  id: '',
  type: 'query',
  rootStateKey: null,
  global: false,
  index: -1,
  hide: VariableHide.dontHide,
  skipUrlSync: false,
  state: LoadingState.NotStarted,
  error: null,
  description: null,
};
export const assetIdVariableConstant: ConstantVariableModel = {
  ...baseModel,
  id: 'assetIdConstant',
  name: 'assetIdConstant',
  type: 'constant',
  query: '',
  current: {
    value: 'valueConstant',
    text: 'valueConstant',
    selected: true,
  },
  options: [],
};

export const assetIdVariableArray: CustomVariableModel = {
  ...baseModel,
  id: 'assetIdArray',
  name: 'assetIdArray',
  type: 'custom',
  query: 'array1,array2",array3',
  current: {
    value: ['array1', 'array2', 'array3'],
    text: ['array1', 'array2', 'array3'],
    selected: true,
  },
  options: [
    {
      value: 'array1',
      text: 'array1',
      selected: true,
    },
    {
      value: 'array2',
      text: 'array2',
      selected: true,
    },
    {
      value: 'array3',
      text: 'array3',
      selected: true,
    },
  ],
  multi: true,
  includeAll: true,
};
