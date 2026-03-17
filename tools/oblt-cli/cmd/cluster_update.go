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
	"github.com/spf13/cobra"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Command to update a cluster.",
	Long:  `Command to update a cluster.`,
	Run:   runUpdate,
}

func init() {
	ClusterCmd.AddCommand(UpdateCmd)

	UpdateCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to destroy. (Required)")
	UpdateCmd.Flags().String(config.ParametersFlag, "", "Parameters values defined in JSON '{ \"var_name1\": \"value\",\"var_name2\": \"value\"}'. (Required)")
	UpdateCmd.Flags().String(config.ParametersFileFlag, "", "Absolute parameters JSON file path to use. It is incompatible with \""+config.ParametersFlag+"\" (Optional)")
}

func runUpdate(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runUpdate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	parameters, err := loadParameters(cmd)

	apm.CobraCheckErr(err, tx, ctx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)

	results, err := obltTestEnvironments.UpdateCluster(clusterName, parameters)
	apm.CobraCheckErr(err, tx, ctx)
	saveResults(results, outputFile)
}
