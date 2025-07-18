import { languages, IRange } from 'monaco-editor/esm/vs/editor/editor.api';
import { MACROS } from './macros';
import { Monaco } from '@grafana/ui';
import { getTemplateSrv } from '@grafana/runtime';

interface KeyValue {
  [key: string]: string[];
}

interface SuggestionDefinition extends Omit<languages.CompletionItem, 'range' | 'insertText'> {
  insertText?: languages.CompletionItem['insertText'];
}

enum SuggestionType {
  'all',
  'macros',
  'tables',
  'fields',
  'variables',
}

const tableColumns: KeyValue = {
  asset: ['asset_id', 'asset_name', 'asset_description', 'asset_model_id'],
  asset_property: ['property_id', 'asset_id', 'property_name', 'property_alias', 'asset_composite_model_id'],
  raw_time_series: [
    'asset_id',
    'property_id',
    'property_alias',
    'event_timestamp',
    'quality',
    'boolean_value',
    'int_value',
    'double_value',
    'string_value',
  ],
  latest_value_time_series: [
    'asset_id',
    'property_id',
    'property_alias',
    'event_timestamp',
    'quality',
    'boolean_value',
    'int_value',
    'double_value',
    'string_value',
  ],
  precomputed_aggregates: [
    'asset_id',
    'property_id',
    'property_alias',
    'event_timestamp',
    'resolution',
    'sum_value',
    'count_value',
    'average_value',
    'maximum_value',
    'minimum_value',
    'stdev_value',
  ],
};

interface SitewiseCompletionProviderType extends languages.CompletionItemProvider {
  fetchSuggestions(range: IRange, types: SuggestionType, table: null | string): languages.CompletionItem[];
  tableDefinitions(): SuggestionDefinition[];
  variableDefinitions(): SuggestionDefinition[];
  fieldDefinitions(table: string): SuggestionDefinition[];
  macroDefinitions(range: IRange): SuggestionDefinition[];
  allDefinitions(range: IRange, table: null | string): SuggestionDefinition[];
  buildAutocompleteSuggestion(definition: SuggestionDefinition, range: IRange): languages.CompletionItem;
  monaco: null | Monaco;
  currentToken: string;
}

// TODO: Check out getStandardSQLCompletionProvider to get standard SQL completion
export const SitewiseCompletionProvider: SitewiseCompletionProviderType = {
  triggerCharacters: ['.', ' ', '$', ',', '(', "'"],

  monaco: null,

  currentToken: 'start',

  provideCompletionItems(model, position, context): languages.ProviderResult<languages.CompletionList> {
    // Setup
    const word = model.getWordUntilPosition(position);
    const range = {
      startLineNumber: position.lineNumber,
      endLineNumber: position.lineNumber,
      startColumn: word.startColumn,
      endColumn: word.endColumn,
    };
    let suggestionType = [SuggestionType.all];

    const last_chars = model.getValueInRange({
      startLineNumber: position.lineNumber,
      startColumn: 0,
      endLineNumber: position.lineNumber,
      endColumn: position.column,
    });
    const words = last_chars.trim().replace('\t', '').split(' ');

    const lastWord = words[words.length - 1].toLowerCase();

    const regResult = /from\s(\w+)/g.exec(last_chars);
    const currentTable = regResult === null ? null : regResult[1];
    const lineText = model.getValueInRange({
      startLineNumber: position.lineNumber,
      startColumn: 0,
      endLineNumber: position.lineNumber,
      endColumn: position.column,
    });

    const isVariableTrigger = lineText.endsWith('${') || lineText.endsWith('$');

    // Check the last word first (before the current space)
    if (isVariableTrigger) {
      suggestionType = [SuggestionType.variables];
    } else if (lastWord === 'from') {
      this.currentToken = 'from';
      suggestionType = [SuggestionType.tables];
    } else if (['where', 'and', 'or'].includes(lastWord)) {
      this.currentToken = 'where';
      if (currentTable === null) {
        suggestionType = [SuggestionType.macros];
      } else {
        suggestionType = [SuggestionType.fields, SuggestionType.macros];
      }
      // If the last word doesn't match any of the above, check the Current Space
    } else if (this.currentToken === 'from') {
      suggestionType = [SuggestionType.tables];
    } else if (this.currentToken === 'where') {
      suggestionType = [SuggestionType.macros];
      // Otherwise suggest everything
    } else {
      suggestionType = [SuggestionType.all];
    }

    let suggestions: languages.CompletionItem[] = [];
    suggestionType.forEach((value) => {
      suggestions = suggestions.concat(this.fetchSuggestions(range, value, currentTable));
    });

    return { suggestions: suggestions };
  },

  buildAutocompleteSuggestion(
    { label, detail, documentation, kind, insertText }: SuggestionDefinition,
    range: IRange
  ): languages.CompletionItem {
    const insertFallback = typeof label === 'string' ? label : label.label;
    const labelObject = typeof label === 'string' ? { label: label, description: detail } : { ...label };

    labelObject.description ??= detail;

    return {
      label: labelObject,
      kind: kind,
      insertText: insertText ?? insertFallback,
      range,
      documentation: documentation,
      detail: detail,
    };
  },

  fetchSuggestions(range: IRange, types: SuggestionType, table: null | string): languages.CompletionItem[] {
    let definitions: SuggestionDefinition[] = [];

    switch (types) {
      case SuggestionType.macros:
        definitions = this.macroDefinitions(range);
        break;
      case SuggestionType.tables:
        definitions = this.tableDefinitions();
        break;
      case SuggestionType.fields:
        if (table != null) {
          definitions = this.fieldDefinitions(table);
        }
        break;
      case SuggestionType.variables:
        definitions = definitions.concat(this.variableDefinitions());
        break;
      default:
        definitions = this.allDefinitions(range, table);
        break;
    }

    return definitions.map((definition) => {
      return this.buildAutocompleteSuggestion(definition, range);
    });
  },

  tableDefinitions(): SuggestionDefinition[] {
    return Object.keys(tableColumns).map((table) => {
      return {
        label: table,
        detail: 'Table',
        kind: this.monaco?.languages.CompletionItemKind.Text || 0,
      };
    });
  },

  fieldDefinitions(table: string): SuggestionDefinition[] {
    return tableColumns[table].map((column) => {
      return {
        label: column,
        detail: 'Field',
        kind: this.monaco?.languages.CompletionItemKind.Field || 0,
      };
    });
  },

  macroDefinitions(range: IRange): SuggestionDefinition[] {
    return MACROS.map((macro) => {
      return {
        label: macro.id,
        kind: this.monaco?.languages.CompletionItemKind.Function || 0,
        documentation: macro.description,
        insertText: macro.id.substring(1, macro.id.length),
        range,
      };
    });
  },

  variableDefinitions(): SuggestionDefinition[] {
    const templateSrv = getTemplateSrv();
    const variables = templateSrv.getVariables();
    return variables.map((v) => {
      return {
        label: v.name,
        detail: 'Grafana Variable',
        kind: this.monaco?.languages.CompletionItemKind.Variable || 0,
        insertText: '{' + v.name + '}',
      };
    });
  },

  allDefinitions(range: IRange, table: string): SuggestionDefinition[] {
    let definitions = this.tableDefinitions().concat(this.macroDefinitions(range)).concat(this.variableDefinitions());
    if (table != null) {
      definitions = definitions.concat(this.fieldDefinitions(table));
    }
    return definitions;
  },
};
