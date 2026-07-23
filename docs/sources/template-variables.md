---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/variables/
description: Use template variables with the AWS IoT SiteWise data source to build dynamic, reusable dashboards.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - template variables
  - variables
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Template variables
title: AWS IoT SiteWise template variables
weight: 300
review_date: 2026-07-23
---

# AWS IoT SiteWise template variables

Use template variables to build dynamic, reusable dashboards. Instead of hard-coding asset IDs, property IDs, or Regions, you can use variables that users select from drop-down menus at the top of the dashboard.

For general information about Grafana variables, refer to [Templates and variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/).

## Before you begin

Before you create a variable, ensure you have:

- Configured the [AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/).
- Reviewed the [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/), because query variables reuse the same visual builder.

## Supported variable types

The AWS IoT SiteWise data source supports the following Grafana variable types.

| Variable type | Supported |
| --- | --- |
| Query | Yes |
| Custom | Yes |
| Data source | Yes |
| Constant | Yes |
| Text box | Yes |

## Create a query variable

A query variable runs an AWS IoT SiteWise query and uses the results to populate a drop-down menu. Query variables use the same visual query builder as panel queries.

To create a query variable:

1. Navigate to **Dashboard settings** > **Variables**.
1. Click **Add variable**.
1. Select **Query** as the variable type.
1. Select the **AWS IoT SiteWise** data source.
1. Select a query type and complete the query fields.
1. Click **Run query** to preview the values, then click **Apply**.

## Query types that return variable options

Not every query type returns values that populate a variable drop-down. The following query types return an ID and name pair for each result.

| Query type | Returns |
| --- | --- |
| **List asset models** | The asset models in the Region. |
| **List assets** | The assets, filtered by model or hierarchy level. |
| **List associated assets** | The child or parent assets associated with an asset. |

Other query types run successfully but don't map their results to variable options. Use the query types in this table when you build query variables.

## Use variables in queries

After you create a variable, reference it in a query field with the `${variable_name}` syntax. You can use variables in the following places:

- **Region:** Switch the Region for a query.
- **Asset and property fields:** Select assets, properties, or property aliases dynamically.
- **Resolution:** Change the aggregate or interpolated resolution.
- **SQL queries:** Insert variable values into `WHERE` conditions and other clauses.

For example, create a **List assets** variable named `asset`, then reference it in the **Asset** field of a property query as `${asset}`. When a user selects a different asset from the drop-down, the panel updates automatically.

## Next steps

- [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshooting/)
