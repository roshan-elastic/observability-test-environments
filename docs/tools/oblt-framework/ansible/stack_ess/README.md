---
render_macros: false
---
# stack_ess role

## Overview

This Role deploy a the Elastic Stack on ESS.
It uses the [Terraform Provider for Elastic Cloud][] to make the deployment.
If k8s is enabled will create the Elastic Stack secret related in the k8s cluster.

The credentials to access the cluster will be stored in the Vault on the following path:

* observability-team/ci/test-clusters/{{ cluster_name }}/tfstate
This secrets also contains the Terraform state info.

## Requirements

It requires a Elastic Cloud credential stored in the Google Cloud Secret Manager (elastic-cloud-observability-team-ENVIRONMENT),
the vault secret must have the following fields:

```json
{
  'apiKey': 'essu_AAAAAAAAAAAAAAAAAAAAAAAAAA',
  'apiKeyId': '3ZbwTo4BGlaLKR2-GoqA',
  'date': '20240318T002219',
  'description': 'User to access to the Elastic Cloud service as a regular user',
  'endpoint': 'https://global.elastic.cloud',
  'ess_endpoint': 'https://cloud.elastic.co',
  'organizationID': '2870499056',
  'serverless_endpoint':
  'https://global.elastic.cloud',
  'region': 'us-east-1',
  'username': 'observability-robots@elastic.co',
  'password': 'PaSWoRD',
  'console': 'https://admin.found.no',
  'consoleApiKey': 'essa_AAAAAAAAAAAAAAAAAAAAAAAAAA'
}
```

## Examples

```yaml
stack:
  version: "8.14.0-SNAPSHOT"
  mode: "ess"
  template: "elasticsearch"
  target: "production"
  ess:
    provider: gcp
    elasticsearch:
      zones: 1
```

```yaml
stack:
  version: "8.14.0-SNAPSHOT"
  mode: "ess"
  template: "elasticsearch"
  target: "staging"
  ess:
    elasticsearch:
      zones: 1
    provider: aws
```

```yaml
stack:
  version: "8.14.0-SNAPSHOT"
  mode: "ess"
  template: "elasticsearch"
  target: "staging"
  ess:
    elasticsearch:
      zones: 1
    provider: azure
```


```yaml
stack:
  version: "8.14.0-SNAPSHOT"
  mode: "ess"
  template: "elasticsearch"
  target: "production"
  observability: true
  ess:
    elasticsearch:
      zones: 1
```

## Dependencies

It includes the [common][] roles.

## License

Apache License 2.0

[common]:../common/README.md
[Terraform Provider for Elastic Cloud]: https://github.com/elastic/terraform-provider-ec

## Parameters

### stack_ess Elasticsearch Service (ESS) configuration

| Name                                | Description                                                                                                                     | Value                                                                                                                                                               |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ess_tf_dir`                        | The directory where the Terraform files are located.                                                                            | `{{ build_dir }}/terraform/elastic_cloud`                                                                                                                           |
| `ess_stack_version`                 | The version of the stack to deploy.                                                                                             | `{{ default_stack_version }}`                                                                                                                                       |
| `ess_ccs_enabled`                   | True to enable Cross Cluster Search (CCS).                                                                                      | `{{ ess_ccs_remote_cluster | length > 1 }}`                                                                                                                         |
| `ess_ccs_remote_cluster_name_vault` | The name of the secret in Google Cloud Secret Manager that contains the remote cluster name.                                    | `oblt-clusters_{{ ess_ccs_remote_cluster }}_cluster-state`                                                                                                          |
| `ess_ccs_remote_cluster_id`         | The ID of the remote cluster.                                                                                                   | `{{ (lookup('oblt.framework.gcp_secret_manager', key=ess_ccs_remote_cluster_name_vault, on_error='warn')|safe|from_yaml).ess_deployment_id | default(None,true) }}` |
| `ess_not_allow_settings`            | True to disable some settings not allowed in staging.                                                                           | `true`                                                                                                                                                              |
| `ess_credentials`                   | The credentials to access the Elastic Cloud service.                                                                            | `{{ lookup('oblt.framework.gcp_secret_manager', key=secret_ess_user)|safe()|from_yaml }}`                                                                           |
| `ess_api_key`                       | The API key to access the Elastic Cloud service                                                                                 | `{{ ess_credentials.apiKey | mandatory }}`                                                                                                                          |
| `ess_org_id`                        | The organization ID to access the Elastic Cloud service.                                                                        | `{{ ess_credentials.organizationID | default(None,true) }}`                                                                                                         |
| `ess_custom_settings`               | True to enable custom settings.                                                                                                 | `true`                                                                                                                                                              |
| `ess_gcs_credentials`               | The name of the secret in Google Cloud Secret Manager that contains the credentials to access the Google Cloud Storage service. | `oblt-clusters_gcs-credentials`                                                                                                                                     |
