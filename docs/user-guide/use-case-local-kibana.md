# Connecting a local Kibana to a oblt cluster using oblt-cli

It is no longer possible to connect several Kibana instances to the same Elasticsearch,
this affect the way we connect our local Kibana instances to our oblt clusters.
Now, we have to use a proxy Elasticsearch instance configured with Cross Cluster Search (CCS) to access the oblt cluster.
Every developer will deploy an ephemeral CCS cluster per environment that the developer want to connect to.
To make this easy, we have develop a CLI tool that generates the configuration of the cluster,
then push the changes to a repository, and finally a pipeline is trigger to create the cluster and send the credentials to Slack.
The only requirement is to have the tool and Git configured to access to Elastic GitHub org.

{% include 'tools/oblt-cli/local-kibana.md' ignore missing %}
