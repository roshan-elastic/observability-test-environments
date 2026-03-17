---
render_macros: false
---
# stack_eck Role

## Overview

This Role deploys and configures an Elastic Cloud on Kubernetes (ECK) cluster,
The kubernetes cluster should exist.
It has also tasks to uninstall the services installed on a Kubernetes cluster.

The credentials to access the cluster will stored in the Vault on the following paths:

* observability-team/ci/test-clusters/{{ cluster_name }}/eck-elasticsearch
* observability-team/ci/test-clusters/{{ cluster_name }}/eck-kibana
* observability-team/ci/test-clusters/{{ cluster_name }}/eck-apm

## Requirements

It requires to have kubectl, and Helm CLI installed.
It requires an Elastic Cloud credential stored in the Google Cloud Secret Manager,
see default variables for more details.

## Dependencies

It includes [common][], and [k8s][] role for testing.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    eck:
      version: 1.0.0
      apiVersion: v1
    elasticsearch:
      enabled: true
      version: 7.4.0
      image: elastic/elasticsearch:7.4.0
      nodes: 1
      mem: 2
      storage: 10Gi
    kibana:
      enabled: true
      version: 7.4.0
      image: elastic/kibana:7.4.0
      mem: 2
    apm:
      enabled: true
      version: 7.4.0
      image: elastic/apm-server:7.4.0
      mem: 0.5
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.k8s
    - role: stack_eck
```

## License

Apache License 2.0

[common]:../common/README.md
[k8s]:../k8s/README.md

## Parameters

### General

| Name                      | Description                                            | Value  |
| ------------------------- | ------------------------------------------------------ | ------ |
| `eck_version`             | The version of the Elastic Cloud on Kubernetes         | `""`   |
| `eck_api_version`         | The API version of the Elastic Cloud on Kubernetes API | `""`   |
| `eck_domain`              | The domain used for URLs                               | `""`   |
| `elasticsearch_host`      | The external URL for Elasticsearch                     | `""`   |
| `elasticsearch_host_int`  | The internal URL for Elasticsearch                     | `""`   |
| `elasticsearch_port_int`  | The internal port for Elasticsearch                    | `9200` |
| `kibana_host`             | The external URL for Kibana                            | `""`   |
| `kibana_host_int`         | The internal URL for Kibana                            | `""`   |
| `kibana_port_int`         | The internal port for Kibana                           | `5601` |
| `apm_host`                | The external URL for APM                               | `""`   |
| `apm_host_int`            | The internal URL for APM                               | `""`   |
| `apm_port_int`            | The internal port for APM                              | `8200` |
| `fleet_host`              | The external URL for Fleet                             | `""`   |
| `fleet_host_int`          | The internal URL for Fleet                             | `""`   |
| `fleet_port_int`          | The internal port for Fleet                            | `8220` |
| `apm_secret_token`        | The secret token for APM                               | `""`   |
| `stack_monitoring_secret` | The k8s secret name for the stack monitoring           | `""`   |
