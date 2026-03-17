---
render_macros: false
---
# k8s_helm Role

## Overview

This role installs and configures Helm on a Kubernetes cluster.
It has also import_tasks to install the Helm CLI,
and to uninstall the installed services on a Kubernetes cluster.

## Requirements

It requires to have Helm CLI installed.

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
    cluster_name: oblt-test
    helm:
      enabled: true
      version: 3.0.2
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.tools
    - role: oblt.framework.k8s_kind
    - role: oblt.framework.k8s_helm
```

## License

Apache License 2.0

[common]:../common/README.md

## Parameters
