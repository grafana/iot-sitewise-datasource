import { languages, IRange } from 'monaco-editor/esm/vs/editor/editor.api';
import { MACROS } from './macros';
import { Monaco } from '@grafana/ui';

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
}

const tableColumns: KeyValue = {
  asset: ['asset_id', 'asset_name', 'asset_description', 'asset_model_id'],
  asset_property: [
    'property_id',
    'asset_id',
    'property_name',
    'property_data_type',
    'property_alias',
    'asset_composite_model_id',
  ],
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
  fetchSuggestions(range: IRange, types: SuggestionType): languages.CompletionList;
  tableDefinitions(): SuggestionDefinition[];
  fieldDefinitions(table: string): SuggestionDefinition[];
  macroDefinitions(): SuggestionDefinition[];
  allDefinitions(): SuggestionDefinition[];
  buildAutocompleteSuggestion(definition: SuggestionDefinition, range: IRange): languages.CompletionItem;
  monaco: null | Monaco;
  currentSpace: string;
}

export const SitewiseCompletionProvider: SitewiseCompletionProviderType = {
  triggerCharacters: ['.', ' ', '$', ',', '(', "'"],

  monaco: null,

  currentSpace: 'start',

  provideCompletionItems(model, position, context): languages.ProviderResult<languages.CompletionList> {
    // Setup
    const word = model.getWordUntilPosition(position);
    const range = {
      startLineNumber: position.lineNumber,
      endLineNumber: position.lineNumber,
      startColumn: word.startColumn,
      endColumn: word.endColumn,
    };
    let suggestionType = SuggestionType.all;

    var last_chars = model.getValueInRange({
      startLineNumber: position.lineNumber,
      startColumn: 0,
      endLineNumber: position.lineNumber,
      endColumn: position.column,
    });
    var words = last_chars.trim().replace('\t', '').split(' ');

    // TODO Find the table so that we can return the fields as well

    const lastWord = words[words.length - 1].toLowerCase();
    // First the last word
    if (lastWord === 'from') {
      this.currentSpace = 'from';
      suggestionType = SuggestionType.tables;
    } else if (lastWord === 'where') {
      this.currentSpace = 'where';
      suggestionType = SuggestionType.macros;
      // Then the current space
    } else if (this.currentSpace === 'from') {
      suggestionType;
    } else if (this.currentSpace === 'where') {
      suggestionType = SuggestionType.macros;
      // Then everything
    } else {
      suggestionType = SuggestionType.all;
    }

    return {
      suggestions: this.fetchSuggestions(range, suggestionType).suggestions,
    };
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

  fetchSuggestions(range: IRange, types: SuggestionType): languages.CompletionList {
    let definitions: SuggestionDefinition[] = [];

    switch (types) {
      case SuggestionType.macros:
        definitions = this.macroDefinitions();
        break;
      case SuggestionType.tables:
        definitions = this.tableDefinitions();
        break;
      default:
        definitions = this.allDefinitions();
        break;
    }

    const suggestions = definitions.map((definition) => {
      return this.buildAutocompleteSuggestion(definition, range);
    });

    return {
      suggestions,
    };
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

  macroDefinitions(): SuggestionDefinition[] {
    return MACROS.map((macro) => {
      return {
        label: macro.id,
        kind: this.monaco?.languages.CompletionItemKind.Function || 0,
        documentation: macro.description,
        insertText: macro.id,
      };
    });
  },

  allDefinitions(): SuggestionDefinition[] {
    return this.tableDefinitions().concat(this.macroDefinitions());
  },
};
