# Access Managed Kubernetes Infrastructure (MKI)

## Overview

The `Managed Kubernetes Infrastructure`, aka [MKI][], is a Kubernetes cluster that is managed by Elastic.
The [MKI][] is used to deploy the Elastic [Serverless][] solution.
To access to [MKI][] there are some [requirements](https://docs.elastic.dev/mki/operations/setup) that need to be met, and some process to follow to [connect to the environments](https://docs.elastic.dev/mki/user-guide/cluster-access).

### Pre-requisites

To access the [MKI][] you need to have the following:

* [Access to the Elastic VPN](https://cloud.elastic.dev/sre/how-to/accessing-ess-production-service/)
* Access to the Elastic-observability GCP project (ask in the [#observablt-robots][] Slack channel)
* [Access to the MKI environments](https://cloud.elastic.dev/teleport/USAGE/#more-information)

{% include '/user-guide/mki-login.md' ignore missing %}

{% include '/user-guide/mki-gen-kibana-yaml.md' ignore missing %}

{% include '/user-guide/mki-run-kibana.md' ignore missing %}

[MKI]: https://docs.elastic.dev/mki
[Serverless]: https://docs.elastic.dev/serverless
[#observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
