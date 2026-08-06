---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/setup/
  - /docs/plugins/grafana-iot-sitewise-datasource/setup/
description: Configure the AWS IoT SiteWise data source in Grafana, including authentication, SiteWise Edge, private data source connect, provisioning, and Terraform.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - configure
  - authentication
  - provisioning
  - terraform
  - private data source connect
  - iam
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Configure
title: Configure the AWS IoT SiteWise data source
weight: 100
review_date: 2026-07-23
---

# Configure the AWS IoT SiteWise data source

This document explains how to configure the AWS IoT SiteWise data source and provides links to related documentation.

## Before you begin

Before you configure the data source, ensure you have:

- **Grafana permissions:** The `Organization administrator` role. Only organization administrators can add data sources.
- **An AWS account** with AWS IoT SiteWise enabled in at least one Region, or a configured SiteWise Edge gateway.
- **AWS credentials or an IAM identity** with read access to AWS IoT SiteWise. At a minimum, grant `iotsitewise:List*`, `iotsitewise:Describe*`, and `iotsitewise:Get*`. To use the SQL query editor, also grant `iotsitewise:ExecuteQuery`.

## Key concepts

If you're new to AWS, these terms are used throughout the configuration.

| Term | Description |
| --- | --- |
| **IAM policy** | A JSON document attached to an identity that grants AWS API permissions. |
| **Assume role** | An AWS mechanism that lets one identity take on temporary credentials for another IAM role, often used for cross-account access. |
| **External ID** | An optional identifier that a role in another account requires when you assume it, which adds a layer of protection for cross-account access. |
| **Region** | The AWS Region, such as `us-east-1`, where your AWS IoT SiteWise data is stored. |
| **Endpoint** | The service URL that the data source connects to. Set a custom endpoint for private networks or SiteWise Edge. |
| **AWS IoT SiteWise Edge** | An on-premises gateway that runs AWS IoT SiteWise locally on your own hardware. |

## Add the data source

To add the AWS IoT SiteWise data source:

1. Click **Connections** in the left-side menu.
1. Click **Add new connection**.
1. Type `AWS IoT SiteWise` in the search bar.
1. Select **AWS IoT SiteWise** from the search results.
1. Click **Add new data source**.

## Authentication

The AWS IoT SiteWise data source uses the same authentication system as the other AWS data sources in Grafana. In the data source settings, these options appear under **Connection Details**. Choose the method that matches your deployment.

| Method | Best for | Grafana Cloud | Server configuration required |
| --- | --- | --- | --- |
| **AWS SDK Default** | Self-managed Grafana with credentials in the environment or default chain | No | Yes |
| **Workspace IAM Role** | Grafana running on Amazon EC2 with an instance profile | No | Yes |
| **Grafana Assume Role** | Grafana Cloud users who want temporary credentials | Yes | No |
| **Access & secret key** | Any deployment | Yes | No |
| **Credentials file** | Self-managed Grafana with an AWS credentials file | No | Yes |

Select the method from the **Authentication Provider** drop-down. The available options depend on the providers your Grafana administrator allows.

### AWS SDK Default

This method uses the default AWS SDK credential chain, which resolves credentials from environment variables, shared configuration, or the container or instance role. Use it when Grafana runs on AWS infrastructure that already has AWS credentials available.

### Workspace IAM Role

This method uses the IAM role attached to the Amazon EC2 instance that runs Grafana. Use it when Grafana runs on Amazon EC2 and you attach an instance profile with access to AWS IoT SiteWise.

### Grafana Assume Role

This method lets Grafana assume an IAM role for temporary credentials without storing long-lived keys. It's available in Grafana Cloud when your administrator enables it. Select **Grafana Assume Role** as the **Authentication Provider**, then expand **How to create an IAM role for grafana to assume** for the trust policy and steps. Create the role so it trusts the Grafana account, and enter its Amazon Resource Name in the **Assume Role ARN** field.

### Access & secret key

This method uses a long-lived AWS access key ID and secret access key. Provide the following values.

| Setting | Description |
| --- | --- |
| **Access Key ID** | The AWS access key ID for an IAM user with access to AWS IoT SiteWise. |
| **Secret Access Key** | The AWS secret access key that pairs with the access key ID. Grafana stores this value as a secure setting. |

### Credentials file

