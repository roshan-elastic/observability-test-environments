## Run local Kibana

Oblt-cli offers a command to run a local Kibana instance with the proper configuration, without to have to generate the `kibana.serverless.yml` file. The command will generate the `KIBANA_SOURCE_DIR/config/kibana.serverless.yaml` kibana configuration file,
make the bootstrap of the Kibana source code, and run the Kibana instance with the correct Node.js version.

```shell
KIBANA_SRC_FOLDER="${HOME}/src/kibana"
KIBANA_YML=${KIBANA_SRC_FOLDER}/config/kibana.serverless.yaml
CLUSTER_NAME=serverless-my-cluster

oblt-cli services mki run-local-kibana \
  --cluster-name "${CLUSTER_NAME}" \
  --kibana-src "${KIBANA_SRC_FOLDER}" \
  --kibana-yaml-path "${KIBANA_YML}"
```

The following logs show the output of the command.
We can see that the command access the VPN, login in the teleport proxy, and get the project details to generate the `kibana.serverless.yml` file. Finally, it install Node.js if it is needed, download the dependencies and run the Kibana instance.

!!! Warning

    Stop the VPN connection before running the command, otherwise the process will fail.
    Aviatrix VPN Client and oblt-cli uses the same port to to receive the credentials, so you need to stop the VPN connection before running the command.

!!! Warning

    The process will ask for authentication several times, so you need to be ready to provide the required information. It ask for authentication in the browser to access the VPN,
    then for permissions to import TLS certificates, and to launch the VPN client as admin user.
    Finally it authenticates against the teleport proxy to access the k8s cluster.

