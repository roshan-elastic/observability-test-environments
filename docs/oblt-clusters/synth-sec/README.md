# Overview

`synth-sec` is a long running ESS cluster deployed in the Cloud First test region in production.
It uses the latest SNAPSHOT, therefore it is `unstable` and contains the features of the
current development branch.
The ESS deployment has Elasticsearch, Kibana, and Integrations server.

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Kind of data

Synthetic data generated for specific purposes.  Currently data is generated on an ad hoc basis into specified indices within a specified date range.  The goal is to have regular generation of data so data is consistently indexed that is close to the current date.

A data set generated for the Security Explore dashboards:  Network, Users and Hosts.

- Index: logs-synthetic-explore-default
- Date Range: 2024-05-01 to 2024-05-30
- Documents: 10000000
- Specification: https://github.com/elastic/guts/blob/main/scripts/gen_synth/document_spec/explore.yml

### Data Caching
Elasticsearch will cache results to improve performance.  This will result in very low query response times when queries are repeated, documented here https://github.com/elastic/guts/issues/44.  Clearing the cache before tests is good practice and can be done using the elasticsearch API `/_cache/clear`.

## Update

The update is automatic and happens on every Tuesday.
The version to update is the latest snapshot.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('synth-sec') }}

## Service map

{{ common.service_map('synth-sec') }}

## Deployed versions

{% include 'synth-sec/cluster-info.md' ignore missing %}

{{ common.common_links() }}
