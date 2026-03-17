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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/k8s"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// PackagesVSCodeCmd represents the packages vscode command
var PackagesVSCodeCmd = &cobra.Command{
	Use:   "vscode",
	Short: "Command to open a vscode from the workspace container deployed.",
	Long: `Command to open a vscode from the workspace container deployed.
	This VSCode uses the shared volume to persist the data.`,
	Run: runPackagesVSCode,
}

func init() {
	PackagesCmd.AddCommand(PackagesVSCodeCmd)

	PackagesVSCodeCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")

	cobra.MarkFlagRequired(PackagesVSCodeCmd.Flags(), config.ClusterNameFlag)
}

// runVSCode open a vscode from the workspace container deployed
func runPackagesVSCode(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	tx, ctx := apm.StartTransaction("runPackagesVSCode", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	script := `oblt-vscode`
	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)
	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}
