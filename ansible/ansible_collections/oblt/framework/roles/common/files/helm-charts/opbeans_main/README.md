# Opbeans Helm chart

This is the common Helm Chart definition for all the Opbeans.
The rest of Opbeans will use this Helm Chart as dependency and add its own settings.

## Parameters

### image Docker image settings

| Name               | Description                                                                              | Value                                                          |
| ------------------ | ---------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/it_opbeans-frontend_nginx` |
| `image.tag`        | Tag of the Docker image (latest).                                                        | `latest`                                                       |
| `image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                                 |

### imagePullSecrets Docker secrets used to pull the Docker images.

| Name                       | Description                                                                                    | Value               |
| -------------------------- | ---------------------------------------------------------------------------------------------- | ------------------- |
| `imagePullSecrets[0].name` | Name of the Docker secret.                                                                     | `docker.elastic.co` |
| `nameOverride`             | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.     | `""`                |
| `fullnameOverride`         | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeans`. | `""`                |

### Service Group settings, allow to group opbenas ussing the service-group label.

| Name                      | Description                                                                      | Value     |
| ------------------------- | -------------------------------------------------------------------------------- | --------- |
| `serviceGroup.name`       | Service group name, this name chain all the services deployed.                   | `opbeans` |
| `serviceGroup.deploy`     | True to deploy a global service to group all opbean with the same service group. | `false`   |
| `serviceGroup.member`     | True if the Opbeans is a service group member.                                   | `true`    |
| `serviceGroup.port`       | Source port of the service group.                                                | `3000`    |
| `serviceGroup.targetPort` | Target port to expose the service group.                                         | `3000`    |
| `annotations`             | Additional annotations to add to the kubernetes deployment.                      | `{}`      |

### serviceAccount Kubernetes Service Account settings

| Name                         | Description                                           | Value  |
| ---------------------------- | ----------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a service account should be created | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account             | `{}`   |
| `serviceAccount.name`        | The name of the service account to use.               | `""`   |

### ingress Ingress settings

| Name                        | Description                                                       | Value                        |
| --------------------------- | ----------------------------------------------------------------- | ---------------------------- |
| `ingress.enabled`           | Enable or disable the ingress resource.                           | `false`                      |
| `ingress.tls`               | Enable or disable the tls configuration for the ingress resource. | `true`                       |
| `ingress.certificateIssuer` | Name of the certificates issuer.                                  | `letsencrypt-staging`        |
| `ingress.host`              | Hostname exposed by the ingress resource.                         | `opbeans.127.0.0.1.ip.es.io` |
| `ingress.annotations`       | Additional annotation to the ingress resource.                    | `{}`                         |

### resources requests and limits

| Name                        | Description                         | Value   |
| --------------------------- | ----------------------------------- | ------- |
| `resources.limits.memory`   | The memory limit for the Opbean     | `512Mi` |
| `resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |
| `volumeMounts`              | Volume mounts                       | `[]`    |
| `volumes`                   | Volumes                             | `[]`    |

### apm APM configuration settings.

| Name                   | Description                                                                              | Value                            |
| ---------------------- | ---------------------------------------------------------------------------------------- | -------------------------------- |
| `apm.type`             | Type of reporting API [apm|otel|annotations]                                             | `apm`                            |
| `apm.url`              | URL to the APM service.                                                                  | `https://apm.127.0.0.1.ip.es.io` |
| `apm.token`            | Token to authenticate against the APM service.                                           | `SuP3RSeCr3T`                    |
| `apm.apikey`           | API Key to authenticate against the APM service. The apikey has precedence to the token. | `""`                             |
| `apm.verifyServerCert` | Enable Certificates verification on TLS connections.                                     | `true`                           |
| `apm.serviceName`      | Name of the service deploy for the Helm chart (APM spans).                               | `opbeans-python`                 |
| `apm.serviceVersion`   | Version of the service deploy for the Helm chart (APM spans).                            | `1.0`                            |
| `apm.logLevel`         | Level of debug.                                                                          | `debug`                          |
| `apm.logFile`          | File to debug to.                                                                        | `stderr`                         |
| `apm.inferredSpans`    | Enable inferred spans.                                                                   | `true`                           |
| `apm.logCorrelation`   | Enable log correlation.                                                                  | `true`                           |
| `apm.transactionRate`  | Configure the APM transaction rate (0.9).                                                | `0.9`                            |
| `apm.environment`      | Environment where the service is deployed.                                               | `production`                     |
| `apm.annotationType`   | Type of annotations to use [none|nodejs|java].                                           | `none`                           |
| `apm.enabled`          | Enable to send data to the APM service.                                                  | `true`                           |

### loadgen Load generator settings.

