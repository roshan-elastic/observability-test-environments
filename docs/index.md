# Observability test Clusters (oblt clusters)

Observability test clusters, aka `oblt clusters`, is a project to facilitate the Elastic Stack deployments
along with a set of applications to generate observability data (logs, metrics, and traces).
The aim of oblt cluster are the developers and CI.
Oblt clusters can deploy the Elastic Stack on Serverless, Elastic Cloud ([ESS][]) or using [ECK][] in a [Kubernetes][] cluster.
The applications to generate data are deployed into a [Kubernetes][] cluster.
The observability data generated comes from the APM Agents, Beats, and Elastic Agent,
every application deployed is monitored and instrumented to report observability data.
The Elastic Stack is also configured to report stack monitoring data and APM data.

* [Serverless](./user-guide/serverless-quick-start.md)
* [Create a ESS cluster](./user-guide/use-case-create-cluster.md)
* [Local development](./user-guide/use-case-local-kibana.md)
* [Deploy Kibana PRs](./user-guide/deploy-my-kibana-pr.md)
* [Elastic Cloud](./user-guide/ess.md)
* [Troubleshooting](./user-guide/troubleshooting.md)
* [Oblt tools](./tools/)
* [Check clusters status](https://status.obs.elastic.dev)

## How to use it

The [oblt-cli][] is the CLI tool used by the developers to manage their personal
oblt clusters to test their new features, troubleshoot bugs, demos, and so on.
You simply download the [oblt-cli][] execute a command like the following,
and you will received your credentials to use the cluster by Slack.
For more detail check [oblt-cli][].

If you preferred to use Slack to manage your oblt clusters,
it is also possible to operate them using [oblt-robot][].
In essence, you talk with a a Slack bot that will use [oblt-cli][] for you.

## How it works

The project uses [Ansible][] to orchestrate the whole process and create [Kubernetes][] resources,
Elastic Terraform provider to deploy the Elastic Stack on ESS,
and [ECK][] Operator to deploy the Elastic Stack in [GKE][] or [Kind][].
All the rest of applications are deployed using [Kubernetes][] manifests or [helm chart][].

A oblt cluster needs a cluster configuration file, this file is a YAML file
that contains a set of settings that define what we want to deploy.
The cluster configuration file defines the kind of Elastic Stack we want [ESS][] or [ECK][],
the version, memory, zones, and other settings for the Elastic Stack deployment,
if we want a [Kubernetes][] and which type [GKE][] or [Kind][],
and finally the application we want to deploy in [Kubernetes][].

All the functionality is wrapped in [Ansible][] roles and playbooks,
a `Makefile` is responsible to prepare the environment to run the [Ansible][] playbooks.
Each [Ansible][] playbook will use the cluster configuration file to perform the deployment you have configured in it.
You have some configuration examples at [test configurations](https://github.com/elastic/observability-test-environments/tree/main/tests/environments).

To know more about the configuration settings check [Oblt clusters Ansible Roles and playbooks](/Ansible/) documentation.

## What applications we deploy

Our environment will be composed of several services interconnected each other,
we want to keep the deploy and update each service as easy as possible,
we do not want to maintain a complex infrastructure and we want to reuse services/deploys
provide by our projects, if it is possible we want to use external projects deploys for
those components we do not maintain (Apache, PostgreSQL, MySQL, HAProxy, Redis, ...).
Because all of these reasons we will choose [Kubernetes][] as platform and Docker as containers.
Each service will have its own deploy. To make it easy,
to change versions and configuration. When it was possible
services are deployed with a [helm chart][].
Each Kubernetes node will report Beats,
and each service that can report Beats will do so. Finally, each service
that can report APM data would do it.
The environment will composed by the following services:

* Apache httpd
* APM Server
* Apm Synthtraces
* Auditbeat
* Elastic Agent
* Elasticsearch
* Filebeat
* Fleet Server
* HAProxy
* Heartbeat
* Kibana
* Metricbeat
* MySQL
* Nginx
* Opbeans application
* Opentelemetry Demo
* Packetbeat
* PostgreSQL
* Redis

* [Projects Big Picture](https://docs.google.com/drawings/d/1MJ-haufU06JnUHYxYDwP-EH-xBpyMDrpDn2UEdS2WgA)

!!! Note

    the `edge-oblt` and `dev-oblt` clusters use the `staging` [`package-storage`](https://github.com/elastic/package-storage) storage cluster branch, so Integrations promoted to 'staging' will show, but those that are only on `snapshot` storage will not; `release-oblt` will use `production` released integrations unless otherwise noted to mirror what shipped stack versions access.

## Elastic Agent Deploys (SIEM)

Elastic Agents are deployed to the 3 pre-defined environments. The hosts reside in GCP in 'elastic-siem' bucket.
The Ansible & Terraform that controls their deployment is managed via a Security Solutions side Jenkins infrastructure.
The very first iteration of this is entirely 'bolt-on' to the existing clusters for sake of speed, and is deployed manually.
[This](https://protections-ci.elastic.co/job/security-sre-rebuild-o11y-test-environment-deploys/) Jenkins job is manually kicked off and is dependent on the given environment's Fleet Server being stable and available.
The setup of the Security Rules, and Agent policies is being worked as a separate phase of integration, as is the daisy chaining of the Jenkins jobs (and other improvements)
The [Elastic Agent Deploy](https://github.com/elastic/siem-team/blob/master/cm/scripts/o11y-cluster-deploys.sh) code lives in the elastic/siem-team repo as noted.
The Security Engineering Productivity group is the POC for this section of the o11y environment services.

## Elasticsearch Data Resilience :bomb:

Observability Test Environments are development environments.

**:warning: We do not make any kind of backup of the data in these environments!**

The data sets for these are environments are huge and most-everything is generated by load generators.
We upgrade frequently with develop versions of the stack that sometimes break the data
or introduce breaking changes on the schema.

We make our best effort to keep the data between upgrades but sometimes it is too much effort or it is not possible.

In these cases we simply delete all indices :radioactive:, though we do try to preserve `.kibana*` indices if it is possible
to do so.

If you need some special requirement about to keep data, backup or something else,
please [open an issue](https://github.com/elastic/observability-robots/issues/new), and we will try to our best to accommodate you.

## Other references

* [K8s Demo](https://github.com/elastic/demos/tree/master/k8s_demo)
* [Beats Dev](https://github.com/elastic/beats-dev)
* [Robot Shop Demo](https://github.com/michaelhyatt/robot-shop-apm)
* [Stack Docker](https://github.com/elastic/stack-docker)
* [Ecommerce Demo](https://github.com/elastic/demos/tree/master/cyclops/misc/docker_cyclops/src)
* [Eden GCP terraform config](https://github.com/elastic/infra/tree/master/terraform/providers/gcp/env/eden)
* [k8s infra config](https://github.com/elastic/infra/tree/master/k8s)
* [GKE - Creating a cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-cluster)
* [Get password from the Vault](https://github.com/elastic/infra/blob/master/docs/vault/README.md)
* [ingress-nginx basic Authentication](https://github.com/kubernetes/contrib/tree/master/ingress/controllers/nginx/examples/auth)
* [Eden dev Beats deployment](https://github.com/elastic/infra/tree/master/k8s/eden/dev/namespaces/kube-system)

[Ansible]: https://github.com/ansible/ansible
[ECK]: https://www.elastic.co/guide/en/cloud-on-k8s/current/k8s-quickstart.html
[ESS]: https://www.elastic.co/cloud/
[GKE]: https://cloud.google.com/kubernetes-engine
[helm chart]: https://helm.sh/
[kind]: https://kind.sigs.k8s.io/
[Kubernetes]: https://kubernetes.io
[oblt-cli]: tools/oblt-cli
[oblt-robot]: tools/oblt-robot
