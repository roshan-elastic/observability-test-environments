# Overview

`edge-sec` is a long running ESS cluster deployed in the Cloud First test region in production.
It uses the latest SNAPSHOT, therefore it is `unstable` and contains the features of the
current development branch.
The ESS deployment has Elasticsearch, Kibana, and Integrations server.

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Kind of data

With a total of 11 endpoints (9 managed via ESTEC) continuously streaming Security events, the data is mainly in the area of Packetbeat, Endpoint Security and Filebeat.

## Update

The update is automatic and happens on every Tuesday.
The version to update is the latest snapshot.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('edge-sec') }}

## Service map

{{ common.service_map('edge-sec') }}

## Deployed versions

{% include 'edge-sec/cluster-info.md' ignore missing %}

{{ common.common_links() }}
