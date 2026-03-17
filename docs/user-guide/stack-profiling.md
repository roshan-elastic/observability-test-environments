# Elastic Stack profiling

## Overview

This guide will walk you through the process of profiling the Elastic Stack using oblt clusters and Observability.
Observability is included in the Elastic Stack by default. It is a collection of tools that can be used to monitor and profile the Elastic Stack.
Stack monitoring is a feature of Observability that provides logs and metrics for the Elastic Stack itself.
APM is a feature of Observability that provides Application performance monitoring data for applications that use the Elastic Stack.
Finally, Profiling is a feature of Observability that provides a way to profile the Elastic Stack execution.
Using these tools, you can profile the Elastic Stack and identify bottlenecks and performance issues.

The process to profile the Elastic Stack using the oblt cluster is as follows:

* Create an oblt cluster using the stack-profiling template.

```bash
oblt-cli cluster create custom \
    --template stack-profiling \
    --parameters '{"StackVersion":"8.11.3"}'
```

* Run the query you want to profile.

![Kibana de query](/images/stack-profiling-dev-query.png){: style="width:600px"}

* Check the [Query Analyzer Dashboard][] to inspect the results.

![Query Analyzer Dashboard](/images/stack-profiling-dashboard.png){: style="width:600px"}

The result of this process is a Kubernetes cluster with the Elastic Stack deployed with ECK and Observability enabled.
The observability data is stored in the [monitoring-oblt][] cluster.
The [monitoring-oblt][] has the [Query Analyzer Dashboard][] that summarizes all the observability data we have collected.
From the [Query Analyzer Dashboard][] you can navigate to logs, metrics inspector, APM, Stack monitoring and Profiling.

## Create a Stack Profiling cluster

To create a Stack Profiling cluster, you need to use the `stack-profiling` template.
The `stack-profiling` template is a custom template that deploys the Elastic Stack with ECK and Observability enabled.

To create a Stack Profiling cluster, run the following command:

```bash
oblt-cli cluster create custom \
    --template stack-profiling \
    --parameters '{"StackVersion":"8.11.3"}'
```

This command will deploy a cluster with Elasticsearch, Kibana, and APM Server and Observability enabled.
By default the Opbeans deployment is also enabled, this will generate Observability data in the cluster to perform queries with.

You can deploy a only Elasticsearch cluster by running the following command:

```bash
oblt-cli cluster create custom \
  --template stack-profiling \
  --parameters '{
      "StackVersion": "8.11.3",
      "KibanaEnabled": "false",
      "ElasticAgentEnabled": "false",
      "OpbeansEnabled": "false"
      }'
```

We will received the credentials to access the cluster in a Slack message. For more details check [retrieve credentials][]

## Ingest Data

By default the Opbeans deployment is enabled, this will generate Observability data in the cluster to perform queries with, but sometimes this is no convenient. Those Opbeans deployment can introduce noise in the data and make it difficult to profile the Elastic Stack. In those cases, you can disable the Opbeans deployment by setting the `"OpbeansEnabled"` parameter to `"false"`.

