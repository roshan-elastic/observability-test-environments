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

| Name                | Description                                                                                          | Value                                   |
| ------------------- | ---------------------------------------------------------------------------------------------------- | --------------------------------------- |
| `image.name`        | Docker image of the Elastic Agent.                                                                   | `docker.elastic.co/beats/elastic-agent` |
| `image.tag`         | Tag for the Docker image.                                                                            | `7.16.2`                                |
| `image.pullSecrets` | List of secret names used to pull Docker images.                                                     | `[]`                                    |
| `image.pullPolicy`  | Pull Docker images policy [IfNotPresent, Always, Never].                                             | `IfNotPresent`                          |
| `imagePullSecrets`  | List of secret names used to pull Docker images.                                                     | `[]`                                    |
| `nameOverride`      | Overrides the the Helm Chart name used to naming kubernetes objects, by default `elastic-agent`.     | `""`                                    |
| `fullnameOverride`  | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `elastic-agent`. | `""`                                    |

### serviceAccount settings

| Name                    | Description                                                         | Value                     |
| ----------------------- | ------------------------------------------------------------------- | ------------------------- |
| `serviceAccount.create` | Specifies whether a ServiceAccount should be created.               | `true`                    |
| `serviceAccount.name`   | The name of the ServiceAccount to use.                              | `""`                      |
| `mode`                  | Elastic Agent deployment type [managed, standalone]                 | `managed`                 |
| `metricsService`        | Name of the kube-state-metrics metrics service for standalone mode. | `kube-state-metrics:8080` |

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

### elasticsearch Elasticsearch settings.

| Name                       | Description                                                                                                  | Value |
| -------------------------- | ------------------------------------------------------------------------------------------------------------ | ----- |
| `elasticsearch.host`       | URL of the Elasticsearch (e.g. https://elasticsearch:9200).                                                  | `""`  |
| `elasticsearch.username`   | Elasticsearch username.                                                                                      | `""`  |
| `elasticsearch.password`   | Elasticsearch password.                                                                                      | `""`  |
| `elasticsearch.secretName` | Kubernetes secret name to use for the Elasticsearch configuration. By default `ELACTIC_AGENT_NAME-es-creds`. | `""`  |
| `elasticsearch.apiToken`   | service token to use for communication with Elasticsearch.                                                   | `""`  |

### kibana Kibana settings.

| Name                | Description                                         | Value |
| ------------------- | --------------------------------------------------- | ----- |
| `kibana.host`       | URL of the Kibana (e.g. https://kibana:5601).       | `""`  |
| `kibana.username`   | Elasticsearch username.                             | `""`  |
| `kibana.password`   | Elasticsearch password.                             | `""`  |
| `kibana.secretName` | service token to use for communication with Kibana. | `""`  |