```log
[2024-05-07T14:34:39.909Z]:[info]:Using config file: /Users/my-user/.oblt-cli/config.yaml
[2024-05-07T14:34:39.91Z]:[info]:SlackChannel: '#observablt-bots'
[2024-05-07T14:34:39.91Z]:[info]:User: 'my-user'
[2024-05-07T14:34:39.91Z]:[info]:Git mode: 'ssh'
[2024-05-07T14:34:52.591Z]:[info]:out: [2024-05-07T16:34:52+0200]: Checking if Aviatrix VPN Client is installed
[2024-05-07T16:34:52+0200]: Adding Aviatrix VPN OpenVPN Client to the PATH
[2024-05-07T14:34:52.592Z]:[info]:Starting
[2024-05-07T14:35:02.603Z]:[info]:VPN credentials received
[2024-05-07T14:35:02.603Z]:[info]:Starting OpenVPN
[2024-05-07T14:35:02.603Z]:[warn]:We need elevated privileges to run OpenVPN, so you will be asked for your password.
[2024-05-07T14:35:03.67Z]:[info]:out: [2024-05-07T16:35:02+0200]: Starting VPN client it requires your password
tab 1 of window id 3274
[2024-05-07T14:35:22.683Z]:[info]:VPN connection established
[2024-05-07T14:35:22.683Z]:[info]:Done
[2024-05-07T16:35:24+0200]: Checking if Teleport is installed
[2024-05-07T16:35:24+0200]: Checking if kubectl is installed
[2024-05-07T16:35:24+0200]: Login in qa with teleport
> Profile URL:        https://teleport-proxy.staging.getin.cloud:3080
  Logged in as:       my-user@elastic.co
  Cluster:            staging
  Roles:              cloud-sre, kube-developer
  Logins:             ubuntu, centos, elastic, root
  Kubernetes:         enabled
  Kubernetes cluster: "qa-awseuw1-cp-internal-app-2"
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:34:49 +0200 CEST [valid for 9h59m0s]
  Extensions:         login-ip, permit-agent-forwarding, permit-port-forwarding, permit-pty, private-key-policy

  Profile URL:        https://teleport-proxy.secops.elstc.co:3080
  Logged in as:       my-user@elastic.co
  Cluster:            production
  Roles:              kube-developer
  Kubernetes:         enabled
  Kubernetes groups:  kube-developer
  Valid until:        2024-05-08 02:25:25 +0200 CEST [valid for 9h50m0s]
  Extensions:         login-ip, permit-port-forwarding, permit-pty, private-key-policy

[2024-05-07T16:35:27+0200]: Getting MKI project details
[2024-05-07T16:35:27+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:35:27+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:35:27+0200]: ENVIRONMENT: qa
[2024-05-07T16:35:27+0200]: CONSOLE: https://adminconsole.qa.cld.elstc.co
[2024-05-07T16:35:27+0200]: PROJECT_TYPE: observability
[2024-05-07T16:35:28+0200]: CLUSTER_NAME: my-project-qa-oblt
[2024-05-07T16:35:28+0200]: PROJECT_ID: b4ba071a21b545d4b690a5f529a11783
[2024-05-07T16:35:28+0200]: ENVIRONMENT: qa
[2024-05-07T16:35:28+0200]: KIBANA_K8S_CLUSTER: qa-awseuw1-cp-internal-app-2
Logged into Kubernetes cluster "qa-awseuw1-cp-internal-app-2". Try 'kubectl version' to test the connection.
Context "staging-qa-awseuw1-cp-internal-app-2" modified.
[2024-05-07T16:35:35+0200]: Kibana yaml config file generated at /Users/my-user/src/kibana/config/kibana.serverless.yaml
[2024-05-07T16:35:35+0200]: Running Kibana from /Users/my-user/src/kibana with config /Users/my-user/src/kibana/config/kibana.serverless.yaml
[2024-05-07T16:35:35+0200]: Run Kibana
[2024-05-07T16:35:35+0200]: Checking if Node.js (18.18.2) is installed
[2024-05-07T16:35:35+0200]: Node.js (18.18.2) is not installed
[2024-05-07T16:35:36+0200]: You have Node.js v20.11.1 installed/default
[2024-05-07T16:35:36+0200]: Do you want to install/activate Node.js (18.18.2) with nvm?
[y/N]: y
[2024-05-07T16:35:48+0200]: Checking if nvm is installed
v18.18.2 is already installed.
Now using node v18.18.2 (npm v9.8.1)
yarn run v1.22.19
$ node scripts/kbn bootstrap
warn updated package map
warn updated tsconfig.json paths
[bazel] Extracting Bazel installation...
[bazel] Starting local Bazel server and connecting to it...
[bazel] INFO: Invocation ID: 84b76d6f-2e4b-4b23-ab9e-8a26f8e794ca
[bazel] $ node ./preinstall_check
[bazel] [1/5] Validating package.json...
[bazel] [2/5] Resolving packages...
[bazel] [3/5] Fetching packages...
[bazel] [4/5] Linking dependencies...
[bazel] [5/5] Building fresh packages...
[bazel] INFO: Analyzed 3 targets (778 packages loaded, 3252 targets configured).
[bazel] INFO: Found 3 targets...
[bazel] INFO: From Action packages/kbn-monaco/target_workers:
[bazel] Browserslist: caniuse-lite is outdated. Please run:
[bazel]   npx update-browserslist-db@latest
[bazel]   Why you should do it regularly: https://github.com/browserslist/update-db#readme
[bazel] INFO: From Action packages/kbn-ui-shared-deps-src/shared_built_assets:
[bazel] Browserslist: caniuse-lite is outdated. Please run:
[bazel]   npx update-browserslist-db@latest
[bazel]   Why you should do it regularly: https://github.com/browserslist/update-db#readme
[bazel] INFO: Elapsed time: 296.612s, Critical Path: 45.28s
[bazel] INFO: 380 processes: 317 disk cache hit, 10 internal, 53 darwin-sandbox.
[bazel]
success shared bundles built
Browserslist: caniuse-lite is outdated. Please run:
  npx update-browserslist-db@latest
  Why you should do it regularly: https://github.com/browserslist/update-db#readme
success yarn.lock analysis completed without any issues
success vscode config updated
✨  Done in 640.82s.
yarn run v1.22.19
$ node scripts/kibana --dev --serverless=oblt --config /Users/my-user/src/kibana/config/kibana.serverless.yaml
```
