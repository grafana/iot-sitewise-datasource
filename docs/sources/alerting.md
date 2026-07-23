---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/alerting/
description: Set up Grafana Alerting with the AWS IoT SiteWise data source.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - aws
  - alerting
  - alerts
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Alerting
title: AWS IoT SiteWise alerting
weight: 360
review_date: 2026-07-23
---

# AWS IoT SiteWise alerting

Use [Grafana Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) with the AWS IoT SiteWise data source to receive notifications when your industrial data meets specific conditions. For example, you can create alerts that trigger when a temperature exceeds a threshold or when an equipment anomaly score rises.

## Before you begin

Before you set up alerting, ensure you have the following prerequisites.

- [Configure the AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/).
- Understand how to use the [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/).
- Understand [Grafana Alerting concepts](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/fundamentals/).

{{< admonition type="note" >}}
Alert rule queries can't include template variables. Use explicit asset IDs, property IDs, or property aliases in queries that back alert rules.
{{< /admonition >}}

## Create an alert rule

To create an alert rule that uses the AWS IoT SiteWise data source:

1. Click **Alerting** in the left-side menu.
1. Click **Alert rules**.
1. Click **New alert rule**.
1. Enter a name for the alert rule.
1. Select the **AWS IoT SiteWise** data source.
1. Build a query that returns the numeric value you want to evaluate.
1. Configure the alert condition by selecting a reducer, such as **Last** or **Mean**, and a threshold.
1. Set the evaluation interval and pending period.
1. Configure notification settings and click **Save rule and exit**.

## Query requirements for alerting

Alert rule queries must return numeric data that Grafana can evaluate against a threshold. Follow these guidelines when you write queries for alert rules.

- **Return a numeric value:** The query must return at least one numeric column that Grafana can evaluate.
- **Use the time series format:** Select the **Time series** format in the query editor. The query needs a time column in ascending order and a numeric value column.
- **Avoid template variables:** Alert queries can't include template variables. Specify assets, properties, and property aliases explicitly.
- **Use macros for SQL time filtering:** In SQL queries, use `$__timeFilter(column)` to filter data to the alert evaluation window.

### Example visual query

To evaluate an asset property against a threshold with the visual builder:

1. Set **Query type** to **Get property value aggregates**.
1. Select the **Asset** and **Property** explicitly. Don't use a template variable.
1. Set **Aggregate** to **Average** and choose a fixed **Resolution**, such as **Hour**.
1. Set **Format** to **Time series**.

### Example SQL query

The following query returns the hourly average value for a specific property from precomputed aggregates, which you can evaluate against a threshold. Filter to a single property so the alert evaluates one series:

```sql
select event_timestamp, average_value
from precomputed_aggregates
where property_alias = '<YOUR_PROPERTY_ALIAS>'
  and resolution = '1h'
  and $__timeFilter(event_timestamp)
order by event_timestamp asc
```

## Performance considerations

Keep the following considerations in mind when you configure alert rules for AWS IoT SiteWise.

- **Evaluation interval:** Set an evaluation interval that accounts for query time and the resolution of your data.
- **No pagination for alert queries:** When a property value history or aggregate query runs for an alert rule, it fetches all results without pagination so the rule can evaluate synchronously. A wide time range or a fine resolution can pull a large volume of data. Use a fixed, coarse resolution and a bounded evaluation window to limit the data returned.
- **API limits:** Each evaluation runs a query against AWS IoT SiteWise. Frequent evaluations across many assets can lead to API throttling. Widen the evaluation interval or increase the aggregate resolution to reduce API calls.
- **Query caching:** Enable [query caching](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/#query-caching) to reduce the number of queries sent to AWS IoT SiteWise. Query caching is available in Grafana Enterprise and Grafana Cloud.

## Next steps

- [Configure notification policies](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/configure-notifications/)
- [Create contact points](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/configure-notifications/manage-contact-points/)
- [Grafana Alerting documentation](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/)
