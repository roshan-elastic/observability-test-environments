// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package interactions

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

type SafeArrayFiles struct {
	sync.RWMutex
	files []files.YamlFile
}

type SafeMap struct {
	sync.RWMutex
	data []map[string]interface{}
}

type SafeMapString struct {
	sync.RWMutex
	data map[string]string
}

// Clusters is the list of clusters configuration files
var Clusters SafeArrayFiles

// ClustersExtra is the list of clusters configuration files enriched with extra information
var ClustersExtra SafeMap

// Templates is the list of templates configuration files
var Templates SafeArrayFiles

// lastExecution is the last time the clusters data was updated
var lastExecution time.Time

// Users is the list of Slack ID to OBLT username
var Users SafeMapString = SafeMapString{data: make(map[string]string)}

// UpdateClustersData will update the clusters data from the repository.
func UpdateClustersData() (executed bool) {
	executed = false
	now := time.Now()
	if now.Sub(lastExecution) > time.Minute {
		lastExecution = now

		obltTestEnvironments, err := OnMemoryUserConf("cache", "cache", true)
		logger.LogError("Error updating clusters Data", err)
		if err == nil {
			UpdateClusters(*obltTestEnvironments)
			UpdateTemplates(*obltTestEnvironments)
			UpdateSlackUsers()
			executed = true
		}
	}
	return executed
}

// updateSlackUsers will update the list of Slack ID to OBLT username
func UpdateSlackUsers() {
	ClustersExtra.Lock()
	Users.Lock()
	defer ClustersExtra.Unlock()
	defer Users.Unlock()
	for _, cluster := range ClustersExtra.data {
		if cluster["slackChannel"] != nil && cluster["obltUsername"] != nil {
			Users.data[cluster["slackChannel"].(string)] = cluster["obltUsername"].(string)
		}
	}
}

// updateTemplates will update the list of templates
func UpdateTemplates(obltTestEnvironments clusters.ObltEnvironmentsRepository) {
	Templates.Lock()
	defer Templates.Unlock()
	Templates.files = obltTestEnvironments.ListTemplates()
}

// updateClusters will update the list of clusters
func UpdateClusters(obltTestEnvironments clusters.ObltEnvironmentsRepository) {
	ClustersExtra.Lock()
	defer ClustersExtra.Unlock()
	Clusters.Lock()
	defer Clusters.Unlock()
	Clusters.files = obltTestEnvironments.ListClusters(true)
	ClustersExtra.data = enrichClusterList(Clusters.files)
}

// FilterClustersByUser will filter the clusters by the user
func FilterClustersByUser(slackID string) []files.YamlFile {
	Clusters.Lock()
	defer Clusters.Unlock()
	var filteredClusters []files.YamlFile
	userID := SlackIDToUsername(slackID)
	for _, cluster := range Clusters.files {
		if cluster.Owner == userID || cluster.Data["slack_channel"] == "@"+slackID {
			filteredClusters = append(filteredClusters, cluster)
		}
	}
	return filteredClusters
}

// SlackIDToUsername will return the OBLT username from the Slack ID
func SlackIDToUsername(slackID string) (obltUser string) {
	obltUser = slackID
	Users.Lock()
	defer Users.Unlock()
	if customUser, ok := Users.data["@"+slackID]; ok {
		obltUser = customUser
	}
	return obltUser
}

// SlackToUsername will return the OBLT username from the Slack ID/Slack name
func SlackToUsername(slackID, slackName string) (obltUser string) {
	obltUser = slackName
	Users.Lock()
	defer Users.Unlock()
	if customUser, ok := Users.data["@"+slackID]; ok {
		obltUser = customUser
	}
	return obltUser
}

// enrichClusterList will enrich the list of clusters with extra information
func enrichClusterList(clusters []files.YamlFile) (items []map[string]interface{}) {
	for _, file := range clusters {
		if file.Data["cluster_name"] != nil {
			clusterName := file.Data["cluster_name"].(string)
			clusterConfigPath := file.Path
			clusterOwner := file.Owner
			isGoldenCluster := false
			if file.Data["golden_cluster"] != nil && file.Data["golden_cluster"].(bool) {
				isGoldenCluster = true
			}
			slackChannel := file.Data["slack_channel"]
			obltUsername := file.Data["oblt_username"]

			items = append(items, map[string]interface{}{
				"clusterName":       clusterName,
				"clusterOwner":      clusterOwner,
				"clusterConfigPath": clusterConfigPath,
				"isGoldenCluster":   isGoldenCluster,
				"slackChannel":      slackChannel,
				"obltUsername":      obltUsername,
			})

		}
	}
	return items
}

// GoldenClusters will return the list of golden clusters
func GoldenClusters() (clusters []string) {
	ClustersExtra.Lock()
	defer ClustersExtra.Unlock()
	for _, cluster := range ClustersExtra.data {
		if cluster["isGoldenCluster"] == true {
			clusters = append(clusters, cluster["clusterName"].(string))
		}
	}
	sort.Strings(clusters)
	return clusters
}

// OnMemoryUserConf will create a fake user configuration file to be able to use the oblt-cli
func OnMemoryUserConf(userID, obltUser string, dryRun bool) (obltTestEnvironments *clusters.ObltEnvironmentsRepository, err error) {
	configFile := config.ForUser(userID)

	userConfig := config.ObltConfiguration{
		ConfigFile:   configFile,
		SlackChannel: "@" + userID,
		Username:     obltUser,
		GitHttpMode:  true,
		Verbose:      true,
	}

	logger.Debugf("Writing configuration file %s", configFile)
	err = os.MkdirAll(filepath.Dir(configFile), 0700)
	if err == nil {
		obltTestEnvironments = clusters.NewObltTestEnvironments(userConfig, false)
		obltTestEnvironments.Repo(dryRun).Sync()
	}
	return obltTestEnvironments, err
}

// ServerlessEnvironments will return the list of environments
func ServerlessEnvironments() (clusters []string) {
	return []string{"qa", "staging", "production"}
}

// ServerlessProjects will return the list of supported projects
func ServerlessProjects() (projects []string) {
	return []string{"observability", "security", "elasticsearch"}
}
