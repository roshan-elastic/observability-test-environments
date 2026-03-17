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
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// CustomCmd represents the create command
var CustomCmd = &cobra.Command{
	Use:   "custom",
	Short: "Command to create a cluster from a template.",
	Long: `Command to create a cluster from a template.

The following example create a cluster using the 'ccs' template:

	oblt-cli cluster create custom --template ccs --cluster-name-prefix 'prefix' --parameters '{ "ClusterName": "oblt", "StackVersion": "7.15.0", "SlackChannel": "#foo" }'`,
	Args: validateCreateCustomFlags,
	Run:  runCustom,
}

func init() {
	CreateCmd.AddCommand(CustomCmd)

	CustomCmd.Flags().String(config.TemplateNameFlag, "", "Template name to use, templates can be listed with 'oblt-cli cluster templates'. (Required if no template file is provided)")
	CustomCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	CustomCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. (Optional)")
	CustomCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to create. This is only supported when running on the CI to help with redeploying clusters. (Optional)")
	CustomCmd.Flags().String(config.ParametersFlag, "", "Parameters values defined in JSON '{ \"var_name1\": \"value\",\"var_name2\": \"value\"}'. (Required)")
	CustomCmd.Flags().String(config.TemplateFileFlag, "", "Absolute template file path to use. (Required if no template name is provided)")
	CustomCmd.Flags().String(config.ParametersFileFlag, "", "Absolute parameters JSON file path to use. It is incompatible with \""+config.ParametersFlag+"\" (Optional)")
}

/*
runCustom process a oblt cluster template with the parameters values passed,
then push a new oblt cluster configuration file to the repository.
*/
func runCustom(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runCustom", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	templateName, _ := cmd.Flags().GetString(config.TemplateNameFlag)
	clusterNamePrefix, _ := cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ := cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	templateFilePath, _ := cmd.Flags().GetString(config.TemplateFileFlag)
	wait, _ := cmd.Flags().GetInt(config.WaitFlag)
	parameters, err := loadParameters(cmd)

	apm.CobraCheckErr(err, tx, ctx)

	err = validateCIMinimumArguments(cmd, args)
	apm.CobraCheckErr(err, tx, ctx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	customCluster := &clusters.CustomCluster{
		TemplateName:      templateName,
		ClusterName:       clusterName,
		ClusterNamePrefix: clusterNamePrefix,
		ClusterNameSuffix: clusterNameSuffix,
		TemplatePath:      templateFilePath,
		Parameters:        parameters,
		Username:          userConfig.Username,
		SlackChannel:      userConfig.SlackChannel,
		ObltRepo:          obltTestEnvironments,
	}

	parametersMap, err := customCluster.Create()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}
