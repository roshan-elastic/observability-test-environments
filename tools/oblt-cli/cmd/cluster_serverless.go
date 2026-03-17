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
	"time"

	"github.com/spf13/cobra"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
)

// ServerlessCmd represents the serverless command
var ServerlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "Command to create an serverless cluster.",
	Long:  `Command to create a custom cluster using the serverless template. If used, any of the docker images must use the same version as the stack version, otherwise the build of the cluster will fail in Elastic Cloud`,
	Run:   runServerless,
}

func init() {
	CreateCmd.AddCommand(ServerlessCmd)

	ServerlessCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	ServerlessCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. This parameter ensures that the name of the cluster is unique, use it with caution. (Optional)")
	ServerlessCmd.Flags().String(config.ProjectTypeFlag, "", "Type of project to deploy [elasticsearch, observability, security] (default: observability) (Optional)")
	ServerlessCmd.Flags().String(config.ElasticsearchDockerImageFlag, "", " Docker image to use for Elasticsearch. This version should be the same build of the remote cluster if the version is an SNAPSHOT(Optional).")
	ServerlessCmd.Flags().String(config.KibanaDockerImageFlag, "", "Docker image to use for Kibana. This version should be the same build of the remote cluster if the version is an SNAPSHOT(Optional).")
	ServerlessCmd.Flags().String(config.FleetDockerImageFlag, "", "Docker image to use for Elastic Agent. This version should be the same build of the remote cluster if the version is an SNAPSHOT(Optional).")
	ServerlessCmd.Flags().String(config.TargetFlag, "", "Target environment to deploy (qa, staging, production) (default: qa)")
}

/*
It will create a cluster configuration file based on the parameters passed.
Then this cluster configuration file is commit and pushed to the oblt test environments repo.
*/
func runServerless(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	tx, ctx := apm.StartTransaction("runServerless", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	clusterNamePrefix, _ := cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ := cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	projectType, _ := cmd.Flags().GetString(config.ProjectTypeFlag)
	elasticsearchDockerImage, _ := cmd.Flags().GetString(config.ElasticsearchDockerImageFlag)
	kibanaDockerImage, _ := cmd.Flags().GetString(config.KibanaDockerImageFlag)
	fleetDockerImage, _ := cmd.Flags().GetString(config.FleetDockerImageFlag)
	target, _ := cmd.Flags().GetString(config.TargetFlag)

	if ciMode && clusterNamePrefix == "" {
		apm.CobraCheckErr(fmt.Errorf("in CI mode you must specify a prefix to identify the source"), tx, ctx)
	}

	wait, _ := cmd.Flags().GetInt(config.WaitFlag)

	essCluster := &clusters.ServerlessCluster{
		TemplateName:             clusters.ServerlessTemplateName,
		ClusterNamePrefix:        clusterNamePrefix,
		ClusterNameSuffix:        clusterNameSuffix,
		Username:                 userConfig.Username,
		SlackChannel:             userConfig.SlackChannel,
		ObltRepo:                 obltTestEnvironments,
		ProjectType:              projectType,
		Target:                   target,
		ElasticsearchDockerImage: elasticsearchDockerImage,
		KibanaDockerImage:        kibanaDockerImage,
		FleetDockerImage:         fleetDockerImage,
	}

	parametersMap, err := essCluster.Create()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}
