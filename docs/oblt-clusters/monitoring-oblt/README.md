# Overview

Monitoring oblt cluster is a long term ESS cluster deployed in the Cloud First test region in production.
We use the next `minor` release, therefore it is `stable` and contains new features for the new release.
The ESS deployment has Elasticsearch, Kibana, and Integrations server.
The rest of oblt clusters send logs, metrics, and APM data to monitoring oblt,
it uses Elastic Stack monitoring and APM for that,
thus this cluster is the place to check the behaviour if the rest of oblt clusters.
there are some Synthetic monitors configured to make a basic functional test on some of the oblt clusters.
The use Elastic Agent standalone to Collect log and metrics from the GKE cluster.
There are also some SIEM instances reporting security data connected to the cluster.

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site](https://elastic.github.io/observability-test-environments/)

## Update

The update of the clusters is automatic and it happens every Tuesday.
The version to update is the latest stable version of the next release Elastic Stack SNAPSHOTs.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('monitoring-oblt') }}

## Deployed versions

{% include 'monitoring-oblt/cluster-info.md' ignore missing %}

{{ common.common_links() }}
