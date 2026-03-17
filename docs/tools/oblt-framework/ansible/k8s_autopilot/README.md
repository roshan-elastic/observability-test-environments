---
render_macros: false
---
# Ansible Role: k8s_autopilot

This Role provision a Google Cloud Platform (GCP) Kubernetes Autopilot cluster.

## Requirements

It requires Google default credentials to be set in the environment.

## Example Playbook

```yaml
- name: Converge
  hosts: localhost
  connection: local
  gather_facts: true
  vars_files:
    - "/tmp/cluster-config.yml"
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.k8s_helm
    - role: oblt.framework.k8s_autopilot
```

## License

Apache License 2.0

## Parameters

### K8s Autopilot Role to create a GKE Autopilot cluster

| Name                    | Description                                             | Value                                                                                                                     |
| ----------------------- | ------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `cluster_context`       | The context name of the cluster to use.                 | `gke_{{ k8s_project }}_{{ k8s_region }}_{{ gke_cluster_name }}`                                                           |
| `k8s_autopilot_channel` | The GCE update channel to use for updating the cluster. | `stable`                                                                                                                  |
| `k8s_autopilot_version` | The version of the GKE Autopilot cluster to use.        | `1.27.11-gke.1062003`                                                                                                     |
| `k8s_autopilot_labels`  | The labels to apply to the cluster.                     | `name={{ cluster_name }},owner={{ oblt_username }},division=engineering,org=obs,team=observability,project=oblt-clusters` |
| `password_apr1`         | The password hash for the user to access the cluster.   | `{{ password_plaintext | password_hash('apr_md5_crypt') }}`                                                               |
