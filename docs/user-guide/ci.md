# oblt-cli in CI

The `oblt-cli` is a tool to manage the observability test environments.
Like we use it for our users we can use it in the same way in CI.
The goal is to isolate the provision of environments and Elastic Stack from the CI tests.
The `oblt-cli` can deploy the Elastic Stack in ESS, ECK, and serverless.
`Oblt-cli` can also deploy k8s clusters in GKE.
It is possible to deploy k8s resources in the cluster using the cluster configuration file.
In the folder [deployments][] you can find some preconfigured deployments.

## GitHub Actions

The `oblt-cli` is used in the GitHub Actions workflows.
We have created some GitHub Actions to make it easier to use it.
The actions take care of the installation and configuration of the `oblt-cli`.

You can check the GitHub Actions documentation for more details.

* [setup-oblt-cli][]: Install and configure the `oblt-cli`.
* [oblt-cli-create-ccs][]: Create a CCS cluster.
* [oblt-cli-create-custom][]: Create a custom cluster.
* [oblt-cli-create-serverless][]: Create a serverless cluster.
* [oblt-cli-cluster-credentials][]: Read the credentials for the given cluster.
* [oblt-cli-destroy-cluster][]: Destroy the given cluster.

### Examples

The following example shows how to use it in a GitHub Actions workflow.
The workflow will configure the `oblt-cli` and create an ESS cluster, then run some API calls and destroy the cluster.

```yaml
{% include '/samples/gh-wf.yml' ignore missing %}
```

