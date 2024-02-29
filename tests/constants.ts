import { type CreateDataSourceArgs } from '@grafana/plugin-e2e';

export const SITE_WISE_DATA_SOURCE_CONFIG = {
  type: 'grafana-iot-sitewise-datasource',
  name: 'IoT SiteWise',
} satisfies CreateDataSourceArgs;
