# Overview

Overview-oblt is an observability test cluster deployed everyday with the latest SNAPSHOT version.
Overview cluster is connected to the CFT region monitoring clusters using CCS (https://overview.elastic-cloud.com/),
so it has access to all the logs and metrics from the CFT region.

The goal of this clusters is provide the latest features to the SRE team and collect feedback about them.

the ESS deployment uses the latest Stack Pack available on ESS and the latest SNAPSHOT version
available in the [artifactory][artifactory] at the time of the creation.
In case of broken SNAPSHOT versions the update of the environment would have a delay.
The Elastic Stack logs and metrics are sent to [logs monitoring-oblt][],
the APM traces of Kibana are sent to [APM monitoring-oblt][].

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Update

The update of the clusters is automatic and it happens every Tuesday.
The version to update is the latest stable version of the latest SNAPSHOT created by `main` branches.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('overview-oblt') }}

## Deployed versions

{% include 'overview-oblt/cluster-info.md' ignore missing %}

{{ common.common_links() }}
