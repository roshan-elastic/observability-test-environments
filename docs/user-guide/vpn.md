# VPN

## Overview

A Elastic [VPN][] connection is required to access to some Elastic Cloud environments.
We have added a wrapper to `oblt-cli` to manage the VPN configuration and connection.
The following command will allow you to connect to a VPN.

```shell
TARGET_ENVIRONMENT=staging

oblt-cli services vpn --environment "${TARGET_ENVIRONMENT}"
```

Oblt-cli will configure the VPN for you, and then you will be asked to connect to the VPN.

```shell
2023/10/25 13:14:49 Using config file: /Users/myuser/.oblt-cli/config.yaml
2023/10/25 13:14:49 SlackChannel: '#observablt-bots'
2023/10/25 13:14:49 User: 'myuser'
2023/10/25 13:14:49 Git mode: 'http'
[2023-10-25T13:14:59+0200]: Checking if Aviatrix VPN Client is installed
[2023-10-25T13:14:59+0200]: Checking if Teleport is installed
[2023-10-25T13:15:00+0200]: Login in staging VPN
[2023-10-25T13:15:00+0200]: Opening Aviatrix VPN Client.../Users/myuser/.oblt-cli/vpn/elastic-cloud-staging-us.ovpn
[2023-10-25T13:15:00+0200]: When the VPN is connected (green), Press enter to continue
```

[VPN]: http://ela.st/vpn}
