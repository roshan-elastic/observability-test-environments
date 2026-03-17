---
render_macros: false
---
# k8s_ingress Role

## Overview

This role installs and configures Ingress on a Kubernetes cluster.
It has also tasks to uninstall the installed services on a Kubernetes cluster.

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
    ingress:
      enabled: true
      version: 1.27.1
    password_plaintext: "supper_secret_password"
  roles:
    - role: oblt.framework.k8s_kind
    - role: oblt.framework.k8s_helm
    - role: oblt.framework.k8s_ingress
```

## License

Apache License 2.0

[common]:../common/README.md

## Parameters

### General

| Name                      | Description                                                                           | Value |
| ------------------------- | ------------------------------------------------------------------------------------- | ----- |
| `ingress_dir`             | The directory where the ingress deploy files are stored                               | `""`  |
| `password_apr1`           | The password hash for the ingress, it is generated with the password_plaintext value. | `""`  |
| `ingress_version_default` | The default version for the ingress                                                   | `""`  |
| `ingress_version`         | The version for the ingress                                                           | `""`  |
