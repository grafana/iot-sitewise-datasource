# Change Log

All notable changes to this project will be documented in this file.

## v1.11.1
- Revert "Replace deprecated setVariableQueryEditor with CustomVariableSupport" in [#229](https://github.com/grafana/iot-sitewise-datasource/pull/229)

## v1.11.0

- Update backend grafana-aws-sdk to v0.19.1 to add `il-central-1` to the opt-in region list
- Update frontend grafana/aws-sdk to v0.1.2 to limit `grafana_assume_role` only to enabled datasources

## v1.10.3

- Update grafana/aws-sdk-react dependency https://github.com/grafana/iot-sitewise-datasource/pull/20

## v1.10.2

- Fix: Fix scoped variables replacement in assetids such as repeat panels by @ahom https://github.com/grafana/iot-sitewise-datasource/pull/205

## v1.10.1

- Fix: Property aggregate queries returning duplicated data https://github.com/grafana/iot-sitewise-datasource/pull/203
- Fix: Query with expression only returns partial data https://github.com/grafana/iot-sitewise-datasource/pull/206

## v1.10.0

- Include propertyName in data frame name for 'raw' queries https://github.com/grafana/iot-sitewise-datasource/pull/199

## v1.9.2

- Fetch asset property info if asset id and property id are available https://github.com/grafana/iot-sitewise-datasource/pull/192
- Handle expression queries with more than 250 data points https://github.com/grafana/iot-sitewise-datasource/pull/194

## v1.9.1

- Replace deprecated setVariableQueryEditor with CustomVariableSupport
  in https://github.com/grafana/iot-sitewise-datasource/pull/184

## v1.9.0

- Add ability to perform property queries by only specifying a property
  alias ([#179](https://github.com/grafana/iot-sitewise-datasource/pull/179))

## 1.8.1

- Update grafana-aws-sdk version to include new region in opt-in region
  list https://github.com/grafana/grafana-aws-sdk/pull/80
- Security: Upgrade Go in build process to 1.20.4
- Update grafana-plugin-sdk-go version to 0.161.0 to avoid a potential http header problem. https://github.com/grafana/athena-datasource/issues/233

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
