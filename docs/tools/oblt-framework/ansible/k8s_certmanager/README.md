---
render_macros: false
---
# k8s_certmanager Role

## Overview

This role installs and configures cert-manager on a Kubernetes cluster.
It has also tasks to uninstall the installed services on a Kubernetes cluster.

## Passwords/Certs

All passwords, secrets and other sensitive information will be stored as secrets in the cluster,
we will get this data from the Vault when the CI builds the environment, so you can always get the
passwords from the Vault using your personal account.

The access details to the clusters are stored in the vault on the following secrets:

* secret/observability-team/ci/test-clusters/CLUSTER_NAME/credentials

The test clusters can use real certificates provided by [let's Encrypt][],
we are limited to 50 certificates per week so for non-permanent environments,
we will use the [let's Encrypt staging][] certificates (self-signed).

* [Rate Limits](https://letsencrypt.org/docs/rate-limits/)

## Requirements

It requires to have Helm CLI installed.

## Dependencies

Include [common][].

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    certmanager:
        enabled: true
        version: v0.10.1
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.k8s_kind
    - role: oblt.framework.k8s_helm
    - role: oblt.framework.k8s_certmanager
```

## License

Apache License 2.0

[common]:../common/README.md
[let's Encrypt staging]: https://letsencrypt.org/docs/staging-environment/
[let's Encrypt]: https://letsencrypt.org/

## Parameters

### General

| Name                 | Description                              | Value |
| -------------------- | ---------------------------------------- | ----- |
| `certificate_issuer` | Default issuer name for the certificates | `""`  |
