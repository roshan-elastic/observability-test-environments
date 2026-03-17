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
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/k8s"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// PackagesListCmd represents the packages list command
var PackagesLoginCmd = &cobra.Command{
	Use:   "auth",
	Short: "Command to authenticate in the Kubernetes cluster",
	Long: `Command to authenticate in the Kubernetes cluster.
	With this command you will retrieve the credentials to authenticate in the Kubernetes cluster.`,
	Run: runPackagesLogin,
}

func init() {
	PackagesCmd.AddCommand(PackagesLoginCmd)

	PackagesLoginCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")

	cobra.MarkFlagRequired(PackagesLoginCmd.Flags(), config.ClusterNameFlag)
}

// runList shows the list of integrations available
func runPackagesLogin(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	tx, ctx := apm.StartTransaction("runPackagesLogin", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	removeLoginMark(userConfig)
	k8sShell := k8s.NewK8sShell(config, "", userConfig.GetDir(), dryRun)
	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}

// getLoginMarkFile returns the path to the file that marks the login
func getLoginMarkFile(userConfig config.ObltConfiguration) string {
	cfgDir := userConfig.GetDir()
	return filepath.Join(cfgDir, "gke-login")
}

// removeLoginMark removes the file that marks the login
func removeLoginMark(userConfig config.ObltConfiguration) {
	loginMark := getLoginMarkFile(userConfig)
	if _, err := os.Stat(loginMark); err == nil {
		os.Remove(loginMark)
	}
}
