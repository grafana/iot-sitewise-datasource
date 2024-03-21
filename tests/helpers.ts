import { test as base, expect } from '@grafana/plugin-e2e';
import { QueryEditor } from './queryEditor.model';

interface Fixtures {
  queryEditor: QueryEditor;
}

const test = base.extend<Fixtures>({
  /** Isolated `QueryEditorPage` instance. */
  queryEditor: async ({ page, grafanaVersion }, use) => {
    const queryEditor = new QueryEditor(page, grafanaVersion);

    await use(queryEditor);
  },
});

export { test, expect };
