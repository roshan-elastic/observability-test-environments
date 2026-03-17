### Project types and environments

The default project type is `observability` so you can create a serverless cluster with the default project type using:

```bash
    oblt-cli cluster serverless
```

You can also create a serverless cluster with a specific project type using:

```bash
    oblt-cli cluster create serverless --project-type "observability"
```

```bash
    oblt-cli cluster create serverless --project-type "elasticsearch"
```

```bash
    oblt-cli cluster create serverless --project-type "security"
```

Finally, to select the target environment:

```bash
    oblt-cli cluster create serverless --project-type "observability" --target "qa"
```

```bash
    oblt-cli cluster create serverless --project-type "observability" --target "staging"
```

```bash
    oblt-cli cluster create serverless --project-type "observability" --target "production"
```

!!! Note

    Run `oblt-cli cluster create serverless --help` to see all the available options.

!!! Warning

    The Staging and QA environments are on development mode so they can be unstable.

### Custom Docker images

In some cases, you may want to use custom Docker images for your serverless cluster. You can do this by specifying the Docker image for each component in the `oblt-cli` command.

Using a custom Kiban Docker image:

```shell
    oblt-cli cluster create serverless \
      --project-type "observability" \
      --kibana-docker-image "my-kibana-image:latest"
```

Using a custom Elasticsearch Docker image:

```shell
    oblt-cli cluster create serverless \
      --project-type "observability" \
      --elasticsearch-docker-image "my-elasticsearch-image:latest"
```

Using a custom Elastic Agent Docker image:

```shell
    oblt-cli cluster create serverless \
      --project-type "observability" \
      --fleet-docker-image "my-elastic-agent-image:latest"
```

### Environment variables

One of the secrets created for each cluster contains the credentials and URLs for the different components of the cluster, defined as environment variables. You can retrieve that secret using the following command:

```shell
    oblt-cli cluster secrets env --cluster-name my-project-oblt --output-file ${PWD}/.env
```

```shell
ELASTICSEARCH_HOSTS=https://my-project-ad85d6.es.us-east-1.aws.elastic.cloud
ELASTICSEARCH_HOST=https://my-project-ad85d6.es.us-east-1.aws.elastic.cloud
ELASTICSEARCH_USERNAME=testing-internal
ELASTICSEARCH_PASSWORD=MyPaSsWord
KIBANA_HOST=https://my-project-ad85d6.kb.us-east-1.aws.elastic.cloud
KIBANA_HOSTS=https://my-project-ad85d6.kb.us-east-1.aws.elastic.cloud
KIBANA_FLEET_HOST=https://my-project-ad85d6.kb.us-east-1.aws.elastic.cloud
KIBANA_USERNAME=testing-internal
KIBANA_PASSWORD=MyPaSsWord
ELASTIC_APM_SERVER_URL=https://my-project-ad85d6.apm.us-east-1.aws.elastic.cloud
ELASTIC_APM_JS_SERVER_URL=https://my-project-ad85d6.apm.us-east-1.aws.elastic.cloud
ELASTIC_APM_JS_BASE_SERVER_URL=https://my-project-ad85d6.apm.us-east-1.aws.elastic.cloud
ELASTIC_APM_SECRET_TOKEN=
ELASTIC_APM_API_KEY=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==
FLEET_ELASTICSEARCH_HOST=https://my-project-ad85d6.es.us-east-1.aws.elastic.cloud
FLEET_ENROLLMENT_TOKEN=
FLEET_SERVER_SERVICE_TOKEN=
FLEET_SERVER_POLICY_ID=
FLEET_TOKEN_POLICY_NAME=
FLEET_URL=https://my-project-ad85d6.fleet.us-east-1.aws.elastic.cloud
ELASTICSEARCH_API_TOKEN=MY-TOKEN
```

### Health check

Oblt-cli provides a health check command to verify the status of the cluster. You can run the following command to check the health of the cluster:

```shell
    oblt-cli cluster healthcheck --cluster-name my-project-oblt --output-file healthcheck.json
```

