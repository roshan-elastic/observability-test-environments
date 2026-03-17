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

// PackagesShellCmd represents the packages shell command
var PackagesShellCmd = &cobra.Command{
	Use:   "exec",
	Short: "Command to open a shell to the containers deployed in the cluster.",
	Long: `Command to open a shell to the containers deployed in the cluster.
	The command will opean an iteractive shell to the containers deployed in the cluster.`,
}

// PackagesShellServiceCmd represents the packages shell service command
var PackagesShellServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Command to open a shell to the service container.",
	Long:  `Command to open a shell to the service container.`,
	Run:   runPackagesShellService,
}

// PackagesShellElasticAgentCmd represents the packages shell elastic-agent command
var PackagesShellElasticAgentCmd = &cobra.Command{
	Use:   "elastic-agent",
	Short: "Command to open a shell to the Elastic Agent container.",
	Long:  `Command to open a shell to the Elastic Agent container.`,
	Run:   runPackagesShellElasticAgent,
}

// PackagesShellPackageRegistryCmd represents the packages shell package-registry command
var PackagesShellPackageRegistryCmd = &cobra.Command{
	Use:   "package-registry",
	Short: "Command to open a shell to the package registry container.",
	Long:  `Command to open a shell to the package registry container.`,
	Run:   runPackagesShellPackageRegistry,
}

// PackagesShellWorkspaceCmd represents the packages shell workspace command
var PackagesShellWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Command to open a shell to the workspace container.",
	Long:  `Command to open a shell to the workspace container.`,
	Run:   runPackagesShellWorkspace,
}

func init() {
	PackagesCmd.AddCommand(PackagesShellCmd)
	PackagesShellCmd.AddCommand(PackagesShellServiceCmd)
	PackagesShellCmd.AddCommand(PackagesShellElasticAgentCmd)
	PackagesShellCmd.AddCommand(PackagesShellPackageRegistryCmd)
	PackagesShellCmd.AddCommand(PackagesShellWorkspaceCmd)

	PackagesShellCmd.PersistentFlags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")

	cobra.MarkFlagRequired(PackagesShellCmd.PersistentFlags(), config.ClusterNameFlag)
}

// runShell open a shell to the containers deployed in the cluster
func runPackagesShell(cmd *cobra.Command, args []string, container string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	tx, ctx := apm.StartTransaction("runPackagesShell", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	script := `oblt-shell-` + container
	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)

	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}

// runPackagesShellService open a shell to the service container
func runPackagesShellService(cmd *cobra.Command, args []string) {
	runPackagesShell(cmd, args, "service")
}

// runPackagesShellElasticAgent open a shell to the elastic-agent container
func runPackagesShellElasticAgent(cmd *cobra.Command, args []string) {
	runPackagesShell(cmd, args, "elastic-agent")
}

// runPackagesShellPackageRegistry open a shell to the package-registry container
func runPackagesShellPackageRegistry(cmd *cobra.Command, args []string) {
	runPackagesShell(cmd, args, "package-registry")
}

// runPackagesShellWorkspace open a shell to the workspace container
func runPackagesShellWorkspace(cmd *cobra.Command, args []string) {
	runPackagesShell(cmd, args, "workspace")
}
