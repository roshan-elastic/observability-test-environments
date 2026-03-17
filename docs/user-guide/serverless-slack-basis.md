## From Slack

For the developers that prefer Slack, we have the `oblt-robot` tool.
This tool provision and update oblt clusters along with other operations.

### Prerequisites

Before you start, make sure you have the following prerequisites:

* [Slack][] installed.
* Access to the [oblt-robot][]. It is available in a few channels:
  * `#observablt-robots`

### Create

To create a serverless cluster, go to the [#observablt-robots](https://elastic.slack.com/archives/CJMURHEHX) channel and type the following command:

```shell
/create-serverless-cluster
```

Slack will show a dialog to configure the environment, project type, and cluster name prefix/suffix.
After you fill in the information, the serverless project will be requested.
The process will take a few minutes to complete. You will receive a Slack message on your Slack user.

![serverless create](/images/oblt-robot-create-serverless-cluster.png){: style="width:600px"}

### Destroy

To destroy a serverless cluster, go to the [#observablt-robots](https://elastic.slack.com/archives/CJMURHEHX) channel and type the following command:

```shell
/destroy-cluster
```

Slack will show a dialog to confirm the cluster name.

![destroy cluster](/images/oblt-robot-destroy-cluster.png){: style="width:600px"}

[oblt-robot]: https://ela.st/oblt-robot
[Slack]: https://slack.com/issues/new?assignees=&labels=cluster%2Cserverless-cluster&projects=&template=cluster-serverless-issue.yaml&title=%5BServerless+Cluster%5D%3A+
