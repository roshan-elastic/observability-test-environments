# Troubleshotting

{% include '/user-guide/how-to-report-issues.md' ignore missing %}

## 401 error accessing MKI

The production environment needs to pass you own API key to access the Serverless admin API.
You will need to enter in the ESS admin console of production environment and create a new API key.
Then you can pass the API key to the `oblt-cli` using the flag `--api-key`.

```log
[2024-05-07T14:25:30+0200]: Getting MKI project details
[2024-05-07T14:25:30+0200]: CLUSTER_NAME: my-project-oblt
[2024-05-07T14:25:30+0200]: PROJECT_ID: ad85d674bb1544aa92178a7aa4dd0400
[2024-05-07T14:25:30+0200]: ENVIRONMENT: production
[2024-05-07T14:25:30+0200]: CONSOLE: https://admin.found.no
[2024-05-07T14:25:30+0200]: PROJECT_TYPE: observability
curl: (22) The requested URL returned error: 401
```

```shell
API_KEY=AAAAAA==
CLUSTER_NAME=my-project-oblt
KIBANA_YML_PATH=${PWD}/kibana.yml
oblt-cli services mki generate-kibana-yaml \
  --cluster-name "${CLUSTER_NAME}" \
  --kibana-yaml-path "${KIBANA_YML_PATH}" \
  --api-key "${API_KEY}"
```

## VPN connection issues with oblt-cli

You can use the `oblt-cli` to connect to the VPN, to do that you must be disconnected from the VPN and stop the VPN client you use.
If you are still having issues, you can try to use your VPN client to connect to the VPN and tell `oblt-cli` to not launch the VPN connection. The flag `--vpn=false` will disable the VPN launch.
The following command will login in MKI without launch the VPN:

```shell
CLUSTER_NAME=serverless-my-cluster
oblt-cli services mki login --cluster-name "${CLUSTER_NAME}" --vpn=false
```

## After launching the VPN connection, the internet connection stuck

This happens because the VPN client is not started correctly, you need to kill the VPN client.

```shell
sudo killall openvpn
```

```shell
sudo killall $(cat "${HOME}/.oblt-cli/vpn/openvpn.pid")
```

## Teleport connection timeout

If you experience a timeout error when trying to connect to the Teleport proxy, you can try to reconnect to the VPN or restart the VPN connection. To do that, you have to relaunch the command again.

```log
ERROR: Teleport proxy not available at teleport-proxy.staging.getin.cloud:3080.
Get "https://teleport-proxy.staging.getin.cloud:3080/webapi/ping/okta": net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
```
