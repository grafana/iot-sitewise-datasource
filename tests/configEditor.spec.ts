import { test, expect } from '@grafana/plugin-e2e';
import { SitewiseOptions, SitewiseSecureJsonData } from '../src/types';
import { handleMocks } from './handleMocks';

test.describe('ConfigEditor', () => {
  test('invalid credentials should return a 400 status code', async ({
    createDataSourceConfigPage,
    readProvisionedDataSource,
    page,
  }) => {
    const provisionedDatasource = await handleMocks(page, '/health', 'e2e-sitewise-invalid-credentials');

    // create a new datasource
    const configPage = await createDataSourceConfigPage({
      type: 'grafana-iot-sitewise-datasource',
    });
    await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
    await page.keyboard.press('Enter');

    // get the provisioned datasource options
    const ds = await readProvisionedDataSource<SitewiseOptions, SitewiseSecureJsonData>(provisionedDatasource);

    // fill in the config form
    await page.getByLabel('Name').fill(ds.name || '');
    await page.getByLabel('Access Key ID').fill(ds.secureJsonData?.accessKey || '');
    await page.getByLabel('Secret Access Key').fill(ds.secureJsonData?.secretKey || '');
    await page.getByLabel('Default Region').fill('us-east-1');
    await page.keyboard.press('Enter');

    // click save and test
    const response = await configPage.saveAndTest();

    // expect network response have error (this is only a meaningful test when we run this with live (not mocked) data)
    const body = await response.json();
    expect(body).toHaveProperty('status', 'ERROR');
    expect(body.message).toContain('invalid');

    // expect error to be shown in the UI
    const errorMessage = await page.getByText('The security token included in the request is invalid');
    expect(errorMessage).toBeVisible();
  });

  test('valid credentials should return a 200 status code', async ({
    createDataSourceConfigPage,
    readProvisionedDataSource,
    page,
  }) => {
    const provisionedDatasource = await handleMocks(page, '/health', 'e2e-sitewise-valid-credentials');

    // create a new datasource
    const configPage = await createDataSourceConfigPage({ type: 'grafana-iot-sitewise-datasource' });
    await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
    await page.keyboard.press('Enter');

    // get the provisioned datasource options
    const ds = await readProvisionedDataSource<SitewiseOptions, SitewiseSecureJsonData>(provisionedDatasource);

    // fill in the config form
    await page.getByLabel('Name').fill(ds.name || '');
    await page.getByLabel('Access Key ID').fill(ds.secureJsonData?.accessKey || '');
    await page.getByLabel('Secret Access Key').fill(ds.secureJsonData?.secretKey || '');
    await page.getByLabel('Default Region').fill('us-east-1');
    await page.keyboard.press('Enter');

    // click save and test
    const response = await configPage.saveAndTest();

    // expect network response have error (this is only a meaningful test when we run this with live (not mocked) data)
    const body = await response.json();
    expect(body).toHaveProperty('status', 'OK');

    // expect success message to be shown in the UI
    const successMessage = page.getByText('OK');
    expect(successMessage).toBeVisible();
  });
});
