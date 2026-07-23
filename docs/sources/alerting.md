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

Use the **Get property value aggregates** query type with the **Average** aggregate and the **Time series** format to evaluate an asset property against a threshold. Select the asset and property explicitly rather than using a template variable.

### Example SQL query

The following query returns the average value for a property from precomputed aggregates, which you can evaluate against a threshold.

```sql
select event_timestamp, average_value
from precomputed_aggregates
where resolution = '1h' and $__timeFilter(event_timestamp)
order by event_timestamp asc
```

## Performance considerations

Keep the following considerations in mind when you configure alert rules for AWS IoT SiteWise.

- **Evaluation interval:** Set an evaluation interval that accounts for query time and the resolution of your data.
- **API limits:** Each evaluation runs a query against AWS IoT SiteWise. Frequent evaluations across many assets can lead to API throttling. Widen the evaluation interval or increase the aggregate resolution to reduce API calls.
- **Query caching:** Enable [query caching](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/#query-caching) to reduce the number of queries sent to AWS IoT SiteWise. Query caching is available in Grafana Enterprise and Grafana Cloud.

## Next steps

- [Configure notification policies](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/configure-notifications/)
- [Create contact points](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/configure-notifications/manage-contact-points/)
- [Grafana Alerting documentation](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/)
