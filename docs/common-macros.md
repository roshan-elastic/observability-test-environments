
{% macro credentials_url(cluster_name) -%}https://ela.st/{{ cluster_name }}-credentials{% endmacro -%}

{% macro cluster_readme_url(cluster_name) -%}https://elastic.github.io/observability-test-environments/{{ cluster_name }}/cluster-info/{% endmacro -%}

{% macro es_url(cluster_name) -%}https://{{ cluster_name }}.es.us-west2.gcp.elastic-cloud.com:443{% endmacro -%}

{% macro apm_url(cluster_name) -%}https://{{ cluster_name }}.apm.us-west2.gcp.elastic-cloud.com:443{% endmacro -%}

{% macro kibana_url(cluster_name) -%}https://{{ cluster_name }}.kb.us-west2.gcp.elastic-cloud.com:443{% endmacro -%}

{% macro docs_site_url() -%}https://elastic.github.io/observability-test-environments{% endmacro -%}

{% macro cluster_urls(cluster_name) -%}
You can check the URLs, user and passwords needed to access any of services in the following URLs

* [Credentials and URLs]({{ credentials_url(cluster_name) }})
* [Elasticsearch]({{ es_url(cluster_name) }})
* [APM]({{ apm_url(cluster_name) }})
* [Kibana]({{ kibana_url(cluster_name) }})
{% endmacro -%}

{% macro common_links() -%}
[artifactory]: https://artifacts-api.elastic.co/v1/versions/
[docs-site]: {{ docs_site_url() }}
[GKE]: https://cloud.google.com/kubernetes-engine
[Kubernetes]: https://kubernetes.io
[monitoring-oblt]: /Monitoring-oblt/
{% endmacro -%}

{% macro service_map(cluster_name) -%}
![{{ cluster_name }}-observability-overview](/images/{{ cluster_name }}-observability-overview.png){: style="width:600px"}

![{{ cluster_name }}-service-map](/images/{{ cluster_name }}-service-map.png){: style="width:600px"}

![{{ cluster_name }}-apm-services](/images/{{ cluster_name }}-apm-services.png){: style="width:600px"}
{% endmacro -%}
