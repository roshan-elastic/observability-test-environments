---
render_macros: false
---
# k8s_kind Role

## Overview

This role creates a K8s cluster in Docker by using [Kind][].

## Requirements

[Docker][], [Kind][], and [kubectl][] installed.

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
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.k8s_kind
```

## License

Apache License 2.0

[common]:../common/README.md
[Kind]: https://kind.sigs.k8s.io/
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[Docker]: https://docker.io

## Parameters

### General

https://hub.docker.com/r/kindest/node/tags?page=1&name=v1.27.
https://github.com/kubernetes-sigs/kind/issues/2376
gh release view --repo kubernetes-sigs/kind v0.13.0 |sed -nE 's/ - ([0-9\.]+): `(kindest\/node:.+)`/{"\1": "\2"}/p' | jq -s add

| Name                   | Description                             | Value                                         |
| ---------------------- | --------------------------------------- | --------------------------------------------- |
| `kindest_node_version` | Version of the Kind nodes to use        | `""`                                          |
| `cluster_context`      | The context name of the cluster to use. | `kind-{{ cluster_name }}`                     |
| `kind_config_file`     | The path to the Kind configuration file | `{{ build_dir }}/kind/config_single_node.yml` |
