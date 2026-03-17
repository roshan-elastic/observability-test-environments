# feature_integrations role

## Overview

This role is used to install and configure `feature_integrations` on a oblt-cluster.
The integrations feature deploy any package that have `_dev/deploy/docker` folder in package folder
in the [integrations repository][].
It uses Docker In Docker to run the integration in a container.
It configures a Fleet policy for the integration and launch and Elastic Agent to run the integration.

The process to deploy the integrations is divided in several steps:

* [Bootstrap the integration][]
  * Configure the Elastic Agent policy
    * Get the version of the integration
    * Get the version of the docker integration
  * Create an Agent policy for the integration
  * Configure the Docker integration in the Agent policy
  * Configure the integration in the Agent policy
* Deploy the [integrations Helm Chart][]

[integrations repository]: https://github.com/elastic/integrations/tree/main
[Bootstrap the integration]: https://github.com/elastic/observability-test-environments/tree/main/ansible/ansible_collections/oblt/framework/roles/common/files/deployments/bootstrap-integration
[integrations Helm Chart]: https://github.com/elastic/observability-test-environments/tree/main/ansible/ansible_collections/oblt/framework/roles/common/files/helm-charts/integrations

## Requirements

It ness a Elastic Stack and a k8s cluster deployed.

## Dependencies

None.

## Example Playbook

```yaml
- name: Converge
  hosts: localhost
  connection: local
  gather_facts: true
  vars:
    cluster_name: "oblt-test"
    stack:
      mode: "ess"
      template: "observability"
    k8s:
      provider: kind
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.k8s_kind
    - role: oblt.framework.stack_ess
  tasks:
    - name: Set cluster created flag
      set_fact:
        cluster_created: true

    - name: Vault Secrets
      include_role:
        name: oblt.framework.common
        tasks_from: cluster_state.yml

    - name: Docker secrets
      include_role:
        name: oblt.framework.k8s
        tasks_from: docker_secrets.yml

    - name: Deploy Apache integration
      include_role:
        name: oblt.framework.feature_integrations
        tasks_from: deploy.yml
      vars:
        feature_integrations_name: "apache"
```

## License

Apache License 2.0

## Parameters

### General

| Name                                      | Description                     | Value |
| ----------------------------------------- | ------------------------------- | ----- |
| `feature_integrations_name`               | Name of the integration         | `""`  |
| `feature_integrations_helm_chart_version` | Integrations Helm chart version | `""`  |
