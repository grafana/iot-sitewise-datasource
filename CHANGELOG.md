# Change Log

All notable changes to this project will be documented in this file.

## v1.8.0

- Update backend dependencies

## v1.7.0
- Variables: Fix assetId field variable replacement ([#172](https://github.com/grafana/iot-sitewise-datasource/pull/172))
- Bump coverage to version 0.1.19 ([#173](https://github.com/grafana/iot-sitewise-datasource/pull/173))
- Update aws-sdk-go for the opt-in region list ([#168](https://github.com/grafana/iot-sitewise-datasource/pull/168))
- Modify templates and add workflows for AWS Datasources squad ([#163](https://github.com/grafana/iot-sitewise-datasource/pull/163))
- Migrate to create-plugin (#159)
  ([#159](https://github.com/grafana/iot-sitewise-datasource/pull/159))

## v1.6.0

- Add Batch API support

## v1.5.1

- Add response format selection to time series queries

## v1.5.0

- Renamed last observed value feature to 'Expand Time Range'
- The expand time range toggle now queries for the previous known value before the start of the current time range, and the next known value after the current time range.

## v1.4.1

- Update Grafana AWS SDK dependencies to the latest versions
- Update Grafana dependencies to 8.5.0

## v1.4.0

- Add support for interpolated property value queries
- Add support for last observed value in property value queries
- Switch from long to wide series to support alerting

## v1.3.0

- Add support to define template variables using iot-sitewise datasource queries
- Add dashboard variable support in query editor

## v1.2.6

- Make asset/model descriptions optional.

## v1.2.5

- Fixes issue with asset explorer.
- Adds support for query by property alias.

## v1.2.4

- Add linux/LDAP based authentication for Edge region.

## v1.2.3

- Update `AUTO` aggregation to better select the resolution, and switch to the raw asset property value data when higher than 1m resolution is needed.

## v1.2.2

- Adds resource cache for describe calls in the plugin back-end

## v1.2.1

- Updates shared aws configuration library
- Bumps min version to 7.5

## v1.2.0

- Shares auth configuration with cloudwatch
- Bumps min version to 7.4

## v1.1.0

- Allowing 'Edge' region
- Support nil values in response #82
- Update aws libraries

## v1.0.0

- Initial Release
