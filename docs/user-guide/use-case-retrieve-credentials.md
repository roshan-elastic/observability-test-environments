# Retrieve credentials

The credentials of all oblt clusters are stores in our Vault.
To retrieve the credentials you can user [oblt-cli][], [oblt-robot][], or [vault-access][Vault CLI].

## oblt-cli

Oblt-cli implement a command to retrieve credentials and configuration template files.
Use the following command to retrieve the credentials of a cluster

```bash
CLUSTER_NAME=edge-oblt
oblt-cli cluster secrets credentials --cluster-name ${CLUSTER_NAME}
```

```bash
CLUSTER_NAME=edge-oblt
oblt-cli cluster secrets credentials --cluster-name ${CLUSTER_NAME} --output-file "${CLUSTER_NAME}-credentials.md"
```

for more secrets check the help `oblt-cli cluster secrets --help`

## oblt-robot

The oblt-robot is available in our Slack channel [observablt-robots][].
You can request the secrets of a cluster by using the command `/cluster-secret`

## Google Cloud Secret Manager

Check that you have installed and configured [gcloud][] or access [Google Cloud Secret Manager][] UI.
You must have access to the `elastic-observability` GCP project.
then you can list the secrets in [Google Cloud Secret Manager][] by executing:

```bash
CLUSTER_NAME=edge-oblt
gcloud secrets list --filter="name:oblt-clusters_${CLUSTER_NAME}_"
```

To retrieve the credentials and useful links use the following command

```bash
CLUSTER_NAME=edge-oblt
gcloud secrets versions access latest --secret="oblt-clusters_${CLUSTER_NAME}_credentials"
```

[oblt-cli]: http://ela.st/oblt-cli
[oblt-robot]: https://ela.st/oblt-robot
[observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
[gcloud]: https://cloud.google.com/sdk/gcloud
[Google Cloud Secret Manager]: https://console.cloud.google.com/security/secret-manager?project=elastic-observability&pli=1
