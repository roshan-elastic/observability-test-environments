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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/bootstrap"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/console"
	"github.com/spf13/cobra"
)

// BootstrapListCmd List the recipes available.
var BootstrapListCmd = &cobra.Command{
	Use:   "list",
	Short: "Command to List the bootstrap recipes available.",
	Long:  "Command to List the bootstrap recipes available.",
	Run:   runBootstrapList,
}

func init() {
	BootstrapCmd.AddCommand(BootstrapListCmd)
	BootstrapListCmd.Flags().String(config.BootstrapFolderFlag, "", "Full path to the bootstarp folder root. (Optional, default: REPO/bootstrap)")
}

func runBootstrapList(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runBootstrapList", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	bootstrapFolder, _ := cmd.Flags().GetString(config.BootstrapFolderFlag)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)
	if bootstrapFolder == "" {
		bootstrapFolder = obltTestEnvironments.GetBootstrapRecipesDir()
	}
	results := bootstrap.ListRecipes(obltTestEnvironments, bootstrapFolder)
	console.PrintYamlFiles(bootstrap.TypeElasticsearch, results[bootstrap.TypeElasticsearch])
	console.PrintYamlFiles(bootstrap.TypeKibana, results[bootstrap.TypeKibana])
	saveResults(results, outputFile)
}
