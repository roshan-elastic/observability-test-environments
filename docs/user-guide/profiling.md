# Profiling Cluster templates

Profiling templates allow to deploy the Elastic Stack configure profiling and deploy some applications to test the profiling in a k8s cluster.

## ESS template

The ESS profiling template deploy the Elastic Stack in ESS and deploy some applications to test the profiling.

```bash
oblt-cli cluster create custom \
    --template ess-profiling.yml \
    --parameters '{"StackVersion":"8.7.0"}'
```

## ECK template

The ESS profiling template deploy the Elastic Stack with ECK and deploy some applications to test the profiling.
All deployments are in the same k8s cluster so profiling is reported for all the applications and the Elastic Stack.

```bash
oblt-cli cluster create custom \
    --template eck-profiling.yml \
    --parameters '{"StackVersion":"8.7.0"}'
```

## Provide custom Docker images

It is possible to provide custom Docker images for the Elastic Stack components.
Each component has its own parameter:

```bash
oblt-cli cluster create custom \
    --template ess-profiling.yml \
    --parameters '{
        "StackVersion":"8.7.0-SNAPSHOT",
        "ElasticsearchDockerImage":"docker.elastic.co/observability-ci/elasticsearch-cloud:8.7.0-SNAPSHOT",
        "KibanaDockerImage":"docker.elastic.co/observability-ci/kibana-cloud:8.7.0-SNAPSHOT",
        "ElasticAgentDockerImage":"docker.elastic.co/observability-ci/elastic-agent-cloud:8.7.0-SNAPSHOT"
        }'
```

```bash
oblt-cli cluster create custom \
    --template eck-profiling.yml \
    --parameters '{
        "StackVersion":"8.7.0-SNAPSHOT",
        "ElasticsearchDockerImage":"docker.elastic.co/observability-ci/elasticsearch:8.7.0-SNAPSHOT",
        "KibanaDockerImage":"docker.elastic.co/observability-ci/kibana:8.7.0-SNAPSHOT",
        "ElasticAgentDockerImage":"docker.elastic.co/observability-ci/elastic-agent:8.7.0-SNAPSHOT"
        }'
```
