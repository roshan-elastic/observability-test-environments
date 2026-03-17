---
render_macros: false
---
# k8s_beats Role

## Overview

This Role deploys [Beats][] into the k8s cluster.

It has also tasks to uninstall the services installed on a Kubernetes cluster.

## Requirements

It requires gcloud, kubectl, and Helm CLI to be installed.

## Dependencies

It includes [common][], [k8s][] and [stack_ess][] or [stack_eck][]  role.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    beats:
      enabled: true
      auditbeat:
        enabled: true
        image: docker.elastic.co/observability-ci/auditbeat:8.0.0-SNAPSHOT
      filebeat:
        enabled: true
        image: docker.elastic.co/observability-ci/filebeat:8.0.0-SNAPSHOT
      heartbeat:
        enabled: true
        image: docker.elastic.co/observability-ci/heartbeat:8.0.0-SNAPSHOT
      metricbeat:
        enabled: true
        image: docker.elastic.co/observability-ci/metricbeat:8.0.0-SNAPSHOT
      packetbeat:
        enabled: false
        image: docker.elastic.co/observability-ci/packetbeat:8.0.0-SNAPSHOT
    apache2:
      enabled: true
      image: docker.io/httpd:2.4.38-alpine
    haproxy:
      enabled: true
      image: docker.io/haproxy:1.9.3-alpine
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.common
    - role: oblt.framework.k8s
    - role: stack_eck
    - role: oblt.framework.k8s_beats

```

## License

Apache License 2.0

[common]:../common/README.md
[k8s]:../k8s/README.md
[stack_ess]:../stack_ess/README.md
[stack_eck]:../stack_eck/README.md
[Beats]: https://www.elastic.co/products/beats

## Parameters

### Beats

| Name                       | Description               | Value   |
| -------------------------- | ------------------------- | ------- |
| `beats`                    | Beats configuration       |         |
| `beats.enabled`            | Enable Beats services     | `false` |
| `beats.auditbeat`          | Auditbeat configuration   |         |
| `beats.auditbeat.enabled`  | Enable Auditbeat service  | `true`  |
| `beats.auditbeat.image`    | Auditbeat image           | `""`    |
| `beats.filebeat`           | Filebeat configuration    |         |
| `beats.filebeat.enabled`   | Enable Filebeat service   | `true`  |
| `beats.filebeat.image`     | Filebeat image            | `""`    |
| `beats.heartbeat`          | Heartbeat configuration   |         |
| `beats.heartbeat.enabled`  | Enable Heartbeat service  | `true`  |
| `beats.heartbeat.image`    | Heartbeat image           | `""`    |
| `beats.metricbeat`         | Metricbeat configuration  |         |
| `beats.metricbeat.enabled` | Enable Metricbeat service | `true`  |
| `beats.metricbeat.image`   | Metricbeat image          | `""`    |
| `beats.packetbeat`         | Packetbeat configuration  |         |
| `beats.packetbeat.enabled` | Enable Packetbeat service | `false` |
| `beats.packetbeat.image`   | Packetbeat image          | `""`    |
| `beats_enabled`            | Enable Beats services     | `""`    |
| `apache_version`           | Apache version to deploy  | `""`    |
| `apache_enabled`           | Enable Apache service     | `""`    |
| `haproxy_version`          | HAProxy version to deploy | `""`    |
| `haproxy_enabled`          | Enable HAProxy service    | `""`    |
