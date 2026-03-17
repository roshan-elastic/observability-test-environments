## Logs, Metrics, and APM

The serverless projects send logs, metrics, and traces to some regional clusters.
Every environments has it own regional cluster to store the logs, metrics, and traces.
It is possibel to check the logs of every component of the project by filtering the logs by the project ID, and component type.

You have detailed info about the regional clusters in the [observability production][], [observability QA][], and [observability staging][] documentation.

### filters

Here you have some useful filters to check the logs, metrics, and traces of the serverless projects.

For logs and metrics you can use the following filters:

* `serverless.project.id`: The project ID of the serverless project.
* `kubernetes.labels.common_k8s_elastic_co/type`: The component type of the project.
  * `elasticsearch`: Elasticsearch logs
  * `kibana`: Kibana logs
  * `fleet`: Elastic Agent logs
* `kubernetes.labels.k8s_elastic_co/application-id`: The component type of the project.
  * `es`: Elasticsearch logs
  * `kibana`: Kibana logs
  * `fleet`: Elastic Agent logs
* `log.level`: The log level of the logs.
  * `DEBUG`: Debug logs
  * `debug`: Debug logs
  * `INFO`: Info logs
  * `info`: Info logs
  * `WARN`: Warning logs
  * `warn`: Warning logs
  * `ERROR`: Error logs
  * `error`: Error logs
  * `FATAL`: Fatal logs
  * `fatal`: Fatal logs
  * `TRACE`: Trace logs
  * `trace`: Trace logs

For APM you can use the following filters:

* `labels.project_id`: The project ID of the serverless project.
* `labels.projectId`: The project ID of the serverless project.
* `labels.project_type`: The project type of the serverless project.
* `labels-projectType`: The project type of the serverless project.
* `labels.project_name`: The project name of the serverless project.
* `labels.projectName`: The project name of the serverless project.
* `service.name`: The service name of the APM project.
  * `elasticsearch`: Elasticsearch APM
  * `kibana`: Kibana APM
  * `fleet-server`: Elastic Agent APM

[Observability Production]: https://docs.elastic.dev/serverless/observability/production
[Observability QA]: https://docs.elastic.dev/serverless/observability/qa
[Observability staging]: https://docs.elastic.dev/serverless/observability/staging