!!! Note

    The example uses the [oblt-actions/google/auth](https://github.com/elastic/oblt-actions/tree/main/google/auth) action to authenticate to Google Cloud through Workload Identity Federation.
    See the [Workload Identity Federation](#workload-identity-federation) section on how to set it up.

## Buildkite Pipeline

You can also use `oblt-cli` in Buildkite pipelines.
The [oblt-cli buildkite plugin][oblt-cli-buildkite-plugin]  helps you to set up the `oblt-cli` in your pipeline.

The following example shows how to use it in a Buildkite pipeline.
It showcases how to create an ESS cluster, how to export the credentials as environment variables,
and how to destroy the cluster.

```yaml
{% include '/samples/bk-pipeline.yml' ignore missing %}
```

!!! Note

    The above example uses the [oblt-google-auth buildkite plugin](https://github.com/elastic/oblt-google-auth-buildkite-plugin) buildkite plugin to authenticate to Google Cloud through Workload Identity Federation.
    See the [Workload Identity Federation](#workload-identity-federation) section on how to set it up.


## Workload Identity Federation

The following example shows how to set up the Workload Identity Federation in [elastic/oblt-infra][] so it can be used
both in Buildkite and GitHub Actions workflows.

```shell
{% include '/samples/setup-workload-identity-federation.hcl' ignore missing %}
```

## Install

The `oblt-cli` is a Go binary, so you can download it from the [releases][] page or build it from source.
There is a sample installation script in the [install.sh][] file.

```bash
PATH=${HOME}/bin:${PATH}
GITHUB_AUTH=my-username:$(gh auth token)
SCRIPT_URL=https://raw.githubusercontent.com/elastic/observability-test-environments/main/tools/oblt-cli/scripts/install.sh
curl -fsSOL -u "${GITHUB_AUTH}" "${SCRIPT_URL}"
chmod ugo+x install.sh
./install.sh "${HOME}/bin"
```

!!! Note

    There are more ways to install the `oblt-cli` see [oblt-cli documentation](/tools/oblt-cli/#installation)

## Configure

The `oblt-cli` needs to be configured to access the observability test environments.
It needs a username to create the configuration folder and a Slack Channel to report.
The configuration is stored in `~/.oblt/config.yaml`.

```bash
OBLT_USERNAME=your-username
SLACK_CHANNEL=#SlackCkannel
oblt-cli configure --username ${OBLT_USERNAME} --slack-channel ${SLACK_CHANNEL}
```

In the CI it is convenient to make it with less commands, so we have added some global flags to the `oblt-cli` to configure it.
The following command will execute a command and configure `oblt-cli` with the username and slack channel.

```bash
OBLT_USERNAME=your-username
SLACK_CHANNEL=#SlackCkannel
oblt-cli cluster create ess --stack-version 8.7.0-SNAPSHOT --username ${OBLT_USERNAME} --slack-channel ${SLACK_CHANNEL} --save-config --git-http-mode
```

For more details check `oblt-cli configure --help` and `oblt-cli --help`

## Create a cluster

There are a few clusters types and templates available.

```bash
oblt-cli cluster create --help
Command to create a cluster.

Usage:
  oblt-cli cluster create [command]

Available Commands:
  ccs         Command to create a CCS cluster.
  custom      Command to create a cluster from a template.
  ess         Command to create an ESS cluster.
```

### ESS cluster

Command to create a custom cluster using the ESS template. If used, any of the docker images must use the same version as the stack version, otherwise the build of the cluster will fail in Elastic Cloud.

```bash
oblt-cli cluster create ess --stack-version 8.7.0-SNAPSHOT --cluster-name-prefix my-job
```

for more details check `oblt-cli cluster create ess --help`

### CCS cluster

Command to create a Cross Cluster Search, aka CCS, cluster. This cluster configures an oblt cluster as remote cluster to use CCS.

```bash
REMOTE_CLUSTER=golden-cluster
oblt-cli cluster create ccs --template-name ccs --remote-cluster=${REMOTE_CLUSTER} --cluster-name-prefix my-job
```

!!! Note

    `golden-cluster` is just a random name,
    so you can see the list of golden-clusters with `oblt-cli cluster list --filter golden_cluster=true`.

For more details check `oblt-cli cluster create ccs --help`

### Custom cluster

Finally, you can create a cluster from a template.
The `oblt-cli cluster custom` command allows you to use any of the [cluster templates][].

```bash
oblt-cli cluster create custom \
  --template ess \
  --parameters '{"StackVersion": "8.10.0", "ProductType": "elasticsearch"}' \
  --cluster-name-prefix my-job
```

```bash
oblt-cli cluster create custom \
  --template-name ccs \
  --parameters '{"RemoteClusterName": "golden-cluster"}' \
  --cluster-name-prefix my-job
```

As you can see `ess`and `ccs` are just templates we can use with the `custom` command.

for more details check `oblt-cli cluster create custom --help`

### Save cluster info

The `oblt-cli` can save the cluster info in a file to be used later.
This is useful to retrieve credentials and make other operations over the cluster.

```bash
oblt-cli cluster create ess \
  --stack-version 8.7.0-SNAPSHOT \
  --output-file ${PWD}/cluster-info.json \
  --cluster-name-prefix my-job
```

This is how `cluster-info.json` looks like:

```JSON
{
  "ClusterConfigFile": "/Users/myusername/.oblt-cli/observability-test-environments/environments/users/myusername/ess-oqqhr.yml",
  "ClusterName": "ess-oqqhr",
  "CommitMessage": "oblt-cli(myusername): Create ess-oqqhr cluster [template=ess]",
  "CommitSha": "fbbdd91138663b5eb7da46fb74a72ab5b09c5321",
  "CommitURL": "https://github.com/elastic/observability-test-environments/commit/fbbdd91138663b5eb7da46fb74a72ab5b09c5321",
  "Date": "2023-09-26T17:20:40+0200",
  "ElasticAgentDockerImage": "docker.elastic.co/observability-ci/elastic-agent-cloud:8.7.0-SNAPSHOT",
  "ElasticsearchDockerImage": "docker.elastic.co/observability-ci/elasticsearch-cloud-ess:8.7.0-SNAPSHOT",
  "GitOps": false,
  "KibanaDockerImage": "docker.elastic.co/observability-ci/kibana-cloud:8.7.0-SNAPSHOT",
  "SlackChannel": "@UCKPL50JY",
  "StackVersion": "8.7.0-SNAPSHOT",
  "TemplateName": "ess",
  "TemplatePath": "/Users/myusername/.oblt-cli/observability-test-environments/environments/users/ess.yml.tmpl",
  "Username": "myusername"
}
```

### Wait for the cluster

The `oblt-cli` can wait for the cluster to be ready.
The cluster creation, update and destroy operations are caused by GitHub pull requests. So those changes are reflected in some files and hence GitHub pull requests drive those operations.
It is possible to monitor the Pull Request created waiting for a merge.
We only need a GITHUB_TOKEN to access the GitHub API, and set the `--wait` flag to the number of minutes to wait.

```bash
GITHUB_TOKEN=your-github-token
oblt-cli cluster create ess \
  --stack-version 8.7.0-SNAPSHOT \
  --output-file ${PWD}/cluster-info.json \
  --cluster-name-prefix my-job \
  --wait 15
```

### Credentials

The `oblt-cli` can retrieve the credentials of a cluster.
The cluster is identified by the name so you need the [cluster info](#save-cluster-info) to retrieve the credentials.

!!! Note

    The credentials are stored in Google Cloud Secret Manager.
    It requires authentication to access the secrets.
    The `oblt-cli` will use the credentials from the environment variables `GOOGLE_APPLICATION_CREDENTIALS`.
    For more details check [Google Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials).

```bash
# Authenticate to Google Cloud
SERVICE_ACCOUNT_EMAIL=my-service-account-email@elastic-observability.iam.gserviceaccount.com
GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/application_default_credentials.json
gcloud auth activate-service-account "${SERVICE_ACCOUNT_EMAIL}" --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"

# Retrieve the credentials
CLUSTER_NAME=$(jq -r .ClusterName cluster-info.json)
oblt-cli cluster secrets cluster-state --cluster-name ${CLUSTER_NAME} --output-file ${PWD}/cluster-state.yaml
```

```yaml
apm_apikey: AAAAAAAAAAAAAAAAAAAAA
apm_token: ABABAGBAGBAGBAGBA
apm_url: https://golden-cluster.apm.us-west2.gcp.elastic-cloud.com:443
certmanager_available: true
cluster_created: true
elasticsearch_apikey: VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHAyYXhUTm1zeWFrdzl0dk5udw==
elasticsearch_password: PPAASSWWOORRDD
elasticsearch_url: https://golden-cluster.es.us-west2.gcp.elastic-cloud.com:443
elasticsearch_username: elastic
ess_cloud_id: golden-cluster:dXMtd2VzdDIuZ2NwLmVsYXN0aWMtY2xvdWQuY29tOjQ0Mzg1ODg4YzE5NjAyNDM3ZmI3MWFhOTU1NzQ3MDAwMGIK==
ess_deployment_id: 02494605cd8e2d62ab6af1690719afaa
fleet_url: https://golden-cluster.fleet.us-west2.gcp.elastic-cloud.com:443
ingress_available: true
ingress_password: PPAASSWWOORRDD
ingress_username: admin
k8s_default_namespace: default
k8s_project: elastic-observability
k8s_provider: gcp
k8s_region: us-central1-c
kibana_password: PPAASSWWOORRDD
kibana_url: https://golden-cluster.kb.us-west2.gcp.elastic-cloud.com:443
kibana_username: elastic
kube_system_namespace: kube-system
serverless_deployment_id: null
```

Or to take a bash script to export the credentials as environment variables.

```bash
CLUSTER_NAME=$(jq -r .ClusterName cluster-info.json)
oblt-cli cluster secrets env --cluster-name ${CLUSTER_NAME} --output-file ${PWD}/env.sh
```

```YAML
ELASTICSEARCH_HOST: https://golden-cluster.es.us-west2.gcp.elastic-cloud.com:443
ELASTICSEARCH_HOSTS: https://golden-cluster.es.us-west2.gcp.elastic-cloud.com:443
ELASTICSEARCH_PASSWORD: PPAASSWWOORRDD
ELASTICSEARCH_USERNAME: elastic

KIBANA_FLEET_HOST: https://golden-cluster.kb.us-west2.gcp.elastic-cloud.com:443
KIBANA_HOST: https://golden-cluster.kb.us-west2.gcp.elastic-cloud.com:443
KIBANA_HOSTS: https://golden-cluster.kb.us-west2.gcp.elastic-cloud.com:443
KIBANA_PASSWORD: PPAASSWWOORRDD
KIBANA_USERNAME: elastic

ELASTIC_APM_API_KEY: PPAASSWWOORRDD
ELASTIC_APM_JS_BASE_SERVER_URL: https://golden-cluster.apm.us-west2.gcp.elastic-cloud.com:443
ELASTIC_APM_JS_SERVER_URL: https://golden-cluster.apm.us-west2.gcp.elastic-cloud.com:443
ELASTIC_APM_SECRET_TOKEN: PPAASSWWOORRDD
ELASTIC_APM_SERVER_URL: https://golden-cluster.apm.us-west2.gcp.elastic-cloud.com:443

FLEET_ELASTICSEARCH_HOST: https://golden-cluster.es.us-west2.gcp.elastic-cloud.com:443
FLEET_URL: https://golden-cluster.fleet.us-west2.gcp.elastic-cloud.com:443
```

## Destroy a cluster

The `oblt-cli` can destroy a cluster.
The cluster is identify by the name so you need the [cluster info](#save-cluster-info) to destroy it.

```bash
CLUSTER_NAME=$(jq -r .ClusterName cluster-info.json)
oblt-cli cluster destroy --cluster-name ${CLUSTER_NAME}
```

## Scripts

The following script shows how to use the `oblt-cli` in a CI.
The scripts use all the features described above.
The scripts will install and configure oblt-cli, create a cluster, wait for the cluster, run some API calls, destroy the cluster, and finally cleanup the environment.

```bash
{% include '/samples/stack-test.sh' ignore missing %}
```

[cluster templates]: /user-guide/cluster-templates/
[releases]: https://github.com/elastic/observability-test-environments/releases
[install.sh]: https://github.com/elastic/observability-test-environments/blob/main/tools/oblt-cli/scripts/install.sh
[deployments]: https://github.com/elastic/observability-test-environments/tree/main/deployments
[oblt-cli-cluster-credentials]: https://github.com/elastic/apm-pipeline-library/tree/main/.github/actions/oblt-cli-cluster-credentials
[oblt-cli-create-ccs]: https://github.com/elastic/oblt-actions/tree/v1/oblt-cli/cluster-create-ccs
[oblt-cli-create-custom]: https://github.com/elastic/oblt-actions/tree/v1/oblt-cli/cluster-create-custom
[oblt-cli-create-serverless]: https://github.com/elastic/oblt-actions/tree/v1/oblt-cli/cluster-create-serverless
[oblt-cli-destroy-cluster]: https://github.com/elastic/oblt-actions/tree/v1/oblt-cli/cluster-destroy
[setup-oblt-cli]: https://github.com/elastic/oblt-actions/tree/v1/oblt-cli/setup
[elastic/oblt-infra]: https://github.com/elastic/oblt-infra
[oblt-cli-buildkite-plugin]: https://github.com/elastic/oblt-cli-buildkite-plugin
