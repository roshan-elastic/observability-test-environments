---
render_macros: false
---
# tools Role

## Overview

This role installs tools that are required by other roles.

## Requirements

None

## Dependencies

None

## Example Playbook

```yaml
- name: Converge
  hosts: localhost
  connection: local
  gather_facts: true
  vars:
    build_dir: "{{ lookup('env', 'MOLECULE_SCENARIO_DIRECTORY') }}/build"
    hermit_enabled: true
    k8s_provider: kind
    k8s_enabled: true
  roles:
    - role: oblt.framework.tools
```

## License

Apache License 2.0

## Parameters

### General

| Name                         | Description                                        | Value  |
| ---------------------------- | -------------------------------------------------- | ------ |
| `build_dir`                  | The directory where the build is stored            | `""`   |
| `bin_dir`                    | The directory where the binaries are stored        | `""`   |
| `hermit_enabled`             | Enable Hermit to install the tools                 | `true` |
| `k8s_provider`               | The k8s provider name                              | `""`   |
| `k8s_enabled`                | Enable K8s                                         | `""`   |
| `save_secrets`               | Save the secrets in the build directory            | `true` |
| `terraform_version`          | The Terraform version to install                   | `""`   |
| `ess_terraform_version`      | Elastic terraform provider version                 | `""`   |
| `build_tf_provider_from_src` | Build the Terraform provider from source           | `""`   |
| `use_snapshot_tf_provider`   | Use the snapshot version of the Terraform provider | `""`   |
| `gcloud_version`             | The GCloud version to install                      | `""`   |
| `kind_version`               | The Kind version to install                        | `""`   |
| `kubectl_version`            | The Kubectl version to install                     | `""`   |
| `helm_version`               | The Helm version to install                        | `""`   |
