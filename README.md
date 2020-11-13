# AWS IoT Sitewise Datasource Development Guide

Please add any feedback to the [issues](https://github.com/grafana/iot-sitewise-datasource/issues) folder, and we will follow up shortly.

For configuraiton options, see: [src/README.md](src/README.md)


To get the latest build artifacts for a branch, click on the green checkbox after the build, and then link to the "package" details:

![package](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/package.png)

Once in circleci, look at the artifacts tab and find the download for your platform:

![artifacts](https://raw.githubusercontent.com/grafana/iot-sitewise-datasource/main/docs/artifacts.png)


## Developer Guide

### Build

#### Getting started
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

4. Build backend for all platforms
```BASH
mage buildAll
```

5. Run tests

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
