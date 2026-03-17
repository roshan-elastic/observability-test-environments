## From Command Line

For the developers that prefer the command line, we have the `oblt-cli` tool.
This tool provision and update oblt clusters along with other operations.

### Prerequisites

Before you start, make sure you have the following prerequisites:

* [oblt-cli][] installed.

### Create

To create a serverless cluster, run the following command:

```shell
oblt-cli cluster create serverless
```

This command will deploy a Serverless project with the default project type `observability`.
The project will receive logs, metrics, and traces from applications deployed on a k8s cluster.
After [oblt-cli][] resquested the project, process will take a few minutes to complete, you will received an Slack message on the User/Channel you configured in oblt-cli.
By default, the cluster will be deployed on the [QA environment][].
Projects are automatically destroyed after three days.
In the following example, we are creating a serverless cluster with the default project type `observability`,
we can see a banner with some news about the oblt-cli tool, a warning about the version of the tool, the template used to create the cluster, and the cluster configuration was created successfully.
The cluster name is `serverless-tupxw`. All the notifications are sent to the `#observablt-bots` Slack channel.

```log
[2024-05-06T15:26:23.32Z]:[info]:Git mode: 'ssh'
[2024-05-06T15:26:36.355Z]:[info]:Using template : /Users/my-user/.oblt-cli/observability-test-environments/environments/users/serverless.yml.tmpl
[2024-05-06T15:26:36.355Z]:[info]:Template parameters: map[ClusterConfigFile:/Users/my-user/.oblt-cli/observability-test-environments/environments/users/my-user/serverless-tupxw.yml ClusterName:serverless-tupxw Date:2024-05-06T17:26:36+0200 GitOps:false SlackChannel:#observablt-bots TemplateName:serverless TemplatePath:/Users/my-user/.oblt-cli/observability-test-environments/environments/users/serverless.yml.tmpl Username:my-user]
[2024-05-06T15:26:39.968Z]:[info]:Cluster created successfully: serverless-tupxw
```

{% include '/user-guide/serverless-local-development-basis.md' ignore missing %}

### destroy

To destroy a serverless cluster, run the following command:

```shell
CLUSTER_NAME=serverless-tupxw
oblt-cli cluster destroy serverless --cluster-name "${CLUSTER_NAME}"
```

This command will destroy the serverless cluster with the name `serverless-tupxw`.
The command will request confirmation before destroying the cluster.
After the confirmation, the request to destroy the cluster will be sent, it will take a few minutes to complete.

```log
[2024-05-06T15:37:53.922Z]:[info]:Using config file: /Users/my-user/.oblt-cli/config.yaml
[2024-05-06T15:37:53.922Z]:[info]:SlackChannel: '#observablt-bots'
[2024-05-06T15:37:53.922Z]:[info]:User: 'my-user'
[2024-05-06T15:37:53.922Z]:[info]:Git mode: 'ssh'
? Do you want to destroy the cluster : serverless-tupxw? (yes/no) Yes
[2024-05-06T15:40:13.047Z]:[info]:Destroying cluster : serverless-tupxw
[2024-05-06T15:40:13.053Z]:[info]:Removing cluster file : /Users/my-user/.oblt-cli/observability-test-environments/environments/users/my-user/serverless-tupxw.yml
```

### Project details

To get the project details, like endpoints, credentials, and other information, run the following command:

```shell
CLUSTER_NAME=serverless-tupxw
oblt-cli cluster secrets credentials --cluster-name "${CLUSTER_NAME}"
```

```log
[2024-05-07T15:35:17.469Z]:[info]:Using config file: /Users/inifc/.oblt-cli/config.yaml
[2024-05-07T15:35:17.469Z]:[info]:SlackChannel: '#observablt-bots'
[2024-05-07T15:35:17.469Z]:[info]:User: 'kuisathaverat'
[2024-05-07T15:35:17.469Z]:[info]:Git mode: 'ssh'
*Cluster Details*

* Oblt username: robots
* Kibana:
  * URL: https://my-project-oblt-ad85d6.kb.us-east-1.aws.elastic.cloud
* Elasticsearch:
  * URL: https://my-project-oblt-ad85d6.es.us-east-1.aws.elastic.cloud

*Cluster Management Links*


``
# Retrieve cluster redentials
oblt-cli cluster secrets credentials --cluster-name my-project-oblt

# Destroy cluster
oblt-cli cluster destroy --cluster-name my-project-oblt

# Retrieve Kibana YAML sample file
oblt-cli cluster secrets kibana-config --cluster-name my-project-oblt
``

To retrieve the secrets you can use oblt-robot, the command is `/cluster-secret`

*help:*
  * https://ela.st/oblt-cli
  * https://ela.st/oblt-robot
  * https://ela.st/oblt-clusters
  * https://ela.st/oblt-deploy
"

*Credentials*

ingress:
  username: admin
  password: MyPaSsWord
elasticsearch:
  login url: https://my-project-oblt-ad85d6.es.us-east-1.aws.elastic.cloud
  username: testing-internal
  password: MyPaSsWord-2

  curl -kv -u testing-internal:MyPaSsWord-2 https://my-project-oblt-ad85d6.es.us-east-1.aws.elastic.cloud

  user: testing-internal
  password: MyPaSsWord-2

kibana:
  url: https://my-project-oblt-ad85d6.kb.us-east-1.aws.elastic.cloud
  username: testing-internal
  password: MyPaSsWord-2

    This Kibana deployment has Okta login enabled, you can use Okta to login to Kibana.

  Check oblt-cli documentation to <https://elastic.github.io/observability-test-environments/tools/oblt-cli/local-kibana/|know hot to connect a local Kibana instance>
apm:
  url: https://my-project-oblt-ad85d6.apm.us-east-1.aws.elastic.cloud
  token: APM_TOKEN

  curl -kv -H "Authorization: Bearer APM_TOKEN" https://my-project-oblt-ad85d6.apm.us-east-1.aws.elastic.cloud
opbeans:
  username: admin
  password: MyPaSsWord
```

[oblt-cli]: https://ela.st/oblt-cli
[QA environment]: https://docs.elastic.dev/serverless/qa
[MKI]: https://docs.elastic.dev/mki
