# Push Licenses

Every instance of Elastic Stack deployed locally or remotely has a trial license that is valid for 30 days.
On ESS the SNAPSHOT deployments uses also a trial license that is valid for 30 days too.
In both cases, you can push a license to an instance of the Elastic Stack using the [license API][] or the ESS Console/API.
You can take a valid license from the [Internal License][] page in the wiki.

To make this process easy oblt-cli provides a command to push a license to an instance of the Elastic Stack.

```bash
ES_URL=http://localhost:9200
ES_USERNAME=elastic
ES_PASSWORD=changeme
oblt-cli license --url ${ES_URL} --username ${ES_USERNAME} --password ${ES_PASSWORD} --type dev
```

For the clusters deployed with oblt-cli there is also a command that only require the cluster name to push a license.

```bash
CLUSTER_NAME=my-cluster
oblt-cli cluster license --cluster-name ${CLUSTER_NAME} --type dev
```

The license type can be `release`, `dev`, `orchestration`, `orchestration-dev`.
This licenses are stored in the Google Cloud Secrets manager and updated every year:

* elastic-stack-license-dev
* elastic-stack-license-orchestration
* elastic-stack-license-orchestration-dev
* elastic-stack-license-release

The following command can be used to get the latest version of orchestration-dev license:

```bash
gcloud secrets versions access latest --secret=elastic-stack-license-orchestration-dev
```

[Internal License]: https://elasticco.atlassian.net/wiki/spaces/PM/pages/46802910/Internal+License+-+X-Pack+and+Endgame
[License API]: https://www.elastic.co/guide/en/elasticsearch/reference/current/get-license.html