This method reads credentials from an AWS shared credentials file on the Grafana server, typically at `~/.aws/credentials`. Provide the profile name.

| Setting | Description |
| --- | --- |
| **Credentials Profile Name** | The profile name in the shared credentials file. Leave blank to use the default profile. |

### Assume Role

The **Assume Role** section appears for every authentication provider. Grafana can assume an IAM role instead of using the selected provider's credentials directly.

| Setting | Description |
| --- | --- |
| **Assume Role ARN** | Optional. The Amazon Resource Name of an IAM role to assume. Grafana uses the selected authentication provider to assume this role instead of using the credentials directly. |
| **External ID** | Optional. The external ID required by a role in another account. This field doesn't appear when you use Grafana Assume Role. |

### Additional settings

Set the following options in the **Additional Settings** section.

| Setting | Description |
| --- | --- |
| **Endpoint** | Optional. A custom endpoint for the AWS IoT SiteWise service, in the form `https://{service}.{region}.amazonaws.com`. Required for SiteWise Edge. This field doesn't appear when you use Grafana Assume Role. |
| **Default Region** | The AWS Region that queries use by default, such as `us-west-2` for US West (Oregon). Select **Edge** to connect to a SiteWise Edge gateway. |

## Connect to private networks

AWS IoT SiteWise is usually reachable over public AWS endpoints. If you need to reach a private endpoint, such as a SiteWise Edge gateway or a service behind a private network that isn't publicly accessible, use one of the following options.

### Private data source connect (Grafana Cloud)

If you use Grafana Cloud, use [Private data source connect (PDC)](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) to securely connect Grafana Cloud to a private network without opening the network to inbound traffic. PDC establishes an outbound connection from an agent that runs inside your network.

When PDC is enabled for your Grafana Cloud instance, a **Private data source connect** section appears in the data source settings, which lets you query data that lives within a secured network without opening the network to inbound traffic from Grafana Cloud. Select the **Private data source connect network** where your endpoint is available.

{{< admonition type="note" >}}
Private data source connect is available exclusively in Grafana Cloud. On self-managed Grafana, use the secure Socks proxy, a VPN, or AWS PrivateLink to reach private endpoints.
{{< /admonition >}}

For setup instructions, refer to [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/).

### Secure Socks proxy (self-managed Grafana)

