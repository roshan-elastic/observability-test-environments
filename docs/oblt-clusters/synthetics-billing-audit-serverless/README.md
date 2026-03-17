# Overview

`synthetics-billing-audit-serverless` is a long term ESS cluster deployed in the Cloud First test region in production.
The Serverless deployment has Elasticsearch, Kibana, and APM managed service.
It's part of our SOX compliance, and is to ensure accuracy of billing data.

## Access the cluster

The serverless project have only the ESS SSO enabled, and it is not possible to configure a different Realm for authentication.
Because of that the only way to access to the deployed cluster is using the ESS SSO.
To access to any cluster deployed with `oblt-cli` you must be a member of the robots organization, you can self-enroll to the robots organization either typing `/oblt-onboarding-cloud` in the [observablt-robots][] Slack channel or filling [this GitHub issue][].

If for some reason you are not able to change of organization, we have a [shared user][] to access to QA and staging environments.

!!! Note

    The shared user is only for internal use. We do not have a shared user in production organization.

## URLs and Credentials

### oblt-robot

To retrieve the credentials of the cluster you can use the oblt-robot Slack bot.
It is installed in the slack channels: [observablt-robots][], #observablt-bots, and a few others.
To ask the bot for the credentials you should write the Slack command `/cluster-secret`
and choose the cluster you want to obtain the credentials.

### oblt-cli

To obtain the credentials of the cluster you can use the following command.

```bash
oblt-cli cluster secrets credentials --cluster-name serverless-production-oblt
```

## Update

The update of the clusters is automatic when the [Production environment][] push a new version of the stack.
The version to update is the latest stable version of the latest on the [Production environment][].

[Production environment]: https://docs.elastic.dev/serverless/production
[shared user]: https://p.elstc.co/paste/7O4qBEtr#4DoUgqXVWBGZODN0SiZKDO4RoMvoISOW+xkGP6shk7s
[observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
[this GitHub issue]: https://ela.st/self-service-oblt-robots-org