| Name               | Description                             | Value         |
| ------------------ | --------------------------------------- | ------------- |
| `loadgen.enabled`  | Enable the load generator job deploy.   | `false`       |
| `loadgen.schedule` | Schedule to run the load generator job. | `*/2 * * * *` |

### loadgen.api API load generator settings.

| Name                           | Description                                                                                                  | Value                                                |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------- |
| `loadgen.api.enabled`          | Enable the API load generator job deploy.                                                                    | `true`                                               |
| `loadgen.api.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx).                     | `docker.elastic.co/observability-ci/opbeans-loadgen` |
| `loadgen.api.image.tag`        | Tag of the Docker image (latest).                                                                            | `daily`                                              |
| `loadgen.api.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                                         | `IfNotPresent`                                       |
| `loadgen.api.urls`             | List of URLs to load generate (opbeans-any:http://opbeans:3000,http://opbeans-1:3000,http://opbeans-2:3000). | `opbeans-any:http://opbeans:3000`                    |
| `loadgen.api.runLength`        | Number of run length per service to generate (opbeans-any:30).                                               | `opbeans-any:30`                                     |
| `loadgen.api.requestMinute`    | Number of request por minute per service to generate (opbeans-any:50).                                       | `opbeans-any:200`                                    |

### loadgen.frontend Frontend load generator settings.

| Name                                | Description                                                                              | Value                                                   |
| ----------------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| `loadgen.frontend.enabled`          | Enable the frontend load generator job deploy.                                           | `true`                                                  |
| `loadgen.frontend.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/functional-opbeans` |
| `loadgen.frontend.image.tag`        | Tag of the Docker image (latest).                                                        | `latest`                                                |
| `loadgen.frontend.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                          |
| `loadgen.frontend.url`              | URL of the frontend service to load generate (http://opbeans-frontend:3000).             | `http://opbeans-frontend:3000`                          |
| `loadgen.parallelism`               | Number of parallel jobs to run.                                                          | `1`                                                     |
| `loadgen.completions`               | Number of completions per job.                                                           | `1`                                                     |
| `loadgen.concurrencyPolicy`         | Concurrency policy for the job.                                                          | `Forbid`                                                |

### loadgen.molotov Molotov load generator settings.

| Name                        | Description                                                | Value                 |
| --------------------------- | ---------------------------------------------------------- | --------------------- |
| `loadgen.molotov.enabled`   | Enable the molotov load generator job deploy.              | `false`               |
| `loadgen.molotov.url`       | URL of the service to load generate (http://opbeans:3000). | `http://opbeans:3000` |
| `loadgen.molotov.duration`  | Duration of each round.                                    | `120`                 |
| `loadgen.molotov.tolerance` | Error tolerance percent.                                   | `75`                  |

### opbeans Opbeans service settings.

| Name                    | Description                              | Value                 |
| ----------------------- | ---------------------------------------- | --------------------- |
| `opbeans.apiServiceUrl` | URL of the opbeans service               | `http://opbeans:3000` |
| `opbeans.port`          | Opbeans service port.                    | `3000`                |
| `opbeans.dtProvability` | Distributed tracing provability.         | `0.5`                 |
| `opbeans.annotations`   | Opbeans deployment aditions annotations. | `{}`                  |

### Database settings.

| Name          | Description                         | Value        |
| ------------- | ----------------------------------- | ------------ |
| `db.enabled`  | Enable the database service deploy. | `false`      |
| `db.host`     | Database host.                      | `opbeans-db` |
| `db.name`     | Database name.                      | `opbeans`    |
| `db.port`     | Database port.                      | `5432`       |
| `db.type`     | Database type [mysql|postgresql]    | `postgresql` |
| `db.username` | Database Username.                  | `elastic`    |
| `db.password` | Database Password.                  | `none`       |

### redis Redis settings.

| Name                                    | Description                                                          | Value                       |
| --------------------------------------- | -------------------------------------------------------------------- | --------------------------- |
| `redis.url`                             | Redis URL.                                                           | `redis://redis-master:6379` |
| `extraEnvVars`                          | Array to add extra environment variables                             | `[]`                        |
| `podSecurityContext.enabled`            | Enabled Apache Server pods' Security Context                         | `false`                     |
| `podSecurityContext.fsGroup`            | Set Apache Server pod's Security Context fsGroup                     | `1001`                      |
| `containerSecurityContext.enabled`      | Enabled Apache Server containers' Security Context                   | `false`                     |
| `containerSecurityContext.runAsUser`    | Set Apache Server containers' Security Context runAsUser             | `1001`                      |
| `containerSecurityContext.runAsNonRoot` | Set Controller container's Security Context runAsNonRoot             | `true`                      |
| `command`                               | Override default container command (useful when using custom images) | `[]`                        |
| `args`                                  | Override default container args (useful when using custom images)    | `[]`                        |
| `labels`                                | to add to all deployed objects                                       | `{}`                        |
