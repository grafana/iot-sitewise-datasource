---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/query/
description: Use the AWS IoT SiteWise query editor to build visual and SQL queries against asset properties, aggregates, and time series.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - query editor
  - sql
  - aggregates
  - macros
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Query editor
title: AWS IoT SiteWise query editor
weight: 200
review_date: 2026-07-23
---

# AWS IoT SiteWise query editor

This document explains how to use the AWS IoT SiteWise query editor. The query editor supports a visual builder for the AWS IoT SiteWise APIs, an SQL clause builder, and a raw SQL code editor.

For general information about the Grafana query editor, refer to [Query and transform data](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/).

## Before you begin

Before you build a query, ensure you have:

- Configured the [AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/).
- Verified that your credentials have read access to the AWS IoT SiteWise assets and properties that you want to query.

## Key concepts

If you're new to AWS IoT SiteWise, these terms are used throughout the query editor.

| Term | Description |
| --- | --- |
| **Asset** | A representation of a device, piece of equipment, or process in AWS IoT SiteWise. |
| **Asset model** | A reusable template that defines the properties and hierarchies for a group of assets. |
| **Property** | A measurement, metric, transform, or attribute that belongs to an asset, such as temperature or pressure. |
| **Property alias** | An alternative name for a property, often an OPC-UA path, that you can use instead of an asset and property ID. |
| **Time series** | A stream of timestamped values for a property. |
| **Aggregate** | A computed value over a time interval, such as the average or maximum. |

## Editor modes

Use the mode selector to switch between the three ways to build a query.

| Mode | Description |
| --- | --- |
| **Builder** | A visual editor for the AWS IoT SiteWise APIs. Select a query type, assets, and properties. |
| **Builder (SQL)** | A visual editor that builds an AWS IoT SiteWise SQL query clause by clause. |
| **Code (SQL)** | A raw SQL editor with autocompletion for tables, columns, and macros. |

The **Builder (SQL)** and **Code (SQL)** modes both run an AWS IoT SiteWise `ExecuteQuery` request and disable the client cache.

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor-switch.png" max-width="800px" class="docs-image--no-shadow" caption="Switch between the Builder and SQL editor modes" >}}

## Query types

In **Builder** mode, select one of the following query types.

| Query type | Description |
| --- | --- |
| **Get property value** | Returns the latest value for one or more asset properties. |
| **Get property value history** | Returns the historical values for one or more asset properties over the dashboard time range. |
| **Get property value aggregates** | Returns aggregated values, such as averages, for one or more asset properties. |
| **Get interpolated property values** | Returns interpolated values at a fixed resolution for one or more asset properties. |
| **List assets** | Returns a paginated list of assets, filtered by asset model or hierarchy level. |
| **List asset models** | Returns the asset models in the Region. |
| **List associated assets** | Returns the child or parent assets associated with an asset. |
| **List time series** | Returns the time series (data streams) in the Region. |

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor.png" max-width="800px" class="docs-image--no-shadow" caption="The AWS IoT SiteWise visual query builder" >}}

## Build a visual query

To build a query in **Builder** mode:

1. Select the **AWS IoT SiteWise** data source.
1. Select a **Query type**.
1. Select a **Region**, or use **Default** to use the data source's default Region.
1. Select the assets and properties, or enter a property alias.
1. Set any query-specific options, such as aggregates or resolution.

### Select assets and properties

Property queries need either an asset and property, or a property alias.

| Field | Description |
| --- | --- |
| **Property Alias** | One or more property aliases, such as an OPC-UA path. Use this field instead of selecting an asset and property. |
| **Asset** | One or more assets. Enter an asset ID, an `externalId:` prefixed value, or a template variable. |
| **Property** | One or more properties for the selected assets. Property names support the `*` wildcard. |

Click **Explore** to open the Asset Browser and navigate assets visually.

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/explorer.png" max-width="800px" class="docs-image--no-shadow" caption="The Asset Browser for navigating assets and models" >}}

### Set aggregates

For the **Get property value aggregates** query type, select one or more aggregate functions. You can display multiple aggregates for a single property.

| Aggregate | Description |
| --- | --- |
| **Average** | The average of the values in the interval. Numeric properties only. |
| **Count** | The number of values in the interval. |
| **Max** | The maximum value in the interval. Numeric properties only. |
| **Min** | The minimum value in the interval. Numeric properties only. |
| **Sum** | The sum of the values in the interval. Numeric properties only. |
| **Stddev** | The standard deviation of the values in the interval. Numeric properties only. |

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor2.png" max-width="800px" class="docs-image--no-shadow" caption="Multiple aggregates shown for a single property" >}}

### Set the resolution

For aggregate and interpolated queries, select a **Resolution** to control the interval that AWS IoT SiteWise uses to compute values. Select **Auto** to let Grafana choose a resolution based on the panel width and time range. Resolution options support template variables.

### Set query options

Use the **Query options** section to refine the results.

| Option | Description |
| --- | --- |
| **Quality** | Filter values by quality: `GOOD`, `BAD`, or `UNCERTAIN`. The default is `GOOD`. |
| **Time** | The time ordering of the results: `ASCENDING` or `DESCENDING`. |
| **Format** | The response format: **Table** or **Time series**. |
| **Expand Time Range** | Include the last value before the time range and the next value after it. Applies to history and aggregate queries. |
| **Format L4E Anomaly Result** | Parse Amazon Lookout for Equipment anomaly results into separate fields. Enabled by default. |
| **Client cache** | Cache query results older than 15 minutes for relative time ranges. Enabled by default. |

### Query-specific fields

Some query types have additional fields.

| Query type | Fields |
| --- | --- |
| **List assets** | **Model ID** to filter by asset model, and **Filter** set to **Top Level** or **All**. The **All** filter requires a model ID. |
| **List associated assets** | **Asset Hierarchy** to select a hierarchy, **Parent** to list parent assets, or **All** to list all child hierarchies. |
| **List time series** | **Time Series Type** set to `ALL`, `ASSOCIATED`, or `DISASSOCIATED`, plus **Alias Prefix** or **Asset Id**. |

## Explore assets with the Asset Browser

The Asset Browser helps you find assets and properties without knowing their IDs. Click **Explore** next to the **Asset** field to open it.

The browser has two tabs.

- **Hierarchy:** Navigate the asset tree. View parent assets, the selected asset, and child hierarchies, and search within each hierarchy.
- **By Model:** Select an asset model, then search and select from the assets that use that model.

## Build an SQL query

AWS IoT SiteWise supports an SQL dialect through the `ExecuteQuery` API. Use **Builder (SQL)** mode to construct a query clause by clause, or **Code (SQL)** mode to write raw SQL.

### SQL builder clauses

In **Builder (SQL)** mode, construct the query with the following clauses. A query must include at least a `SELECT` and a `FROM` clause to be valid. A preview of the generated SQL appears as you build the query.

| Clause | Description |
| --- | --- |
| **From** | The view to query. Select one of the AWS IoT SiteWise views. |
| **Select** | The columns to return. Choose columns and optional aggregation functions, or select all columns. |
| **Where** | Filter conditions. Combine conditions with `AND` or `OR`. |
| **Group By** | The columns to group by. |
| **Having** | Filter conditions on aggregated results. Available when a `Group By` clause is set. |
| **Order By** | The column and direction, either **Ascending** or **Descending**. |
| **Limit** | The maximum number of rows to return. The default is `100`. |

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/query-builder.png" max-width="800px" class="docs-image--no-shadow" caption="The AWS IoT SiteWise SQL query builder" >}}

### Available views

The `FROM` clause supports the following AWS IoT SiteWise views.

| View | Description |
| --- | --- |
| `asset` | Asset metadata, such as name, description, model, and parent asset. |
| `asset_property` | Property metadata, such as name, alias, data type, and attributes. |
| `raw_time_series` | Raw timestamped property values. |
| `latest_value_time_series` | The latest value for each time series. |
| `precomputed_aggregates` | Precomputed aggregate values, such as average, sum, and count, at a resolution. |

{{< figure src="https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/query-builder2.png" max-width="800px" class="docs-image--no-shadow" caption="A data preview for an SQL query" >}}

### Write raw SQL

In **Code (SQL)** mode, write the query directly. The editor provides autocompletion for views, columns, macros, and template variables. Select the **AWS Region** from the editor header.

The default query is:

```sql
select $__selectAll from raw_time_series where $__timeFilter(event_timestamp)
```

## Macros

Use macros in SQL queries to reference the dashboard time range and other dynamic values. Grafana replaces each macro before sending the query to AWS IoT SiteWise.

| Macro | Description |
| --- | --- |
| `$__selectAll` | Expands to all columns of the current `FROM` view. |
| `$__timeFrom` | The start of the dashboard time range as a timestamp. |
| `$__timeTo` | The end of the dashboard time range as a timestamp. |
| `$__timeFilter(column)` | Filters the specified column to the dashboard time range. |
| `$__autoResolution()` | Returns a resolution based on the panel interval when querying `precomputed_aggregates`. |

## Format Amazon Lookout for Equipment results

AWS IoT SiteWise can store Amazon Lookout for Equipment (L4E) anomaly detection results as a composite model property named `AWS/L4E_ANOMALY_RESULT`. When you enable **Format L4E Anomaly Result** on a **Get property value** or **Get property value history** query, Grafana parses the JSON result into separate fields, including:

- `anomaly_score`: The anomaly score for the prediction.
- `prediction_reason`: The reason for the prediction.
- A diagnostic contribution field for each property that contributes to the anomaly.

The original JSON value is retained in the `AWS/L4E_ANOMALY_RESULT` column.

## Query examples

Use the following examples as starting points.

Return the raw values for a property over the dashboard time range:

```sql
select event_timestamp, double_value
from raw_time_series
where $__timeFilter(event_timestamp)
order by event_timestamp asc
```

Return hourly averages from precomputed aggregates:

```sql
select event_timestamp, average_value
from precomputed_aggregates
where resolution = '1h' and $__timeFilter(event_timestamp)
```

List the assets that use a specific asset model:

```sql
select asset_id, asset_name
from asset
where asset_model_id = '<YOUR_ASSET_MODEL_ID>'
```

## Use cases

The query editor supports a range of industrial monitoring scenarios, for example:

- **Equipment monitoring:** Visualize property values, such as temperature or pressure, across a fleet of assets.
- **Aggregated trends:** Chart averages, maximums, and standard deviations over time to spot performance changes.
- **Anomaly detection:** Surface Amazon Lookout for Equipment anomaly scores alongside the properties that contribute to them.

## Next steps

- [Template variables](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/template-variables/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshooting/)
