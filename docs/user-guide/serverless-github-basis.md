## From GitHub

For the developers that prefer GitHub, we have a issue drive process on the [observability-test-environments][] repository.
This process can provision clusters only.

### Prerequisites

Before you start, make sure you have the following prerequisites:

* Access to the repository [observability-test-environments][]

### Create

To create a serverless cluster, go to the [Create a Create Serverless Cluster GitHub issue][] and fill in the template.

You have to choose a title for the issue and fill in the following fields:

* **Environment**: The environment where the cluster will be deployed.
* **Project Type**: The project type of the cluster.
* **Cluster Prefix**: The prefix of the cluster name.
* **Cluster Suffix**: The suffix of the cluster name.

![github issue create serverless](/images/github-issue-serverless.png){: style="width:600px"}

[observability-test-environments]: https://github.com/elastic/observability-test-environments
[Create a Create Serverless Cluster GitHub issue]: https://github.com/elastic/observability-test-environments/issues/new?assignees=&labels=cluster%2Cserverless-cluster&projects=&template=cluster-serverless-issue.yaml&title=%5BServerless+Cluster%5D%3A+
