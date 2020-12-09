# AWS IoT SiteWise Datasource

This datasource supports reading data from [AWS IoT SiteWise](https://aws.amazon.com/iot-sitewise/) and showing it in a Grafana dashboard.


## Add the data source

1. In the side menu under the **Configuration** link, click on **Data Sources**.
1. Click the **Add data source** button.
1. Select **IoT sitewise** in the **Industrial & IoT** section.


## Authentication

The IoT SiteWise plugin authentication matches the standard Cloudwatch plugin system.  See the [grafana cloudwatch documentation](https://grafana.com/docs/grafana/latest/datasources/cloudwatch/#authentication) for authentication options and setup.


Once authentication is configured, click "Save and Test" to verify the service is working. Once this is configured, you can specify default values for the configuration.


## Query editor

Use the "query type" selector to pick an appropriate query.
![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor.png)

Click on the "Explore" button to open an asset/model navigation interface:
![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/explorer.png)

Multiple aggregations can be showin for a single property:
![query-editor](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/editor2.png)


### Alerting

Standard grafana alertings is support with this plugin, however note that alert queries may not include template variables.
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
    type: datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
```

### Using `accessKey` and `secretKey`

```yaml
apiVersion: 1

datasources:
  - name: IoT Sitewise
    type: datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: '<your access key>'
      secretKey: '<your secret key>'
```
