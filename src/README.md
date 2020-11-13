# AWS IoT SiteWise Datasource

This datasource supports reading data from [AWS IoT SiteWise](https://aws.amazon.com/iot-sitewise/) and showing it in a Grafana dashboard.


## Add the data source

1. In the side menu under the **Configuration** link, click on **Data Sources**.
1. Click the **Add data source** button.
1. Select **IoT sitewise** in the **Industrial & IoT** section.

| Name                     | Description                                                                                                             |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| Name                     | The data source name. This is how you refer to the data source in panels and queries.                                   |
| Auth Provider            | Specify the provider to get credentials.                                                                                |
| Default Region           | Used in query editor to set region. (can be changed on per query basis)                                                 |
| Credentials profile name | Specify the name of the profile to use (if you use `~/.aws/credentials` file), leave blank for default.                 |
| Assume Role Arn          | Specify the ARN of the role to assume.                                                                                  |


## Authentication

In this section we will go through the different type of authentication you can use for IoT sitewise data source.

### Example AWS credentials

If the Auth Provider is `Credentials file`, then Grafana tries to get credentials in the following order:

- Hard-code credentials
- Environment variables (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`)
- Existing default config files
- ~/.aws/credentials
- IAM role for Amazon EC2

Refer to [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) in the AWS documentation for more information.

### AWS credentials file

Create a file at `~/.aws/credentials`. That is the `HOME` path for user running grafana-server.

> **Note:** If the credentials file in the right place, but it is not working, then try moving your .aws file to '/usr/share/grafana/'. Make sure your credentials file has at most 0644 permissions.

Example credential file:

```bash
[default]
aws_access_key_id = <your access key>
aws_secret_access_key = <your access key>
region = us-west-2
```

Once authentication is configured, click "Save and Test" to verify the service is working. Once this is configured, you can specify default values for the configuration.

## Query editor

TODO: query editor docs

![query-editor](https://github.com/grafana/iot-sitewise-datasource/blob/user-readme/docs/editor.png?raw=true)


### Alerting

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
