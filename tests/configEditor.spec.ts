import { test, expect, ReadProvisionedDataSourceArgs, DataSourceSettings } from '@grafana/plugin-e2e';
import { SitewiseOptions, SitewiseSecureJsonData } from '../src/types';
import dotenv from 'dotenv';
import path from 'path';

// Read from ".env" file.
dotenv.config({ path: path.resolve(__dirname, '.env') });

test.describe('ConfigEditor', () => {
  test('invalid credentials should return a 400 status code', async ({
    createDataSourceConfigPage,
    page,
    selectors,
  }) => {
    // create a new datasource and navigate to config page
    const configPage = await createDataSourceConfigPage({ type: 'grafana-iot-sitewise-datasource' });

    // fill in the config form
    await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
    await page.keyboard.press('Enter');
    await page.getByLabel('Access Key ID').fill('bad1credentials');
    await page.getByLabel('Secret Access Key').fill('very-bad-credentials');
    await page.getByRole('combobox', { name: 'Default Region' }).click();
    await configPage.getByGrafanaSelector(selectors.components.Select.option).getByText('us-east-1').click();
    await expect(configPage.saveAndTest()).not.toBeOK();
    await expect(configPage).toHaveAlert('error', { hasText: 'The security token included in the request is invalid' });
  });

  test('valid credentials should return a 200 status code', async ({
    createDataSourceConfigPage,
    readProvisionedDataSource,
    page,
    selectors,
  }) => {
    const { accessKey, secretKey } = await getTestCredentials(readProvisionedDataSource);

    // create a new datasource and navigate to config page
    const configPage = await createDataSourceConfigPage({ type: 'grafana-iot-sitewise-datasource' });

    // fill in the config form
    await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
    await page.keyboard.press('Enter');
    await page.getByLabel('Access Key ID').fill(accessKey);
    await page.getByLabel('Secret Access Key').fill(secretKey);
    await page.getByRole('combobox', { name: 'Default Region' }).click();
    await configPage.getByGrafanaSelector(selectors.components.Select.option).getByText('us-east-1').click();
    await expect(configPage.saveAndTest()).toBeOK();
  });
});

async function getTestCredentials(
  readProvisionedDataSource: <T = {}, S = {}>(args: ReadProvisionedDataSourceArgs) => Promise<DataSourceSettings<T, S>>
) {
  // get access key from env (in ci) or from provisioning repo (if running e2e test locally)
  let accessKey = '';
  let secretKey = '';
  if (process.env.AWS_ACCESS_KEY && process.env.AWS_SECRET_KEY) {
    accessKey = process.env.AWS_ACCESS_KEY;
    secretKey = process.env.AWS_SECRET_KEY;
  } else {
    try {
      const ds = await readProvisionedDataSource<SitewiseOptions, SitewiseSecureJsonData>({
        fileName: 'iot-sitewise.yaml',
      });
      if (!ds.secureJsonData || !ds.secureJsonData.accessKey || !ds.secureJsonData.secretKey) {
        throw new Error('Provisioned datasource does not have valid credentials');
      }
      accessKey = ds.secureJsonData.accessKey;
      secretKey = ds.secureJsonData.secretKey;
    } catch (err) {
      throw new Error(
        'Missing valid credentials for e2e tests. Please provide AWS_ACCESS_KEY and AWS_SECRET_KEY in the environment variables or provision a datasource with valid credentials in the provisioning repo.'
      );
    }
  }
  return { accessKey, secretKey };
}
