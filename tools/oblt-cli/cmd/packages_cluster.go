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
	"strings"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// PackagesClusterCmd represents the cluster command
var PackagesClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Command to manage oblt clusters.",
	Long: `Command to manage oblt clusters.
	With this command you can create oblt clusters running integrations.
	`,
}

var PackagesClusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Command to create an oblt cluster.",
	Long: `Command to create an oblt cluster.
	With this command you can create an oblt cluster running the Elastic Stack and integrations.`,
	Run: runPackagesCreate,
}

var PackagesClusterServiceCmd = &cobra.Command{
	Use:   "create-service",
	Short: "Command to create an oblt cluster that runs a service.",
	Long: `Command to create an oblt cluster that runs a service.
	With this command you can create an oblt cluster running integrations (no Stack).
	`,
	Run: runPackagesService,
}

// This allow to generate the documentation for both commands
var PackagesClusterSecretsCmd = &cobra.Command{
	Use:   ClusterSecretsCmd.Use,
	Short: ClusterSecretsCmd.Short,
	Long:  ClusterSecretsCmd.Long,
	Run:   ClusterSecretsCmd.Run,
}

var PackagesDestroyCmd = &cobra.Command{
	Use:   DestroyCmd.Use,
	Short: DestroyCmd.Short,
	Long:  DestroyCmd.Long,
	Run:   DestroyCmd.Run,
}

var PackagesClusterListCmd = &cobra.Command{
	Use:   ListCmd.Use,
	Short: ListCmd.Short,
	Long:  ListCmd.Long,
	Run:   ListCmd.Run,
}

func init() {
	PackagesClusterCmd.AddCommand(PackagesClusterCreateCmd)
	PackagesClusterCmd.AddCommand(PackagesClusterSecretsCmd)
	PackagesClusterCmd.AddCommand(PackagesDestroyCmd)
	PackagesClusterCmd.AddCommand(PackagesClusterListCmd)
	PackagesClusterCmd.AddCommand(PackagesClusterServiceCmd)
	PackagesCmd.AddCommand(PackagesClusterCmd)

	PackagesClusterCreateCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	PackagesClusterCreateCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. This parameter ensures that the name of the cluster is unique, use it with caution. (Optional)")
	PackagesClusterCreateCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to create. This is only supported when running on the CI to help with redeploying clusters. (Optional)")
	PackagesClusterCreateCmd.Flags().String(config.StackVersionFlag, "", "Stack version to use for the deployment. (Required)")
	PackagesClusterCreateCmd.Flags().Bool(config.IsReleaseFlag, false, "True is the Elastic Stack version is a release. The default value is false. (Optional)")
	PackagesClusterCreateCmd.Flags().String(config.IntegrationFlag, "", "Integration to be deployed. (Required)")
	PackagesClusterCreateCmd.Flags().String(config.RepositoryFlag, "elastic/integrations", "Integrations repository to use for the `_dev` integration source. (elastic/integrations by default)(Optional)")
	PackagesClusterCreateCmd.Flags().String(config.BranchFlag, "main", "Branch to use for the `_dev` integration source. (main by default)(Optional)")

	PackagesClusterServiceCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	PackagesClusterServiceCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. This parameter ensures that the name of the cluster is unique, use it with caution. (Optional)")
	PackagesClusterServiceCmd.Flags().String(config.IntegrationFlag, "", "Integration to be deployed. (Required)")
	PackagesClusterServiceCmd.Flags().String(config.RepositoryFlag, "elastic/integrations", "Integrations repository to use for the `_dev` integration source. (elastic/integrations by default)(Optional)")
	PackagesClusterServiceCmd.Flags().String(config.BranchFlag, "main", "Branch to use for the `_dev` integration source. (main by default)(Optional)")

	cobra.MarkFlagRequired(PackagesClusterCreateCmd.Flags(), config.StackVersionFlag)
	cobra.MarkFlagRequired(PackagesClusterCreateCmd.Flags(), config.IntegrationFlag)

	cobra.MarkFlagRequired(PackagesClusterServiceCmd.Flags(), config.IntegrationFlag)
}

// runPackagesCreate is the function executed to create a cluster
func runPackagesCreate(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	var clusterName string
	var clusterNamePrefix string
	var clusterNameSuffix string
	var stackVersion string
	var isRelease bool
	var integration string
	var repository string
	var branch string

	tx, ctx := apm.StartTransaction("runPackagesCreate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)

	clusterName, _ = cmd.Flags().GetString(config.ClusterNameFlag)
	clusterNamePrefix, _ = cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ = cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	stackVersion, _ = cmd.Flags().GetString(config.StackVersionFlag)
	isRelease, _ = cmd.Flags().GetBool(config.IsReleaseFlag)
	integration, _ = cmd.Flags().GetString(config.IntegrationFlag)
	repository, _ = cmd.Flags().GetString(config.RepositoryFlag)
	branch, _ = cmd.Flags().GetString(config.BranchFlag)

	wait, _ := cmd.Flags().GetInt(config.WaitFlag)

	if strings.EqualFold(stackVersion, "") {
		apm.CobraCheckErr(fmt.Errorf(`required "stack-version" not set`), tx, ctx)
	}

	clusterInfo := &clusters.IntegrationCluster{
		ClusterName:       clusterName,
		ClusterNamePrefix: clusterNamePrefix,
		ClusterNameSuffix: clusterNameSuffix,
		Integration:       integration,
		Repository:        repository,
		Branch:            branch,
		StackVersion:      stackVersion,
		IsRelease:         isRelease,
		Username:          userConfig.Username,
		SlackChannel:      userConfig.SlackChannel,
		ObltRepo:          obltTestEnvironments,
	}

	parametersMap, err := clusterInfo.Create()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}

// runPackagesService is the function executed to create a service cluster
func runPackagesService(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	var clusterNamePrefix string
	var clusterNameSuffix string
	var integration string
	var repository string
	var branch string

	tx, ctx := apm.StartTransaction("runPackagesCreate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)

	clusterNamePrefix, _ = cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ = cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	integration, _ = cmd.Flags().GetString(config.IntegrationFlag)
	repository, _ = cmd.Flags().GetString(config.RepositoryFlag)
	branch, _ = cmd.Flags().GetString(config.BranchFlag)

	wait, _ := cmd.Flags().GetInt(config.WaitFlag)

	integrationClusters := &clusters.IntegrationCluster{
		ClusterNamePrefix: clusterNamePrefix,
		ClusterNameSuffix: clusterNameSuffix,
		Integration:       integration,
		Repository:        repository,
		Branch:            branch,
		Username:          userConfig.Username,
		SlackChannel:      userConfig.SlackChannel,
		ObltRepo:          obltTestEnvironments,
	}
	parametersMap, err := integrationClusters.CreateService()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}
