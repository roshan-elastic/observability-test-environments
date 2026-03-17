# Integrations Helm Chart

Helm chart to deploy integrations.

## Parameters

### kibana Kibana connection configuration

| Name              | Description                | Value                |
| ----------------- | -------------------------- | -------------------- |
| `kibana.url`      | Kibana host to connect to. | `http://kibana:5601` |
| `kibana.username` | Kibana username to use.    | `elastic`            |
| `kibana.password` | Kibana password to use.    | `changeme`           |

### fleet Fleet Server connection configuration.

| Name           | Description                                       | Value                      |
| -------------- | ------------------------------------------------- | -------------------------- |
| `fleet.url`    | Fleet Server host to connect to.                  | `http://fleet-server:8220` |
| `fleet.policy` | Fleet Server policy to configure the integration. | `Default Integration`      |

### elasticAgent Elastic Agent configuration

| Name                   | Description                          | Value  |
| ---------------------- | ------------------------------------ | ------ |
| `elasticAgent.enabled` | Enable or disable the Elastic Agent. | `true` |

### workspace Workspace configuration

| Name                | Description                      | Value  |
| ------------------- | -------------------------------- | ------ |
| `workspace.enabled` | Enable or disable the workspace. | `true` |

### packageRegistry Package Registry configuration

| Name                      | Description                             | Value  |
| ------------------------- | --------------------------------------- | ------ |
| `packageRegistry.enabled` | Enable or disable the package registry. | `true` |

### integration System integration configuration

| Name                  | Description                                                                                          | Value    |
| --------------------- | ---------------------------------------------------------------------------------------------------- | -------- |
| `integration.name`    | name of the integration.                                                                             | `system` |
| `integration.version` | version/branch of the integration.                                                                   | `main`   |
| `imagePullSecrets`    | list of imagePullSecrets to use.                                                                     | `[]`     |
| `nameOverride`        | Overrides the the Helm Chart name used to naming kubernetes objects, by default `elastic-agent`.     | `""`     |
| `fullnameOverride`    | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `elastic-agent`. | `""`     |

### serviceAccount Kubernetes Service Account settings

| Name                         | Description                                                                                                            | Value  |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a service account should be created                                                                  | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account                                                                              | `{}`   |
| `serviceAccount.name`        | The name of the service account to use. If not set and create is true, a name is generated using the fullname template | `""`   |

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

### service service Service settings

| Name           | Description  | Value       |
| -------------- | ------------ | ----------- |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `8080`      |

### ingress Ingress settings

| Name                        | Description                                    | Value                                 |
| --------------------------- | ---------------------------------------------- | ------------------------------------- |
| `ingress.enabled`           | Enable or disable the ingress resource.        | `false`                               |
| `ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`                   |
| `ingress.host`              | Hostname exposed by the ingress resource.      | `package-registry.127.0.0.1.ip.es.io` |
| `ingress.annotations`       | Additional annotation to the ingress resource. | `{}`                                  |
