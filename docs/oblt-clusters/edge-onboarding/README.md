# Overview

To ensure we can offer users end-to-end Observability for their entire stack, we will be focusing
on providing an end-to-end moon landing experience around a limited technology stack with a core set of use cases.

This stack will comprise of:

* Hosts (Linux)
* Containers (docker, containerd and any runtime detected via the Kubernetes integration)
* Services (at least a java application)
* Kubernetes (self-managed cluster)
* An ingress controller (Nginx ingress controller)
* A database (PostgreSQL)
* A web server (Tomcat)

Scenarios / data to contain

* Infra metrics (host + k8s)
* Custom (Prometheus) metrics
* Metrics from integrations (e.g. PostgreSQL, NGINX, Tomcat, etc.)
* System logs (e.g. k8s logs)
* Custom application logs (e.g. through container logs)
* Logs for “known” components / integrations (e.g. NGINX ingress controller)
* Logs-only services
* Services with traces, logs and metrics
* Services in different languages (Java, Node.js, .NET, Python, Go, etc.)

## Update

The update of the clusters is automatic and it happens once a week every Tuesday.
The version to update is the latest stable version of the latest SNAPSHOT created by `main` branches.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('edge-onboarding') }}

## Deployed versions

{% include 'edge-onboarding/cluster-info.md' ignore missing %}

{{ common.common_links() }}
