# Opbeans Python Helm chart

Opbeans Python Helm chart deploy the [Opbeans Python](https://github.com/elastic/opbeans-dotnet/)
and all services needed (PostgreSQL, Redis, and Elasticseacrh).

## Parameters

### opbeans Opbeans Python Helm chart configuration

| Name                       | Description                                                                                    | Value  |
| -------------------------- | ---------------------------------------------------------------------------------------------- | ------ |
| `opbeans.enabled`          | Enable or disable the whole kubernetes deployments, this allow to use it on parent charts.     | `true` |
| `opbeans.nameOverride`     | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.     | `""`   |
| `opbeans.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeans`. | `""`   |

### opbeans.image Docker image settings

| Name                       | Description                                                                              | Value                                                             |
| -------------------------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| `opbeans.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/opbeans-dotnet`               |
| `opbeans.image.tag`        | Tag of the Docker image (latest).                                                        | `daily`                                                           |
| `opbeans.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                                    |
| `opbeans.command`          | Override default container command (useful when using custom images)                     | `["dotnet"]`                                                      |
| `opbeans.args`             | Override default container args (useful when using custom images)                        | `["opbeans-dotnet.dll","--urls=http://*:$(OPBEANS_SERVER_PORT)"]` |

### apm APM configuration settings.

| Name                           | Description                                                   | Value                            |
| ------------------------------ | ------------------------------------------------------------- | -------------------------------- |
| `opbeans.apm.type`             | Type of reporting API [apm|otel|annotations]                  | `apm`                            |
| `opbeans.apm.url`              | URL to the APM service.                                       | `https://apm.127.0.0.1.ip.es.io` |
| `opbeans.apm.token`            | Token to authenticate against the APM service.                | `SuP3RSeCr3T`                    |
| `opbeans.apm.verifyServerCert` | Enable Certificates verification on TLS connections.          | `true`                           |
| `opbeans.apm.serviceName`      | Name of the service deploy for the Helm chart (APM spans).    | `opbeans-go`                     |
| `opbeans.apm.serviceVersion`   | Version of the service deploy for the Helm chart (APM spans). | `1.0`                            |
| `opbeans.apm.logLevel`         | Level of debug.                                               | `debug`                          |
| `opbeans.apm.logFile`          | File to debug to.                                             | `stderr`                         |
| `opbeans.apm.inferredSpans`    | Enable inferred spans.                                        | `true`                           |
| `opbeans.apm.logCorrelation`   | Enable log correlation.                                       | `true`                           |
| `opbeans.apm.transactionRate`  | Configure the APM transaction rate (0.9).                     | `0.9`                            |
| `opbeans.apm.environment`      | Environment where the service is deployed.                    | `production`                     |
| `opbeans.apm.enabled`          | Enable to send data to the APM service.                       | `true`                           |

### ingress Ingress settings

| Name                                | Description                                    | Value                               |
| ----------------------------------- | ---------------------------------------------- | ----------------------------------- |
| `opbeans.ingress.enabled`           | Enable or disable the ingress resource.        | `false`                             |
| `opbeans.ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`                 |
| `opbeans.ingress.host`              | Hostname exposed by the ingress resource.      | `opbeans-dotnet.127.0.0.1.ip.es.io` |
| `opbeans.ingress.annotations`       | Additional annotation to the ingress resource. | `undefined`                         |

### opbeans.resources requests and limits

| Name                                | Description                         | Value   |
| ----------------------------------- | ----------------------------------- | ------- |
| `opbeans.resources.limits.memory`   | The memory limit for the Opbean     | `1.5Gi` |
| `opbeans.resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `opbeans.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |
