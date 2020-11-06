# AWS IoT Sitewise Datasource

WORK IN PROGRESS


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

### Docker

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
docker run -e AWS_SHARED_CREDENTIALS_FILE="/Users/grafana/.aws/credentials" -d -p 3000:3000 -v ~/.aws/credentials:/Users/grafana/.aws/credentials -v "$(pwd)"/dist:/var/lib/grafana/plugins --name=grafana grafana/grafana:latest
```

3. Reload plugin

```BASH
docker restart grafana
```

Access from `http://localhost:3000`. 
First time login will be user:**admin** password:**admin**
