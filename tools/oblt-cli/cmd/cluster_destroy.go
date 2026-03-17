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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// DestroyCmd represents the destroy command
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Command to destroy a cluster.",
	Long: `Command to destroy a cluster.
you can use 'oblt-cli cluster list' to list the clusters available to destroy.`,
	Run: runDestroy,
}

func init() {
	ClusterCmd.AddCommand(DestroyCmd)

	DestroyCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to destroy. (Required)")
	DestroyCmd.Flags().Bool(config.ForceFlag, false, "Do not ask for comfirmation.")
	DestroyCmd.Flags().Bool(config.AllFlag, false, "Destroy any cluster not only users clusters.")
	DestroyCmd.Flags().Bool(config.WipeupFlag, false, "Wipeup the cluster all user clusters destroy.")
}

/*
runDestroy Searches for the file of the cluster you select and remove the file from the repository.
*/
func runDestroy(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	all, _ := cmd.Flags().GetBool(config.AllFlag)
	wipeup, _ := cmd.Flags().GetBool(config.WipeupFlag)

	labelClusterName := apm.Label{Key: "cluster-name", Value: clusterName}
	tx, ctx := apm.StartTransaction("runDestroy", "request", []apm.Label{labelVersion, labelClusterName}, userConfig)
	defer apm.Flush(tx)

	if wipeup && all {
		apm.ReportError(tx, "wipeup and all flags are not compatible")
	}

	if wipeup && clusterName != "" {
		apm.ReportError(tx, "wipeup and clusterName flags are not compatible")
	}

	if clusterName == "" && !wipeup {
		apm.ReportError(tx, "clusterName or wipeup flag is required")
	}

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	force, _ := cmd.Flags().GetBool(config.ForceFlag)
	if !force {
		prompt.Confirm("Do you want to destroy the cluster : " + clusterName + "?")
	}

	if wipeup {
		results, err := obltTestEnvironments.Wipeup()
		apm.CobraCheckErr(err, tx, ctx)
		for _, result := range results {
			saveResults(result, outputFile)
		}
	} else {
		results, err := obltTestEnvironments.DestroyCluster(clusterName, all)
		apm.CobraCheckErr(err, tx, ctx)
		saveResults(results, outputFile)
	}
}
