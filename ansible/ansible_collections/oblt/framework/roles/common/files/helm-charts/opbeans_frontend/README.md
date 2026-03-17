# Opbeans Frontend Helm chart

Opbeans Frontend Helm chart deploy the [Opbeans Frontend](https://github.com/elastic/opbeans-frontend/).

## Parameters

### opbeans Opbeans Frontend Helm chart configuration

| Name                       | Description                                                                                     | Value  |
| -------------------------- | ----------------------------------------------------------------------------------------------- | ------ |
| `opbeans.enabled`          | Enable or disable the whole kubernetes deployments, this allow to use it on parent charts.      | `true` |
| `opbeans.nameOverride`     | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.      | `""`   |
| `opbeans.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `""`   |

### opbeans.image Docker image settings

| Name                       | Description                                                                              | Value                                                 |
| -------------------------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `opbeans.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/opbeans-frontend` |
| `opbeans.image.tag`        | Tag of the Docker image (latest).                                                        | `daily`                                               |
| `opbeans.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                        |

### apm APM configuration settings.

| Name                         | Description                                   | Value         |
| ---------------------------- | --------------------------------------------- | ------------- |
| `opbeans.apm.type`           | Type of reporting API [apm|otel|annotations]  | `annotations` |
| `opbeans.apm.annotationType` | Type of annotations to use [none|nodejs|java] | `nodejs`      |
| `opbeans.apm.enabled`        | Enable to send data to the APM service.       | `true`        |

### ingress Ingress settings

| Name                                | Description                                    | Value                                 |
| ----------------------------------- | ---------------------------------------------- | ------------------------------------- |
| `opbeans.ingress.enabled`           | Enable or disable the ingress resource.        | `true`                                |
| `opbeans.ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`                   |
| `opbeans.ingress.host`              | Hostname exposed by the ingress resource.      | `opbeans-frontend.127.0.0.1.ip.es.io` |
| `opbeans.ingress.annotations`       | Additional annotation to the ingress resource. | `undefined`                           |

### resources requests and limits

| Name                                | Description                         | Value   |
| ----------------------------------- | ----------------------------------- | ------- |
| `opbeans.resources.limits.memory`   | The memory limit for the Opbean     | `1Gi`   |
| `opbeans.resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `opbeans.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |
