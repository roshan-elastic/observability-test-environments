# Contributing

## How to test the changes in the cluster
Please follow these instructions:

- Edit the `deploy-config.yaml` to select the Elastic Stack deploy flavor
and enable/disable the services that you need. This is the default configuration:

```YAML
#Select where the Elastic Stack (Elasticsearch, Kibana, and APM server) would be deployed,
#on Elastic Cloud or in a k8s cluster. ['none','Elastic Cloud', 'k8s']
elastic_stack_flavor: k8s
#Destroy the infrastructure created.
destroy: true
#Deploy Helm service into the k8s cluster.
deploy_helm: true
#Deploy ingress service into the k8s cluster.
deploy_ingress: true
#Deploy Elasticsearch service into the k8s cluster.
deploy_elasticsearch: true
#Deploy Kibana service into the k8s cluster.
deploy_kibana: true
#Deploy Betas services into the k8s cluster.
deploy_beats: true
#Deploy Curator service into the k8s cluster.
deploy_curator: true
#Deploy APM service into the k8s cluster.
deploy_apm: true
#Deploy Apache sample app into the k8s cluster.
deploy_apache: true
#Deploy PetClinic sample app into the k8s cluster.
deploy_petclinic: false
#Deploy Opbeans sample app into the k8s cluster.
deploy_opbeans: true
#Deploy Grafana and Prometheus into the k8s cluster.
deploy_grafana_prometheus: true
#Deploy AWS infrastructure.
deploy_aws: true
```
`elastic_stack_flavor: k8s` will enable the creation of the cluster on a k8s cluster.

- Create a branch in the `elastic/observability-test-environments` repository with your changes following the syntax `oblt-cli/<something>`. The cluster will be created.
- Once you are done with testing, you can create revert the changes as long as you follow the above-mentioned syntax.


## How to lint the changes locally

### Precommit

This particular process will help to evaluate some linting before committing any changes. Therefore you need the pre-commit.

#### Installation

Follow <https://pre-commit.com/#install> and `pre-commit install`

Some hooks might require some extra tools such as:

* [shellcheck](https://github.com/koalaman/shellcheck#installing)
* [yamllint](https://yamllint.readthedocs.io/en/stable/quickstart.html)

#### Enabled hooks

See the `.pre-commit-config.yaml` file.
