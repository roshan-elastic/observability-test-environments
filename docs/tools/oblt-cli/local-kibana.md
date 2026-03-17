# Local Kibana instance

When the process to create the CCS cluster ends a Slack message is sent
with the content of a kibana.yml file ready to use.
If you copy the content of the Slack message into a `config/kibana.dev.yml` file.
There is a another way to get the kibana.yml file, you can use the `oblt-cli` tool to get the kibana.yml file from the Vault.

```shell
oblt-cli cluster secrets kibana-config --cluster-name=my-cluster --output-file ${PWD}/kibana.yml
```

you have to check that there is no configuration settings which have been duplicated or are obsolete.
If you find obsolete settings please report them by creating an issue in [oblt test environments repo](https://github.com/elastic/observability-test-environments/issues).

After save your `kibana.yml` file you can start your local Kibana instance from code by using the following commands:

* if your source code is new

```shell
NODE_OPTIONS=" --max-old-space-size=4096"
FORCE_COLOR=1
nvm install $(cat .node-version)
yarn kbn bootstrap
yarn start
```

* If you had already ran the `yarn kbn bootstrap`

```shell
NODE_OPTIONS=" --max-old-space-size=4096" yarn start
```

* If you do not need to run Kibana from code you can use the Kibana Docker image
that matches with the Elastic Stack version in the oblt cluster.

```shell
docker run -it -v ${HOME}/kibana.yml:/usr/share/kibana/config/kibana.yml -p 5601:5601 docker.elastic.co/kibana/kibana:8.0.0-SNAPSHOT
```

## Post-configuration

There are some settings that is not possible to set in the `kibana.yml` file,
for those cases we need to make manual operations to configure those settings (automatic way is WIP).

### Logs and Metrics

Logs and Metrics need a couple of settings to use CCS:
* Go to the `Inventory` page under `Infrastructure` and click on `Settings` in the upper right hand corner of the screen.
Change the `Metric indices` to the the value of:

```yaml
remote_cluster:metricbeat-*,metricbeat-*,remote_cluster:metrics-*,metrics-*
```

* Go to the `Stream` page under `Logs` and click on `Settings` in the upper right hand corner of the screen.
Change the `Log indices` to the the value of:

```yaml
remote_cluster:filebeat-*,remote_cluster:logs-*,remote_cluster:kibana_sample_data_logs*,filebeat-*,kibana_sample_data_logs*,logs-*
```

It is possible to use oblt-cli bootstrap feature to configure the settings for Logs and Metrics:

```bash
ADMIN_PASSWORD=changeme
oblt-cli bootstrap kibana \
  --verbose \
  --bootstrap-folder "${HOME}/.oblt-cli/observability-test-environments/deployments/bootstrap-ccs" \
  --username admin \
  --password ${ADMIN_PASSWORD} \
  --url http://localhost:5601
```

!!! Note

      When you use the template `ess-overview` there are four remote clusters configured.
      You will have to replace `remote_cluster:` with `metrics-*:` in the metrics settings and `logs-*:` in the logs settings.

If you are using the `ess-overview` template you can use the following command to configure the settings:

```bash
ADMIN_PASSWORD=changeme
oblt-cli bootstrap kibana \
  --verbose \
  --bootstrap-folder "${HOME}/.oblt-cli/observability-test-environments/deployments/bootstrap-ess-overview" \
  --username admin \
  --password ${ADMIN_PASSWORD} \
  --url http://localhost:5601
```

### Security Solution

Security Solution needs some additional settings to use CCS:

In Kibana, navigate to `Stack Management` and in the Kibana section select `Advanced Settings`.  If you select **Security Solution** in the `Category` drop down you will be able to quickly find the entries that need to be updated.

Duplicate the entries in the Elasticsearch indices field, `securitySolution:defaultIndex`, and prefix each of them with `remote_cluster:`.

Repeat the steps with the Threat indices field, `securitySolution:defaultThreatIndex`.  Duplicate the entries and prefix each of them with `remote_cluster:`.

Select `Data Views` in the Kibana section of Stack Management and create data views for the following index patterns.

 - remote_cluster:auditbeat-*
 - remote_cluster:filebeat-*
 - remote_cluster:logs-*
 - remote_cluster:metrics-*
 - remote_cluster:packetbeat-*
 - remote_cluster:winlogbeat-*

It is possible to use oblt-cli bootstrap feature to configure the settings for security solution:

```bash
CLUSTER_NAME=edge-oblt-ccs-ohcnl
ADMIN_PASSWORD=changeme
oblt-cli bootstrap kibana \
  --verbose \
  --bootstrap-folder "${HOME}/.oblt-cli/observability-test-environments/deployments/bootstrap-ccs" \
  --username admin \
  --password ${ADMIN_PASSWORD} \
  --url http://localhost:5601
```

For additional details please refer to the Security Solution tool documentation https://github.com/elastic/security-team/blob/main/tools/oblt-ccs-configuration.md.

## Serverless

We can also use [serverless](https://docs.elastic.dev/serverless) deployments to connect our development local Kibana instance.
To do that we need to create a new serverless deployment with `oblt-cli` tool,
then generate a `config/kibana.serverless.yml` file with the proper values at `elasticsearch.hosts`, `elasticsearch.serviceAccountToken`, `xpack.encryptedSavedObjects.encryptionKey`, `xpack.security.encryptionKey`, and `xpack.reporting.encryptionKey` settings from the credentials of the cluster.

!!! Warning

    We will run two Kibana instances against the same Elasticsearch cluster,
    so we need to have special care to not create objects in both instances.
    Create objects in both can cause versions conflicts and other issues.
    In case or issues, deleting the `.kibana*` indices in the Elasticsearch cluster can help.

* [Create a serverless cluster](/user-guide/serverless.md)
* [Get credentials cluster](/docs/user-guide/use-case-retrieve-credentials.md)
* Use the `oblt-cli` tool to generate the kibana.yml file by connecting to the MKI k8s cluster to retrieve some settings from the Kibana configuration.

```shell
KIBANA_SRC_FOLDER="${HOME}/src/kibana"
KIBANA_YML=${KIBANA_SRC_FOLDER}/config/kibana.serverless.yaml
CLUSTER_NAME=serverless-my-cluster
ENVIRONMENT=staging

oblt-cli services mki generate-kibana-yaml \
  --cluster-name "${CLUSTER_NAME}" \
  --environment "${ENVIRONMENT}" \
  --kibana-yaml-path "${KIBANA_YML}"
```

```yaml
elasticsearch.hosts: https://serverless-my-cluster-fffe89.es.us-east-1.aws.staging.elastic.cloud
elasticsearch.serviceAccountToken: ABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABABA
elasticsearch.ssl.verificationMode: none
elasticsearch.ignoreVersionMismatch: true
migrations:
  skip: true
server.host: 0.0.0.0
xpack:
  reporting:
    encryptionKey: FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
  security:
    authc:
      providers:
        basic:
          cloud-basic:
            order: 150
    encryptionKey: CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC
    loginAssistanceMessage: 'Credentials: elastic / MyEsPaSsWoRd'
  encryptedSavedObjects:
    encryptionKey: DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD
```

* Start your local Kibana instance

```shell
KIBANA_SRC_FOLDER="${HOME}/src/kibana"
KIBANA_YML=${KIBANA_SRC_FOLDER}/config/kibana.serverless.yaml
NODE_OPTIONS=" --max-old-space-size=4096"

yarn start --serverless=oblt --config="${KIBANA_YML}"
```
