### Local Development

It is possible to connect a local development environment to the serverless cluster.
The goal is to test the changes in Kibana against the serverless cluster.

!!! Warning

    You need access to the MKI environment to run the following commands.
    Check [MKI requirements][] for more details.

!!! Warning

    Stop the VPN connection before running the command, otherwise the process will fail.
    Aviatrix VPN Client and oblt-cli uses the same port to to receive the credentials, so you need to stop the VPN connection before running the command.

To connect a local development environment to the serverless cluster, run the following command:

```shell
CLUSTER_NAME=serverless-tupxw
KIBANA_SRC_DIR="${HOME}/src/kibana"
KIBANA_GENERATED_YAML_OUTPUT="${PWD}/kibana_serverless.yaml"
oblt-cli services mki run-local-kibana \
  --cluster-name "${CLUSTER_NAME}" \
    --kibana-src "${KIBANA_SRC_DIR}" \
  --kibana-yaml-path "${KIBANA_GENERATED_YAML_OUTPUT}"
```

If you need to edit the Kibana YAML file, you can use the following command:

```shell
CLUSTER_NAME=serverless-tupxw
KIBANA_SRC_DIR="${HOME}/src/kibana"
KIBANA_GENERATED_YAML_OUTPUT="${PWD}/kibana_serverless.yaml"
oblt-cli services mki generate-kibana-yaml \
  --cluster-name "${CLUSTER_NAME}" \
  --kibana-yaml-path "${KIBANA_GENERATED_YAML_OUTPUT}"

# After editing the Kibana YAML file, you can start Kibana with the following command:
yarn start --serverless=oblt --config="${KIBANA_GENERATED_YAML_OUTPUT}"
```

The process will connect to the VPN, login in the [MKI][] environment, and start the Kibana configured to use the serverless cluster.

[MKI]: https://docs.elastic.dev/mki
[MKI requirements]: ./mki.md
