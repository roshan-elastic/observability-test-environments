
# Enable SSO

The ESS deployment can use the Okta authentication,
the oblt clusters are preconfigured to use Okta,
so you only have to create the cluster and register the cluster URL in the [oktanaut app][].

* Create a new cluster
* register the cluster URL at [oktanaut app][]
* Update your cluster configuration file with the `idp.entity_id` returned by [oktanaut app][]

```yaml
elastic_cloud:
  provider: gcp
  region: gcp-us-west2
  endpoint: https://cloud.elastic.co
  zones: 3
  template: gcp-io-optimized-v2
  autoscale: true
  # idp.entity_id returned by oktanaut app
  saml_id: 01234567890123456789
```

Check [how to create a cluster](../user-guide/use-case-create-cluster.md) and [how to update a cluster](../user-guide/use-case-update-cluster.md) for more details.
To configure the Okta authentication you have to follow the [Oktanaut instructions][]
The configuration settings are in the file [es_settings.yml.j2](https://github.com/elastic/observability-test-environments/blob/main/ansible/ansible_collections/oblt/framework/roles/stack_ess/templates/terraform/elastic_cloud/es_settings.yml.j2).

[oktanaut app]: https://oktanaut.app.elastic.dev/
[Oktanaut instructions]: https://github.com/elastic/infra/blob/master/flavortown/oktanaut/README.md
