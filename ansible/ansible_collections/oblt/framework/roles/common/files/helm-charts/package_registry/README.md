# Package registry Helm chart

This is the common Helm Chart definition for the package registry.

## Parameters

### image Docker image settings

| Name               | Description                                                                                  | Value                                                              |
| ------------------ | -------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/integrations/package-registry). | `docker.elastic.co/observability-ci/integrations/package-registry` |
| `image.tag`        | Tag of the Docker image (latest).                                                            | `latest`                                                           |
| `image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                         | `IfNotPresent`                                                     |


### imagePullSecrets Docker secrets used to pull the Docker images.

| Name                       | Description                                                                                    | Value               |
| -------------------------- | ---------------------------------------------------------------------------------------------- | ------------------- |
| `imagePullSecrets[0].name` | Name of the Docker secret.                                                                     | `docker.elastic.co` |
| `nameOverride`             | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.     | `""`                |
| `fullnameOverride`         | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeans`. | `""`                |
| `annotations`              | Additional annotations to add to the kubernetes deployment.                                    | `{}`                |


### serviceAccount Kubernetes Service Account settings

| Name                         | Description                                           | Value  |
| ---------------------------- | ----------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a service account should be created | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account             | `{}`   |
| `serviceAccount.name`        | The name of the service account to use.               | `""`   |


### securityContext Enabled Apache Server containers' Security Context




### service service Service settings

| Name           | Description  | Value       |
| -------------- | ------------ | ----------- |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80`        |


### ingress Ingress settings

| Name                        | Description                                    | Value                                 |
| --------------------------- | ---------------------------------------------- | ------------------------------------- |
| `ingress.enabled`           | Enable or disable the ingress resource.        | `false`                               |
| `ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`                   |
| `ingress.host`              | Hostname exposed by the ingress resource.      | `package-registry.127.0.0.1.ip.es.io` |
| `ingress.annotations`       | Additional annotation to the ingress resource. | `{}`                                  |


### resources requests and limits

| Name                        | Description                              | Value       |
| --------------------------- | ---------------------------------------- | ----------- |
| `resources.limits.memory`   | The memory limit for the Opbean          | `512Mi`     |
| `resources.requests.memory` | The requested memory for the Opbean      | `256Mi`     |
| `resources.requests.cpu`    | The requested cpu for the Opbean         | `0.1`      |
| `extraEnvVars`              | Array to add extra environment variables | `undefined` |
