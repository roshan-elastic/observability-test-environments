# Overview

All APM Managed services should have an APM_SERVER_URL configured to be able to send APM data, metrics and potentially logs to an APM Server.

To reduce the requirements on local development footprint we should ideally just send data to a remote APM Server which we all have access to.

The cluster URL and info can be retrieved using the following oblt-cli command:

```shell
oblt-cli cluster secrets info --cluster-name apm-oblt
```

The credentials can be retrieved using the following oblt-cli commands:

For general credentials document use

```shell
oblt-cli cluster secrets credentials --cluster-name apm-oblt
```

For a environment variables file use

```shell
oblt-cli cluster secrets env --cluster-name apm-oblt --output-file .env
```

## Update

Apm-oblt is updated every Tuesday to the latest BC/Release version.

## Deployed versions

{% include 'apm-oblt/cluster-info.md' ignore missing %}
