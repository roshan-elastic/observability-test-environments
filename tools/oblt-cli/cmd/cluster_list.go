// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	apmLib "go.elastic.co/apm/v2"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Command to list the cluster configurations.",
	Long:  `Command to list the cluster configurations created by the user. Configurations which are listed are not guaranteed to have been deployed.`,
	Run:   runList,
}

func init() {
	ClusterCmd.AddCommand(ListCmd)

	ListCmd.Flags().Bool(config.AllFlag, false, "Lists oblt and users clusters.")
	ListCmd.Flags().String(config.FilterFlag, "", "Filter by YAML key01=value01,key02=value02 pair.")
}

/*
It will checkout the observability test environments repository,
and check all the files in the user environments folder.
*/
func runList(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	var files []files.YamlFile
	var results = make(map[string]interface{})
	var items []interface{}
	data := [][]string{}
	all, _ := cmd.Flags().GetBool(config.AllFlag)
	filtersValue, _ := cmd.Flags().GetString(config.FilterFlag)

	tx, ctx := apm.StartTransaction("runList", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	if filtersValue != "" {
		filtersMap := processFilter(ctx, tx, filtersValue)
		files, err = obltTestEnvironments.FindClusterByYamlPath(filtersMap)
		apm.CobraCheckErr(err, tx, ctx)
	} else {
		files = obltTestEnvironments.ListClusters(all)
	}

	for _, file := range files {
		if file.Data["cluster_name"] != nil {
			clusterName, _ := file.Data["cluster_name"].(string)
			clusterConfigPath := file.Path
			clusterOwner := file.Owner
			isGoldenCluster, _ := file.Data["golden_cluster"].(bool)
			slackChannel, _ := file.Data["slack_channel"].(string)
			obltUsername, _ := file.Data["oblt_username"].(string)
			updated_at, _ := file.Data["updated_at"].(string)

			items = append(items, map[string]interface{}{
				"clusterName":       clusterName,
				"clusterOwner":      clusterOwner,
				"clusterConfigPath": clusterConfigPath,
				"isGoldenCluster":   isGoldenCluster,
				"slackChannel":      slackChannel,
				"obltUsername":      obltUsername,
				"updatedAt":         updated_at,
			})

			data = append(data, []string{clusterName, clusterOwner, updated_at, clusterConfigPath})
		}
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Cluster Name", "Owner", "Updated at", "File path"})
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	results["items"] = items
	saveResults(results, outputFile)
}

// processFilter will process the filter flag and return a map with the key value pairs.
func processFilter(ctx context.Context, tx *apmLib.Transaction, filtersValue string) map[string]string {
	var filtersMap = make(map[string]string)
	filters := strings.Split(filtersValue, ",")
	for _, filter := range filters {
		items := strings.Split(filter, "=")
		if len(items) != 2 {
			apm.CobraCheckErr(errors.New("filter must be in the format key=value"), tx, ctx)
		}
		filtersMap[items[0]] = items[1]
	}
	return filtersMap
}
