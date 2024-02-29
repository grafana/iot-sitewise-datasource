import { test as setup } from '@grafana/plugin-e2e';
import { SITE_WISE_DATA_SOURCE_CONFIG } from './constants';

setup('data source', async ({ createDataSource }) => {
  await createDataSource(SITE_WISE_DATA_SOURCE_CONFIG);
});
