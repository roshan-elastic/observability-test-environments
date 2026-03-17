# Destroy a Cluster

* [How to destroy a cluster with oblt robot](#how-to-destroy-a-cluster-with-oblt-robot)
* [How to destroy a cluster with oblt CLI](#how-to-destroy-a-cluster-with-oblt-cli)
* [How to destroy a cluster Manually](#how-to-destroy-a-cluster-manually)
* [How to destroy a cluster locally](#how-to-destroy-a-cluster-locally)

## How to destroy a cluster with oblt robot

The fastest and easiest way to destroy an oblt cluster is using [oblt robot][oblt-robot],
you do not need to install anything, you interact with a Slack robot directly in the Slack app.

* Go to the `#observablt-bots` Slack channel
* Write the comment `/destroy-cluster`
* Choose the cluster to destroy.

![destroy-cluster](../images/oblt-robot-destroy-cluster.png){: style="width:450px"}

* Wait for the bot to send you the confirmation By Slack.

For more commands and info check the documentation at [oblt robot][oblt-robot]

## How to destroy a cluster with oblt CLI

If you preferred the command line to destroy an oblt cluster,
you can use [oblt CLI][oblt-cli].

* List your clusters

```bash
oblt-cli cluster list
```

* Destroy your cluster using the name you got in the previous step

```bash
oblt-cli cluster destroy --cluster-name my-cluster-name
```

* Wait for the CI to send you the confirmation By Slack.

For more commands and info check the documentation at [oblt CLI][oblt-cli]

## How to destroy all my clusters

If you want to destroy all your clusters at once you can.
The command following command will destroy all your clusters.

```bash
oblt-cli cluster destroy --wipeup
```

## How to destroy a cluster Manually

It is possible to make the destroy operation manually
by deleting the cluster configuration file from the repository.

* Pre-requisites
  * Git installed
* Checkout the repo `git@github.com:elastic/observability-test-environments.git`
* delete the `environments/users/USERNAME/my-config-cluster.yml` file
* Create a feature branch and push your changes to the repo
* Create a new PR
* The [CI will trigger a new build][oblt-manager] to destroy the cluster.

## How to destroy a cluster locally

* Pre-requisites
  * Git installed
  * Python3 installed
  * [Elastic Vault service][] access
* Checkout the repo `git@github.com:elastic/observability-test-environments.git`
* Execute the `destroy-cluster` target

```bash
CLUSTER_CONFIG_FILE=$(pwd)/environments/USERNAME/my-config-cluster.yml \
make -C ansible destroy-cluster
```

* delete the `environments/users/USERNAME/my-config-cluster.yml` file
* Create a feature branch and push your changes to the repo
* Create a new PR
* The CI will fail because the cluster was destroyed so you can merge the PR directly.

[oblt-cli]: http://ela.st/oblt-cli
[oblt-robot]: https://ela.st/oblt-robot
[oblt-manager]: https://github.com/elastic/observability-test-environments/actions/workflows/cluster-manager.yml
