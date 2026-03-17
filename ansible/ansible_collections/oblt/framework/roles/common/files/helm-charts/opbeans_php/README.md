# Opbeans Go Helm chart

Opbeans Go Helm chart deploy the [Opbeans Go](https://github.com/elastic/opbeans-php/)
and all services needed (PostgreSQL, and Redis).

## Parameters

### opbeans Opbeans Go Helm chart configuration

| Name                       | Description                                                                                     | Value  |
| -------------------------- | ----------------------------------------------------------------------------------------------- | ------ |
| `opbeans.enabled`          | Enable or disable the whole kubernetes deployments, this allow to use it on parent charts.      | `true` |
| `opbeans.nameOverride`     | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.      | `""`   |
| `opbeans.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `""`   |

### opbeans.image Docker image settings

| Name                       | Description                                                                              | Value                                            |
| -------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------ |
| `opbeans.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/opbeans-php` |
| `opbeans.image.tag`        | Tag of the Docker image (latest).                                                        | `daily`                                          |
| `opbeans.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                   |
| `opbeans.extraEnvVars`     | Array to add extra environment variables                                                 | `undefined`                                      |

### apm APM configuration settings.

| Name                           | Description                                                   | Value                            |
| ------------------------------ | ------------------------------------------------------------- | -------------------------------- |
| `opbeans.apm.type`             | Type of reporting API [apm|otel|annotations]                  | `apm`                            |
| `opbeans.apm.url`              | URL to the APM service.                                       | `https://apm.127.0.0.1.ip.es.io` |
| `opbeans.apm.token`            | Token to authenticate against the APM service.                | `SuP3RSeCr3T`                    |
| `opbeans.apm.verifyServerCert` | Enable Certificates verification on TLS connections.          | `true`                           |
| `opbeans.apm.serviceName`      | Name of the service deploy for the Helm chart (APM spans).    | `opbeans-php`                    |
| `opbeans.apm.serviceVersion`   | Version of the service deploy for the Helm chart (APM spans). | `1.0`                            |
| `opbeans.apm.logLevel`         | Level of debug.                                               | `debug`                          |
| `opbeans.apm.logFile`          | File to debug to.                                             | `stderr`                         |
| `opbeans.apm.inferredSpans`    | Enable inferred spans.                                        | `true`                           |
| `opbeans.apm.logCorrelation`   | Enable log correlation.                                       | `true`                           |
| `opbeans.apm.transactionRate`  | Configure the APM transaction rate (0.9).                     | `0.9`                            |
| `opbeans.apm.environment`      | Environment where the service is deployed.                    | `production`                     |
| `opbeans.apm.enabled`          | Enable to send data to the APM service.                       | `true`                           |

### db Database settings.

| Name                  | Description                         | Value            |
| --------------------- | ----------------------------------- | ---------------- |
| `opbeans.db.enabled`  | Enable the database service deploy. | `true`           |
| `opbeans.db.host`     | Database host.                      | `opbeans-php-db` |
| `opbeans.db.name`     | Database name.                      | `opbeans-php`    |
| `opbeans.db.port`     | Database port.                      | `5432`           |
| `opbeans.db.type`     | Database type [mysql|postgresql]    | `postgresql`     |
| `opbeans.db.username` | Database Username.                  | `elastic`        |
| `opbeans.db.password` | Database Password.                  | `none`           |

### ingress Ingress settings

| Name                                | Description                                    | Value                            |
| ----------------------------------- | ---------------------------------------------- | -------------------------------- |
| `opbeans.ingress.enabled`           | Enable or disable the ingress resource.        | `true`                           |
| `opbeans.ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`              |
| `opbeans.ingress.host`              | Hostname exposed by the ingress resource.      | `opbeans-php.127.0.0.1.ip.es.io` |
| `opbeans.ingress.annotations`       | Additional annotation to the ingress resource. | `undefined`                      |

### opbeans.redis Redis settings.

| Name                | Description | Value                                   |
| ------------------- | ----------- | --------------------------------------- |
| `opbeans.redis.url` | Redis URL.  | `redis://opbeans-php-redis-master:6379` |

### opbeans.resources requests and limits

| Name                                | Description                         | Value   |
| ----------------------------------- | ----------------------------------- | ------- |
| `opbeans.resources.limits.memory`   | The memory limit for the Opbean     | `1.5Gi` |
| `opbeans.resources.requests.memory` | The requested memory for the Opbean | `256Mi` |
| `opbeans.resources.requests.cpu`    | The requested cpu for the Opbean    | `0.1`   |

### postgresql PostgreSQL Helm chart configuration

| Name                          | Description                                                                                     | Value            |
| ----------------------------- | ----------------------------------------------------------------------------------------------- | ---------------- |
| `postgresql.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-php-db` |

### postgresql.auth PostgreSQL authentication settings

| Name                             | Description                                                                                                                                                                                                                                                                                                                                   | Value                  |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------- |
| `postgresql.auth.username`       | Name for a custom user to create                                                                                                                                                                                                                                                                                                              | `elastic`              |
| `postgresql.auth.existingSecret` | Name of existing secret to use for PostgreSQL credentials. `auth.postgresPassword`, `auth.password`, and `auth.replicationPassword` will be ignored and picked up from this secret. The secret might also contains the key `ldap-password` if LDAP is enabled. `ldap.bind_password` will be ignored and picked from this secret in this case. | `opbeans-php-db-creds` |

### postgresql.primary

| Name                                     | Description                  | Value   |
| ---------------------------------------- | ---------------------------- | ------- |
| `postgresql.primary.persistence.enabled` | Enable persistence using PVC | `false` |

### postgresql.primary.initdb Scripts to run at first boot.

| Name                                         | Description                                               | Value                |
| -------------------------------------------- | --------------------------------------------------------- | -------------------- |
| `postgresql.primary.initdb.scriptsConfigMap` | Name of existing ConfigMap containing the initdb scripts. | `opbeans-php-initdb` |
