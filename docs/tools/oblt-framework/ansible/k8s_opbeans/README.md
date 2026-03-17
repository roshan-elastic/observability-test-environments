---
render_macros: false
---
# k8s_opbeans Role

## Overview

This Role deploy [Opbeans][] sample applications into the k8s cluster,
it also install [MySQL][], [PostgreSQL][], and [Redis][] as auxiliary services.
It has also tasks to uninstall the services installed on a Kubernetes cluster.

Opbeans is the sample application used to test APM, it is composed by a frontend made in
Node.js and different backends each one for each APM Agent (.NET, Go, Java, Node.js, PHP, Python, Ruby, and RUM).
We will deploy several services, frontend, backend, APM server, database, and cache. The frontend will be a NGINX as proxy
of the Node.js frontend application in one pod with two replicas. The backend will be a service
pointing to pods with different versions on the APM Agent (.NET, Go, Java, Node.js, PHP, Python, and Ruby).
Finally, we will need the database service made with a single PostgreSQL pod, the cache service
 made by a Redis pod with 3 replicas, and an APM Server pod also with 3 replicas.

* [Opbeans frontend][]
* [Opbeans DotNet][]
* [Opbeans Go][]
* [Opbeans Java][]
* [Opbeans Node.js][]
* [Opbeans PHP][]
* [Opbeans Python][]
* [Opbeans Ruby][]

![beats-environment-04](/images/beats-environment-04.png){: style="width:600px"}

![beats-environment-05](/images/beats-environment-05.png){: style="width:600px"}

This section configures the Opbeans settings and its auxiliary services.
It has a job executed every couple of minutes to generate APM traces.
There is an endpoint exposed in the cluster per Opbean to allow
connections and generate APM traces.

### Data generators

We will use k8s scheduled task to launch some process to generate data by acceding to the
demo applications.

[Opbeans loadgen](https://github.com/elastic/opbeans-loadgen)

### Opbeans Scheduled tasks

* Opbeans load generator [k8s/opbeans/opbeans-loadgen-job.yaml](https://github.com/elastic/observability-test-environments/blob/main/ansible/ansible_collections/oblt/framework/roles/k8s-opbeans/templates/k8s/opbeans/opbeans-loadgen-job.yaml.j2)

## MySQL, PostgreSQL, and Redis

There are other indirect related auxiliary services like
MySQL, PostgreSQL, and Redis that are deployed by the Opbeans.

## Requirements

It requires to have gcloud, kubectl, and Helm CLI installed.

## License

Apache License 2.0

[common]:../common/README.md
[k8s]:../k8s/README.md
[stack_ess]:../stack_ess/README.md
[stack_eck]:../stack_eck/README.md
[Opbeans]: https://hub.docker.com/u/opbeans/
[MySQL]: https://github.com/helm/charts/tree/master/stable/mysql
[PostgreSQL]: https://github.com/helm/charts/tree/master/stable/postgresql
[Redis]: https://github.com/helm/charts/tree/master/stable/redis

## Parameters

### General

It cannot be deploy on the default namespace https://github.com/elastic/apm-mutating-webhook/issues/44
NOTE: the orther matter to make the updates with updatecli do no change the order of the charts.
see [Helm Apps Role](../helm_apps/README.md)

| Name                       | Description                                               | Value  |
| -------------------------- | --------------------------------------------------------- | ------ |
| `opbeans_dir`              | The directory where the opbeans Helm Charts are stored    | `""`   |
| `apm_attacher_dir`         | The directory where the apm-attacher Helm Chart is stored | `""`   |
| `opbeans_transaction_rate` | The rate of transactions for the opbeans                  | `""`   |
| `opbeans_go_enabled`       | Enable the opbeans-go chart                               | `""`   |
| `opbeans_dotnet_enabled`   | Enable the opbeans-dotnet chart                           | `""`   |
| `opbeans_frontend_enabled` | Enable the opbeans-frontend chart                         | `""`   |
| `opbeans_java_enabled`     | Enable the opbeans-java chart                             | `""`   |
| `opbeans_node_enabled`     | Enable the opbeans-node chart                             | `""`   |
| `opbeans_php_enabled`      | Enable the opbeans-php chart                              | `""`   |
| `opbeans_python_enabled`   | Enable the opbeans-python chart                           | `""`   |
| `opbeans_ruby_enabled`     | Enable the opbeans-ruby chart                             | `""`   |
| `opbeans_apm_enabled`      | Enable sent traces to the APM service                     | `true` |
| `opbeans_helm_charts`      | List of opbeans Helm charts to install                    | `{}`   |

## Dependencies

It includes [common][], [k8s][] and [stack_ess][] or [stack_eck][] role.

## Example Playbook

```yaml
- hosts: localhost
  connection: local
  environment:
    HOME: "{{ build_dir }}"
    PATH: "{{ build_dir }}/bin:{{ lookup('env','PATH') }}"
  vars:
    cluster_name: oblt-test
    opbeans:
      enabled: true
  roles:
    - role: oblt.framework.tools
    - role: oblt.framework.k8s
    - role: stack_eck
    - role: oblt.framework.k8s_opbeans
```
