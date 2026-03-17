# Overview template cluster

The oberview-oblt cluster is a clusters connected to the overview clusters inside the CFT region in production.
This clusters contains monitoring data from all the clusters in the region.
Because this cluster has a Kibana instance deployed it is not possible to use it for development.
To allow developers to use the monitoring data from the region we have created the `ess-overview` template cluster.
This template deploy a Elasticsearch only cluster configured to connect to the overview clusters in the region,
as overview-oblt does.
The clusters created from this template can be used for development and testing running a local Kibana instance.

## Create a cluster from the template

Use the following command to create a cluster from the template for the Elastic Stack version 8.7.0:

```bash
oblt-cli cluster create custom \
    --template ess-overview \
    --parameters '{"StackVersion":"8.7.0"}'
```

It is possible to pass custom Docker images for Elasticsearch by using the `ElasticsearchDockerImage` parameter:

```bash
oblt-cli cluster create custom \
    --template ess-overview \
    --parameters '{"StackVersion":"8.7.0", "ElasticsearchDockerImage":"docker.elastic.co/observability-ci/elasticsearch-cloud-ess:8.7.0-16046737"}'
```

{% include 'tools/oblt-cli/local-kibana.md' ignore missing %}
