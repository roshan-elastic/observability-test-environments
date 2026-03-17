---
render_macros: false
---
# feature_siem_linux role

## Overview

This role is used to install and configure `feature_siem_linux` on a oblt-cluster.
This `feature_siem_linux` feature deploys a VM with a SIEM solution on a oblt-cluster.

## Requirements

A oblt cluster deployed.
A Elastic Agent policy created with the ID `security-linux-policy`.
There is a bootstrap folder at `/deployments/bootstrap-defend` with the REST calls to create the Elastic Agent policy.

## Dependencies

none.

## Example Playbook

```yaml
- name: Test feature_siem_linux role
  hosts: localhost
  connection: local
  gather_facts: true
  vars:
    feature_siem_linux_name: "test-instance-vm-gcp-00"
    feature_siem_linux_user: "admin"
    feature_siem_linux_fleet_url: "https://fleet.example.com"
    feature_siem_linux_kibana_url: "https://kibana.example.com"
    feature_siem_linux_kibana_username: "elastic"
    feature_siem_linux_kibana_password: "changeme"
  roles:
    - role: feature_siem_linux
```

## License

Apache License 2.0

## Parameters

### General

| Name                                           | Description                      | Value                                                                                                       |
| ---------------------------------------------- | -------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `feature_siem_linux_user`                      | User to connect to the VM        | `""`                                                                                                        |
| `feature_siem_linux_name`                      | Name of the VM                   | `""`                                                                                                        |
| `feature_siem_linux_fleet_url`                 | URL of the Fleet server          | `{{ fleet_url | mandatory }}`                                                                               |
| `feature_siem_linux_kibana_url`                | URL of the Kibana server         | `{{ kibana_url | mandatory }}`                                                                              |
| `feature_siem_linux_kibana_username`           | Username to connect to Kibana    | `{{ kibana_username | mandatory }}`                                                                         |
| `feature_siem_linux_kibana_password`           | Password to connect to Kibana    | `{{ kibana_password | mandatory }}`                                                                         |
| `feature_siem_linux_elastic_agent_version`     | Version of the Elastic Agent     | `8.13.0`                                                                                                    |
| `feature_siem_linux_elastic_agent_package`     | URL of the Elastic Agent package | `elastic-agent-{{ feature_siem_linux_elastic_agent_version }}-amd64.deb`                                    |
| `feature_siem_linux_elastic_agent_package_url` | URL of the Elastic Agent package | `https://artifacts.elastic.co/downloads/beats/elastic-agent/{{ feature_siem_linux_elastic_agent_package }}` |
