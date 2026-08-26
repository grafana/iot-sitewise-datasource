---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/
description: Use the AWS IoT SiteWise data source to query and visualize industrial equipment data from AWS IoT SiteWise in Grafana.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - iot
  - aws
  - amazon
  - data source
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: AWS IoT SiteWise
title: AWS IoT SiteWise data source
weight: 10
review_date: 2026-07-23
---

# AWS IoT SiteWise data source

The AWS IoT SiteWise data source lets you query and visualize data from [AWS IoT SiteWise](https://aws.amazon.com/iot-sitewise/) in Grafana. AWS IoT SiteWise is a managed service that collects, stores, organizes, and monitors data from industrial equipment at scale, so you can build dashboards for asset properties, aggregates, and time series without moving your data.

{{< admonition type="warning" >}}
Use Grafana version 10.4.0 or later with the AWS IoT SiteWise data source. Grafana instances earlier than 10.4.0 can't use AWS IoT SiteWise data source versions later than 1.25.2.
{{< /admonition >}}

## Supported features

The following table lists the features available with this data source.

| Feature | Supported |
| --- | --- |
| Metrics | Yes |
| Logs | No |
| Traces | No |
| Alerting | Yes |
| Annotations | Yes |
| Template variables | Yes |

## Requirements

To use the AWS IoT SiteWise data source, you need:

- A Grafana instance running version 10.4.0 or later.
- An AWS account with AWS IoT SiteWise enabled in at least one Region, or a configured SiteWise Edge gateway.
- AWS credentials or an IAM identity with read access to AWS IoT SiteWise.

## Get started

The following pages help you get started with the AWS IoT SiteWise data source.

- [Configure the AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/)
- [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/)
- [Template variables](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/template-variables/)
- [Annotations](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/annotations/)
- [Alerting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/alerting/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshooting/)

## Additional features

After you configure the data source, you can use the following Grafana features.

- Use [Explore](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/explore/) to query data without building a dashboard.
- Add [Transformations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/transform-data/) to manipulate query results.
- Set up [Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) rules to get notified when data changes.
- Configure and use [Template variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/) to build dynamic dashboards.

## Plugin updates

Always ensure that your plugin version is up-to-date so you have access to all current features and improvements. Navigate to **Plugins and data** > **Plugins** to check for updates. Grafana recommends upgrading to the latest Grafana version, and this applies to plugins as well.

{{< admonition type="note" >}}
Plugins are automatically updated in Grafana Cloud.
{{< /admonition >}}

## Related resources

- [AWS IoT SiteWise documentation](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/what-is-sitewise.html)
- [AWS IoT SiteWise data source plugin GitHub repository](https://github.com/grafana/iot-sitewise-datasource)
- [Grafana community forum](https://community.grafana.com/)