```json
{
  "elasticsearch": {
    "cluster_name": "ad85d674bb1544aa92178a7aa4dd0465",
    "cluster_uuid": "Iw6lDLCaTSGdNinRyzgnEA",
    "name": "serverless",
    "tagline": "You Know, for Search",
    "version": {
      "build_date": "2023-10-31",
      "build_flavor": "serverless",
      "build_hash": "00000000",
      "build_snapshot": false,
      "build_type": "docker",
      "lucene_version": "9.7.0",
      "minimum_index_compatibility_version": "8.11.0",
      "minimum_wire_compatibility_version": "8.11.0",
      "number": "8.11.0"
    }
  },
  "kibana": {
    "metrics": {
      "collection_interval_in_millis": 5000,
      "concurrent_connections": 6,
      "elasticsearch_client": {
        "totalActiveSockets": 0,
        "totalIdleSockets": 4,
        "totalQueuedRequests": 0
      },
      "last_updated": "2024-05-07T15:53:39.121Z",
      "os": {
        "cpu": {
          "cfs_period_micros": 100000,
          "cfs_quota_micros": 200000,
          "control_group": "/",
          "stat": {
            "number_of_elapsed_periods": 814921,
            "number_of_times_throttled": 1863,
            "time_throttled_nanos": 422263932114
          }
        },
        "cpuacct": {
          "control_group": "/",
          "usage_nanos": 2944230568970
        },
        "distro": "Ubuntu",
        "distroRelease": "Ubuntu-20.04",
        "load": {
          "15m": 3.32,
          "1m": 2.1,
          "5m": 3.15
        },
        "memory": {
          "free_in_bytes": 104335679488,
          "total_in_bytes": 133452976128,
          "used_in_bytes": 29117296640
        },
        "platform": "linux",
        "platformRelease": "linux-5.10.210-201.855.amzn2.aarch64",
        "uptime_in_millis": 688259130
      },
      "process": {
        "event_loop_delay": 11.583487,
        "event_loop_delay_histogram": {
          "exceeds": 0,
          "fromTimestamp": "2024-05-07T15:53:34.122Z",
          "lastUpdatedAt": "2024-05-07T15:53:39.120Z",
          "max": 11.583487,
          "mean": 10.075406997979798,
          "min": 9.29792,
          "percentiles": {
            "50": 10.387455,
            "75": 10.420223,
            "95": 10.575871,
            "99": 10.706943
          },
          "stddev": 0.48435921920702013
        },
        "event_loop_utilization": {
          "active": 97.15028692781925,
          "idle": 4901.52854500711,
          "utilization": 0.019435192816781455
        },
        "memory": {
          "array_buffers_in_bytes": 970812,
          "external_in_bytes": 4548491,
          "heap": {
            "size_limit": 764411904,
            "total_in_bytes": 359981056,
            "used_in_bytes": 311762472
          },
          "resident_set_size_in_bytes": 496357376
        },
        "pid": 7,
        "uptime_in_millis": 81492113.936818
      },
      "processes": [
        {
          "event_loop_delay": 11.583487,
          "event_loop_delay_histogram": {
            "exceeds": 0,
            "fromTimestamp": "2024-05-07T15:53:34.122Z",
            "lastUpdatedAt": "2024-05-07T15:53:39.120Z",
            "max": 11.583487,
            "mean": 10.075406997979798,
            "min": 9.29792,
            "percentiles": {
              "50": 10.387455,
              "75": 10.420223,
              "95": 10.575871,
              "99": 10.706943
            },
            "stddev": 0.48435921920702013
          },
          "event_loop_utilization": {
            "active": 97.15028692781925,
            "idle": 4901.52854500711,
            "utilization": 0.019435192816781455
          },
          "memory": {
            "array_buffers_in_bytes": 970812,
            "external_in_bytes": 4548491,
            "heap": {
              "size_limit": 764411904,
              "total_in_bytes": 359981056,
              "used_in_bytes": 311762472
            },
            "resident_set_size_in_bytes": 496357376
          },
          "pid": 7,
          "uptime_in_millis": 81492113.936818
        }
      ],
      "requests": {
        "disconnects": 0,
        "statusCodes": {
          "200": 1
        },
        "status_codes": {
          "200": 1
        },
        "total": 1
      },
      "response_times": {
        "avg_in_millis": 2,
        "max_in_millis": 2
      }
    },
    "name": "kb",
    "status": {
      "core": {
        "elasticsearch": {
          "level": "available",
          "meta": {
            "incompatibleNodes": [],
            "warningNodes": []
          },
          "summary": "Elasticsearch is available"
        },
        "savedObjects": {
          "level": "available",
          "meta": {
            "migratedIndices": {
              "migrated": 0,
              "patched": 6,
              "skipped": 0
            }
          },
          "summary": "SavedObjects service has completed migrations and is available"
        }
      },
      "overall": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "plugins": {
        "actions": {
          "level": "available",
          "summary": "All services and plugins are available"
        },
        "advancedSettings": {
          "level": "available",
          "summary": "All services and plugins are available"
        },

        ...

        "visualizations": {
          "level": "available",
          "summary": "All services and plugins are available"
        }
      }
    },
    "uuid": "2a56b6ea-3530-469b-97df-a6f12e7118fd",
    "version": {
      "build_date": "2024-05-03T00:43:48.160Z",
      "build_flavor": "serverless",
      "build_hash": "f7be3ba82cd93c7ece35189105aa279be589b68b",
      "build_number": 74129,
      "build_snapshot": false,
      "number": "8.15.0"
    }
  }
}
```

{% include '/user-guide/mki.md' ignore missing %}
