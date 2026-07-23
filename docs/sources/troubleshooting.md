---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshoot/
  - /docs/plugins/grafana-iot-sitewise-datasource/troubleshoot/
description: Troubleshoot common issues with the AWS IoT SiteWise data source in Grafana.
keywords:
  - grafana
  - aws iot sitewise
  - sitewise
  - troubleshooting
  - errors
  - authentication
  - query
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Troubleshooting
title: Troubleshoot AWS IoT SiteWise data source issues
weight: 400
review_date: 2026-07-23
---

# Troubleshoot AWS IoT SiteWise data source issues

This document provides solutions to common issues you might encounter when you configure or use the AWS IoT SiteWise data source. For configuration instructions, refer to [Configure the AWS IoT SiteWise data source](https://grafana.com/docs/plugins/grafana-iot-sitewise-datasource/latest/configure/).

## Plugin version and compatibility issues

Many issues that appear after a Grafana or plugin upgrade come from a version mismatch between the plugin and your Grafana instance.

### Check the plugin version

To find the installed plugin version:

1. Navigate to **Administration** > **Plugins and data** > **Plugins**.
1. Search for **AWS IoT SiteWise** and select the plugin.
1. Read the installed version on the plugin page. The page also lists the available versions and the Grafana version each one supports.

### Verify version compatibility

The plugin requires a minimum Grafana version. The current plugin releases require Grafana 10.4.0 or later. On Grafana instances earlier than 10.4.0, use plugin version 1.25.2, which is the last version that supports earlier Grafana releases.

| Grafana version | Supported plugin version |
| --- | --- |
| 10.4.0 and later | The latest plugin version |
| Earlier than 10.4.0 | 1.25.2 |

If your Grafana instance is earlier than 10.4.0, pin the plugin to version 1.25.2. To use a later plugin version, upgrade Grafana to 10.4.0 or later first. Always confirm the exact minimum Grafana version on the plugin page before you upgrade.

### Errors after a Grafana or plugin upgrade

**Symptoms:**

- The data source configuration page or query editor fails to load after an upgrade.
- The browser console shows a JavaScript error, such as `TypeError: Cannot read properties of undefined`.
- Panels that previously worked return errors or stop rendering.

**Solutions:**

1. Compare the plugin version with your Grafana version against the compatibility table.
1. If the plugin is too new for your Grafana version, downgrade the plugin to a compatible version, or upgrade Grafana.
1. If the plugin is too old for your Grafana version, update the plugin to the latest compatible version.
1. Clear the browser cache and reload after you change versions.

### Pin or test a plugin version on self-managed Grafana

On self-managed Grafana, you control which plugin version is installed. To install a specific version with `grafana-cli`:

```bash
grafana-cli plugins install grafana-iot-sitewise-datasource 1.25.2
```

To pin the version in a container image, set the version in the `GF_INSTALL_PLUGINS` environment variable:

```bash
GF_INSTALL_PLUGINS=grafana-iot-sitewise-datasource@1.25.2
```

Test plugin and Grafana upgrades in a non-production instance before you upgrade production. This lets you confirm compatibility without affecting live dashboards.

### Grafana Cloud upgrades

In Grafana Cloud, plugins update automatically and you can't self-service a rollback to a previous plugin version.

{{< admonition type="note" >}}
To reduce the risk of an upgrade breaking a dashboard, test changes in a separate Grafana instance, such as a local OSS instance, before you rely on them in Grafana Cloud. If a Grafana Cloud upgrade breaks the plugin and you can't resolve it, contact [Grafana Support](https://grafana.com/help/).
{{< /admonition >}}

## Authentication errors

These errors occur when credentials are invalid, missing, or don't have the required permissions.

### "The security token included in the request is invalid"

**Symptoms:**

- **Save & test** fails with a security token error.
- Queries return authorization errors.

**Possible causes and solutions:**

| Cause | Solution |
| --- | --- |
| Invalid or mistyped credentials | Verify the access key ID and secret access key in the AWS console. Regenerate them if necessary. |
| Expired temporary credentials | Create new credentials and update the data source configuration. |
| Missing session token | For temporary credentials, provide a session token. |
| Wrong Region | Verify that the **Default Region** matches where your AWS IoT SiteWise data is stored. |

### "Access denied" or "not authorized"

**Symptoms:**

- Queries fail with access denied messages.
- Assets, models, or properties don't load in drop-down menus.

**Solutions:**

1. Confirm that the IAM identity has `iotsitewise:List*`, `iotsitewise:Describe*`, and `iotsitewise:Get*` permissions.
1. To use the SQL editor, grant `iotsitewise:ExecuteQuery`.
1. For cross-account access, verify the assume role Amazon Resource Name and external ID.

## Connection errors

These errors occur when Grafana can't reach the AWS IoT SiteWise endpoints.

### Connection refused or timeout errors

**Symptoms:**

- The data source test times out.
- Queries fail with network errors.

**Solutions:**

1. Verify network connectivity from the Grafana server to the AWS IoT SiteWise endpoints.
1. Check that firewall rules allow outbound HTTPS on port `443`.
1. If you set a custom endpoint, verify that the URL is correct and reachable.
1. For Grafana Cloud accessing private resources, configure [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/).

<!-- vale Grafana.Headings = NO -->
<!-- vale Grafana.Gerunds = NO -->

## SiteWise Edge errors

These errors occur when you configure a SiteWise Edge gateway connection.

### "edge region requires an explicit endpoint"

**Symptoms:**

- **Save & test** fails after you select **Edge** as the Region.

**Solution:**

Set the **Endpoint** field to the URL of your SiteWise Edge gateway.

### "edge region requires an SSL certificate"

**Symptoms:**

- **Save & test** fails for an Edge connection, in any authentication mode.

**Solution:**

Provide a valid PEM certificate in the **SSL Certificate** field. An SSL certificate is required for all Edge connections, including the **Standard** authentication mode. The value begins with `-----BEGIN CERTIFICATE-----`.

### "missing edge auth user" or "missing edge auth password"

**Symptoms:**

- **Save & test** fails with a message about a missing Edge authentication user or password.

**Solution:**

For the **Linux** and **LDAP** authentication modes, provide both the **Username** and **Password** for the gateway's local authentication proxy. The **Standard** mode uses the authentication provider configured above and doesn't require these fields.

<!-- vale Grafana.Headings = YES -->
<!-- vale Grafana.Gerunds = YES -->

## Query errors

These errors occur when you run queries against the data source.

### No data or empty results

**Symptoms:**

- A query runs without error but returns no data.
- Panels show a "No data" message.

**Possible causes and solutions:**

| Cause | Solution |
| --- | --- |
| Time range doesn't contain data | Expand the dashboard time range or verify that data exists in AWS IoT SiteWise. |
| Wrong asset or property selected | Verify that you selected the correct asset and property, or the correct property alias. |
| Quality filter excludes values | Change the **Quality** option, because the default `GOOD` filter can exclude values recorded with other qualities. |
| Permissions issue | Verify that the identity has read access to the specific asset and property. |

### API throttling or rate limit errors

**Symptoms:**

- Queries fail intermittently with throttling or "rate exceeded" errors.
- Panels fail to load when a dashboard has many queries.

**Solutions:**

1. Reduce the dashboard refresh frequency.
1. Increase the aggregate resolution to reduce the number of API calls.
1. Keep the **Client cache** option enabled to reuse results for relative time ranges.
1. Request a quota increase from AWS.

## Field name issues

Panel overrides and transformations can break when the field names in a query result change. The data source builds field names from the names it reads from AWS IoT SiteWise: the value field is typically named after the property, and the series after the asset. The plugin resolves these names in the query's Region.

### Why field names change

| Cause | Explanation |
| --- | --- |
| Asset or property renamed in AWS | The field and series names follow the asset and property names in AWS IoT SiteWise, so a rename in AWS changes the names in Grafana. |
| Query type changed | Aggregate queries add a separate field for each selected aggregate, named `avg`, `min`, `max`, `sum`, `count`, or `stddev`, whereas **Get property value** and **Get property value history** queries name the value field after the property. |
| Default Region mismatch | The plugin resolves asset and property names in the query's Region. If a query relies on the default Region and that Region changes, names can differ or fail to resolve. |
| Stale cached results | The client cache can return an earlier result that has a different field structure. |

### Use field overrides safely

Hard-coded field names break when the underlying names change. To make overrides and transformations more resilient:

- Use the **Fields with name matching regex** override matcher instead of **Fields with name**, so a partial match still applies when a name changes.
- Match fields by type, such as all numeric fields, when the override doesn't depend on a specific field.
- Add a **Rename by regex** or **Organize fields** transformation to set stable display names before you apply overrides.
- Set the Region explicitly in the query rather than relying on the default Region, so name resolution stays consistent.

### Check the default Region

A changed or mismatched default Region is a common cause of field names that appear as "not found" in overrides.

1. In **Builder** mode, set the query **Region** explicitly instead of using **Default**.
1. Verify that the **Default Region** in the data source configuration matches where your assets are stored.
1. Confirm that the assets and properties you query exist in the selected Region.

### Clear stale results from the client cache

If field names look stale after you rename assets or properties, change tags, or switch the Region, clear the cached results:

1. In **Builder** mode, expand **Query options** and turn off **Client cache**.
1. Run the query again to fetch fresh results.
1. Re-enable **Client cache** if you want to continue reusing results for relative time ranges.

The **Builder (SQL)** and **Code (SQL)** modes don't use the client cache, so you can also run the same query in an SQL mode to bypass cached results.

## Property alias display

When you query by property alias, you might expect the panel to show the full alias, such as an OPC-UA path. Instead, the data source shows the asset property name.

### Property name takes priority over the alias

The data source builds series and field names from the metadata it reads from AWS IoT SiteWise. When a property resolves to an asset property that has a name, the plugin uses the **asset property name** for the series and field names, not the full property alias. This behavior is intentional.

As a result:

- A **Get property value** or **Get property value history** query names the series after the asset and names the value field after the property.
- The full alias isn't shown as the field name, even when you query by alias.

If you need the alias to appear in the panel, add a **Rename by regex** or **Organize fields** transformation to set the display name you want.

{{< admonition type="note" >}}
An option to make the alias display configurable is an open feature request. To follow the discussion or add your use case, refer to the [plugin issue tracker](https://github.com/grafana/iot-sitewise-datasource/issues).
{{< /admonition >}}

### Alias resolution and property data types

Alias resolution can behave differently depending on the property's data type, such as a boolean property compared with a numeric property. For a property queried by alias, the data source infers the data type from the first returned value when the type isn't otherwise defined, so the resolved name or type can differ between data types. If a property queried by alias shows an unexpected or inconsistent name:

1. Verify that the property alias is correct and associated with the expected asset property.
1. Where possible, query by asset and property instead of by alias, so the plugin resolves a consistent asset property name.
1. Add a **Rename by regex** or **Organize fields** transformation to set a stable display name.

## Template variable errors

These errors occur when you use template variables with the data source.

### Variables return no values

**Solutions:**

1. Verify that the data source connection works by testing it in the data source settings.
1. Use a query type that returns variable options: **List asset models**, **List assets**, or **List associated assets**.
1. For the **List assets** query, verify that the **Model ID** is set when you use the **All** filter.
1. Verify that the identity has permission to list the requested resources.

## Enable debug logging

To capture detailed error information for troubleshooting:

1. Set the Grafana log level to `debug` in the configuration file:

   ```ini
   [log]
   level = debug
   ```

1. Review logs in `/var/log/grafana/grafana.log`, or your configured log location.
1. Look for AWS IoT SiteWise entries that include request and response details.
1. Reset the log level to `info` after troubleshooting to avoid excessive log volume.

## Get additional help

If you've tried the solutions in this document and still encounter issues:

1. Check the [Grafana community forums](https://community.grafana.com/) for similar issues.
1. Review the [AWS IoT SiteWise data source plugin GitHub issues](https://github.com/grafana/iot-sitewise-datasource/issues) for known bugs.
1. Consult the [AWS IoT SiteWise documentation](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/what-is-sitewise.html) for service-specific guidance.
1. When you report an issue, include:
   - Your Grafana version and plugin version.
   - Error messages, with sensitive information redacted.
   - Steps to reproduce the issue.
   - Relevant configuration, with credentials redacted.
