# Overview

Edge-lite is an observability test cluster deployed everyday with the latest SNAPSHOT version.
The edge-lite deployment has two pieces, one is the Elastic Stack and a Kubernetes(k8s) cluster.

The Elastic Stack is deployed on ESS running on the Cloud First Test region(CFT),
the ESS deployment uses the latest Stack Pack available on ESS and the latest SNAPSHOT version
available at the time of the creation,
the deployment is created every day so the SNAPSHOT version could have 24 hours.
In case of broken SNAPSHOT versions the update of the environment would have a delay.
The Elastic Stack logs and metrics are sent to [monitoring-oblt][monitoring-oblt],
the APM traces of Kibana are sent to [monitoring-oblt][monitoring-oblt].

The k8s cluster is deployed on Google Kubernetes Engine(GKE),
this k8s cluster is monitored with The Elastic Agent running on standalone mode.
The k8s metrics and logs are sent to the Elasticsearch running edge-lite.
In the k8s cluster we deploy some applications and services in order to generate data (logs, metrics, and APM).
The environment will composed by the following services:

* Apache httpd
* Auditbeat
* Elastic Agent standalone
* Filebeat
* Grafana
* HAProxy
* Heartbeat
* Metricbeat
* MySQL
* Nginx
* Opbeans application
* Packetbeat
* PostgreSQL
* Prometheus
* Redis

 Edge-lite is destroy and created everyday, the data of the cluster would not be older than 24 hours.
 For other use cases that needs more historic data you should use other cluster.

The Elastic Stack is configured with Stack Monitoring and APM reporting to [monitoring-oblt][monitoring-oblt]

For more details about how to use it and implementation check the [Documentation Site][docs-site]

## Update

The update of the clusters is automatic and it happens once a week every day.
The version to update is the latest stable version of the latest SNAPSHOT created by `main` branches.
Developers clusters are also automatically updated on every update of the cluster.

{% import 'common-macros.md' as common %}

## URLs and Credentials

{{ common.cluster_urls('edge-lite-oblt') }}

## Service map

{{ common.service_map('edge-oblt') }}

## Deployed versions

{% include 'edge-lite-oblt/cluster-info.md' ignore missing %}

{{ common.common_links() }}
