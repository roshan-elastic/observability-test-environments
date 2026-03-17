# APM Synthtrace

[@kbn/apm-synthtrace](https://github.com/elastic/kibana/blob/main/packages/kbn-apm-synthtrace/README.md) is a tool in technical preview to generate synthetic APM data. It is intended to be used for development and testing of the Elastic APM app in Kibana.

## Parameters

### image Image settings

| Name               | Description       | Value                                               |
| ------------------ | ----------------- | --------------------------------------------------- |
| `image.repository` | Image repository  | `docker.elastic.co/observability-ci/apm-synthtrace` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent`                                      |
| `image.tag`        | Image tag         | `1.0.1`                                             |

### imagePullSecrets Docker secrets used to pull the Docker images.

| Name                       | Description                                                                                    | Value               |
| -------------------------- | ---------------------------------------------------------------------------------------------- | ------------------- |
| `imagePullSecrets[0].name` | Name of the Docker secret.                                                                     | `docker.elastic.co` |
| `nameOverride`             | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.     | `""`                |
| `fullnameOverride`         | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeans`. | `""`                |
| `podAnnotations`           | Additional annotations to add to the kubernetes deployment.                                    | `{}`                |

### serviceAccount Kubernetes Service Account settings

| Name                         | Description                                           | Value  |
| ---------------------------- | ----------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a service account should be created | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account             | `{}`   |
| `serviceAccount.name`        | The name of the service account to use.               | `""`   |

### podSecurityContext Enabled pods' Security Context

| Name              | Description                          | Value |
| ----------------- | ------------------------------------ | ----- |
| `securityContext` | Enabled containers' Security Context | `{}`  |

### resources requests and limits

| Name                        | Description                                                                                                        | Value                                  |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------ | -------------------------------------- |
| `resources.limits.memory`   | The memory limit for the Opbean                                                                                    | `2Gi`                                  |
| `resources.requests.memory` | The requested memory for the Opbean                                                                                | `1Gi`                                  |
| `resources.requests.cpu`    | The requested cpu for the Opbean                                                                                   | `200m`                                 |
| `scenarios`                 | List of scenarios to run see https://github.com/elastic/kibana/tree/main/packages/kbn-apm-synthtrace/src/scenarios | `["simple_logs.ts","simple_trace.ts"]` |
| `cmdOpts`                   | Additional command line options to pass to the synthtrace command                                                  | `--live`                               |

### elasticsearch Elasticsearch settings.

| Name                     | Description                                         | Value |
| ------------------------ | --------------------------------------------------- | ----- |
| `elasticsearch.host`     | URL of the Elasticsearch (e.g. elasticsearch:9200). | `""`  |
| `elasticsearch.username` | Elasticsearch username.                             | `""`  |
| `elasticsearch.password` | Elasticsearch password.                             | `""`  |
