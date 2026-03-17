---
render_macros: false
---
# k8s Role

## Overview

This Role deploy and configure a kubernetes cluster,
this cluster would have Helm installed, an ingress controller, and a certificates manager.
It has also tasks to uninstall the services installed on a Kubernetes cluster.

## Requirements

It requires to have gcloud, kubectl, and Helm CLI installed.

## Dependencies

It include k8s_kind, [k8s_gcp][] or [k8s_kind][], [k8s_helm][], [k8s_ingress][], and [k8s_certmanager][] roles.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    k8s:
      enabled: true
      provider: kind
      region: "none"
      domain: "ip.es.io"
      default_namespace: "default"
    password_plaintext: "super_password"
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.k8s
```

## License

Apache License 2.0

[k8s_gcp]:../k8s_gcp/README.md
[k8s_kind]:../k8s_kind/README.md
[k8s_helm]:../k8s_helm/README.md
[k8s_ingress]:../k8s_ingress/README.md
[k8s_certmanager]:../k8s_certmanager/README.md

## Parameters
