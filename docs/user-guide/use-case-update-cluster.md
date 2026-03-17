# Update existing cluster

## Update with oblt-cli

The easiest way to update a cluster is using the [oblt-cli](/tools/oblt-cli.md).
The `oblt-cli` will allow you to update a cluster with a single command.
The command `oblt-cli cluster update` will update a cluster configuration file with the new configuration you provide.
The new configuration will be merged with the existing one.
The new configuration is passed as a JSON file or as a JSON string.

```bash
CLUSTER_NAME=my-oblt
PARAMS_JSON='{"stack":{"version":"8.12.0-SNAPSHOT"}}'
oblt-cli cluster update --cluster-name ${CLUSTER_NAME} --parameters ${PARAMS_JSON}
```

```bash
CLUSTER_NAME=my-oblt
PARAMS_JSON_FILE=/tmp/params.json
oblt-cli cluster update --cluster-name ${CLUSTER_NAME} --parameters-file "${PARAMS_JSON_FILE}"
```

For details about the configuration settings available see [cluster configuration](./cluster-config.md).

!!! Warning

    The name of the cluster cannot be changed.

!!! Note

    The cluster configuration file is a YAML file, but the `oblt-cli` will accept a JSON file or a JSON string as input.

!!! Note

    Trigger a change `updated_at` value will trigger a full restart of the cluster without changes in the configuration.
    ```bash
    CLUSTER_NAME=my-oblt
    PARAMS_JSON='{"stack":{"updated_at":"2023-09-18:00:00Z"}}'
    oblt-cli cluster update --cluster-name ${CLUSTER_NAME} --parameters ${PARAMS_JSON}
    ```

## How to update a cluster Manually

In the git repository `https://github.com/elastic/observability-test-environments`
there is a folder `environments/users/USERNAME` that contains a folder per environment.
Every environment folder will have a cluster configuration file per cluster.
Any change we made on a cluster configuration file in the repository will trigger a update process.
so to edit a cluster we will need to edit a cluster configuration file and push the changes to the repo.
This are the steps to follow:

* Pre-requisites
  * Git installed
* Checkout the repo `git@github.com:elastic/observability-test-environments.git`
* Go to the user's folder `environments/users/USERNAME`.
* Edit the cluster config file `my-config-cluster.yml`
* Create a feature branch and push your changes to the repo
* Create a new PR
* The [CI will trigger a new build][oblt-manager] to create the new cluster based in your cluster configuration file.

You have examples of configuration at [test configurations](../tests/environments)

[oblt-manager]: https://github.com/elastic/observability-test-environments/actions/workflows/cluster-manager.yml
