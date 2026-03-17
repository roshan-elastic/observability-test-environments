
Elastic Stack (deployed at 2023-09-06 10:36:04)

[Open :fontawesome-solid-house:][kibana-url]{ .md-button }
[Console :fontawesome-solid-gears:][console-url]{ .md-button }
[Logs :fontawesome-solid-scroll:][logs-url]{ .md-button }
[APM :material-chart-line:][apm-url]{ .md-button }
[Metrics :material-chart-line:][metrics-url]{ .md-button }

!!! Note

    * ess: Elastic Cloud Deployment (Elastic Terraform provider)
    * eck: Kubernetes Deployment with ECK

### Elasticsearch:
* type: ess
* [09520b59b6bc1057340b55750186466ea715e30e](https://github.com/elastic/elasticsearch/compare/09520b59b6bc1057340b55750186466ea715e30e...8.7)
<details>
<summary>Click for more details</summary>

```json
{
  "cluster_name": "463bdb12ebd6490cafe287dff57da44f",
  "cluster_uuid": "Ff1A7WvtRFK8IG4-dd9vXg",
  "name": "instance-0000000019",
  "tagline": "You Know, for Search",
  "version": {
    "build_date": "2023-03-27T16:31:09.816451435Z",
    "build_flavor": "default",
    "build_hash": "09520b59b6bc1057340b55750186466ea715e30e",
    "build_snapshot": false,
    "build_type": "docker",
    "lucene_version": "9.5.0",
    "minimum_index_compatibility_version": "7.0.0",
    "minimum_wire_compatibility_version": "7.17.0",
    "number": "8.7.0"
  }
}
```

</details>

### Kibana:
* type: ess
* [05f12599523732051b84dde0b8d5610e0db2b06d](https://github.com/elastic/Kibana/compare/05f12599523732051b84dde0b8d5610e0db2b06d...8.7)
<details>
<summary>Click for more details</summary>

```json
{
  "metrics": {
    "collection_interval_in_millis": 5000,
    "concurrent_connections": 16,
    "elasticsearch_client": {
      "totalActiveSockets": 0,
      "totalIdleSockets": 5,
      "totalQueuedRequests": 0
    },
    "last_updated": "2023-09-06T10:36:00.704Z",
    "os": {
      "distro": "Ubuntu",
      "distroRelease": "Ubuntu-20.04",
      "load": {
        "15m": 3.49,
        "1m": 3.4,
        "5m": 3.2
      },
      "memory": {
        "free_in_bytes": 90797862912,
        "total_in_bytes": 147675066368,
        "used_in_bytes": 56877203456
      },
      "platform": "linux",
      "platformRelease": "linux-5.4.0-1049-gcp",
      "uptime_in_millis": 30419713630
    },
    "process": {
      "event_loop_delay": 10.120189922920892,
      "event_loop_delay_histogram": {
        "exceeds": 0,
        "fromTimestamp": "2023-09-06T10:35:55.704Z",
        "lastUpdatedAt": "2023-09-06T10:36:00.701Z",
        "max": 14.385151,
        "mean": 10.120189922920892,
        "min": 7.335936,
        "percentiles": {
          "50": 10.117119,
          "75": 10.199039,
          "95": 10.313727,
          "99": 10.534911
        },
        "stddev": 0.33643433220108254
      },
      "memory": {
        "heap": {
          "size_limit": 3405774848,
          "total_in_bytes": 322809856,
          "used_in_bytes": 268661672
        },
        "resident_set_size_in_bytes": 488189952
      },
      "pid": 19,
      "uptime_in_millis": 659890.373833
    },
    "processes": [
      {
        "event_loop_delay": 10.120189922920892,
        "event_loop_delay_histogram": {
          "exceeds": 0,
          "fromTimestamp": "2023-09-06T10:35:55.704Z",
          "lastUpdatedAt": "2023-09-06T10:36:00.701Z",
          "max": 14.385151,
          "mean": 10.120189922920892,
          "min": 7.335936,
          "percentiles": {
            "50": 10.117119,
            "75": 10.199039,
            "95": 10.313727,
            "99": 10.534911
          },
          "stddev": 0.33643433220108254
        },
        "memory": {
          "heap": {
            "size_limit": 3405774848,
            "total_in_bytes": 322809856,
            "used_in_bytes": 268661672
          },
          "resident_set_size_in_bytes": 488189952
        },
        "pid": 19,
        "uptime_in_millis": 659890.373833
      }
    ],
    "requests": {
      "disconnects": 0,
      "statusCodes": {
        "302": 1
      },
      "status_codes": {
        "302": 1
      },
      "total": 1
    },
    "response_times": {
      "avg_in_millis": 3,
      "max_in_millis": 3
    }
  },
  "name": "instance-0000000002",
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
            "patched": 2,
            "skipped": 0
          }
        },
        "summary": "SavedObjects service has completed migrations and is available"
      }
    },
    "overall": {
      "level": "available",
      "summary": "All services are available"
    },
    "plugins": {
      "actions": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "advancedSettings": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "aiops": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "alerting": {
        "level": "available",
        "summary": "Alerting is (probably) ready"
      },
      "apm": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "banners": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "bfetch": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "canvas": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cases": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "charts": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloud": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudChat": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudDataMigration": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudDefend": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudExperiments": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudFullStory": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudLinks": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "cloudSecurityPosture": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "console": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "contentManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "controls": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "crossClusterReplication": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "customBranding": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "customIntegrations": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dashboard": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dashboardEnhanced": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "data": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dataViewEditor": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dataViewFieldEditor": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dataViewManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dataViews": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "dataVisualizer": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "devTools": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "discover": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "discoverEnhanced": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "ecsDataQualityDashboard": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "embeddable": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "embeddableEnhanced": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "encryptedSavedObjects": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "enterpriseSearch": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "esUiShared": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "eventAnnotation": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "eventLog": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionError": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionGauge": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionHeatmap": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionImage": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionLegacyMetricVis": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionMetric": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionMetricVis": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionPartitionVis": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionRepeatImage": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionRevealImage": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionShape": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionTagcloud": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressionXY": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "expressions": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "features": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "fieldFormats": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "fileUpload": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "files": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "filesManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "fleet": {
        "level": "available",
        "summary": "Fleet is available"
      },
      "ftrApis": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "globalSearch": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "globalSearchBar": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "globalSearchProviders": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "graph": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "grokdebugger": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "guidedOnboarding": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "home": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "imageEmbeddable": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "indexLifecycleManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "indexManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "infra": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "ingestPipelines": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "inputControlVis": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "inspector": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "kibanaOverview": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "kibanaReact": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "kibanaUsageCollection": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "kibanaUtils": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "kubernetesSecurity": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "lens": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "licenseApiGuard": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "licenseManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "licensing": {
        "level": "available",
        "summary": "License fetched"
      },
      "lists": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "logstash": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "management": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "maps": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "mapsEms": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "ml": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "monitoring": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "monitoringCollection": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "navigation": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "newsfeed": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "notifications": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "observability": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "osquery": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "painlessLab": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "presentationUtil": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "profiling": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "remoteClusters": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "reporting": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "rollup": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "ruleRegistry": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "runtimeFields": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedObjects": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedObjectsFinder": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedObjectsManagement": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedObjectsTagging": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedObjectsTaggingOss": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "savedSearch": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "screenshotMode": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "screenshotting": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "searchprofiler": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "security": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "securitySolution": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "sessionView": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "share": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "snapshotRestore": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "spaces": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "stackAlerts": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "stackConnectors": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "synthetics": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "taskManager": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "telemetry": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "telemetryCollectionManager": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "telemetryCollectionXpack": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "telemetryManagementSection": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "threatIntelligence": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "timelines": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "transform": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "translations": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "triggersActionsUi": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "uiActions": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "uiActionsEnhanced": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "unifiedFieldList": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "unifiedHistogram": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "unifiedSearch": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "upgradeAssistant": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "urlDrilldown": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "urlForwarding": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "usageCollection": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "ux": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visDefaultEditor": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeGauge": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeHeatmap": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeMarkdown": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeMetric": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypePie": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeTable": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeTagcloud": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeTimelion": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeTimeseries": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeVega": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeVislib": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visTypeXy": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "visualizations": {
        "level": "available",
        "summary": "All dependencies are available"
      },
      "watcher": {
        "level": "available",
        "summary": "All dependencies are available"
      }
    }
  },
  "uuid": "8d4fa3e0-4e55-463d-b4dd-9041acb53e84",
  "version": {
    "build_hash": "05f12599523732051b84dde0b8d5610e0db2b06d",
    "build_number": 61109,
    "build_snapshot": false,
    "number": "8.7.0"
  }
}
```

</details>

### APM:
* type: ess
* [80446fbb7881463fa549a8d669055a6f4b897f70](https://github.com/elastic/apm-server/compare/80446fbb7881463fa549a8d669055a6f4b897f70...8.7)
<details>
<summary>Click for more details</summary>

```json
{
  "build_date": "2023-03-27T18:03:39-04:00",
  "build_sha": "80446fbb7881463fa549a8d669055a6f4b897f70",
  "publish_ready": true,
  "version": "8.7.0"
}
```

</details>

[kibana-url]: https://autopilot-oblt.kb.us-west2.gcp.elastic-cloud.com:443
[logs-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/logs/stream?logFilter=(language:kuery,query:%27%22service.name%22:%22autopilot-oblt%22%27)
[metrics-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/monitoring
[apm-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/apm/services?kuery=labels.deploymentName:%20%22autopilot-oblt%22
[console-url]: https://admin.found.no/deployments?q=alias%3Aautopilot-oblt
