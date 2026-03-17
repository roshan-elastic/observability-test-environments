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
	"fmt"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/k8s"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// PackagesLogsCmd represents the packages logs command
var PackagesLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Command to show the logs from the containers deployed in the cluster.",
	Long:  `Command to show the logs from the containers deployed in the cluster.`,
}

// PackagesLogsServiceCmd represents the packages logs service command
var PackagesLogsServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Command to show the service containers logs.",
	Long:  `Command to show the service containers logs.`,
	Run:   runPackagesLogsService,
}

// PackagesLogsElasticAgentCmd represents the packages logs elastic-agent command
var PackagesLogsElasticAgentCmd = &cobra.Command{
	Use:   "elastic-agent",
	Short: "Command to show the Elastic Agent containers logs.",
	Long:  `Command to show the Elastic Agent containers logs.`,
	Run:   runPackagesLogsElasticAgent,
}

// PackagesLogsPackageRegistryCmd represents the packages logs package-registry command
var PackagesLogsPackageRegistryCmd = &cobra.Command{
	Use:   "package-registry",
	Short: "Command to show the package registry containers logs.",
	Long:  `Command to show the package registry containers logs.`,
	Run:   runPackagesLogsPackageRegistry,
}

// PackagesLogsWorkspaceCmd represents the packages logs workspace command
var PackagesLogsWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Command to show the workspace containers logs.",
	Long:  `Command to show the workspace containers logs.`,
	Run:   runPackagesLogsWorkspace,
}

// PackagesLogsWorkspaceCmd represents the packages logs workspace command
var PackagesLogsRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Command to sync a local folder from the workspace container logs folder.",
	Long: `Command to sync a local folder from the workspace container logs folder.
		The command will sync a local folder with the logs from the workspace container every 30 seconds.
		It uses kubectl + rsync to sync the files.`,
	Run: runPackagesLogsRemote,
}

func init() {
	PackagesCmd.AddCommand(PackagesLogsCmd)
	PackagesLogsCmd.AddCommand(PackagesLogsServiceCmd)
	PackagesLogsCmd.AddCommand(PackagesLogsElasticAgentCmd)
	PackagesLogsCmd.AddCommand(PackagesLogsPackageRegistryCmd)
	PackagesLogsCmd.AddCommand(PackagesLogsWorkspaceCmd)
	PackagesLogsCmd.AddCommand(PackagesLogsRemoteCmd)

	PackagesLogsCmd.PersistentFlags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	PackagesLogsRemoteCmd.PersistentFlags().String(config.OutputFolderFlag, "", "Folder to sync from the remote folder. (Required)")

	cobra.MarkFlagRequired(PackagesLogsCmd.PersistentFlags(), config.ClusterNameFlag)
	cobra.MarkFlagRequired(PackagesLogsRemoteCmd.PersistentFlags(), config.OutputFolderFlag)
}

// runLogs show the logs a container deployed in the cluster
func runPackagesLogs(cmd *cobra.Command, args []string, script, request string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	tx, ctx := apm.StartTransaction(request, "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)
	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}

// runPackagesLogsService show the logs from the service container deployed in the cluster
func runPackagesLogsService(cmd *cobra.Command, args []string) {
	runPackagesLogs(cmd, args, "oblt-logs-service", "runPackagesLogsService")
}

// runPackagesLogsElasticAgent show the logs from the elastic-agent container deployed in the cluster
func runPackagesLogsElasticAgent(cmd *cobra.Command, args []string) {
	runPackagesLogs(cmd, args, "oblt-logs-elastic-agent", "runPackagesLogsElasticAgent")
}

// runPackagesLogsPackageRegistry show the logs from the package-registry container deployed in the cluster
func runPackagesLogsPackageRegistry(cmd *cobra.Command, args []string) {
	runPackagesLogs(cmd, args, "oblt-logs-package-registry", "runPackagesLogsPackageRegistry")
}

// runPackagesLogsWorkspace show the logs from the workspace container deployed in the cluster
func runPackagesLogsWorkspace(cmd *cobra.Command, args []string) {
	runPackagesLogs(cmd, args, "oblt-logs-workspace", "runPackagesLogsWorkspace")
}

func runPackagesLogsRemote(cmd *cobra.Command, args []string) {
	localFolder, _ := cmd.Flags().GetString(config.OutputFolderFlag)
	script := fmt.Sprintf("oblt-ksync-task %s", localFolder)
	runPackagesLogs(cmd, args, script, "runPackagesLogsRemote")
}
