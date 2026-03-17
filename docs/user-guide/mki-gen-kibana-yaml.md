## Generate Kibana configuration

The main reason to add this feature to `oblt-cli` is to allow developers to connect a local Kibana instance to a real serverless project.
To do that we need a serverless project up and running created with `oblt-cli` and a `kibana.yml` file configured to access to the serverless project.
The Kibana configuration file would have the proper values at `elasticsearch.hosts`, `elasticsearch.serviceAccountToken`, `xpack.encryptedSavedObjects.encryptionKey`, `xpack.security.encryptionKey`, and `xpack.reporting.encryptionKey` settings from the credentials of the cluster.
The following command will generate a `kibana.serverless.yml` file with the proper values for the `serverless-my-cluster` oblt cluster.

```shell
KIBANA_SRC_FOLDER="${HOME}/src/kibana"
KIBANA_YML=${KIBANA_SRC_FOLDER}/config/kibana.serverless.yaml
CLUSTER_NAME=serverless-my-cluster

oblt-cli services mki generate-kibana-yaml \
  --cluster-name "${CLUSTER_NAME}" \
  --kibana-yaml-path "${KIBANA_YML}"
```

The following logs show the output of the command.
We can see that the command access the VPN, login in the teleport proxy, and get the project details to generate the `kibana.serverless.yml` file.

!!! Warning

    Stop the VPN connection before running the command, otherwise the process will fail.
    Aviatrix VPN Client and oblt-cli uses the same port to to receive the credentials, so you need to stop the VPN connection before running the command.

!!! Warning

    The process will ask for authentication several times, so you need to be ready to provide the required information. It asks for authentication in the browser to access the VPN,
    then for permissions to import TLS certificates, and to launch the VPN client as admin user.
    Finally it authenticates against the teleport proxy to access the k8s cluster.

```log
[2024-05-07T14:18:18.683Z]:[info]:Using config file: /Users/my-user/.oblt-cli/config.yaml
[2024-05-07T14:18:18.684Z]:[info]:SlackChannel: '#observablt-bots'
[2024-05-07T14:18:18.684Z]:[info]:User: 'my-user'
[2024-05-07T14:18:18.684Z]:[info]:Git mode: 'ssh'
[2024-05-07T14:18:32.445Z]:[info]:out: [2024-05-07T16:18:32+0200]: Checking if Aviatrix VPN Client is installed
[2024-05-07T16:18:32+0200]: Adding Aviatrix VPN OpenVPN Client to the PATH
[2024-05-07T14:18:32.447Z]:[info]:Starting
[2024-05-07T14:18:52.466Z]:[info]:VPN credentials received
[2024-05-07T14:18:52.466Z]:[info]:Starting OpenVPN
[2024-05-07T14:18:52.466Z]:[warn]:We need elevated privileges to run OpenVPN, so you will be asked for your password.
[2024-05-07T14:18:53.497Z]:[info]:out: [2024-05-07T16:18:52+0200]: Starting VPN client it requires your password
tab 1 of window id 3017
[2024-05-07T14:19:12.519Z]:[info]:VPN connection established
[2024-05-07T14:19:12.519Z]:[info]:Done
[2024-05-07T16:19:14+0200]: Checking if Teleport is installed
[2024-05-07T16:19:14+0200]: Checking if kubectl is installed
[2024-05-07T16:19:14+0200]: Login in qa with teleport
> Profile URL:        https://teleport-proxy.staging.getin.cloud:3080
  Logged in as:       my-user@elastic.co
  Cluster:            staging
  Roles:              cloud-sre, kube-developer
  Logins:             ubuntu, centos, elastic, root
  Kubernetes:         enabled
  Kubernetes cluster: "qa-awseuw1-cp-internal-app-2"
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:34:49 +0200 CEST [valid for 10h16m0s]
  Extensions:         login-ip, permit-agent-forwarding, permit-port-forwarding, permit-pty, private-key-policy

  Profile URL:        https://teleport-proxy.secops.elstc.co:3080
  Logged in as:       my-user@elastic.co
  Cluster:            production
  Roles:              kube-developer
  Kubernetes:         enabled
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:25:25 +0200 CEST [valid for 10h6m0s]
  Extensions:         login-ip, permit-port-forwarding, permit-pty, private-key-policy

[2024-05-07T16:19:16+0200]: Getting MKI project details
[2024-05-07T16:19:16+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:19:16+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:19:16+0200]: ENVIRONMENT: qa
[2024-05-07T16:19:16+0200]: CONSOLE: https://adminconsole.qa.cld.elstc.co
[2024-05-07T16:19:16+0200]: PROJECT_TYPE: observability
[2024-05-07T16:19:17+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:19:17+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:19:17+0200]: ENVIRONMENT: qa
[2024-05-07T16:19:17+0200]: KIBANA_K8S_CLUSTER: qa-awseuw1-cp-internal-app-2
Logged into Kubernetes cluster "qa-awseuw1-cp-internal-app-2". Try 'kubectl version' to test the connection.
Context "staging-qa-awseuw1-cp-internal-app-2" modified.
[2024-05-07T16:19:24+0200]: Kibana yaml config file generated at /Users/my-user/src/kibana/config/kibana.serverless.yaml
```

Once you have the `kibana.serverless.yml` file, you can run a local Kibana instance with a yarn command. For further details about developing with Kibana please go to [Kibana Development Getting Started](https://www.elastic.co/guide/en/kibana/current/development-getting-started.html)```

```shell
KIBANA_SRC_FOLDER="${HOME}/src/kibana"
KIBANA_YML=${KIBANA_SRC_FOLDER}/config/kibana.serverless.yaml
NODE_OPTIONS=" --max-old-space-size=4096"

yarn start --serverless=oblt --config="${KIBANA_YML}"
```
