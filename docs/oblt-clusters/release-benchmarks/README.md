# Release Benchmarks

## Overview

This ESS cluster is similar to [release-oblt][], but without apps deployed.
edge-benchmarks is used to run the benchmarks for the edge versions.

## Update

The update of the clusters is automatic and it happens every day.
The version to update is the latest stable version of the next release Elastic Stack,
this include Build Candidates versions(BC).
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

{{ common.common_links() }}

[release-oblt]: ../release-oblt/README.md
