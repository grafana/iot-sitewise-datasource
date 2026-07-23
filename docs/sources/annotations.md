---
description: Add annotations from AWS IoT SiteWise data in Grafana.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - aws
  - annotations
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Annotations
title: AWS IoT SiteWise annotations
weight: 350
review_date: 2026-07-23
---

# AWS IoT SiteWise annotations

[Annotations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/annotate-visualizations/) let you overlay event information on top of graphs. You can add annotations from AWS IoT SiteWise data through **Dashboard settings** > **Annotations**.

## Before you begin

Before you create an annotation query, ensure you have:

- Configured the [AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/).
- Reviewed the [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/), because annotation queries use the same query editor.

## Create an annotation query

To create an annotation query:

1. Open the dashboard where you want to add annotations.
1. Click the dashboard settings icon.
1. Select **Annotations** from the left-side menu.
1. Click **Add annotation query**.
1. Select the **AWS IoT SiteWise** data source.
1. Build a query that returns the values you want to display as events.

## Annotation columns

Grafana maps query result columns to annotation fields by name. The column names are matched case-insensitively. The following table describes the columns that AWS IoT SiteWise annotation queries use to render annotations.

| Column | Required | Description |
| --- | --- | --- |
| **time** | Yes | The event start time. Alias a timestamp column as `time`, such as `event_timestamp as time`. |
| **timeend** | No | The event end time. When present, Grafana renders a region annotation that spans from **time** to **timeend**. |
| **text** | No | The event description shown in the annotation tooltip. |
| **tags** | No | A comma-separated string used as event tags for filtering annotations. |

## Example annotation queries

The following SQL query creates an annotation for each property value that exceeds a threshold. It aliases `event_timestamp` as `time`, uses the property alias as the annotation text, and uses the value quality as a tag:

```sql
select
  event_timestamp as time,
  property_alias as text,
  quality as tags
from raw_time_series
where $__timeFilter(event_timestamp) and double_value > 95
order by event_timestamp asc
```

You can also use a visual query, such as **Get property value history**, to return timestamped values for an asset property. Grafana renders each returned value as an annotation using the timestamp and value.

## Next steps

- [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/)
- [Grafana annotations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/annotate-visualizations/)
