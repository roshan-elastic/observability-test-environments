# Overview

Observability CI is a long term ESS cluster deployed in the Cloud First test region in production.
We use the latest release or current BC,
therefore it is `stable` and contains the features of the current release or Build Candidate (BC).
The ESS deployment has Elasticsearch, Kibana, and Integrations server.

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Update

The update of the clusters is automatic and it happens once a week every Tuesday.
The version to update is the latest stable version of the next release Elastic Stack,
this include Build Candidates versions(BC).
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('observability-ci') }}

## Deployed versions

{% include 'observability-ci/cluster-info.md' ignore missing %}

{{ common.common_links() }}
