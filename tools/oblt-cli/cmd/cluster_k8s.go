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

// CreateCmd represents the create command
var ClusterK8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Command to access to the k8s cluster resources.",
	Long: `Command to access to the k8s cluster resources.
		The command open a subshell with the kubectl command configured to access to the k8s cluster.
		In that subshell you can execute any kubectl, helm and other k8s commands in the context of the k8s cluster.`,
	Run: runK8s,
}

func init() {
	ClusterCmd.AddCommand(ClusterK8sCmd)

	ClusterK8sCmd.PersistentFlags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")

	cobra.MarkFlagRequired(ClusterK8sCmd.PersistentFlags(), config.ClusterNameFlag)
}

func runK8s(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	tx, ctx := apm.StartTransaction("runK8s", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	script := `echo "Welcome to the k8s cluster shell"
	`

	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)

	err = k8sShell.ExecScript()
	apm.CobraCheckErr(err, tx, ctx)
}
