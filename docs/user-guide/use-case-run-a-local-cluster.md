# Run a local cluster

You can run a Elastic Stack cluster locally using the oblt test environments,
this configuration will deploy a [Kind k8s cluster](https://kind.sigs.k8s.io/),
then, inside this k8s cluster you can deploy the Elastic Stack with ECK or using Helm Charts.
Both options will deploy Elasticsearch, Kibana, and APM.

Check the [How to create a cluster locally](./use-case-create-cluster.md#how-to-create-a-cluster-locally) before start.

## Run a local ECK cluster

This deploy uses ECK to deploy the stack into the k8s cluster.
To evaluate this approach for your needs, do the following:

* Checkout the oblt test environment repo
`git clone git@github.com:elastic/observability-test-environments.git`
* Start some of the pre configured test environment we have
`CLUSTER_CONFIG_FILE=$(pwd)/tests/environments/eck-kind.yml make -C ansible create-cluster`
* To operate with the k8s cluster you have to go to the folder `ansible/build`
and change the HOME and PATH env vars with the following command
`export HOME=$(pwd) && export PATH=${HOME}/bin:${PATH}`
then you can make `kubectl get po -A && kubectl get ingress`
* Get the elastic user password with the command
`kubectl get secret elastic-stack-es-elastic-user -o=jsonpath='{.data.elastic}'`
* Connect to the local Elasticsearch and Kibana instances https://kibana.127.0.0.1.ip.es.io and https://elasticsearch.127.0.0.1.ip.es.io

When you end your work you can destroy the environment by using the destroy command

`CLUSTER_CONFIG_FILE=$(pwd)/tests/environments/eck-kind.yml make -C ansible create-cluster`
