# Overview

Edge logs cluster is a long term ESS cluster deployed in the Cloud First test region in production.
This environment is logs-only oriented and is used to test the service detection from logs.
We use the next `major` version, therefore it's `unstable` by definition and contains all the new features for the next major.
The ESS deployment has Elasticsearch, Kibana, and Integrations server.
To generate data edge logs oblt cluster add a [Kubernetes][] [GKE][] cluster to deploy some applications.
The GKE cluster contains the following applications:

* Opbeans application
* MySQL
* PostgreSQL
* Redis

The use Elastic Agent standalone to Collect log and metrics from the GKE cluster.
There are also some SIEM instances reporting security data connected to the cluster.

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site][docs-site]

Requested at [#2182](https://github.com/elastic/observability-robots/issues/2182)

## Update

The update of the clusters is automatic and it happens once a week every Tuesday.
The version to update is the latest stable version of the latest SNAPSHOT created by `main` branches.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('edge-logs') }}

## Deployed versions

{% include 'edge-logs/cluster-info.md' ignore missing %}

{{ common.common_links() }}
