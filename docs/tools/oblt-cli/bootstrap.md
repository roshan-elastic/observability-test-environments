
# Bootstrap Elasticsearch and Kibana

There are many configuration which we share across developers. We have unified some of those shared configuration in `oblt-cli`
to make it possible to run some recipes to initialize and configure settings in Elasticsearch and Kibana.

## List bootstrap recipes

This command will list all the recipes available for Elasticsearch and Kibana.

```bash
oblt-cli bootstrap list
```

## Apply recipes to a cluster

The easiest way to bootstrap an oblt cluster is using this command, you only have to pass the name of the cluster to apply all the recipes.

```bash
oblt-cli bootstrap cluster --cluster-name edge-oblt
```

In the case you only want to run some recipes it is possible to choose which recipes to apply to Elasticsearch and Kibana

```bash
oblt-cli bootstrap cluster --cluster-name edge-oblt --recipes-elasticsearch '["test"]'
oblt-cli bootstrap cluster --cluster-name edge-oblt --recipes-kibana '["test"]'
```

## Apply recipes to a Elasticsearch instance

In some cases we want to initialize an Elasticsearch that is not part of a oblt cluster, for this task we have to run  `oblt-cli bootstrap elasticsearch` command.
You need to pass the Elastisearch URL and a method to authenticate.

```bash
oblt-cli bootstrap elasticsearch --url https://es.example.com --username elastic --password changeme
```

In the case you only want to run some recipes it is possible to choose which recipes to apply to Elasticsearch

```bash
oblt-cli bootstrap elasticsearch --url https://es.example.com --username elastic --password changeme --recipes '["test"]'
```

## Apply recipes to a Kibana instance

In some cases we want to initialize a Kibana that is not part of a oblt cluster, for this task we have to run  `oblt-cli bootstrap kibana` command.
You need to pass the Kibana URL and a method to authenticate.

```bash
oblt-cli bootstrap kibana --url https://kibana.example.com --username elastic --password changeme
```

In the case you only want to run some recipes it is possible to choose which recipes to apply to Elasticsearch

```bash
oblt-cli bootstrap kibana --url https://kibana.example.com --username elastic --password changeme --recipes '["test"]'
```

## Add your own recipe

The recipes are simple YAML files to describe a REST API call.
The recipes are in the [observability-test-environments repo](https://github.com/elastic/observability-test-environments) in the folder `deployments/bootstrap`.
They are two types of recipes, for Elasticsearch at `deployments/bootstrap/elasticsearch`, and for Kibana at `deployments/bootstrap/kibana`. The structure of a recipe file is the following:

```YAML
---
# Description of what the recipe does.
description: Basic test request
# Path to the REST API call.
api: "/"
#HTTP method to use it can be GET, POST, PUT, DELETE, ...
method: GET
# HTTP header we need to pass.
headers:
  Content-Type: application/json
  kbn-xsrf: true
# The payload to send to the REST API.
body: ""
# The expected return code.
return_code: 200
```

A new file in those folders is a new recipe available to use by everyone.
