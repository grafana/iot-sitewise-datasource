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

After you create a variable, reference it in a query field with the `${variable_name}` syntax. The data source interpolates variables in the following fields:

- **Region:** Switch the Region for a query.
- **Asset and property fields:** Set the asset, property, or property alias dynamically.
- **Model ID:** Filter a **List assets** query by a selected asset model.
- **Resolution:** Change the aggregate or interpolated resolution.
- **SQL queries:** Insert variable values into `WHERE` conditions and other clauses. String values are quoted and escaped automatically.

The **Asset**, **Property**, and **Property Alias** fields accept multi-value variables. When a variable holds more than one value, the query expands to include every selected value.

## Examples

### Chain a model and asset variable

Build a dashboard where the user first selects an asset model, then an asset that uses that model:

1. Create a **List asset models** variable named `model`.
1. Create a **List assets** variable named `asset`. Set the query **Filter** to **All** and the **Model ID** to `${model}`.
1. In your panels, reference `${asset}` in the **Asset** field.

When the user selects a different model, the asset drop-down updates to show only the assets for that model, and the panels update automatically.

### Use a variable in an SQL query

Reference a variable in a `WHERE` condition. The data source quotes and escapes string values, so you don't add quotation marks:

```sql
select event_timestamp, double_value
from raw_time_series
where asset_id = ${asset} and $__timeFilter(event_timestamp)
order by event_timestamp asc
```

## Next steps

- [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshooting/)