On self-managed Grafana, if your administrator enables the secure Socks proxy, a **Secure Socks Proxy** section appears in the data source configuration so you can route requests through the proxy. This option requires Grafana version 10.0.0 or later. To enable it, configure the `[secure_socks_datasource_proxy]` section in the Grafana configuration file. For more information, refer to [Configure Grafana](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/setup-grafana/configure-grafana/#secure_socks_datasource_proxy).

## Configure SiteWise Edge

SiteWise Edge lets you run AWS IoT SiteWise on an on-premises gateway. To connect to a gateway, select **Edge** as the **Default Region**. An explicit endpoint is required for Edge connections.

Select the **Authentication Mode** for the gateway.

| Mode | Description |
| --- | --- |
| **Standard** | Uses the AWS authentication provider that you configured for the data source. |
| **Linux** | Uses Linux-based authentication against the gateway's local authentication proxy. |
| **LDAP** | Uses LDAP-based authentication against the gateway's local authentication proxy. |

For Linux and LDAP modes, provide the following values.

| Setting | Description |
| --- | --- |
| **Username** | The username sent to the local authentication proxy. |
| **Password** | The password sent to the local authentication proxy. Grafana stores this value as a secure setting. |
| **SSL Certificate** | The PEM certificate used for SSL-enabled authentication. The value begins with `-----BEGIN CERTIFICATE-----`. Grafana stores this value as a secure setting. |

To replace a saved certificate, click **Reset** and enter a new certificate.

## Verify the connection

Click **Save & test** to verify the configuration. On success, Grafana returns `OK`. The data source runs a `ListAssetModels` request against AWS IoT SiteWise to confirm that the credentials and Region are valid.

If the test fails, refer to [Troubleshooting](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshooting/).

## Provision the data source

You can define the data source in YAML files as part of the Grafana provisioning system. For more information, refer to [Provision Grafana](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/provisioning/#data-sources).

### Access and secret key

```yaml
apiVersion: 1

datasources:
  - name: AWS IoT SiteWise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
```

### Credentials file

```yaml
apiVersion: 1

datasources:
  - name: AWS IoT SiteWise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
      profile: default
```

### Assume role

```yaml
apiVersion: 1

datasources:
  - name: AWS IoT SiteWise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
      assumeRoleArn: arn:aws:iam::123456789012:role/grafana-sitewise
      externalId: <YOUR_EXTERNAL_ID>
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
```

### Edge gateway

```yaml
apiVersion: 1

datasources:
  - name: AWS IoT SiteWise Edge
    type: grafana-iot-sitewise-datasource
    jsonData:
      defaultRegion: Edge
      endpoint: https://<YOUR_EDGE_GATEWAY_HOST>
      edgeAuthMode: linux
      edgeAuthUser: <YOUR_EDGE_USERNAME>
    secureJsonData:
      edgeAuthPass: <YOUR_EDGE_PASSWORD>
      cert: |
        -----BEGIN CERTIFICATE-----
        <YOUR_CERTIFICATE>
        -----END CERTIFICATE-----
```

The following table describes the provisioning keys.

| Key | Description |
| --- | --- |
| `authType` | The authentication method: `keys`, `credentials`, `default`, `ec2_iam_role`, or `grafana_assume_role`. |
| `defaultRegion` | The default AWS Region. Set to `Edge` for a SiteWise Edge gateway. |
| `profile` | The credentials file profile name. |
| `assumeRoleArn` | The Amazon Resource Name of an IAM role to assume. |
| `externalId` | The external ID for cross-account role assumption. |
| `endpoint` | A custom service endpoint. Required for Edge. |
| `edgeAuthMode` | The Edge authentication mode: `default`, `linux`, or `ldap`. |
| `edgeAuthUser` | The Edge local proxy username. |
| `accessKey` | The AWS access key ID. Store in `secureJsonData`. |
| `secretKey` | The AWS secret access key. Store in `secureJsonData`. |
| `sessionToken` | An optional session token for temporary credentials. Store in `secureJsonData`. |
| `edgeAuthPass` | The Edge local proxy password. Store in `secureJsonData`. |
| `cert` | The PEM SSL certificate for Edge. Store in `secureJsonData`. |
| `enableSecureSocksProxy` | Set to `true` to route requests through the secure Socks proxy on self-managed Grafana. |

### Private data source connect

To route a provisioned data source through the secure Socks proxy, set `enableSecureSocksProxy` to `true`.

```yaml
apiVersion: 1

datasources:
  - name: AWS IoT SiteWise
    type: grafana-iot-sitewise-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
      enableSecureSocksProxy: true
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
```

## Provision with Terraform

You can use the [Grafana Terraform provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs) to provision the AWS IoT SiteWise data source as code. The following examples use the `grafana_data_source` resource.

### Access and secret key with Terraform

```hcl
resource "grafana_data_source" "sitewise" {
  type = "grafana-iot-sitewise-datasource"
  name = "AWS IoT SiteWise"

  json_data_encoded = jsonencode({
    authType      = "keys"
    defaultRegion = "us-east-1"
  })

  secure_json_data_encoded = jsonencode({
    accessKey = var.aws_access_key
    secretKey = var.aws_secret_key
  })
}
```

To use temporary AWS credentials from STS, add `sessionToken` to `secure_json_data_encoded`:

```hcl
  secure_json_data_encoded = jsonencode({
    accessKey    = var.aws_access_key
    secretKey    = var.aws_secret_key
    sessionToken = var.aws_session_token
  })
```

### Assume role with Terraform

```hcl
resource "grafana_data_source" "sitewise" {
  type = "grafana-iot-sitewise-datasource"
  name = "AWS IoT SiteWise"

  json_data_encoded = jsonencode({
    authType      = "default"
    defaultRegion = "us-east-1"
    assumeRoleArn = "arn:aws:iam::123456789012:role/grafana-sitewise"
    externalId    = var.external_id
  })
}
```

For more information, refer to the [`grafana_data_source` resource](https://registry.terraform.io/providers/grafana/grafana/latest/docs/resources/data_source) in the Grafana Terraform provider documentation.

## Next steps

- [AWS IoT SiteWise query editor](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/query-editor/)
- [Template variables](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/template-variables/)
