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

	"github.com/spf13/cobra"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/artifacts"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/prompt"
)

// EssCmd represents the ess command
var EssCmd = &cobra.Command{
	Use:   "ess",
	Short: "Command to create an ESS cluster.",
	Long:  `Command to create a custom cluster using the ESS template. If used, any of the docker images must use the same version as the stack version, otherwise the build of the cluster will fail in Elastic Cloud`,
	Run:   runEss,
}

func init() {
	CreateCmd.AddCommand(EssCmd)

	EssCmd.Flags().BoolVar(&interactive, config.InteractiveFlag, false, "If true, it will ask for the different versions in a interactive way.")
	EssCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	EssCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. This parameter ensures that the name of the cluster is unique, use it with caution. (Optional)")
	EssCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to create. This is only supported when running on the CI to help with redeploying clusters. (Optional)")
	EssCmd.Flags().String(config.StackVersionFlag, "", "Stack version to use for the deployment. (Required)")
	EssCmd.Flags().Bool(config.IsReleaseFlag, false, "True is the Elastic Stack version is a release. The default value is false. (Optional)")

	cobra.MarkFlagRequired(EssCmd.Flags(), config.StackVersionFlag)
}

/*
It will create a cluster configuration file based on the parameters passed.
Then this cluster configuration file is commit and pushed to the oblt test environments repo.
*/
func runEss(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	var clusterName = ""
	var clusterNamePrefix string
	var clusterNameSuffix string
	var stackVersion string
	var dockerImageVersion string
	var isRelease bool

	tx, ctx := apm.StartTransaction("runEss", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	if interactive {
		versions, err := artifacts.GetVersions()
		apm.CobraCheckErr(err, tx, ctx)

		clusterNamePrefix, clusterNameSuffix, stackVersion, dockerImageVersion = prompt.EssConfig(versions)

		if strings.HasPrefix(dockerImageVersion, "---") {
			// keep the ability to not set the Docker images
			dockerImageVersion = ""
		}
	} else {
		clusterName, _ = cmd.Flags().GetString(config.ClusterNameFlag)
		clusterNamePrefix, _ = cmd.Flags().GetString(config.ClusterNamePrefixFlag)
		clusterNameSuffix, _ = cmd.Flags().GetString(config.ClusterNameSuffixFlag)
		stackVersion, _ = cmd.Flags().GetString(config.StackVersionFlag)
		isRelease, _ = cmd.Flags().GetBool(config.IsReleaseFlag)
	}

	err = validateCIMinimumArguments(cmd, args)
	apm.CobraCheckErr(err, tx, ctx)

	wait, _ := cmd.Flags().GetInt(config.WaitFlag)

	if strings.EqualFold(stackVersion, "") {
		apm.CobraCheckErr(fmt.Errorf(`required "stack-version" not set`), tx, ctx)
	}

	essCluster := &clusters.ESSCluster{
		TemplateName:      clusters.ESSTemplateName,
		ClusterName:       clusterName,
		ClusterNamePrefix: clusterNamePrefix,
		ClusterNameSuffix: clusterNameSuffix,
		Username:          userConfig.Username,
		SlackChannel:      userConfig.SlackChannel,
		ObltRepo:          obltTestEnvironments,
		StackVersion:      stackVersion,
		IsRelease:         isRelease,
	}

	parametersMap, err := essCluster.Create()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}
