import { test as base, expect } from '@grafana/plugin-e2e';
import { QueryEditor } from './queryEditor.model';

interface Fixtures {
  queryEditor: QueryEditor;
}

const test = base.extend<Fixtures>({
  /** Isolated `QueryEditorPage` instance. */
  queryEditor: async ({ page }, use) => {
    const queryEditor = new QueryEditor(page);

    await use(queryEditor);
  },
});

export { test, expect };
