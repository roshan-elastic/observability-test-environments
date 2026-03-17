# elastic-agent

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 7.16.2](https://img.shields.io/badge/AppVersion-7.16.2-informational?style=flat-square)

Elastic Agent Helm chart for Kubernetes

**Homepage:** <https://github.com/elastic/beats/tree/master/x-pack/elastic-agent>

## Source Code

* <https://github.com/elastic/beats/tree/master/x-pack/elastic-agent>

## Requirements

Kubernetes: `>=1.16.0-0`

## Parameters

### Docker image settings

| Name                  | Description                                                        | Value                                   |
| --------------------- | ------------------------------------------------------------------ | --------------------------------------- |
| `image.name`          | Docker image of the Elastic Agent.                                 | `docker.elastic.co/beats/elastic-agent` |
| `image.tag`           | Tag for the Docker image.                                          | `7.16.2`                                |
| `image.pullSecrets`   | List of secret names used to pull Docker images.                   | `[]`                                    |
| `image.pullPolicy`    | Pull Docker images policy [IfNotPresent, Always, Never].           | `IfNotPresent`                          |
| `commandArgs`         | Command arguments to pass to the Elastic Agent.                    | `["-e","-d","*"]`                       |
| `mode`                | Elastic Agent deployment type [deploymnent, daemonset, standalone] | `deployment`                            |
| `kubeSystemNamespace` | Kubernetes namespace where the kube-system pods are running.       | `kube-system`                           |
| `metricsService`      | Kubernetes service name for kube-state-metrics.                    | `kube-state-metrics:8080`               |

### Elastic Agent Standalone packages configuration

| Name                            | Description                                          | Value    |
| ------------------------------- | ---------------------------------------------------- | -------- |
| `standalonePackages.kubernetes` | Kubernetes package versions in used standalone mode. | `1.17.2` |
| `standalonePackages.log`        | log package versions in used standalone mode.        | `1.0.0`  |
| `standalonePackages.system`     | system package versions in used standalone mode.     | `1.11.0` |

### standaloneInputs standalone inputs configuration

| Name                                        | Description                         | Value   |
| ------------------------------------------- | ----------------------------------- | ------- |
| `standaloneInputs.kubernetesClusterMetrics` | Collect Kubernetes cluster metrics. | `true`  |
| `standaloneInputs.kubernetesNodeMetrics`    | Collect Kubernetes nodes metrics.   | `true`  |
| `standaloneInputs.systemLogs`               | Collect Nodes logs.                 | `true`  |
| `standaloneInputs.systemMetrics`            | Collect Nodes metrics.              | `true`  |
| `standaloneInputs.containerLogs`            | Collect containers logs.            | `true`  |
| `standaloneInputs.extra`                    | Configure additional inputs.        | `[]`    |
| `ecsVersion`                                | Elastic Common Schema versions.     | `1.9.0` |

### fleet Elastic Agent settings.

| Name                    | Description                                                                                                                                                            | Value   |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `fleet.secretName`      | Kubernetes secret name to use for the Elastic Agent configuration. By default `ELACTIC_AGENT_NAME-fleet-creds`.                                                        | `""`    |
| `fleet.enrollToken`     | token to use for the Elastic Agent enrollment. This is not needed in case fleet server is configured and `fleet.enroll` is set. Then the token is fetched from Kibana. | `""`    |
| `fleet.tokenPolicyName` | token policy name to use for fetching token from Kibana. This requires Kibana configs to be set. This setting is ignored if you provide a `enrollToken`.               | `""`    |
| `fleet.force`           | To ensure that enrollment occurs on every start of the container set `fleet.force` to true.                                                                            | `true`  |
| `fleet.enroll`          | set to `true` for enrollment into fleet-server. If not set, Elastic Agent is run in standalone mode.                                                                   | `true`  |
| `fleet.insecure`        | communicate with Fleet with either insecure HTTP or unverified HTTPS.                                                                                                  | `false` |
| `fleet.serverUrl`       | URL of the Fleet Server to enroll into (e.g. https://fleet-server:8220).                                                                                               | `""`    |

### fleetServer Fleet server settings

| Name                       | Description                                                                                                           | Value     |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------- | --------- |
| `fleetServer.enabled`      | set to true enables bootstrapping of Fleet Server inside Elastic Agent (forces `fleet.enroll` enabled).               | `false`   |
| `fleetServer.insecureHttp` | expose Fleet Server over HTTP (not recommended; insecure).                                                            | `false`   |
| `fleetServer.host`         | binding host for Fleet Server HTTP (overrides the policy). By default this is 0.0.0.0.                                | `0.0.0.0` |
| `fleetServer.port`         | binding port for Fleet Server HTTP (overrides the policy).                                                            | `8220`    |
| `fleetServer.serviceToken` | service token to use for communication with elasticsearch.                                                            | `""`      |
| `fleetServer.policyId`     | policy ID for Fleet Server to use for itself ("Default Fleet Server policy" used when undefined).                     | `""`      |
| `fleetServer.secretName`   | Kubernetes secret name to use for the Fleet server configuration. By default `ELACTIC_AGENT_NAME-fleet-server-creds`. | `""`      |

### elasticsearch Elasticsearch settings.

| Name                       | Description                                                                                                  | Value |
| -------------------------- | ------------------------------------------------------------------------------------------------------------ | ----- |
| `elasticsearch.host`       | URL of the Elasticsearch (e.g. https://elasticsearch:9200).                                                  | `""`  |
| `elasticsearch.username`   | Elasticsearch username.                                                                                      | `""`  |
| `elasticsearch.password`   | Elasticsearch password.                                                                                      | `""`  |
| `elasticsearch.secretName` | Kubernetes secret name to use for the Elasticsearch configuration. By default `ELACTIC_AGENT_NAME-es-creds`. | `""`  |
| `elasticsearch.apiToken`   | service token to use for communication with Elasticsearch.                                                   | `""`  |

### kibana Kibana settings.

| Name                      | Description                                                                                          | Value |
| ------------------------- | ---------------------------------------------------------------------------------------------------- | ----- |
| `kibana.host`             | URL of the Kibana (e.g. https://kibana:5601).                                                        | `""`  |
| `kibana.username`         | Elasticsearch username.                                                                              | `""`  |
| `kibana.password`         | Elasticsearch password.                                                                              | `""`  |
| `kibana.secretName`       | service token to use for communication with Kibana.                                                  | `""`  |
| `kibana.agent_policies`   | list of agent policies in YAML to create in Kibana.                                                  | `[]`  |
| `kibana.package_policies` | list of package policies in YAML to create in Kibana.                                                | `[]`  |
| `nameOverride`            | Overrides the the Helm Chart name used to naming kubernetes objects, by default `elastic-agent`.     | `""`  |
| `fullnameOverride`        | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `elastic-agent`. | `""`  |

### serviceAccount Kubernetes Service Account settings

| Name                         | Description                                                                                                            | Value  |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a service account should be created                                                                  | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account                                                                              | `{}`   |
| `serviceAccount.name`        | The name of the service account to use. If not set and create is true, a name is generated using the fullname template | `""`   |

### serviceAccount Allows you to set the [securityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for the pod


### serviceAccount Allows you to set the [securityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for the container


### serviceAccount Additional labels to add to the Kubernetes objects.


### resources Allows you to set the [resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/) for the kubernetes pods.


### nodeSelector Configurable [nodeSelector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) so that you can target specific nodes for your Elasticsearch cluster.

| Name          | Description                                                                                          | Value |
| ------------- | ---------------------------------------------------------------------------------------------------- | ----- |
| `tolerations` | Configurable [tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/). | `[]`  |

### affinity Value for the [node affinity settings](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity-beta-feature).

| Name        | Description                                | Value |
| ----------- | ------------------------------------------ | ----- |
| `extraEnvs` | Additional environment variables for pods. | `[]`  |
