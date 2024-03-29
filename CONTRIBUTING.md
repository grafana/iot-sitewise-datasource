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

To get the best fidelity, we make live requests to AWS for many of our e2e tests. In order to run them you will need to create an AWS User and secret/access keys and add them in either as env variables (`AWS_ACCESS_KEY` and `AWS_SECRET_KEY`) which is how we run our tests in CI or add a yaml file to the provisioning repo like so:

```
apiVersion: 1

deleteDatasources:
  - name: sitewise
    orgId: 1

datasources:
  - name: sitewise
    type: grafana-iot-sitewise-datasource
    editable: true
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: {your access key here}
      secretKey: {your secret key here}
```

### Running e2e tests locally

To run e2e tests locally:

```
yarn run e2e
```

This will then print out a report that can be viewed. Or To run e2e tests locally with [UI mode](https://playwright.dev/docs/test-ui-mode) for easier debugging:

```
yarn run e2e:debug
```

You may also wish to enable "traces" in the playwright.config.ts file which will show screenshots of failures and network requests.