There is many ways to ingest data in the cluster, another option is to use `esrally` at the time you profile the query using a `corpora` in your track, check [Esrally section](#esrally).

Finally, you can restore data from a snapshot, check [restore a snapshot][] for more details.

## Perform a query

To perform a query, the easiest way is to use the Kibana dev tools.
Open `KIBANA_URL/app/dev_tools#/console` in your browser and run the query you want to profile.

This is an example of a query that we can use to profile the Elastic Stack:

```json
GET metrics-system.cpu-default/_search
{
  "query": {
    "match": {
      "data_stream.dataset": "system.cpu"
    }
  }
}
```

The logs, metrics, APM, and profiler data will be stored in the [monitoring-oblt][] cluster automatically.

## Inspect the results

To inspect the results, we can use the [Query Analyzer Dashboard][].
This dashboard is available in the [monitoring-oblt][] cluster.
The dashboard summarizes all the observability data we have collected.

Before you see the proper results you must filter the data by the cluster name, cluster UUID, and transaction of your search.
In our example, we are using the `eck-profiling-ihata` cluster name. The cluster UUID is `HyjNARB7RTmjxDPlVNztag`. finally, the transaction name is `GET /{index}/_search` and the path of that transaction `/metrics-system.cpu-default/_search`.
To filter the data in the dashboard we will need to use the following filters:

An APM filter by the deploymentName label and transaction details:

```console
(labels.deploymentName:eck-profiling-ihata and data_stream.type:traces AND transaction.name:"GET /{index}/_search" AND url.path:"/metrics-system.cpu-default/_search")
```

An APM metrics filter by a label and the type of data stream:

```console
(labels.deploymentName:eck-profiling-ihata and data_stream.type:metrics)
```

Finally, Stack monitoring and Profiler data is filtered by using the cluster UUID:

```console
cluster_uuid:HyjNARB7RTmjxDPlVNztag
```

Joining the three filters we can see all the data in the dashboard:

```console
(labels.deploymentName:eck-profiling-ihata AND data_stream.type:traces AND transaction.name:"GET /{index}/_search" AND url.path:"/metrics-system.cpu-default/_search")
OR (labels.deploymentName:eck-profiling-ihata AND data_stream.type:metrics)
OR cluster_uuid:HyjNARB7RTmjxDPlVNztag
```

!!!Note

    The cluster UUID can be found in the `ELASTICSEARCH_URL/` path. For example, `https://elasticsearch-ihata-0001.es.us-west2.gcp.elastic-cloud.com/` has the cluster UUID `HyjNARB7RTmjxDPlVNztag`.

    ```bash
      CLUSTER_NAME=eck-profiling-ihata
      ENV_FILE="${PWD}/.env"

      # Get the secrets from the cluster
      oblt-cli cluster secrets env --cluster-name ${CLUSTER_NAME} --output-file "${ENV_FILE}"
      source "${ENV_FILE}"

      curl -sL -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" -X GET "${ELASTICSEARCH_HOST}/" | jq -r '.cluster_uuid'
    ```

From the [Query Analyzer Dashboard][] you can navigate to logs, metrics inspector, APM, Stack monitoring and Profiling.

![APM](/images/stack-profiling-apm-link.png){: style="width:200px"}
![Stack Monitoring](/images/stack-profiling-stack-monitoring-link.png){: style="width:200px"}
![Profiler](/images/stack-profiling-profiler-link.png){: style="width:200px"}
![Logs](/images/stack-profiling-logs-link.png){: style="width:200px"}

## Spotting issues

The [Query Analyzer Dashboard][] would help us to spot issues in the Elastic Stack.
We are looking for peaks of resource usage in the visualizations we have in the dashboard.
We can correlate the peaks in the visualizations with the queries we are running.
A peak in resource usage would indicate a bottleneck in the Elastic Stack.
By narrowing down the time range in the dashboard, we can go for more details in the logs, metrics, APM, Stack Monitoring, and Profiler.

Latency, Throughput, and Rate visualizations give us a big picture of the query execution and the time range where we have to focus.

![Latency, Throughput, and Rate](/images/stack-profiling-latency-throughput-rate.png){: style="width:600px"}

Zooming in the Latency, Throughput, and Rate visualizations we can see the time range where we have to focus.

![Latency, Throughput, and Rate Zoom](/images/stack-profiling-latency-throughput-rate-01.png){: style="width:600px"}

At this point, we can review the next groups of metrics to try to identify the bottleneck:

![cpu](/images/stack-profiling-cpu.png){: style="width:600px"}
![memory](/images/stack-profiling-memory.png){: style="width:600px"}
![IO](/images/stack-profiling-io.png){: style="width:600px"}
![GC](/images/stack-profiling-gc.png){: style="width:600px"}

In our example there is nothing really alarming but a 40% CPU usage is a good indicator that we have to look at the CPU metrics.

![cpu range zoom](/images/stack-profiling-cpu-zoom.png){: style="width:600px"}

At this point, we need more details so we can go to stack monitoring to see more details.

![Stack Monitoring](/images/stack-profiling-stack-monitoring-link.png){: style="width:200px"}

![es overview](/images/stack-profiling-es-overview.png){: style="width:600px"}

![es nodes](/images/stack-profiling-es-nodes.png){: style="width:600px"}

The overview and nodes visualizations give us a good idea of the resource usage in the cluster.
We know that the CPU usage was in the node `eck-profiling-ihata-es-default-2` so we can go to the node details.

![es node details](/images/stack-profiling-es-node-details.png){: style="width:600px"}

If we do not find anything in the node details, we can go to the APM data.

![APM](/images/stack-profiling-apm-link.png){: style="width:200px"}

The profiler data ...

![Profiler](/images/stack-profiling-profiler-link.png){: style="width:200px"}

Or even the logs.

![Logs](/images/stack-profiling-logs-link.png){: style="width:200px"}

## Advanced usage

### Individual Components

The `stack-profiling` template deploys allows to deploy individual components of the Elastic Stack.
To deploy individual components, you need to use the `stack-profiling` template and set the `enabled` parameter to `false` for the components you do not want to deploy.

For example, to deploy only Elasticsearch, you can run the following command:

```bash
oblt-cli cluster create custom \
    --template stack-profiling \
    --parameters '{
      "StackVersion":"8.11.3",
      "ElasticsearchEnabled": true,
      "KibanaEnabled": false,
      "ElasticAgentEnabled": false,
      }'
```

### Fixed Docker Images

The `stack-profiling` template allows to use fixed docker images for the Elastic Stack components.
Each part of the stack has a parameter to set the docker image to use.

For example, to use the `8.11.3` version of Elasticsearch, you can run the following command:

```bash
oblt-cli cluster create custom \
    --template stack-profiling \
    --parameters '{
      "StackVersion":"8.11.3",
      "ElasticsearchEnabled": true,
      "ElasticsearchDockerImage":"docker.elastic.co/observability-ci/elasticsearch-cloud:8.11.3-65c4b655",
      "KibanaEnabled": true,
      "KibanaDockerImage":"docker.elastic.co/observability-ci/kibana-cloud:8.11.3-65c4b655",
      "ElasticAgentEnabled": true,
      "ElasticAgentDockerImage":"docker.elastic.co/observability-ci/elastic-agent-cloud:8.11.3-65c4b655"
      }'
```

### Query with curl

In some cases, you may want to run a query directly in Elasticsearch with `curl` instead of using Kibana.
To do that, you can use the `ELASTICSEARCH_URL` and `ELASTICSEARCH_USERNAME` and `ELASTICSEARCH_PASSWORD` environment variables.

For example, to run the query we used in the previous example, you can run the following command:

```bash
CLUSTER_NAME=eck-profiling-ihata
ENV_FILE="${PWD}/.env"

# Get the secrets from the cluster
oblt-cli cluster secrets env --cluster-name ${CLUSTER_NAME} --output-file "${ENV_FILE}"
source "${ENV_FILE}"

curl -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" -X GET "${ELASTICSEARCH_HOST}/metrics-system.cpu-default/_search" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "data_stream.dataset": "system.cpu"
    }
  }
}'
```

### Esrally

In some cases, you may want to profile a query with a specific data and during a specific time range.
To do that, you can use the [esrally][] tool. [esrally][] is a tool that allows you to run benchmarks against Elasticsearch.
The [esrally][] tool simplifies the way to load data and run some operations against the Elastic Stack.
[esrally][] uses `tracks`and each execution of a track is called a `race`.
We can use existing `tracks` from the [rally-tracks][] repository or create our own `tracks`.

!!! Note

    To install [esrally][] you can follow the [official documentation](https://esrally.readthedocs.io/en/stable/install.html).

#### Run a track

To run a `race` you only need [esrally][] and and the Elastic Stack credentials.
An execution of a `race` from the [rally-tracks][] repository would look like this:

```bash
CLUSTER_NAME=eck-profiling-ihata
TRACK_NAME=elastic/logs
CHALLENGE=logging-indexing-querying
RACE_ID=${CLUSTER_NAME}-${TRACK_NAME}-$(date -Iseconds)
ENV_FILE="${PWD}/.env"

# Get the secrets from the cluster
oblt-cli cluster secrets env --cluster-name ${CLUSTER_NAME} --output-file "${ENV_FILE}"
source "${ENV_FILE}"

# run a track race
esrally race \
  --track="${TRACK_NAME}" \
  --challenge="${CHALLENGE}" \
  --target-hosts "${ELASTICSEARCH_HOST}" \
  --client-options "basic_auth_user:'${ELASTICSEARCH_USERNAME}',basic_auth_password:'${ELASTICSEARCH_PASSWORD}'" \
  --race-id "${RACE_ID}" \
  --cluster-name "${CLUSTER_NAME}" \
  --pipeline benchmark-only
```

#### Run custom track

It is possible to run a custom track.
You would need some data and a track file.
You can use existing data in a cluster to create a new track check [Creating a track from data][].
You can also create a track from scratch check, see [creating a track from scratch][]

In the oblt cluster repo we have an example of custom track that uses the geonames data.

```bash
CLUSTER_NAME=eck-profiling-ihata
ENV_FILE="${PWD}/.env"

# Clone the repo
git clone https://github.com/elastic/observability-test-environments.git
cd observability-test-environments

# Get the secrets from the cluster
oblt-cli cluster secrets env --cluster-name ${CLUSTER_NAME} --output-file "${ENV_FILE}"
source "${ENV_FILE}"

# Generate the track
cd deployments/track/sample
./gen-track.sh

# run a track race
./track.sh
```

The `gen-track.sh` script will generate a track file and the data to run the track.
You can check the `track.sh` script to see how to run a track.
Also, the `track.json.tmpl` file is a good example of a track file.

[Query Analyzer Dashboard]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/dashboards#/view/a7d4baea-5956-4749-bffa-d289e5d660b4?_g=(filters:!(),time:(from:now-15m,to:now))&_a=(query:(language:kuery,query:'cluster_uuid:*%20OR%20(labels.deploymentName:*%20and%20data_stream.type:traces%20AND%20transaction.name:*%20AND%20url.path:*)%20OR%20(labels.deploymentName:*%20and%20data_stream.type:metrics)'))
[monitoring-oblt]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/
[retrieve credentials]: /user-guide/use-case-retrieve-credentials.md
[esrally]: https://esrally.readthedocs.io/en/stable/index.html
[rally-tracks]: https://github.com/elastic/rally-tracks
[creating a track from data]: https://esrally.readthedocs.io/en/stable/adding_tracks.html#creating-a-track-from-data-in-an-existing-cluster
[creating a track from scratch]: https://esrally.readthedocs.io/en/stable/adding_tracks.html#creating-a-track-from-scratch
[restore a snapshot]: /user-guide/oblt-cluster-snapshots.md
