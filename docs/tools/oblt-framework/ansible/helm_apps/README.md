---
render_macros: false
---
# Helm Apps Role

## Overview

This Role deploy a Helm chart specified in a configuration.

It also has tasks to uninstall the services installed on a Kubernetes cluster.

## Requirements

It is necessary to have gcloud, and the Helm CLI installed.

## Dependencies

It include [common][], and [k8s_helm][] role.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
    environment:
      HOME: "{{ build_dir }}"
      PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
    vars:
      cluster_name: oblt-test
      save_secrets: false
      apps:
        helm:
          - name: dummy-helm
            values_file: "{{ playbook_dir }}/helm/dummy/helm-values.yml"
            chart: "{{ build_dir }}/helm-charts/hello"
            version: 0.1.0
          - name: dummy-helm-1
            values_file: "{{ playbook_dir }}/helm/dummy1/helm-values.yml"
            chart: "{{ build_dir }}/helm-charts/hello"
            version: 0.1.0
    roles:
      - role: oblt.framework.tools
      - role: oblt.framework.common
      - role: oblt.framework.k8s_kind
      - role: oblt.framework.helm_apps
```

## License

Apache License 2.0

[common]:../common/README.md
[k8s_helm]:../k8s_helm/README.md

## Parameters

### General

| Name        | Description                                                                | Value |
| ----------- | -------------------------------------------------------------------------- | ----- |
| `helm_apps` | List of Helm charts to install, see [K8s helm Role](../k8s_helm/README.md) | `[]`  |
