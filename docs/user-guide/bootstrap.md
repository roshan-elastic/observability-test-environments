---
render_macros: false
---
# Bootstrap settings using the API

Sometimes we need to initialize the Elastic Stack with a bunch of configuration that is not possible to make in the configuration files,
or it is better to use the API to configure them. For those cases oblt cluster added the support to bootstrap the Elastic Stack with REST API calls.

The bootstrap feature is based on a folder that contains a subfolder per each element to initialize,
on each folder you will have a set of recipes files that will be used to initialize the element.
The recipes are evaluated in alphabetical order, so you can use numbers to control the order of the calls.

```
bootstrap
├── apm
│   └── 00 - api call.yml
│   └── 01 - api call.yml
├── elasticsearch
│   └── 00 - api call.yml
│   └── 01 - api call.yml
├── kibana
│   └── 00 - api call.yml
│   └── 01 - api call.yml
└── fleet
│   └── 00 - api call.yml
│   └── 01 - api call.yml
```

The recipe is a yaml file that contains the following fields:

```yaml
---
# Description of what the recipe do.
description: Basic test request
# Path to the REST API call.
api: "/"
#HTTP method to use it can be GET, POST, PUT, DELETE, ...
method: GET
# HTTP header we need to pass. (Optional)
headers:
  Content-Type: application/json
  accept: application/json
  kbn-xsrf: true
# The payload to send to the REST API.
body: |
  {
    "test": "test"
  }
# The format of the payload. (Optional)
body_format: json
# The expected return codes.
return_code:
  - 200
  - 302
# If we want to ignore the errors. (Optional)
ignore_errors: true
# If we want to retry the call in case of error. (Optional)
retries: 3
# The delay between retries. (Optional)
delay: 10
```

`Oblt-cli` implement the bootstrap command, the first step is to list the available recipes:

```bash
oblt-cli bootstrap list
```

if you have your own folder with recipes you can use the `--bootstrap-folder` option to point to the folder.

```bash
oblt-cli bootstrap list --bootstrap-folder /tmp/bootstrap-test
```

To bootstrap a oblt cluster you can use the `cluster` command, this command will bootstrap all the elements in the cluster.

```bash
oblt-cli bootstrap cluster --cluster-name edge-oblt --bootstrap-folder /tmp/bootstrap-test
```

It is possible to bootstrap only one element of the cluster, for example, if you want to bootstrap only the elasticsearch cluster you can use the `--recipes-elasticsearch` and the list of recipes you want to use.

```bash
oblt-cli bootstrap cluster --cluster-name edge-oblt --bootstrap-folder /tmp/bootstrap-test --recipes-elasticsearch '["00-test", "01-test"]'
```

The same for Kibana and the rest of elements.

```bash
oblt-cli bootstrap cluster --cluster-name edge-oblt --bootstrap-folder /tmp/bootstrap-test --recipes-kibana '["00-test", "01-test"]'
```

Finally is possible to bootstrap Elasticsearch or Kibana instances directly.

```bash
oblt-cli bootstrap elasticsearch --bootstrap-folder /tmp/bootstrap-test --recipes '["test"]' --username elastic --password changeme --url https://edge-oblt.es.us-west2.gcp.elastic-cloud.com
```

```bash
oblt-cli bootstrap kibana --bootstrap-folder /tmp/bootstrap-test --recipes '["test"]' --username elastic --password changeme --url https://edge-oblt.kb.us-west2.gcp.elastic-cloud.com
```

## Templating in the recipes

The recipes support templating, you can use [Go templating][] in the body of the recipe.

```yaml
---
# Description of what the recipe do.
description: Basic test request
# Path to the REST API call.
api: "/"
#HTTP method to use it can be GET, POST, PUT, DELETE, ...
method: GET
# HTTP header we need to pass. (Optional)
headers:
  Content-Type: application/json
  kbn-xsrf: true
# The payload to send to the REST API.
# https://pkg.go.dev/text/template
body: |
  {
    "test": "{{ .TestVaue }}"
  }
# The format of the payload. (Optional)
body_format: json
# The expected return codes.
return_code:
  - 200
```

You see that the Body field has a `{{ .TestVaue }}` this is a template that will be replaced by the value of the `TestVaue` field in the recipe.
To pass the value to the recipe you can use the `--parameters` flag.

```shell
oblt-cli bootstrap cluster \
  --cluster-name edge-oblt \
  --bootstrap-folder /tmp/bootstrap-test \
  --parameters '{"TestVaue": "test"}'
```

Finally, you can use the [Spring functions][] to manipulate the values in the recipes.

```yaml
---
# Description of what the recipe do.
description: Basic test request
# Path to the REST API call.
api: "/"
#HTTP method to use it can be GET, POST, PUT, DELETE, ...
method: GET
# HTTP header we need to pass. (Optional)
headers:
  Content-Type: application/json
  kbn-xsrf: true
# The payload to send to the REST API.
# https://pkg.go.dev/text/template
# https://masterminds.github.io/sprig/
body: |
  {
    "test": "{{ .TestVaue }}",
    "date": "{{ now | date "2006-01-02" }}",
    "fromEnvironmentVariables": "{{ env "HOME" }}"
    "interpolateEnvironmentVariables": "{{ expandenv "interpolate $HOME" }}"
  }
# The format of the payload. (Optional)
body_format: json
# The expected return codes.
return_code:
  - 200
```

[Go templating]: https://pkg.go.dev/text/template
[Spring functions]: https://masterminds.github.io/sprig/
