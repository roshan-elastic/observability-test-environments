---
render_macros: false
---
# Stack Serverless Role

## Overview

The Serverless role deploy the serverless project on Elastic Cloud.
It uses the [Project API][] to make the deployment.

## Requirements

Credentials for the serverless account must be stored in the Google Cloud Secret Manager (elastic-cloud-observability-team-ENVIRONMENT),

## Example Playbook

```yaml
---
- name: Converge
  hosts: localhost
  connection: local
  gather_facts: true
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: "oblt-test-{{ cluster_seed }}"
    save_secrets: true
    stack:
      mode: serverless
      template: observability
      target: qa
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.stack_serverless
  tasks:
    - name: Set cluster created flag
      set_fact:
        cluster_created: true
    - name: Vault Secrets
      include_role:
        name: oblt.framework.common
        tasks_from: cluster_state.yml
```

## License

Apache License 2.0

[Project API]: https://backstage.elastic.dev/catalog/default/api/project-api/definition

## Parameters

### General

| Name                     | Description                                          | Value |
| ------------------------ | ---------------------------------------------------- | ----- |
| `serverless_credentials` | The credentials for the serverless account from GCSM | `""`  |
| `serverless_endpoint`    | The endpoint for the serverless API                  | `""`  |
| `serverless_apikey`      | The API key for the serverless API                   | `""`  |
