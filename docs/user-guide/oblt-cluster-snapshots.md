# Oblt clusters snapshots

The oblt clusters are configured to store a snapshot in the ESS snapshot repository every 30 minutes,
this is the default configuration in ESS.
Additionally, the oblt clusters are configured to store a snapshot in a [GCS bucket][] every 24 hours.
The [GCS bucket][] is configured in the `external_gcs_repository` repository.
The [GCS bucket][] stores 5 snapshots, the oldest snapshot is deleted when a new snapshot is created.
Each cluster has its own folder, named as the cluster name.
Inside that folder, the snapshots are stored.
The developers oblt clusters are configured to attach the snapshots repositories of the production oblt clusters.
This allows the developers to restore the production oblt clusters snapshots in their own clusters.
Each cluster has its own snapshot repository, the name is like `external_gcs_repository_CLUSTER_NAME`.

## Restore a snapshot

To restore a snapshot in the developers oblt clusters, you will do it as on any ESS cluster using the Kibana UI.
As alternative, you can create a "Restore Snapshot" GitHub issue or use the Elasticsearch API.

Check the following links for more details:

* [restore a snapshot via GH issue][]
* [restore a snapshot][]

## Local Elasticsearch

In some cases, developers may want to restore a snapshot in their local Elasticsearch instance.

### From the Bucket

The first option is to configure the local Elasticsearch instance to use the [GCS bucket][] as snapshot repository.
To do that, you have to configure the [repository-gcs][] plugin in your Elasticsearch instance.

```bash
DOCKER_ES_ID=$(docker ps --filter name=elasticsearch --quiet )
ELASTICSEARCH_HOST=http://localhost:9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=changeme
GCP_ACCOUNT_NAME=user-account@elastic.co

# install the repository-gcs plugin
docker exec -it ${DOCKER_ES_ID} bin/elasticsearch-plugin install https://artifacts.elastic.co/downloads/elasticsearch-plugins/repository-gcs/repository-gcs-8.0.0-alpha2.zip

# add GCP credentials to the keystore
docker cp  ${HOME}/.config/gcloud/legacy_credentials/${GCP_ACCOUNT_NAME}/adc.json ${DOCKER_ES_ID}:/tmp
docker exec -it ${DOCKER_ES_ID} bin/elasticsearch-keystore add-file gcs.client.external.credentials_file /tmp/adc.json
docker exec -it ${DOCKER_ES_ID} bin/elasticsearch-keystore list
docker stop ${DOCKER_ES_ID}
docker start ${DOCKER_ES_ID}

# create the snapshot repository
curl -X PUT \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/external_gcs_repository" \
  -d '{
    "type": "gcs",
    "settings":{
      "bucket": "oblt-clusters",
      "client": "external",
      "base_path": "edge-oblt"
      }
    }'

# list the snapshots
curl -X GET \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/external_gcs_repository/_all?pretty"
```

### From the local filesystem

The second option is to restore the snapshot from the local filesystem,
To do that, you will need to download the snapshot from the [GCS bucket][] and restore it in your local Elasticsearch instance.

The following example shows how to download the snapshot from the [GCS bucket][] and configure the local Elasticsearch instance to use it.

```bash
CLUSTER_NAME=oblt-cluster
DOCKER_ES_ID=$(docker ps --filter name=elasticsearch --quiet )
ELASTICSEARCH_HOST=http://localhost:9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=changeme
LOCAL_ES_SNAPSHOT_DIR=${HOME}/Downloads/${CLUSTER_NAME}

# download the snapshot repository
gcloud storage cp --recursive "${CLUSTER_NAME}" "gs://oblt-clusters" "${LOCAL_ES_SNAPSHOT_DIR}"

# create the snapshot repository
curl -X PUT \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/local_repository" \
  -d "{\"type\": \"url\",\"settings\":{\"url\": \"${LOCAL_ES_SNAPSHOT_DIR}\"}}"

# list the snapshots
curl -X GET \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/local_repository/_all?pretty"
```

### Mount or restore a snapshot

Once you have the snapshot ID, you can mount or restore it in your local Elasticsearch instance.
The easy way it is to use the Kibana UI to mount or restore a snapshot.
Another option is to use the Elasticsearch API.

The following example mounts the snapshot in the local Elasticsearch instance.

```bash
# Mount a snapshot
SNAPSHOT_ID=external-snapshots-2021.10.14-c33zfqvattm0e5icrpwb-w
curl -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/external_gcs_repository/${SNAPSHOT_ID}/_mount?wait_for_completion=true" \
  -d '
    {
      "index": "apm-7.15.1-transaction-000001",
      "renamed_index": "apm-7.15.1-transaction-000001-snap",
      "index_settings": {
        "index.number_of_replicas": 0
      },
      "ignored_index_settings": [ "index.refresh_interval" ]
    }
    '
```

The following example restores the snapshot in the local Elasticsearch instance.

```bash
# Restore a snapshot
SNAPSHOT_ID=external-snapshots-2021.10.14-c33zfqvattm0e5icrpwb-w
curl -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_snapshot/external_gcs_repository/${SNAPSHOT_ID}/_restore?wait_for_completion=true" \
  -d '
    {
      "indices": "apm-7.15.1-transaction-000001",
      "ignore_unavailable": true,
      "include_global_state": false,
      "rename_pattern": "apm-7.15.1-transaction-000001",
      "rename_replacement": "apm-7.15.1-transaction-000001-snap",
      "partial": true,
      "index_settings": {
        "index.number_of_replicas": 0,
        "auto_expand_replicas": false,
        "index.routing.allocation.include._tier_preference" : "data_hot"
      }
    }
    '
```

## Optimize the restore time

The restore time depends on the snapshot size and the cluster resources.
It could take a few minutes or hours to restore a snapshot.
To optimize the restore time, you can tune the cluster settings.

The following example shows how to tune the cluster settings to optimize the restore time.
The example increases the limits of resources used to process the restore operation.

```bash
curl -X PUT \
  -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}" \
  -H "Content-Type: application/json" \
  "${ELASTICSEARCH_HOST}/_cluster/settings?pretty" \
  -d '{
        "persistent" : {
          "cluster.routing.allocation.node_concurrent_recoveries" : "20",
          "indices.recovery.max_concurrent_file_chunks" : "8",
          "indices.recovery.max_concurrent_operations" : "4",
          "indices.recovery.max_bytes_per_sec" : "-1"
        }
      }'
```

The above step is not needed if you choose the "Restore Snapshot" GitHub issue automation, it's already taken care of.

Check the following links for more details:

* [restore a snapshot via GH issue][]
* [restore a snapshot][]
* [repository-gcs][]
* [Download GCS objects][]
* [Read-only URL repository][]

[GCS bucket]: https://console.cloud.google.com/storage/browser/oblt-clusters
[restore a snapshot via GH issue]: https://github.com/elastic/observability-test-environments/issues/new?template=cluster-restore-snapshot-issue.yaml&title=%5BRestore+Snapshot%5D%3A+&labels=cluster%2Ccluster-restore-snapshot
[restore a snapshot]: https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-restore-snapshot.html
[repository-gcs]: https://www.elastic.co/guide/en/elasticsearch/reference/current/repository-gcs.html
[Download GCS objects]: https://cloud.google.com/storage/docs/downloading-objects
[Read-only URL repository]: https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-read-only-repository.html
