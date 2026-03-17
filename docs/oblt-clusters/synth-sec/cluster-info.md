
Elastic Stack (deployed at 2024-08-06 02:41:22)

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
* [4c782f7ca883c945716fa97a84e4c35ebe8a4e74](https://github.com/elastic/elasticsearch/compare/4c782f7ca883c945716fa97a84e4c35ebe8a4e74...8.16)
<details>
<summary>Click for more details</summary>

```json
{
  "cluster_name": "45a09618f6a145e8bf4f00ae818c8b77",
  "cluster_uuid": "iCZLG45ZT9qsibgK8a84BA",
  "name": "instance-0000000000",
  "tagline": "You Know, for Search",
  "version": {
    "build_date": "2024-08-05T13:15:35.606117526Z",
    "build_flavor": "default",
    "build_hash": "4c782f7ca883c945716fa97a84e4c35ebe8a4e74",
    "build_snapshot": true,
    "build_type": "docker",
    "lucene_version": "9.11.1",
    "minimum_index_compatibility_version": "7.0.0",
    "minimum_wire_compatibility_version": "7.17.0",
    "number": "8.16.0-SNAPSHOT"
  }
}
```

</details>

### Kibana:
* type: ess
* [477d07f7fc8f7dcc93a2e8701205398860ea6676](https://github.com/elastic/Kibana/compare/477d07f7fc8f7dcc93a2e8701205398860ea6676...8.16)
<details>
<summary>Click for more details</summary>

```json
{
  "metrics": {
    "collection_interval_in_millis": 5000,
    "concurrent_connections": 9,
    "elasticsearch_client": {
      "totalActiveSockets": 0,
      "totalIdleSockets": 7,
      "totalQueuedRequests": 0
    },
    "last_updated": "2024-08-06T02:41:19.021Z",
    "os": {
      "distro": "Ubuntu",
      "distroRelease": "Ubuntu-20.04",
      "load": {
        "15m": 1.4,
        "1m": 2.37,
        "5m": 2.06
      },
      "memory": {
        "free_in_bytes": 39760769024,
        "total_in_bytes": 71758884864,
        "used_in_bytes": 31998115840
      },
      "platform": "linux",
      "platformRelease": "linux-5.15.0-1032-gcp",
      "uptime_in_millis": 1216240580
    },
    "process": {
      "event_loop_delay": 10.911743,
      "event_loop_delay_histogram": {
        "exceeds": 0,
        "fromTimestamp": "2024-08-06T02:41:14.021Z",
        "lastUpdatedAt": "2024-08-06T02:41:19.021Z",
        "max": 10.911743,
        "mean": 10.060965161290323,
        "min": 9.207808,
        "percentiles": {
          "50": 10.297343,
          "75": 10.321919,
          "95": 10.502143,
          "99": 10.559487
        },
        "stddev": 0.45895391661635276
      },
      "event_loop_utilization": {
        "active": 45.01543207780924,
        "idle": 4955.115454000072,
        "utilization": 0.00900285074599707
      },
      "memory": {
        "array_buffers_in_bytes": 3238887,
        "external_in_bytes": 6812126,
        "heap": {
          "size_limit": 3405774848,
          "total_in_bytes": 447713280,
          "used_in_bytes": 393832320
        },
        "resident_set_size_in_bytes": 571441152
      },
      "pid": 21,
      "uptime_in_millis": 687080.632041
    },
    "processes": [
      {
        "event_loop_delay": 10.911743,
        "event_loop_delay_histogram": {
          "exceeds": 0,
          "fromTimestamp": "2024-08-06T02:41:14.021Z",
          "lastUpdatedAt": "2024-08-06T02:41:19.021Z",
          "max": 10.911743,
          "mean": 10.060965161290323,
          "min": 9.207808,
          "percentiles": {
            "50": 10.297343,
            "75": 10.321919,
            "95": 10.502143,
            "99": 10.559487
          },
          "stddev": 0.45895391661635276
        },
        "event_loop_utilization": {
          "active": 45.01543207780924,
          "idle": 4955.115454000072,
          "utilization": 0.00900285074599707
        },
        "memory": {
          "array_buffers_in_bytes": 3238887,
          "external_in_bytes": 6812126,
          "heap": {
            "size_limit": 3405774848,
            "total_in_bytes": 447713280,
            "used_in_bytes": 393832320
          },
          "resident_set_size_in_bytes": 571441152
        },
        "pid": 21,
        "uptime_in_millis": 687080.632041
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
      "avg_in_millis": 1,
      "max_in_millis": 1
    }
  },
  "name": "instance-0000000044",
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
            "patched": 7,
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
      "aiAssistantManagementSelection": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "aiops": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "alerting": {
        "level": "available",
        "reported": true,
        "summary": "Alerting is (probably) ready"
      },
      "apm": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "apmDataAccess": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "assetsDataAccess": {
        "level": "available",
        "summary": "All services are available"
      },
      "banners": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "bfetch": {
        "level": "available",
        "summary": "All services are available"
      },
      "canvas": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cases": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "charts": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloud": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudChat": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudDataMigration": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudDefend": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudExperiments": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudFullStory": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudLinks": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "cloudSecurityPosture": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "console": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "contentManagement": {
        "level": "available",
        "summary": "All services are available"
      },
      "controls": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "crossClusterReplication": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "customBranding": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "customIntegrations": {
        "level": "available",
        "summary": "All services are available"
      },
      "dashboard": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dashboardEnhanced": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "data": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataQuality": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataViewEditor": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataViewFieldEditor": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataViewManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataViews": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "dataVisualizer": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "datasetQuality": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "devTools": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "discover": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "discoverEnhanced": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "discoverShared": {
        "level": "available",
        "summary": "All services are available"
      },
      "ecsDataQualityDashboard": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "elasticAssistant": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "embeddable": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "embeddableEnhanced": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "encryptedSavedObjects": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "enterpriseSearch": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "entityManager": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "esUiShared": {
        "level": "available",
        "summary": "All services are available"
      },
      "esql": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "esqlDataGrid": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "eventAnnotation": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "eventAnnotationListing": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "eventLog": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "exploratoryView": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionError": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionGauge": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionHeatmap": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionImage": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionLegacyMetricVis": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionMetric": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionMetricVis": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionPartitionVis": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionRepeatImage": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionRevealImage": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionShape": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionTagcloud": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressionXY": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "expressions": {
        "level": "available",
        "summary": "All services are available"
      },
      "features": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "fieldFormats": {
        "level": "available",
        "summary": "All services are available"
      },
      "fieldsMetadata": {
        "level": "available",
        "summary": "All services are available"
      },
      "fileUpload": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "files": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "filesManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "fleet": {
        "level": "available",
        "reported": true,
        "summary": "Fleet is available"
      },
      "ftrApis": {
        "level": "available",
        "summary": "All services are available"
      },
      "globalSearch": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "globalSearchBar": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "globalSearchProviders": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "graph": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "grokdebugger": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "guidedOnboarding": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "home": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "imageEmbeddable": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "indexLifecycleManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "indexManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "infra": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "ingestPipelines": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "inputControlVis": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "inspector": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "integrationAssistant": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "kibanaOverview": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "kibanaReact": {
        "level": "available",
        "summary": "All services are available"
      },
      "kibanaUsageCollection": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "kibanaUtils": {
        "level": "available",
        "summary": "All services are available"
      },
      "kubernetesSecurity": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "lens": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "licenseApiGuard": {
        "level": "available",
        "summary": "All services are available"
      },
      "licenseManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "licensing": {
        "level": "available",
        "reported": true,
        "summary": "License fetched"
      },
      "links": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "lists": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "logsDataAccess": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "logsExplorer": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "logsShared": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "logstash": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "management": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "maps": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "mapsEms": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "metricsDataAccess": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "ml": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "monitoring": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "monitoringCollection": {
        "level": "available",
        "summary": "All services are available"
      },
      "navigation": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "newsfeed": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "noDataPage": {
        "level": "available",
        "summary": "All services are available"
      },
      "notifications": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observability": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityAIAssistant": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityAIAssistantApp": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityAiAssistantManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityLogsExplorer": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityOnboarding": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "observabilityShared": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "osquery": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "painlessLab": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "presentationPanel": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "presentationUtil": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "profiling": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "profilingDataAccess": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "remoteClusters": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "reporting": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "rollup": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "ruleRegistry": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "runtimeFields": {
        "level": "available",
        "summary": "All services are available"
      },
      "savedObjects": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "savedObjectsFinder": {
        "level": "available",
        "summary": "All services are available"
      },
      "savedObjectsManagement": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "savedObjectsTagging": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "savedObjectsTaggingOss": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "savedSearch": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "screenshotMode": {
        "level": "available",
        "summary": "All services are available"
      },
      "screenshotting": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "searchConnectors": {
        "level": "available",
        "summary": "All services are available"
      },
      "searchHomepage": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "searchInferenceEndpoints": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "searchNotebooks": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "searchPlayground": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "searchprofiler": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "security": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "securitySolution": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "securitySolutionEss": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "sessionView": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "share": {
        "level": "available",
        "summary": "All services are available"
      },
      "slo": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "snapshotRestore": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "spaces": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "stackAlerts": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "stackConnectors": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "synthetics": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "taskManager": {
        "level": "available",
        "reported": true,
        "summary": "Task Manager is healthy"
      },
      "telemetry": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "telemetryCollectionManager": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "telemetryCollectionXpack": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "telemetryManagementSection": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "threatIntelligence": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "timelines": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "transform": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "translations": {
        "level": "available",
        "summary": "All services are available"
      },
      "triggersActionsUi": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "uiActions": {
        "level": "available",
        "summary": "All services are available"
      },
      "uiActionsEnhanced": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "unifiedDocViewer": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "unifiedHistogram": {
        "level": "available",
        "summary": "All services are available"
      },
      "unifiedSearch": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "upgradeAssistant": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "uptime": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "urlDrilldown": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "urlForwarding": {
        "level": "available",
        "summary": "All services are available"
      },
      "usageCollection": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "ux": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visDefaultEditor": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeGauge": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeHeatmap": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeMarkdown": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeMetric": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypePie": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeTable": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeTagcloud": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeTimelion": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeTimeseries": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeVega": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeVislib": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visTypeXy": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "visualizations": {
        "level": "available",
        "summary": "All services and plugins are available"
      },
      "watcher": {
        "level": "available",
        "summary": "All services and plugins are available"
      }
    }
  },
  "uuid": "226555aa-391e-47e6-876b-60976a17f473",
  "version": {
    "build_date": "2024-08-05T11:09:18.304Z",
    "build_flavor": "traditional",
    "build_hash": "477d07f7fc8f7dcc93a2e8701205398860ea6676",
    "build_number": 77092,
    "build_snapshot": true,
    "number": "8.16.0"
  }
}
```

</details>

### APM:
* type: ess
* [972c4917a78d124da158f9f7d555e98cc00197be](https://github.com/elastic/apm-server/compare/972c4917a78d124da158f9f7d555e98cc00197be...8.16)
<details>
<summary>Click for more details</summary>

```json
{
  "build_date": "2024-08-02T19:20:23Z",
  "build_sha": "972c4917a78d124da158f9f7d555e98cc00197be",
  "publish_ready": true,
  "version": "8.16.0"
}
```

</details>

[kibana-url]: https://synth-sec.kb.us-west2.gcp.elastic-cloud.com:443
[logs-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/logs/stream?logFilter=(language:kuery,query:%27%22service.name%22:%22synth-sec%22%27)
[metrics-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/monitoring
[apm-url]: https://monitoring-oblt.kb.us-west2.gcp.elastic-cloud.com/app/apm/services?kuery=labels.deploymentName:%20%22synth-sec%22
[console-url]: https://admin.found.no/deployments?q=alias%3Asynth-sec
