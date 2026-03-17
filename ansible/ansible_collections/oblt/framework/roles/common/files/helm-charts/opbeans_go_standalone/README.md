# Opbeans Go Helm chart

Opbeans Go Helm chart deploy the [Opbeans Go](https://github.com/elastic/opbeans-go/)
and all services needed (PostgreSQL, and Redis).

## Parameters

### opbeans-go Opbeans Go Helm chart configuration

| Name                       | Description                                                                                     | Value                   |
| -------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------- |
| `opbeans.enabled`          | Enable or disable the whole kubernetes deployments, this allow to use it on parent charts.      | `true`                  |
| `opbeans.nameOverride`     | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.      | `""`                    |
| `opbeans.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-go-standalone` |

### opbeans.image Docker image settings

| Name                       | Description                                                                              | Value                                                                                                                                      |
| -------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `opbeans.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/opbeans-go`                                                                                            |
| `opbeans.image.tag`        | Tag of the Docker image (latest).                                                        | `daily`                                                                                                                                    |
| `opbeans.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                                                                                                             |
| `opbeans.command`          | Override default container command (useful when using custom images)                     | `["/opbeans-go"]`                                                                                                                          |
| `opbeans.args`             | Override default container args (useful when using custom images)                        | `["-log-level","debug","-log-json","-listen=:$(OPBEANS_SERVER_PORT)","-frontend=/opbeans-frontend","-db=postgres:","-cache=$(REDIS_URL)"]` |
| `opbeans.extraEnvVars`     | Array to add extra environment variables                                                 | `undefined`                                                                                                                                |

### apm APM configuration settings.

| Name                           | Description                                                   | Value                                   |
| ------------------------------ | ------------------------------------------------------------- | --------------------------------------- |
| `opbeans.apm.type`             | Type of reporting API [apm|otel|annotations]                  | `apm`                                   |
| `opbeans.apm.url`              | URL to the APM service.                                       | `http://opbeans-go-standalone-apm:8200` |
| `opbeans.apm.token`            | Token to authenticate against the APM service.                | `SuP3RSeCr3T`                           |
| `opbeans.apm.verifyServerCert` | Enable Certificates verification on TLS connections.          | `true`                                  |
| `opbeans.apm.serviceName`      | Name of the service deploy for the Helm chart (APM spans).    | `opbeans-go-standalone`                 |
| `opbeans.apm.serviceVersion`   | Version of the service deploy for the Helm chart (APM spans). | `1.0`                                   |
| `opbeans.apm.logLevel`         | Level of debug.                                               | `debug`                                 |
| `opbeans.apm.logFile`          | File to debug to.                                             | `stderr`                                |
| `opbeans.apm.inferredSpans`    | Enable inferred spans.                                        | `true`                                  |
| `opbeans.apm.logCorrelation`   | Enable log correlation.                                       | `true`                                  |
| `opbeans.apm.transactionRate`  | Configure the APM transaction rate (0.9).                     | `0.9`                                   |
| `opbeans.apm.environment`      | Environment where the service is deployed.                    | `production`                            |

### db Database settings.

| Name                  | Description                         | Value                      |
| --------------------- | ----------------------------------- | -------------------------- |
| `opbeans.db.enabled`  | Enable the database service deploy. | `true`                     |
| `opbeans.db.host`     | Database host.                      | `opbeans-go-standalone-db` |
| `opbeans.db.name`     | Database name.                      | `opbeans-go-standalone`    |
| `opbeans.db.port`     | Database port.                      | `5432`                     |
| `opbeans.db.type`     | Database type [mysql|postgresql]    | `postgresql`               |
| `opbeans.db.username` | Database Username.                  | `elastic`                  |
| `opbeans.db.password` | Database Password.                  | `none`                     |

### ingress Ingress settings


### opbeans.redis Redis settings.

| Name                | Description | Value                                             |
| ------------------- | ----------- | ------------------------------------------------- |
| `opbeans.redis.url` | Redis URL.  | `redis://opbeans-go-standalone-redis-master:6379` |

### postgresql PostgreSQL Helm chart configuration

| Name                          | Description                                                                                     | Value                      |
| ----------------------------- | ----------------------------------------------------------------------------------------------- | -------------------------- |
| `postgresql.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-go-standalone-db` |

### postgresql.auth PostgreSQL authentication settings

| Name                             | Description                                                                                                                                                                                                                                                                                                                                   | Value                            |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| `postgresql.auth.username`       | Name for a custom user to create                                                                                                                                                                                                                                                                                                              | `elastic`                        |
| `postgresql.auth.existingSecret` | Name of existing secret to use for PostgreSQL credentials. `auth.postgresPassword`, `auth.password`, and `auth.replicationPassword` will be ignored and picked up from this secret. The secret might also contains the key `ldap-password` if LDAP is enabled. `ldap.bind_password` will be ignored and picked from this secret in this case. | `opbeans-go-standalone-db-creds` |

### postgresql.primary

| Name                                     | Description                  | Value   |
| ---------------------------------------- | ---------------------------- | ------- |
| `postgresql.primary.persistence.enabled` | Enable persistence using PVC | `false` |

### postgresql.primary.initdb Scripts to run at first boot.

| Name                                         | Description                                               | Value                          |
| -------------------------------------------- | --------------------------------------------------------- | ------------------------------ |
| `postgresql.primary.initdb.scriptsConfigMap` | Name of existing ConfigMap containing the initdb scripts. | `opbeans-go-standalone-initdb` |

### redis Redis Helm chart configuration

| Name                     | Description                                                                                     | Value                         |
| ------------------------ | ----------------------------------------------------------------------------------------------- | ----------------------------- |
| `redis.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-go-standalone-redis` |
| `redis.auth.enabled`     | Enable password authentication                                                                  | `false`                       |

### redis.master Redis master settings


### redis.master.resources requests and limits

| Name                                     | Description                         | Value   |
| ---------------------------------------- | ----------------------------------- | ------- |
| `redis.master.resources.limits.memory`   | The memory limit for the Opbean     | `512Mi` |
| `redis.master.resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `redis.master.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |

### redis.master.persistence Redis replica persistence settings

| Name                               | Description                  | Value   |
| ---------------------------------- | ---------------------------- | ------- |
| `redis.master.persistence.enabled` | Enable persistence using PVC | `false` |

### redis.replica Redis replica settings

| Name                         | Description                  | Value |
| ---------------------------- | ---------------------------- | ----- |
| `redis.replica.replicaCount` | Number of replicas to deploy | `1`   |

### redis.replica.resources requests and limits

| Name                                      | Description                         | Value   |
| ----------------------------------------- | ----------------------------------- | ------- |
| `redis.replica.resources.limits.memory`   | The memory limit for the Opbean     | `512Mi` |
| `redis.replica.resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `redis.replica.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |

### redis.replica.persistence Redis replica persistence settings

| Name                                | Description                  | Value   |
| ----------------------------------- | ---------------------------- | ------- |
| `redis.replica.persistence.enabled` | Enable persistence using PVC | `false` |

### serviceAccount Service Account settings

| Name                    | Description                                          | Value                   |
| ----------------------- | ---------------------------------------------------- | ----------------------- |
| `serviceAccount.create` | Specifies whether a ServiceAccount should be created | `true`                  |
| `serviceAccount.name`   | The name of the ServiceAccount to use.               | `opbeans-go-standalone` |

### apm APM Server settings

| Name          | Description                   | Value  |
| ------------- | ----------------------------- | ------ |
| `apm.enabled` | Enable the APM Server deploy. | `true` |

### apm.image Docker image to deploy.

| Name                   | Description               | Value                              |
| ---------------------- | ------------------------- | ---------------------------------- |
| `apm.image.name`       | Docker image name.        | `docker.elastic.co/apm/apm-server` |
| `apm.image.tag`        | Docker image tag.         | `8.8.1`                            |
| `apm.image.pullPolicy` | Docker image pull policy. | `IfNotPresent`                     |

### apm.command Command to run.

| Name       | Description       | Value                           |
| ---------- | ----------------- | ------------------------------- |
| `apm.args` | Arguments to run. | `["--strict.perms=false","-e"]` |

### apm.elasticsearch Elasticsearch settings.

| Name                         | Description              | Value                               |
| ---------------------------- | ------------------------ | ----------------------------------- |
| `apm.elasticsearch.host`     | Elasticsearch hosts.     | `https://elasticsearch.example.com` |
| `apm.elasticsearch.username` | Elasticsearch username.  | `elastic`                           |
| `apm.elasticsearch.password` | Elasticsearch password.  | `changeme`                          |
| `apm.secretToken`            | APM Server secret token. | `SuP3RSeCr3T`                       |

### apm.service Kubernetes service settings.

| Name                           | Description                                                    | Value       |
| ------------------------------ | -------------------------------------------------------------- | ----------- |
| `apm.service.type`             | Kubernetes service type.                                       | `ClusterIP` |
| `apm.service.port`             | Kubernetes service port.                                       | `8200`      |
| `apm.podSecurityContext`       | Security context policies to add to the APM Server pods.       | `{}`        |
| `apm.containerSecurityContext` | Security context policies to add to the APM Server containers. | `{}`        |
| `apm.extraEnvVars`             | Extra environment variables to add to the APM Server pods.     | `[]`        |

### apm.resources requests and limits


### apm.imagePullSecrets Docker secrets used to pull the Docker images.

| Name                           | Description                | Value               |
| ------------------------------ | -------------------------- | ------------------- |
| `apm.imagePullSecrets[0].name` | Name of the Docker secret. | `docker.elastic.co` |
