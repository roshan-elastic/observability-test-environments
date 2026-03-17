# Overview

The aim of serverless-production-sec is to evaluate Security features
in a serverless environment. The project is created in the [Production environment][]
and therefore its version is kept up-to-date with the latest deployed there.

CCS is not supported *it is not possible to attach developers clusters* if you want a developer
cluster to work on serverless you have to use the serverless template to deploy one see
[Serverless deployments](https://elastic.github.io/observability-test-environments/user-guide/serverless/)

## Access the cluster

The serverless project has only the ESS SSO enabled, and it is impossible to configure a different Realm for authentication.
Because of that, the only way to access the deployed cluster is to use the ESS SSO.
To access any cluster deployed with `oblt-cli` you must be a member of the robots organization, you can self-enroll to the robots organization either typing `/oblt-onboarding-cloud` in the [observablt-robots][] Slack channel or filling [this GitHub issue][].

We have a [shared user][] to access the QA and staging environments if you cannot change organizations.

!!! Note

    The shared user is only for internal use. We do not have a shared user for the production organization.

## URLs and Credentials

### oblt-robot

To retrieve the credentials of the cluster you can use the oblt-robot Slack bot.
It is installed in the slack channels: [observablt-robots][], #observablt-bots, and a few others.
To ask the bot for the credentials you should write the Slack command `/cluster-secret`
and choose the cluster you want to obtain the credentials.

### oblt-cli

To obtain the credentials of the cluster you can use the following command.

```bash
oblt-cli cluster secrets credentials --cluster-name serverless-production-sec
```

[Production environment]: https://docs.elastic.dev/serverless/production
[shared user]: https://p.elstc.co/paste/7O4qBEtr#4DoUgqXVWBGZODN0SiZKDO4RoMvoISOW+xkGP6shk7s
[observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
[this GitHub issue]: https://ela.st/self-service-oblt-robots-org
