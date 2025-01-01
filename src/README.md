# AWS IoT SiteWise Datasource

This datasource supports reading data from [AWS IoT SiteWise](https://aws.amazon.com/iot-sitewise/) and showing it in a Grafana dashboard.

## Add the data source

1. In the side menu under the **Configuration** link, click on **Data Sources**.
1. Click the **Add data source** button.
1. Select **IoT sitewise** in the **Industrial & IoT** section.

## Authentication

The IoT SiteWise plugin authentication matches the standard Cloudwatch plugin system. See the [grafana cloudwatch documentation](https://grafana.com/docs/grafana/latest/datasources/cloudwatch/#authentication) for authentication options and setup.

Once authentication is configured, click "Save and Test" to verify the service is working. Once this is configured, you can specify default values for the configuration.

## Querying data

Users can choose between the UI driven Query builder and the Query code editor, that uses [Iot Sitewise query language](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/sql.html) queries to fetch data. The Query builder offers a guided, less technical interface, whereas the Query code editor is more technical but also more powerful.

### Query builder

Use the "query type" selector to pick an appropriate query.

![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor.png)

Click on the "Explore" button to open an asset/model navigation interface:

![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/explorer.png)

Multiple aggregations can be shown for a single property:

![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor2.png)

### Query code editor

You can run [Iot Sitewise query language](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/sql.html) queries in the code editor:
![raw-query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor-switch.png)

The query editor supports the following macros:

* $__selectAll - Shortcut to select available fields in the current table: `select $__selectAll from raw_time_series`
* $__rawTimeFrom - Lower limit of the time range as a timestamp: `select $__selectAll from raw_time_series where event_timestamp > $__rawTimeFrom`
* $__rawTimeTo - Upper limit of the time range as a timestamp: `select $__selectAll from raw_time_series where event_timestamp <= $__rawTimeTo`
* $__unixEpochFilter(column) - Filter the specified field according to the time range: `select $__selectAll from raw_time_series where $__unixEpochFilter(event_timestamp)`
* $__resolution - Shortcut to the applicable aggregate resolution based on the panel interval: `select $__selectAll from precomputed_aggregates where $__unixEpochFilter(event_timestamp) and resolution = '$__resolution'`

#### Example queries

The queries below provide a simple introduction to the [Iot Sitewise query language](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/sql.html). See the linked documentation for more details.

**Retrieve all raw events**

```sql
select $__selectAll from raw_time_series where $__unixEpochFilter(event_timestamp)
```

**Retrieve asset and property name along with raw events**

```sql
select r.event_timestamp, a.asset_name, p.property_name, r.double_value
from asset a, asset_property p, raw_time_series r
where $__unixEpochFilter(event_timestamp)
```

## Alerting

Standard grafana alerting is support with this plugin, however note that alert queries may not include template variables.
See the [Alerting](https://grafana.com/docs/grafana/latest/alerting/alerts-overview/) documentation for more on Grafana alerts.

## Configure the data source with provisioning

You can configure data sources using config files with Grafana's provisioning system. You can read more about how it works and all the settings you can set for data sources on the [provisioning docs page](https://grafana.com/docs/grafana/latest/administration/provisioning/).

Here are some provisioning examples for this data source.

### Using a credentials file

If you are using Credentials file authentication type, then you should use a credentials file with a config like this.

```yaml
apiVersion: 1

datasources:
  - name: IoT Sitewise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
```

### Using `accessKey` and `secretKey`

```yaml
apiVersion: 1

datasources:
  - name: IoT Sitewise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: '<your access key>'
      secretKey: '<your secret key>'
```
