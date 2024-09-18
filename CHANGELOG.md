# Change Log

All notable changes to this project will be documented in this file.
## 1.24.0

- fix: add check for nil for property value [#352](https://github.com/grafana/iot-sitewise-datasource/pull/352)
- Fix golangci-lint errors [#353](https://github.com/grafana/iot-sitewise-datasource/pull/353)
- fix: migrate asset id on the frontend in the query editor [#350](https://github.com/grafana/iot-sitewise-datasource/pull/350)
- Remove "ANY" as a query quality option [#347](https://github.com/grafana/iot-sitewise-datasource/pull/347)
- Chore: Rename datasource file [#344](https://github.com/grafana/iot-sitewise-datasource/pull/344)
- Add precommit hook [#338](https://github.com/grafana/iot-sitewise-datasource/pull/338)
- Remove unused fields from Get property value query editor [#343](https://github.com/grafana/iot-sitewise-datasource/pull/343)
- Remove ANY option for interpolated property quality [#342](https://github.com/grafana/iot-sitewise-datasource/pull/342)


## 1.23.0

- Feat: add list time series to query editor by @ssjagad in [#336](https://github.com/grafana/iot-sitewise-datasource/pull/336)
- Chore: update dependencies in [#337](https://github.com/grafana/iot-sitewise-datasource/pull/337)
- Migrate to new form styling in config and query editors [#332](https://github.com/grafana/iot-sitewise-datasource/pull/332)

## 1.22.1

- Fix: use ReadAuthSettings to get authSettings in [#333](https://github.com/grafana/iot-sitewise-datasource/pull/333)

## 1.22.0

- Update to use GetSessionWithAuthSettings [#330](https://github.com/grafana/iot-sitewise-datasource/pull/330)

## 1.21.0

- Refactor Paginator [#313](https://github.com/grafana/iot-sitewise-datasource/pull/313)
- Fix e2e tests [#316](https://github.com/grafana/iot-sitewise-datasource/pull/316)
- Added Stalebot [#314](https://github.com/grafana/iot-sitewise-datasource/pull/314)
- Added lint rule [#319](https://github.com/grafana/iot-sitewise-datasource/pull/319)
- Add a frontend cache for relative time range queries [#318](https://github.com/grafana/iot-sitewise-datasource/pull/318)
- Removed unsupported time ordering from interpolated query and fix issue with caching `aggregates` queries [#323](https://github.com/grafana/iot-sitewise-datasource/pull/323)

## 1.20.0

- Perf: Update batch api queries to request maximum number of dependencies in (#310)[https://github.com/grafana/iot-sitewise-datasource/pull/310]

## 1.19.0

- Fetch properties using both AssetId and AssetIds in (#307)[https://github.com/grafana/iot-sitewise-datasource/pull/307]
- Migrate to CustomVariableSupport in (#304)[https://github.com/grafana/iot-sitewise-datasource/pull/304]

## 1.18.0

- Fix fetching asset properties in (#302)[https://github.com/grafana/iot-sitewise-datasource/pull/302]
- Feature: L4E struct support in (#300)[https://github.com/grafana/iot-sitewise-datasource/pull/300]
- Response processing: Add struct data type handling in (#297)[https://github.com/grafana/iot-sitewise-datasource/pull/297]
- Query Editor: Improve menu placement for dropdowns in (#292)[https://github.com/grafana/iot-sitewise-datasource/pull/292]
- Add keywords in (#291)[https://github.com/grafana/iot-sitewise-datasource/pull/291]
- E2E: Add happy path playwright query tests in (#290)[https://github.com/grafana/iot-sitewise-datasource/pull/290]
- Query Editor: Stop running queries on every change in (#274)[https://github.com/grafana/iot-sitewise-datasource/pull/274]

## 1.17.0

- Update grafana-aws-sdk to 0.21.0 and prepare 1.16.2 in (#282)[https://github.com/grafana/iot-sitewise-datasource/pull/282]
- fix: clear selected property when assets are removed (#278)[https://github.com/grafana/iot-sitewise-datasource/pull/278]
- E2E Tests (#285)[https://github.com/grafana/iot-sitewise-datasource/pull/285]
- Add support for composite model properties (#279)[https://github.com/grafana/iot-sitewise-datasource/pull/279]
- Use non-batch APIs at the edge (#281)[https://github.com/grafana/iot-sitewise-datasource/pull/281]

## 1.16.1

- Upgrade aws-sdk-go to v1.49.6 to have access to the `ExecuteQuery` API (#266)
- Fix: Infer data type for disassociated streams for property value queries by alias (#275)

## 1.16.0

- Use query region to get client for queries (#258)
- Feat: implement an "all" option for list associated assets query (#261)

## 1.15.0

- Support multiple assets for interpolated queries in (#256)[https://github.com/grafana/iot-sitewise-datasource/pull/256]

## 1.14.0

- Query and Config editors: Migrate to new form styling under feature toggle in (#244)[https://github.com/grafana/iot-sitewise-datasource/pull/244]

## 1.13.0

- Update dependencies and create-plugin configuration by @idastambuk
  in https://github.com/grafana/iot-sitewise-datasource/pull/243
- Property aggregate processing: Move out ErrorEntries processing from SuccessEntries block by @idastambuk
  in https://github.com/grafana/iot-sitewise-datasource/pull/240
- Bump go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace from 0.37.0 to 0.44.0 by @dependabot
  in https://github.com/grafana/iot-sitewise-datasource/pull/241
- Bump @babel/traverse from 7.17.10 to 7.23.2 by @dependabot
  in https://github.com/grafana/iot-sitewise-datasource/pull/245
- Bump loader-utils from 2.0.2 to 2.0.4 by @dependabot in https://github.com/grafana/iot-sitewise-datasource/pull/248
- Bump semver from 5.7.1 to 5.7.2 by @dependabot in https://github.com/grafana/iot-sitewise-datasource/pull/247
- Bump google.golang.org/grpc from 1.58.2 to 1.58.3 by @dependabot
  in https://github.com/grafana/iot-sitewise-datasource/pull/246
- Upgrade yaml package by @fridgepoet in https://github.com/grafana/iot-sitewise-datasource/pull/249
- Upgrade underscore, debug dependencies by @fridgepoet in https://github.com/grafana/iot-sitewise-datasource/pull/252
- Bump yaml from 2.2.1 to 2.3.4 by @dependabot in https://github.com/grafana/iot-sitewise-datasource/pull/253
- Bump json5 from 2.2.1 to 2.2.3 by @dependabot in https://github.com/grafana/iot-sitewise-datasource/pull/254

**Full Changelog**: https://github.com/grafana/iot-sitewise-datasource/compare/v1.12.1...v1.13.0

## v1.12.1

- Disassociated streams: Hash entryId to fix bug with property aliases longer than 64 characters in [#239](https://github.com/grafana/iot-sitewise-datasource/pull/239)

## v1.12.0

- Query by property alias: Add support for unassociated streams in [#231](https://github.com/grafana/iot-sitewise-datasource/pull/231)

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
