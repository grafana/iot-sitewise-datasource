## Developer Guide

## How to build the Sitewise data source plugin locally

## Dependencies

Make sure you have the following dependencies installed first:

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (see [go.mod](../go.mod#L3) for minimum required version)
- [Mage](https://magefile.org/)
- [Node.js (Long Term Support)](https://nodejs.org)
- [Yarn](https://yarnpkg.com)

## Frontend

1. Install dependencies

```BASH
yarn install
```

2. Build plugin in development mode or run in watch mode

```BASH
yarn dev
```

or

```BASH
yarn watch
```

3. Build plugin in production mode

```BASH
yarn build
```

4. Run tests

```BASH
yarn test
```

## Backend

1. Build backend for all platforms

```BASH
mage buildAll
```

2. Run tests

```BASH
mage test
```

### Install

Instructions to install grafana server locally can be found, here:

- [Grafana Server](https://grafana.com/docs/grafana/latest/installation/)

To install the plugin locally, copy the built plugin to the Grafana plugin directory (usually: `/var/lib/grafana/plugins`)

- https://grafana.com/docs/grafana/latest/plugins/installation/

### Docker development setup

1. Create AWS credentials file:

```BASH
cat << EOF >> ~/.aws/credentials
[DEFAULT]

[default]
aws_access_key_id=<your aws access key id>
aws_secret_access_key=<your aws secret access key>
EOF
```

2. Start Grafana docker

```BASH
cd /Workspace/iot-sitewise-datasource
yarn server
```

OR

```BASH
# Run from directory containing iot-sitewise-datasource clone
cd /Workspace/iot-sitewise-datasource
docker run -e GF_DEFAULT_APP_MODE=development -e AWS_SHARED_CREDENTIALS_FILE="/Users/grafana/.aws/credentials" -d -p 3000:3000 -v ~/.aws/:/Users/grafana/.aws/ -v "$(pwd)"/dist:/var/lib/grafana/plugins --name=grafana grafana/grafana:latest
```

3. Reload plugin

```BASH
docker restart grafana
```

Access from `http://localhost:3000`.
First time login will be user:**admin** password:**admin**

### Build a release

You need to have commit rights to the GitHub repository to publish a release.

1. Update the version number in the `package.json` file.
2. Update the `CHANGELOG.md` by copy and pasting the relevant PRs from [Github's Release drafter interface](https://github.com/grafana/iot-sitewise-datasource/releases/new) or by running `yarn generate-release-notes` (you'll need to install the [gh cli](https://cli.github.com/) and [jq](https://jqlang.github.io/jq/) to run this command).
3. PR the changes.
4. Once merged, follow the Drone release process that you can find [here](https://github.com/grafana/integrations-team/wiki/Plugin-Release-Process#drone-release-process)

## E2E Tests

This plugin uses [playwright](https://playwright.dev/) and [@grafana/plugin-e2e](https://github.com/grafana/plugin-tools/tree/main/packages/plugin-e2e) for e2e end tests.

We support writing/running e2e tests against both mock and live data. Hitting live endpoints in AWS gives us the best fidelity but is potentially difficult for external contributors, as they won't share our internal AWS credentials. Using live endpoints can also sometimes be slow, or run into potential rate limiting issues if we run many tests at once. To get around these issues we allow developers to run their tests against mock data, and also have made it easier to record live data and generate mock data for tests. See helper functions in the handleMocks.ts file to generate and/or call mocks.

To run existing e2e tests we have several scripts to choose from:

- `yarn run test:e2e:debug` will open up playwright's ui mode so you can see screenshots and debug broken tests more easily. This will also use mocks.
- `test:e2e:use-mocks` to run all tests with mocks in the terminal (no ui mode)
- `test:e2e:use-live-data` to run all tests with live data
- `test:e2e:use-live-data:generate-mocks` to run all tests with live data and also generate/overwrite any existing mocks.

To use the live data mode, you'll need to create a yaml file in provisioning/datasources called `iot-sitewise.e2e.yaml` and add the following:

```
apiVersion: 1

deleteDatasources:
  - name: e2e-sitewise-invalid-credentials
    orgId: 1
  - name: e2e-sitewise-valid-credentials
    orgId: 1

datasources:
  - name: e2e-sitewise-invalid-credentials
    type: grafana-iot-sitewise-datasource
    defaultRegion: us-east-1
    editable: true
    secureJsonData:
      accessKey: invalid-mock-access-key
      secretKey: invalid-mock-secret-key

  - name: e2e-sitewise-valid-credentials
    type: grafana-iot-sitewise-datasource
    editable: true
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: { put your access key here that has access to sitewise in aws }
      secretKey: { put your secret key here that has access to sitewise in aws }
```
