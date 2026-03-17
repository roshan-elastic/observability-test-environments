# Opbeans Python Helm chart

Opbeans Python Helm chart deploy the [Opbeans Python](https://github.com/elastic/opbeans-node/)
and all services needed (PostgreSQL, Redis, and Elasticseacrh).

## Parameters

### opbeans Opbeans Python Helm chart configuration

| Name                       | Description                                                                                    | Value  |
| -------------------------- | ---------------------------------------------------------------------------------------------- | ------ |
| `opbeans.enabled`          | Enable or disable the whole kubernetes deployments, this allow to use it on parent charts.     | `true` |
| `opbeans.nameOverride`     | Overrides the the Helm Chart name used to naming kubernetes objects, by default `opbeans`.     | `""`   |
| `opbeans.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeans`. | `""`   |

### opbeans.image Docker image settings

| Name                       | Description                                                                              | Value                                             |
| -------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------- |
| `opbeans.image.name`       | Name of the Docker image (docker.elastic.co/observability-ci/it_opbeans-frontend_nginx). | `docker.elastic.co/observability-ci/opbeans-node` |
| `opbeans.image.tag`        | Tag of the Docker image (latest).                                                        | `daily`                                           |
| `opbeans.image.pullPolicy` | Docker image pull policy (IfNotPresent|Always|Never)                                     | `IfNotPresent`                                    |
| `opbeans.extraEnvVars`     | Array to add extra environment variables                                                 | `undefined`                                       |

### apm APM configuration settings.

| Name                         | Description                                   | Value         |
| ---------------------------- | --------------------------------------------- | ------------- |
| `opbeans.apm.type`           | Type of reporting API [apm|otel|annotations]  | `annotations` |
| `opbeans.apm.annotationType` | Type of annotations to use [none|nodejs|java] | `nodejs`      |
| `opbeans.apm.enabled`        | Enable to send data to the APM service.       | `true`        |

### db Database settings.

| Name                  | Description                         | Value             |
| --------------------- | ----------------------------------- | ----------------- |
| `opbeans.db.enabled`  | Enable the database service deploy. | `true`            |
| `opbeans.db.host`     | Database host.                      | `opbeans-node-db` |
| `opbeans.db.name`     | Database name.                      | `opbeans-node`    |
| `opbeans.db.port`     | Database port.                      | `5432`            |
| `opbeans.db.type`     | Database type [mysql|postgresql]    | `postgresql`      |
| `opbeans.db.username` | Database Username.                  | `elastic`         |
| `opbeans.db.password` | Database Password.                  | `none`            |

### ingress Ingress settings

| Name                                | Description                                    | Value                             |
| ----------------------------------- | ---------------------------------------------- | --------------------------------- |
| `opbeans.ingress.enabled`           | Enable or disable the ingress resource.        | `true`                            |
| `opbeans.ingress.certificateIssuer` | Name of the certificates issuer.               | `letsencrypt-stage`               |
| `opbeans.ingress.host`              | Hostname exposed by the ingress resource.      | `opbeans-node.127.0.0.1.ip.es.io` |
| `opbeans.ingress.annotations`       | Additional annotation to the ingress resource. | `undefined`                       |

### opbeans.redis Redis settings.

| Name                | Description | Value                                    |
| ------------------- | ----------- | ---------------------------------------- |
| `opbeans.redis.url` | Redis URL.  | `redis://opbeans-node-redis-master:6379` |

### postgresql PostgreSQL Helm chart configuration

| Name                          | Description                                                                                     | Value             |
| ----------------------------- | ----------------------------------------------------------------------------------------------- | ----------------- |
| `postgresql.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-node-db` |

### postgresql.auth PostgreSQL authentication settings

| Name                             | Description                                                                                                                                                                                                                                                                                                                                   | Value                   |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| `postgresql.auth.username`       | Name for a custom user to create                                                                                                                                                                                                                                                                                                              | `elastic`               |
| `postgresql.auth.existingSecret` | Name of existing secret to use for PostgreSQL credentials. `auth.postgresPassword`, `auth.password`, and `auth.replicationPassword` will be ignored and picked up from this secret. The secret might also contains the key `ldap-password` if LDAP is enabled. `ldap.bind_password` will be ignored and picked from this secret in this case. | `opbeans-node-db-creds` |

### postgresql.primary

| Name                                     | Description                  | Value   |
| ---------------------------------------- | ---------------------------- | ------- |
| `postgresql.primary.persistence.enabled` | Enable persistence using PVC | `false` |

### postgresql.primary.initdb Scripts to run at first boot.

| Name                                         | Description                                               | Value                 |
| -------------------------------------------- | --------------------------------------------------------- | --------------------- |
| `postgresql.primary.initdb.scriptsConfigMap` | Name of existing ConfigMap containing the initdb scripts. | `opbeans-node-initdb` |

### redis Redis Helm chart configuration

| Name                     | Description                                                                                     | Value                |
| ------------------------ | ----------------------------------------------------------------------------------------------- | -------------------- |
| `redis.fullnameOverride` | Overrides the the Helm Chart fullname used to naming kubernetes objects, by default `opbeanst`. | `opbeans-node-redis` |
| `redis.auth.enabled`     | Enable password authentication                                                                  | `false`              |

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
