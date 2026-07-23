---
aliases:
  - /docs/plugins/grafana-iot-sitewise-datasource/latest/troubleshoot/
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

- **Save & test** fails for a Linux or LDAP Edge connection.

**Solution:**

Provide a valid PEM certificate in the **SSL Certificate** field. The value begins with `-----BEGIN CERTIFICATE-----`.

### "missing edge auth user" or "missing edge auth password"

**Symptoms:**

- **Save & test** fails with a message about a missing Edge authentication user or password.

**Solution:**

For Linux and LDAP authentication modes, provide both the **Username** and **Password** for the gateway's local authentication proxy.

<!-- vale Grafana.Headings = YES -->
<!-- vale Grafana.Gerunds = YES -->

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
