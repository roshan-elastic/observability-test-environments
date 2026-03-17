# Opbeans Cluster template

The oblt cluster template `opbeans` allows to deploy a ESS cluster for APM testing.
An `opbeans` cluster will deploy the Elastic Stack in ESS and Opbeans on a k8s cluster.
The `Opbeans` are configured to report APM data to the ESS cluster using load generators.
Additionally the Elastic Agent running in standalone mode is deployed to grab logs and metrics from k8s.

The following command will deploy a cluster for the Elastic Stack version `8.7.0`

```bash
  oblt-cli cluster create custom --template opbeans --parameters '{"StackVersion": "8.7.0"}'
```

It is possible to use any Stack Pack version deployed on ESS.
The Stack pack usually are the version of the stack or the version plus `-SNAPSHOT`.

```bash
  oblt-cli cluster create custom --template opbeans --parameters '{"StackVersion": "8.7.0-SNAPSHOT"}'
```

Finally, in the case we need specific Docker images we can set the Docker image for each component of the Elastic Stack.
It is important that those Docker images were build to run on ESS, if not they will not work.

```bash
  oblt-cli cluster create custom \
          --template opbeans \
          --parameters '{"StackVersion":"8.3.0-SNAPSHOT","ElasticsearchDockerImage":"docker.elastic.co/observability-ci/elasticsearch-cloud:8.3.0-aaaaaa","KibanaDockerImage":"docker.elastic.co/observability-ci/kibana-cloud:8.3.0-aaaaaa","ElasticAgentDockerImage":"docker.elastic.co/observability-ci/elastic-agent-cloud:8.3.0-aaaaaa"}'
```
