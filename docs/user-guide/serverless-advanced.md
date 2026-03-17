# Serverless Advanced

## Environments

Serverless has three environments where we can deploy our serverless clusters:

* `qa`: Deploy a cluster on the [QA environment][] **(default)**
* `staging`: Deploy a cluster on the [Staging environment][] environment
* `production`: Deploy a cluster on the [Production environment][] environment

The deployment is made using the API described at [Serverless Project API reference][]

The Elastic Stack is deployed on the given environment, the deployment uses the Elastic Stack available on
the environment at the time of the creation.

## Project types

The project types available are:

* `observability`: Deploy a cluster with the observability project (default) (Elasticsearch, Kibana, APM, Fleet)
* `elasticsearch`: Deploy a cluster with the elasticsearch project. (Elasticsearch, Kibana)
* `security`: Deploy a cluster with the security project. (Elasticsearch, Kibana, Fleet)

## Oblt Tools

Observability has different user types. Thus, we provide various tools adapted to users' needs.
You can use three different approaches:

* Use [oblt-robot][] with the slack command `/create-serverless-cluster`.
* [Create a GitHub issue][] to deploy your own serverless using the existing automation.
* Use [oblt-cli][] and the `oblt-cli cluster create serverless`.

{% include '/user-guide/serverless-oblt-cli-basis.md' ignore missing %}

{% include '/user-guide/serverless-oblt-cli-adv.md' ignore missing %}

{% include '/user-guide/serverless-access.md' ignore missing %}

{% include '/user-guide/serverless-slack-basis.md' ignore missing %}

{% include '/user-guide/serverless-access.md' ignore missing %}

{% include '/user-guide/serverless-github-basis.md' ignore missing %}

{% include '/user-guide/serverless-access.md' ignore missing %}

{% include '/user-guide/serverless-access.md' ignore missing %}

{% include '/user-guide/mki.md' ignore missing %}

change expiration date
Add k8s deployments
Import data

[QA environment]: https://docs.elastic.dev/serverless/qa
[Staging environment]: https://docs.elastic.dev/serverless/staging
[Production environment]: https://docs.elastic.dev/serverless/production
[Serverless Project API reference]: https://docs.google.com/document/d/1CRVHisgpU-e1uUbjLySy4lPaFanJ0FwpVXp2Wu2JkGU/edit#
[Serverless Project API]: https://backstage.elastic.dev/catalog/default/api/project-api/definition#/default
[oblt-cli]: https://ela.st/oblt-cli
[oblt-robot]: https://ela.st/oblt-robot
[Create a GitHub issue]: https://ela.st/create-serverless-oblt
[this GitHub issue]: https://ela.st/self-service-oblt-robots-org
[observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
[Serverless environments]: https://docs.elastic.dev/serverless/environments
