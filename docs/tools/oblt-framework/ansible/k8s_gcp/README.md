---
render_macros: false
---
# k8s_gcp Role

## Overview

This role creates a K8s cluster on Google Kubernetes Engine (GKE)
and grab the authentication details to start using `gcloud` and `kubectl`.

## Requirements

It requires a service account credential stored in the Vault,
see default variables for more details.

## Dependencies

Includes [common][].

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-cluster
    k8s:
      enabled: true
      project: "elastic-observability"
      region: "europe-west1-c"
      max_node_count: "3"
      machine_type: "n1-standard-4"
      default_namespace: "default"
      domain: "ip.es.io"
  roles:
    - role: oblt.framework.tools
    - role: gcp
    - role: oblt.framework.k8s_gcp
```

## License

Apache License 2.0

[common]:../common/README.md

## Parameters

### General

| Name               | Description                                         | Value |
| ------------------ | --------------------------------------------------- | ----- |
| `gcp_dir`          | The directory where the GCP files are stored        | `""`  |
| `cluster_context`  | The name of the k8s configuration context to create | `""`  |
| `static_ip_region` | The region to create the static IP                  | `""`  |
