import { test, expect } from './helpers';
import { interceptRequests } from './interceptRequests';
import { type SitewiseOptions, type SitewiseSecureJsonData } from '../src/types';

test.describe('Query Editor', () => {
  test.describe('Queries', () => {
    test.beforeEach(async ({ page, panelEditPage, readProvisionedDataSource }) => {
      await interceptRequests(page);

      /* Configure data source */

      const ds = await readProvisionedDataSource<SitewiseOptions, SitewiseSecureJsonData>({
        fileName: 'mock-iot-sitewise.e2e.yaml',
      });
      await panelEditPage.datasource.set(ds.name);
      await panelEditPage.setVisualization('Table');
    });

    test('Get property value', async ({ page, panelEditPage, queryEditor }) => {
      await expect(queryEditor.assetSelect).not.toBeVisible();
      await expect(queryEditor.propertyAliasInput).not.toBeVisible();

      await queryEditor.selectQueryType('Get property value');

      await expect(queryEditor.assetSelect).toBeVisible();
      await expect(queryEditor.propertyAliasInput).toBeVisible();
      await expect(queryEditor.propertySelect).not.toBeVisible();

      await queryEditor.selectAsset('Demo Wind Farm Asset');

      await expect(queryEditor.propertySelect).toBeVisible();
      await expect(queryEditor.queryOptionsSelect).not.toBeVisible();

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).toBeVisible();

      await queryEditor.selectProperty('Total Average Power');

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).not.toBeVisible();
      await expect(panelEditPage.panel.data).toContainText(['15.6 kW', 'GOOD']);
    });

    test('Get property value history', async ({ page, panelEditPage, queryEditor }) => {
      await expect(queryEditor.propertyAliasInput).not.toBeVisible();
      await expect(queryEditor.assetSelect).not.toBeVisible();

      await queryEditor.selectQueryType('Get property value history');

      await expect(queryEditor.propertyAliasInput).toBeVisible();
      await expect(queryEditor.assetSelect).toBeVisible();

      await expect(queryEditor.queryOptionsSelect).not.toBeVisible();

      await queryEditor.selectAsset('Demo Wind Farm Asset');

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).toBeVisible();

      await queryEditor.selectProperty('Total Average Power');

      await queryEditor.openQueryOptions();
      await expect(queryEditor.qualitySelect).toBeVisible();
      await expect(queryEditor.timeSelect).toBeVisible();
      await expect(queryEditor.formatSelect).toBeVisible();

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).not.toBeVisible();
      await expect(panelEditPage.panel.data).toContainText(['15.6 kW', '14.3 kW', '16.3 kW', 'GOOD']);
    });

    test('Get property value aggregates', async ({ page, panelEditPage, queryEditor }) => {
      await expect(queryEditor.propertyAliasInput).not.toBeVisible();
      await expect(queryEditor.assetSelect).not.toBeVisible();

      await queryEditor.selectQueryType('Get property value aggregates');

      await expect(queryEditor.propertyAliasInput).toBeVisible();
      await expect(queryEditor.assetSelect).toBeVisible();

      await queryEditor.selectAsset('Demo Wind Farm Asset');

      await expect(queryEditor.aggregatePicker).not.toBeVisible();

      await queryEditor.openQueryOptions();
      await expect(queryEditor.resolutionSelect).not.toBeVisible();
      await expect(queryEditor.qualitySelect).not.toBeVisible();
      await expect(queryEditor.timeSelect).not.toBeVisible();
      await expect(queryEditor.formatSelect).not.toBeVisible();

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).toBeVisible();

      await queryEditor.selectProperty('Total Average Power');

      await expect(queryEditor.aggregatePicker).toBeVisible();
      await expect(queryEditor.resolutionSelect).toBeVisible();
      await expect(queryEditor.qualitySelect).toBeVisible();
      await expect(queryEditor.timeSelect).toBeVisible();
      await expect(queryEditor.formatSelect).toBeVisible();

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).not.toBeVisible();
      await expect(panelEditPage.panel.data).toContainText(['15.6 kW', '14.3 kW', '16.3 kW', 'GOOD']);
    });

    test('Get interpolated property values', async ({ page, panelEditPage, queryEditor }) => {
      await expect(queryEditor.propertyAliasInput).not.toBeVisible();
      await expect(queryEditor.assetSelect).not.toBeVisible();

      await queryEditor.selectQueryType('Get interpolated property values');

      await expect(queryEditor.propertyAliasInput).toBeVisible();
      await expect(queryEditor.assetSelect).toBeVisible();
      await expect(queryEditor.qualitySelect).not.toBeVisible();
      await expect(queryEditor.timeSelect).not.toBeVisible();
      await expect(queryEditor.formatSelect).not.toBeVisible();
      await expect(queryEditor.resolutionSelect).not.toBeVisible();

      await queryEditor.selectAsset('Demo Wind Farm Asset');

      await expect(queryEditor.queryOptionsSelect).not.toBeVisible();

      await queryEditor.runQuery();

      await expect(page.getByText('No data')).toBeVisible();

      await queryEditor.selectProperty('Total Average Power');
      await queryEditor.runQuery();

      await expect(page.getByText('No data')).not.toBeVisible();
      await expect(panelEditPage.panel.data).toContainText(['15.6 kW', '14.3 kW', '16.3 kW', 'GOOD']);
    });

    test('List assets', async ({ panelEditPage, queryEditor }) => {
      await expect(queryEditor.modelIdSelect).not.toBeVisible();
      await expect(queryEditor.filterSelect).not.toBeVisible();

      await queryEditor.selectQueryType('List assets');

      await expect(queryEditor.modelIdSelect).toBeVisible();
      await expect(queryEditor.filterSelect).toBeVisible();

      await queryEditor.runQuery();

      await expect(panelEditPage.panel.data).toContainText(['Demo Wind Farm Asset']);
    });

    test('List asset models', async ({ panelEditPage, queryEditor }) => {
      await queryEditor.selectQueryType('List asset models');
      await queryEditor.runQuery();

      await expect(panelEditPage.panel.data).toContainText(['Demo Wind Farm Asset Model']);
    });

    test('List associated assets', async ({ panelEditPage, queryEditor }) => {
      await expect(queryEditor.propertyAliasInput).not.toBeVisible();
      await expect(queryEditor.assetSelect).not.toBeVisible();
      await expect(queryEditor.hierarchySelect).not.toBeVisible();

      await queryEditor.selectQueryType('List associated assets');

      await expect(queryEditor.propertyAliasInput).not.toBeVisible();
      await expect(queryEditor.assetSelect).toBeVisible();
      await expect(queryEditor.hierarchySelect).toBeVisible();

      await queryEditor.selectHierarchy('** All **');
      await queryEditor.runQuery();

      await expect(panelEditPage.panel.data).toContainText(['Demo Turbine Asset']);
    });
  });

  test.describe('Raw Code Editor', () => {
    test.beforeEach(async ({ page, panelEditPage, readProvisionedDataSource }) => {
      await interceptRequests(page);

      /* Configure data source */

      const ds = await readProvisionedDataSource<SitewiseOptions, SitewiseSecureJsonData>({
        fileName: 'mock-iot-sitewise.e2e.yaml',
      });
      await panelEditPage.datasource.set(ds.name);
      await panelEditPage.setVisualization('Table');
    });

    // Uncomment the following code when Builder mode is ready
    // test('Switch to Code Editor', async ({ page, panelEditPage, selectors }) => {
    //   await page.getByRole('radio', { name: 'Code' }).click();
    //   await page.waitForFunction(() => window.monaco);
    //   await expect(panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container)).toBeVisible();
    // });

    // test('Displays the correct initial value', async ({ page, panelEditPage, selectors }) => {
    //   await page.getByRole('radio', { name: 'Code' }).click();
    //   await page.waitForFunction(() => window.monaco);
    //   await expect(panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container)).toContainText(
    //     'select $__selectAll from raw_time_series where $__timeFilter(event_timestamp)'
    //   );
    // });

    // test('Accepts keyboard input', async ({ page, panelEditPage, selectors }) => {
    //   await page.getByRole('radio', { name: 'Code' }).click();
    //   await page.waitForFunction(() => window.monaco);
    //   const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
    //   await editor.click();
    //   await page.keyboard.insertText('SELECT * FROM new_table');
    //   await expect(editor).toContainText('SELECT * FROM new_table');
    // });
  });
});
