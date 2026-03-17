# Edge Benchmarks

## Overview

This ESS cluster is similar to [edge-lite-oblt][], but without apps deployed.
edge-benchmarks is used to run the benchmarks for the edge versions.

## Update

The update of the clusters is automatic and it happens once a week every day.
The version to update is the latest stable version of the latest SNAPSHOT created by `main` branches.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

{{ common.common_links() }}

[edge-lite-oblt]: ../edge-lite-oblt/README.md
