# Overview

`release-sec` is a long running ESS cluster deployed in the Cloud First test region in production.
We use the latest release or current BC, therefore it is `stable` and contains the features of
the current release or Build Candidate (BC).
The ESS deployment has Elasticsearch, Kibana, and Integrations server.

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Kind of data

With a total of 11 endpoints (9 managed via ESTEC) continuously streaming Security events, the data is mainly in the area of Packetbeat, Endpoint Security and Filebeat.

## Update

The update is automatic and happens every day.
The version to update is the latest stable version of the next release Elastic Stack,
this includes Build Candidates (BC) versions.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('release-sec') }}

## Service map

{{ common.service_map('release-sec') }}

## Deployed versions

{% include 'release-sec/cluster-info.md' ignore missing %}

{{ common.common_links() }}
