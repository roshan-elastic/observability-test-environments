### Access the cluster

The serverless project has only the ESS SSO enabled, and it is impossible to configure a different Realm for authentication.
Because of that, the only way to access the deployed cluster is to use the ESS SSO.
We provide two options to access the cluster:

* [Join robots organization](#join-robots-organization).
* Use a [shared user][] to access the QA and staging environments if you cannot change organizations.

!!! Note

    The shared user is only for internal use. We do not have a shared user for the production organization.

!!! Warning

    The Robots organization is full automated with [oblt-cli](https://ela.st/oblt-cli) [oblt-robot](https://ela.st/oblt-robot)
    the manual deployment will be removed.

#### Join robots organization

All the elastic employees can join the robots organization to access the serverless projects.
Every Elastician has an `@elastic.co` account created in the Elastic Cloud Service (ESS) and can join the robots organization. Additionally, Elasticians have a `@elasticsearch.com` email alias that can be used to join the robots organization.

The `@elasticsearch.com` alias does not have an ESS account created by default,
If you use your `@elasticsearch.com` alias, you need to create an account in ESS using the `@elasticsearch.com` alias.
The `@elastisearch.com` alias allows to create email aliases like `my-user+alias@elasticsearch.com` those aliases can be used to create multiple accounts in ESS.

You can use the [sign up](https://cloud.elastic.co/registration) process to create a new user in ESS.
There are three environments where the Robots ESS organization is deployed, each environment is independent.
We send one invitation per environment, you have to accept those you want access to.
you need a user for each environment you want to access:

* [Production](https://cloud.elastic.co/registration)
* [QA](https://console.qa.cld.elstc.co/registration)
* [Staging](https://staging.found.no/registration)

!!! Warning

    ESS user can not be member of more than one organization, so please make sure that the user is not a member of any other organization.

To join the robots organization, you need an invitation and accept it.
You can request an invitation by typing `/oblt-onboarding-cloud` in the [observablt-robots][] Slack channel or filling [this GitHub issue][]. This triggers an automatic process that will send you an invitation for every environment we have (production, QA, staging) to join the robots organization.
Once you accept the invitation, you will have access to the serverless projects.

[shared user]: https://p.elstc.co/paste/xR1ZSG46#tzMtAKrVvR+KsEsF4w8E7-aOAdBxpQxDDGV/pNk3v6+
[this GitHub issue]: https://ela.st/self-service-oblt-robots-org
[observablt-robots]: https://elastic.slack.com/archives/CJMURHEHX
