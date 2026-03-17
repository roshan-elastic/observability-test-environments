---
render_macros: false
---
# Cluster Configuration

The cluster configuration file is a YAML file that contains the settings for the cluster.
It contains the Elastic Stack deployment configuration and the Kubernetes cluster settings.
It also can contains the configuration for the applications to deploy.
The basic structure of the configuration file is as follows:

```yaml
cluster_name: "my-unique-fqdn-name"
slack_channel: "#observablt-bots"
oblt_username: "robots"
stack:
  version: "7.15.0"
```

## Elastic Stack

The Elastic Stack configuration defines the settings for the Elastic Stack deployment.
It contains the settings for the Elastic Cloud deployment and the applications to deploy.

```yaml
stack:
  version: "7.15.0"
  mode: "ess"
  template: "gcp-io-optimized-v3"
  target: "production
```

[Elastic Stack parameters](#elastic-stack_1)

## Kubernetes

`k8s` section defines the settings for the Kubernetes cluster.
By default [Kind][] is used to provision Kubernetes clusters
so the `provider` is set to `kind`.

```yaml
k8s:
  enabled: true
  provider: kind
  region: "none"
  domain: "ip.es.io"
  default_namespace: "default"
  project: "elastic-observability"
```

[Kubernetes parameters](#kubernetes_1)

## Bootstrap cluster

In some cases, we have to make some configuration using the API,
for those cases we have the option to bootstrap the cluster using files (recipes) that define the API call.
The implemented a process that reads the files in a folder and execute the API calls defined in those files.
this folder will contain a folder from elasticsearch, kibana, apm, and fleet API REST calls.
folder:

* elasticsearch
* kibana
* apm
* fleet
* ess

All `yml` files in those folder will be processed.

This is an example of request that return the Elasticsearch version info.

```yaml
---
# Description of what the recipe do.
description: Basic test request
# Path to the REST API call.
api: "/"
#HTTP method to use it can be GET, POST, PUT, DELETE, ...
method: GET
# HTTP header we need to pass.
headers:
  Content-Type: application/json
# The payload to send to the REST API.
body: ""
# The expected return code.
return_code: 200
```

To enable the bootstrap process in a cluster,
we have to add a section where you define which folders contain the bootstrap recipes.

```yaml
bootstrap_cluster: true
bootstrap_cluster_dirs:
  - {{ deployments_dir }}/bootstrap-ess
  - {{ deployments_dir }}/bootstrap-fleet
```

for more information check the [Bootstrap cluster](./bootstrap.md) documentation.

## Kibana index backup/restore

The Kibana index is migrated every time we update to a new version of the Elastic Stack,
due to we use SNAPSHOTs it is common that we hit issues related with the Kibana index migration process,
we are using not final versions of Kibana Object that change during the develop.
To allow keep some of those object or completely ignore them, we have added a couple of settings to backup the Kibana index or purge it on updates.

```yaml
kibana_backup_enabled: false
kibana_purge_enabled: false
```

The process creates a new index `kibana_reindex_temp` to copy all the objects from the current Kibana index,
then delete the  `.kibana*` indices. When Kibana starts again it will create a new `.kibana` index valid for the new version,
then we reindex the objects in `kibana_reindex_temp` to the `.kibana` index.
Finally we removed the `kibana_reindex_temp`.

In purge mode, the process makes a backup of the Kibana objects in the index `kibana_reindex_temp`
and delete all the `.kibana*` indices.

## Observability

Kibana and APM support sending logs, metrics, and APM data to another Elasticsearch cluster.
To enable this feature you have to set `observability_enabled` to `true`,
and `observability_secret` to point to the Vault folder with the oblt cluster secrets you want to use to send the data.
By default `observability_secret` points to the `monitoring-oblt`,
so enabling `observability_enabled` is enough to start sending data to `monitoring-oblt`.

```yaml
observability_enabled: false
observability_secret: "oblt-clusters_monitoring-oblt_cluster-state"
```

## Remote cluster configured using the API

There are some cases where the remote clusters cannot be configured using terraform.
For those cases we need to use the Elasticsearch API.
The following settings allow that configuration.

```yaml
remote_cluster_configure: true
remote_clusters:
  - host: 365b555050af4980872d3e88c6c4cd69.us-central1.gcp.foundit.no
    port: 9400
    name: remote_overview_01
  - host: 365b555050af4980872d3e88c6c4cd50.us-central1.gcp.foundit.no
    port: 9400
    name: remote_overview_02
```

* remote_cluster_configure: Enable the configuration of a remote cluster using the API.
* remote_clusters: It is the list of remote clusters we want to configure. Each remote cluster has the following settings.
  * host: Host name of the remote cluster.
  * port: Port to connect in the remote cluster.
  * name: Name for the remote cluster in the Kibana configuration.

## Kibana Config

The Kibana configuration is stored in a secret in the Vault.

```yaml
kibana_yml_gcsm_def_secret: "oblt-clusters_CLUSTER_NAME_kibana-yml"
```

## Deploy Info

The deploy info is stored in a secret in the Vault.

```yaml
deploy_info_gcsm_def_secret: "oblt-clusters_CLUSTER_NAME_deploy-info"
```

[Kind]: https://kind.sigs.k8s.io/

## Parameters

### Configuration

Cluster configuration file parameters confire the type of cluster to create.
It is a YAML file that contains the settings for the cluster.
It contains the Elastic Stack deployment configuration and the Kubernetes cluster settings.
The configuration file define the applications to the ploy and the configuration for each application.
The basic configuration file contains a name for the cluster an the type of stack to deploy.

```yaml
cluster_name: "test-oblt"
stack:
version: "8.13.2"
```

The main sections are the following:

* `cluster_name`: Name of the cluster (e.g. sample-oblt)
* `stack`: Elastic Stack deployment configuration
* `k8s`: Kubernetes cluster settings
* `features`: List of features the cluster should deploy.

| Name                                | Description                                                                                 | Value              |
| ----------------------------------- | ------------------------------------------------------------------------------------------- | ------------------ |
| `cluster_name`                      | Name of the cluster (e.g. sample-oblt)                                                      | `""`               |
| `create_users`                      | True to create users for the cluster (default: true)                                        | `false`            |
| `grab_cluster_info`                 | True to grab cluster information (default: true)                                            | `false`            |
| `configure_external_snapshots_repo` | True to configure an external snapshot repository (default: false)                          | `false`            |
| `ess_not_allow_settings`            | True to disable some settings not allowed in staging.                                       | `true`             |
| `build_tf_provider_from_src`        | True to build the Elastic terraform provider from source (default: fatse).                  | `false`            |
| `use_snapshot_tf_provider`          | True to use a snapshot version of the Elastic terraform provider (default: false)           | `true`             |
| `ess_terraform_version`             | Version of the Elastic terraform provider to use (default: 0.5.0)                           | `""`               |
| `observability_secret`              | Secret name for the observability cluster (default: oblt-clusters_monitoring-oblt)          | `""`               |
| `secret_ess_user`                   | path to the secret that contains the ESS credentials (elastic-cloud-observability-team-pro) | `""`               |
| `golden_cluster`                    | True to create a golden cluster (default: false)                                            | `true`             |
| `slack_channel`                     | Slack channel to use for notifications                                                      | `#observablt-bots` |
| `oblt_username`                     | Username to use for the observability team (default: robots)                                | `robots`           |
| `expire_date`                       | Date to expire the cluster (e.g. 2023-12-31).                                               | `""`               |
| `bootstrap_cluster`                 | True to bootstrap the cluster (default: false) see [bootstrap](./bootstrap.md)              | `false`            |
| `bootstrap_cluster_dirs`            | Directories where the bootstrap cluster are located.                                        | `[]`               |
| `profiling_agent_deploy`            | True to deploy the profiling agent in the k8s cluster (default: true)                       | `true`             |

### Kubernetes

The Kubernetes configuration defines the settings for the Kubernetes cluster.
It is optional by default gcp with default settings.

| Name                        | Description                                                                                                                            | Value   |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `k8s`                       | Kubernetes cluster settings                                                                                                            |         |
| `k8s.cluster_name`          | Name of the k8s cluster, it replaces de default name (default: {{ cluster_name }})                                                     | `""`    |
| `k8s.certmanager`           | True to enable cert-manager                                                                                                            | `false` |
| `k8s.default_namespace`     | Default namespace for all resources                                                                                                    | `""`    |
| `k8s.domain`                | Domain name for the cluster                                                                                                            | `""`    |
| `k8s.enabled`               | True to enable Kubernetes cluster creation                                                                                             | `true`  |
| `k8s.ingress`               | True to enable ingress controller                                                                                                      | `false` |
| `k8s.kube_system_namespace` | Namespace for the kube-system resources                                                                                                | `""`    |
| `k8s.machine_type`          | Machine type, e.g. n1-standard-4                                                                                                       | `""`    |
| `k8s.max_node_count`        | Maximum number of nodes                                                                                                                | `10`    |
| `k8s.project`               | Project name                                                                                                                           | `""`    |
| `k8s.provider`              | Cloud provider, e.g. gcp, kind                                                                                                         | `""`    |
| `k8s.region`                | Cloud region, e.g. us-east-1, europe-west1-c                                                                                           | `""`    |
| `k8s.shared`                | Name of the shared cluster.This is used to share the same cluster between different deployments, when is set a cluster is not created. | `""`    |
| `k8s.static_ip`             | True to enable static IP                                                                                                               | `false` |
| `k8s.certificate_issuer`    | Issuer of the certificate, e.g. letsencrypt-prod, letsencrypt-staging (default: letsencrypt-staging)                                   | `""`    |
| `k8s.network`               | Network to use for the cluster                                                                                                         | `""`    |

### Elastic Stack

The Elastic Stack configuration defines the settings for the Elastic Stack deployment.

| Name                     | Description                                                                                                       | Value  |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------- | ------ |
| `stack`                  | Elastic Stack deployment settings                                                                                 |        |
| `stack.version`          | Version of the stack                                                                                              | `""`   |
| `stack.mode`             | Stack provider, e.g. ess, eck, serverless, external (default: ess)                                                | `""`   |
| `stack.template`         | Name of the template to use for the deployment creation (e.g. elasticsearch, kibana, observability, ccs)          | `""`   |
| `stack.target`           | Target environment, e.g. production, staging (default: production)                                                | `""`   |
| `stack.observability`    | True to enable observability features to send logs, metrics, and APM to the monitoring cluster (default: false)   | `true` |
| `stack.update_channel`   | Update channel to use for the stack, e.g. update_channel: unstable, development, release                          | `""`   |
| `stack.update_mode`      | Update mode to use for the cluster, e.g. update, recreate (default: update)                                       | `""`   |
| `stack.update_schedule`  | Update schedule to use for the cluster, e.g. daily, Monday, Tuesday, Wednesday, Thursday, Friday (default: daily) | `""`   |
| `stack.package_registry` | Package registry to use for the stack.                                                                            | `""`   |

### ESS

It requires a Elastic Cloud credential stored in the Vault,
the vault secret must have the following fields:

```ini
"apiKey": "OWOWOWOWOOOWOWOWOWOWOWOWOWOWOWWOOWWOOWWOOWWOOWOWOWOWWO==",
"description": "user to access to the Elastic Cloud service as a regular user",
"endpoint": "https://staging.found.no",
"region": "us-east-1",
"username": "user@elastic.co"
```


| Name                                   | Description                                                                                                                                        | Value   |
| -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `ess`                                  | ESS deployment settings                                                                                                                            |         |
| `stack.ess.saml_id`                    | Service Provider (SP) ID configured in ESS to use for SSO authentication (optional)                                                                | `""`    |
| `stack.ess.autoscale`                  | True to enable autoscaling (default: true)                                                                                                         | `true`  |
| `stack.ess.ccs_remote_clusters`        | List of remote clusters to configure CCS (optional)                                                                                                | `[]`    |
| `stack.ess.elasticsearch.enabled`      | True to enable Elasticsearch deployment (default: true)                                                                                            | `true`  |
| `stack.ess.elasticsearch.extra_config` | Extra configuration for Elasticsearch (optional)                                                                                                   | `{}`    |
| `stack.ess.elasticsearch.image`        | Docker image for Elasticsearch (optional)                                                                                                          | `""`    |
| `stack.ess.elasticsearch.mem`          | Memory for Elasticsearch (default: 16)(ignored on autoscaling enabled)                                                                             | `16`    |
| `stack.ess.elasticsearch.zones`        | Number of zones to use for the Elasticsearch deployment (default: 3)                                                                               | `3`     |
| `stack.ess.integrations.enabled`       | True to enable Integrations server deployment (default: true)                                                                                      | `true`  |
| `stack.ess.integrations.extra_config`  | Extra configuration for Integrations server (optional)                                                                                             | `{}`    |
| `stack.ess.integrations.image`         | Docker image for Elastic Agent (optional)                                                                                                          | `""`    |
| `stack.ess.integrations.mem`           | Memory for Integrations server (default: 2)(ignored on autoscaling enabled)                                                                        | `2`     |
| `stack.ess.integrations.zones`         | Number of zones to use for the Integrations serve deployment (default: 1)                                                                          | `1`     |
| `stack.ess.kibana.enabled`             | True to enable Kibana deployment (default: true)                                                                                                   | `true`  |
| `stack.ess.kibana.extra_config`        | Extra configuration for Kibana (optional)                                                                                                          | `{}`    |
| `stack.ess.kibana.further_details`     | Message to show in the Kibana login page (optional)                                                                                                | `""`    |
| `stack.ess.kibana.image`               | Docker image for Kibana (optional)                                                                                                                 | `""`    |
| `stack.ess.kibana.mem`                 | Memory for Kibana (default: 2)(ignored on autoscaling enabled)                                                                                     | `2`     |
| `stack.ess.kibana.zones`               | Number of zones to use for the Kibana deployment (default: 1)                                                                                      | `1`     |
| `stack.ess.kibana.open_ai_secret`      | Vault secret with the OpenaAI connection information '{"apiKey":"XXXXXXX", "service": "azureOpenAI"}', '{"apiKey":"XXXXXXX", "service": "openAI"}' | `""`    |
| `stack.ess.provider`                   | Cloud provider, e.g. gcp, aws, azure (default: gcp)                                                                                                | `""`    |
| `stack.ess.region`                     | Cloud region, e.g. us-east-1, europe-west1-c, gcp-us-west2, azure-eastus2,  (default: us-central1-c)                                               | `""`    |
| `stack.ess.restore_snapshot_from`      | Restore a snapshot from a given deployment ID (optional)                                                                                           | `""`    |
| `stack.ess.restore_snapshot_name`      | Name of the snapshot to restore (optional)(default: __latest_success__)                                                                            | `""`    |
| `stack.ess.template`                   | Name of the ESS hardware template to use for the deployment creation (default: gcp-io-optimized-v3)                                                | `""`    |
| `stack.ess.trusted_org`                | Trusted ESS organization (Optional)                                                                                                                | `""`    |
| `stack.ess.reset_password`             | True to reset the password for the ESS deployment (default: false)                                                                                 | `false` |

### Elastic Cloud on Kubernetes (ECK)

The Elastic Cloud on Kubernetes (ECK) configuration defines the settings for the ECK deployment.

| Name                                   | Description                                                             | Value  |
| -------------------------------------- | ----------------------------------------------------------------------- | ------ |
| `eck`                                  | ECK deployment settings                                                 |        |
| `stack.eck.agent.enabled`              | True to enable Elastic Agent deployment (default: true)                 | `true` |
| `stack.eck.agent.extra_config`         | Extra configuration for Elastic Agent (optional)                        | `{}`   |
| `stack.eck.agent.image`                | Docker image for Elastic Agent (optional)                               | `""`   |
| `stack.eck.agent.mem`                  | Memory for Elastic Agent (default: 2)                                   | `2`    |
| `stack.eck.elasticsearch.enabled`      | True to enable Elasticsearch deployment (default: true)                 | `true` |
| `stack.eck.elasticsearch.extra_config` | Extra configuration for Elasticsearch (optional)                        | `{}`   |
| `stack.eck.elasticsearch.image`        | Docker image for Elasticsearch (optional)                               | `""`   |
| `stack.eck.elasticsearch.mem`          | Memory for Elasticsearch (default: 16)                                  | `16`   |
| `stack.eck.elasticsearch.node_count`   | Number of nodes for the Elasticsearch deployment (default: 3)           | `3`    |
| `stack.eck.elasticsearch.node_storage` | Size of the storage for each node (default: 256)                        | `256`  |
| `stack.eck.kibana.enabled`             | True to enable Kibana deployment (default: true)                        | `true` |
| `stack.eck.kibana.extra_config`        | Extra configuration for Kibana (optional)                               | `{}`   |
| `stack.eck.kibana.further_details`     | Message to show in the Kibana login page (optional)                     | `""`   |
| `stack.eck.kibana.image`               | Docker image for Kibana (optional)                                      | `""`   |
| `stack.eck.kibana.mem`                 | Memory for Kibana (default: 2)                                          | `2`    |
| `stack.eck.license`                    | License to use for ECK (default: trial) ['enterprise', 'trial', 'none'] | `""`   |

### External

The External configuration defines the settings for the external Stack deployment.

| Name                    | Description                                                      | Value |
| ----------------------- | ---------------------------------------------------------------- | ----- |
| `stack.external.secret` | Vault secret with the fields to the external Stack (default: '') | `""`  |

### Serverless

The Serverless configuration defines the settings for the Serverless deployment.

| Name                                   | Description                                                                                                            | Value   |
| -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ------- |
| `serverless`                           | Serverless deployment settings                                                                                         |         |
| `stack.serverless.region`              | Cloud region, e.g. gcp-us-central1, aws-eu-west-1, azure-eastus2,  (default: aws-eu-west-1)                            | `""`    |
| `stack.serverless.elasticsearch.image` | Docker image for Elasticsearch (optional)                                                                              | `""`    |
| `stack.serverless.kibana.image`        | Docker image for Kibana (optional)                                                                                     | `""`    |
| `stack.serverless.fleet.image`         | Docker image for Fleet (optional)                                                                                      | `""`    |
| `stack.serverless.cluster`             | ID of the cluster to use for the deployment creation (e.g prd-awsuse1-cp-internal-app-1, qa-awsuse1-cp-internal-app-1) | `""`    |
| `stack.ess.profiling`                  | Configuration for Profiling                                                                                            | `{}`    |
| `stack.ess.profiling.enabled`          | True to enable profiling (default: false)                                                                              | `false` |
| `stack.ess.profiling.image`            | Docker image for profiling (optional)                                                                                  | `""`    |

### Applications

The applications configuration defines which applications to deploy in the cluster.

| Name                       | Description                                         | Value |
| -------------------------- | --------------------------------------------------- | ----- |
| `apps`                     | Applications to deploy                              |       |
| `apps.helm`                | List of Helm charts to deploy                       |       |
| `apps.helm[0].chart`       | Chart to deploy                                     | `""`  |
| `apps.helm[0].extra_args`  | Extra arguments to pass to the Helm install command | `{}`  |
| `apps.helm[0].name`        | Name of the Helm chart                              | `""`  |
| `apps.helm[0].namespace`   | Namespace to deploy the chart                       | `""`  |
| `apps.helm[0].values_file` | Values file for the Helm chart                      | `""`  |
| `apps.helm[0].version`     | Version of the chart                                | `""`  |
| `apps.k8s`                 | List of Kubernetes manifests to deploy              |       |
| `apps.k8s[0].name`         | Name of the Kubernetes manifest                     | `""`  |
| `apps.k8s[0].src`          | Source directory for the Kubernetes manifest        | `""`  |

### Opbeans

The Opbeans configuration defines the settings for the Opbeans deployment.

| Name              | Description                             | Value   |
| ----------------- | --------------------------------------- | ------- |
| `opbeans.enabled` | True to deploy Opbeans (default: false) | `false` |

### Synthetics monitors

```yaml
- id: release-oblt
folder: "{{ synthetics_deployments_dir }}/{{ cluster_name }}"
locale: "en-US"
timezone: "America/Los_Angeles"
params:
ELASTICSEARCH_URL: "https://elasticsearch.example.com"
APM_URL: "https://apm.example.com"
KIBANA_URL: "https://kibana.example.com"
KIBANA_USERNAME: "elastic"
KIBANA_PASSWORD: "changeme"
```

| Name                                      | Description                            | Value |
| ----------------------------------------- | -------------------------------------- | ----- |
| `synthetics_monitors_folders`             | List of Synthetics monitors to deploy  |       |
| `synthetics_monitors_folders[0].id`       | Id of the monitor                      | `""`  |
| `synthetics_monitors_folders[0].folder`   | Folder where the monitor is located    | `""`  |
| `synthetics_monitors_folders[0].locale`   | Locale to use for the monitor          | `""`  |
| `synthetics_monitors_folders[0].timezone` | Timezone to use for the monitor        | `""`  |
| `synthetics_monitors_folders[0].params`   | Parameters to use for the monitor      | `{}`  |
| `load_balancer_ip`                        | Loadbalancer IP to use for the cluster | `""`  |

### Features

The features configuration defines the features to deploy in the cluster.
List of features to deploy in the cluster.
The features are a list of objects with the following fields:

* `name`: Name of the feature
* `enabled`: True to enable the feature
* `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


The features configuration defines the features to deploy in the cluster.
List of features to deploy in the cluster.
The features are a list of objects with the following fields:

* `name`: Name of the feature
* `enabled`: True to enable the feature
* `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


The features configuration defines the features to deploy in the cluster.
List of features to deploy in the cluster.
The features are a list of objects with the following fields:

* `name`: Name of the feature
* `enabled`: True to enable the feature
* `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


The features configuration defines the features to deploy in the cluster.
List of features to deploy in the cluster.
The features are a list of objects with the following fields:

* `name`: Name of the feature
* `enabled`: True to enable the feature
* `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


The features configuration defines the features to deploy in the cluster.
List of features to deploy in the cluster.
The features are a list of objects with the following fields:

* `name`: Name of the feature
* `enabled`: True to enable the feature
* `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```


List of features to deploy in the cluster.
The features are a list of objects with the following fields:
- `name`: Name of the feature
- `enabled`: True to enable the feature
- `parameters`: Parameters to use for the feature

Each feature define the parameters to use for the feature.
you can find details about the parameters in the feature documentation.

```yaml
features:
- name: "feature01"
enabled: true
- name: "feature02"
enabled: true
parameter00: "value00"
parameter01: "value01"
```
