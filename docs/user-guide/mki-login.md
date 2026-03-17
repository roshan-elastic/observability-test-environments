## Login to MKI

In order to make it accessible to more people we have wrapped the process and some installation in `oblt-cli`.
The following command will allow you to access a MKI k8s cluster where a oblt cluster is running.

```shell
CLUSTER_NAME=serverless-my-cluster
oblt-cli services mki login --cluster-name "${CLUSTER_NAME}"
```

!!! Warning

    Stop the VPN connection before running the command, otherwise the process will fail.
    Aviatrix VPN Client and oblt-cli uses the same port to to receive the credentials, so you need to stop the VPN connection before running the command.

!!! Warning

    The process will ask for authentication several times, so you need to be ready to provide the required information. It ask for authentication in the browser to access the VPN,
    then for permissions to import TLS certificates, and to launch the VPN client as admin user.
    Finally it authenticates against the teleport proxy to access the k8s cluster.

First you have to connect the VPN, oblt-cli will configure the VPN for you, and then you will be asked to connect to the VPN.

```log
[2024-05-07T14:12:32.782Z]:[info]:Using config file: /Users/my-user/.oblt-cli/config.yaml
[2024-05-07T14:12:32.782Z]:[info]:SlackChannel: '#observablt-bots'
[2024-05-07T14:12:32.782Z]:[info]:User: 'my-user'
[2024-05-07T14:12:32.782Z]:[info]:Git mode: 'ssh'
[2024-05-07T14:12:54.184Z]:[info]:out: [2024-05-07T16:12:54+0200]: Checking if Aviatrix VPN Client is installed
[2024-05-07T16:12:54+0200]: Adding Aviatrix VPN OpenVPN Client to the PATH
[2024-05-07T14:12:54.187Z]:[info]:Starting
[2024-05-07T14:13:14.22Z]:[info]:VPN credentials received
[2024-05-07T14:13:14.22Z]:[info]:Starting OpenVPN
[2024-05-07T14:13:14.22Z]:[warn]:We need elevated privileges to run OpenVPN, so you will be asked for your password.
[2024-05-07T14:13:14.901Z]:[info]:out: [2024-05-07T16:13:14+0200]: Starting VPN client it requires your password
tab 1 of window id 2957
[2024-05-07T14:13:33.915Z]:[info]:VPN connection established
[2024-05-07T14:13:33.915Z]:[info]:Done
```

Once you are connected to the VPN, you will login to the MKI cluster.

```log
[2024-05-07T16:13:35+0200]: Checking if Teleport is installed
[2024-05-07T16:13:35+0200]: Checking if kubectl is installed
[2024-05-07T16:13:35+0200]: Login in qa with teleport
> Profile URL:        https://teleport-proxy.staging.getin.cloud:3080
  Logged in as:       my-user@elastic.co
  Cluster:            staging
  Roles:              cloud-sre, kube-developer
  Logins:             ubuntu, centos, elastic, root
  Kubernetes:         enabled
  Kubernetes cluster: "qa-awseuw1-cp-internal-app-2"
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:34:49 +0200 CEST [valid for 10h21m0s]
  Extensions:         login-ip, permit-agent-forwarding, permit-port-forwarding, permit-pty, private-key-policy

  Profile URL:        https://teleport-proxy.secops.elstc.co:3080
  Logged in as:       my-user@elastic.co
  Cluster:            production
  Roles:              kube-developer
  Kubernetes:         enabled
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:25:25 +0200 CEST [valid for 10h12m0s]
  Extensions:         login-ip, permit-port-forwarding, permit-pty, private-key-policy

[2024-05-07T16:13:37+0200]: Getting MKI project details
[2024-05-07T16:13:37+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:13:38+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:13:38+0200]: ENVIRONMENT: qa
[2024-05-07T16:13:38+0200]: CONSOLE: https://adminconsole.qa.cld.elstc.co
[2024-05-07T16:13:38+0200]: PROJECT_TYPE: observability
[2024-05-07T16:13:39+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:13:39+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:13:39+0200]: ENVIRONMENT: qa
[2024-05-07T16:13:39+0200]: KIBANA_K8S_CLUSTER: qa-awseuw1-cp-internal-app-2
Logged into Kubernetes cluster "qa-awseuw1-cp-internal-app-2". Try 'kubectl version' to test the connection.
Context "staging-qa-awseuw1-cp-internal-app-2" modified.
bash-5.2$
```

Finally, you can access to the k8s resources using `kubectl`.

```shell
bash-3.2$ kubectl get pods

NAME                                     READY   STATUS    RESTARTS   AGE
es-es-index-575f9784b7-f4zwv             1/1     Running   0          5m9s
es-es-index-575f9784b7-j5zs6             1/1     Running   0          3m53s
es-es-index-575f9784b7-z7pkb             1/1     Running   0          4m27s
es-es-search-55d9cd8b67-sz2ft            1/1     Running   0          45m
fleet-757dc66bf9-hrhxg                   1/1     Running   0          2d2h
kb-background-tasks-kb-649f444f6-gh2wv   1/1     Running   0          31h
kb-ui-kb-5d97665746-wnxxd                1/1     Running   0          31h
```
