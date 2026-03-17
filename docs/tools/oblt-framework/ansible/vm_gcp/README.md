---
render_macros: false
---
# vm_gcp Role

## Overview

This Role provision a VM in GCP.

## Requirements

It requires Google default credentials to be set in the environment.

## Example Playbook

```yaml
- name: Converge
  hosts: localhost
  connection: local
  gather_facts: true
  vars:
    build_dir: "/tmp/build"
    cluster_name: "vm-gcp-1234"
    vm_gcp_ssh_private_key: "~/.ssh/google_compute_engine"
    vm_gcp_ssh_public_key: "~/.ssh/google_compute_engine.pub"
    vm_gcp_user: "admin"
    vm_gcp_name: "test-instance-1234"
  roles:
    - role: oblt.framework.common
    - role: vm_gcp
```

## License

Apache License 2.0

## Parameters

### vm_gcp Role to create a GCP VM

| Name                            | Description                                                             | Value                                                              |
| ------------------------------- | ----------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `vm_gcp_machine_type`           | The machine type to use for the instance.                               | `e2-micro`                                                         |
| `vm_gcp_name`                   | The name of the instance.                                               | `test-instance`                                                    |
| `vm_gcp_image`                  | The image to use for the instance.                                      | `debian-cloud/debian-11`                                           |
| `vm_gcp_vpc_name`               | The name of the instance.                                               | `{{ cluster_name }}-network`                                       |
| `vm_gcp_terradorm_dir`          | The directory where the terraform files are located.                    | `{{ build_dir }}/terraform/gcp/{{ vm_gcp_name }}`                  |
| `vm_gcp_bucket`                 | The name of the bucket to use for the Terraform state.                  | `oblt-clusters`                                                    |
| `vm_gcp_bucket_prefix`          | The prefix to use for the Terraform state in the bucket.                | `terraform/gcp/{{ vm_gcp_name }}`                                  |
| `vm_gcp_user`                   | The user to access the instance.                                        | `admin`                                                            |
| `vm_gcp_ssh_private_key`        | The path to the private key to use for SSH access.                      | `{{ private_key_file | default('~/.ssh/google_compute_engine') }}` |
| `vm_gcp_ssh_public_key`         | The path to the public key to use for SSH access.                       | `{{ vm_gcp_ssh_private_key }}.pub`                                 |
| `vm_gcp_scripts`                | The list of scripts to run on the instance.                             | `["{{ vm_gcp_terradorm_dir }}/script.sh"]`                         |
| `vm_gcp_open_ports`             | The list of ports to open on the instance.                              | `["22"]`                                                           |
| `vm_gcp_project`                | The GCP project to use.                                                 | `{{ gcp_project | mandatory }}`                                    |
| `vm_gcp_zone`                   | The GCP zone to use.                                                    | `{{ gcp_zone | mandatory }}`                                       |
| `vm_gcp_cluster_name`           | The name of the GCP cluster.                                            | `{{ cluster_name | mandatory }}`                                   |
| `vm_gcp_oblt_username`          | The username owner of the project.                                      | `{{ oblt_username | mandatory }}`                                  |
| `vm_gcp_windows_startup_script` | The path to the Windows startup script.                                 | `{{ vm_gcp_terradorm_dir }}/windows-startup-script.cmd`            |
| `vm_gcp_os`                     | The operating system of the instance.                                   | `linux`                                                            |
| `gcp_vm_upload_files`           | The list of files to upload to the instance.                            | `[]`                                                               |
| `vm_gcp_startup_script`         | Script to run on startup (e.g. "echo 'Hello, World!' > /tmp/hello.txt") | `""`                                                               |
