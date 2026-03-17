# Overview

Dev oblt cluster is a long term ESS cluster deployed in the Cloud First test region in production.
We use the next `minor` release, therefore it is `stable` and contains new features for the new release.
The ESS deployment has Elasticsearch, Kibana, and Integrations server.
To generate data dev oblt cluster add a [Kubernetes][] [GKE][] cluster to deploy some applications.
The GKE cluster contains the following applications:

* Opbeans application
* MySQL
* PostgreSQL
* Redis
* Nginx
* HAProxy
* Apache httpd
* Auditbeat
* Packetbeat

The use Elastic Agent standalone to Collect log and metrics from the GKE cluster.
There are also some SIEM instances reporting security data connected to the cluster.

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Update

The update of the clusters is automatic and it happens every day.
The version to update is the latest stable version of the next release Elastic Stack SNAPSHOTs.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('dev-oblt') }}

## Service map

{{ common.service_map('dev-oblt') }}

## Deployed versions

{% include 'dev-oblt/cluster-info.md' ignore missing %}

{{ common.common_links() }}
