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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/mki"
	"github.com/spf13/cobra"
)

// MkiCmd represents the create command
var MkiCmd = &cobra.Command{
	Use:   "mki",
	Short: "Command to operate to the Elastic MKI environment.",
	Long:  `Command to operate to the Elastic MKI environment.`,
}

var MkiLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Command to login to the Elastic MKI environment.",
	Long: `Command to login to the Elastic MKI environment.
This operation is performed before to access perform any operation with Teleport.

The following example connect to the staging MKI environment:

	oblt-cli services mki login --cluster-name "keep_serverless-qa-oblt" --environment "qa"

For more information about the MKI environment, please check https://docs.elastic.dev/mki/user-guide/cluster-access`,
	Args: validateEnvironmentFlag,
	Run:  runMkiLogin,
}

var MkiGenKibanaYamlCmd = &cobra.Command{
	Use:   "generate-kibana-yaml",
	Short: "Command to generata a Kibana yaml config file to connet to a serverless project.",
	Long: `Command to generata a Kibana yaml config file to connet to a serverless project.
The Kibana YAML generated allows to connect a local Kibana instance to the serverless project.

The following example connect to the staging MKI environment:

	oblt-cli services mki generate-kibana-yaml --cluster-name "servreless-qa-oblt" --kibana-yaml-path "${PWD}/kibana.serverless.yaml"

For more information about the MKI environment, please check https://docs.elastic.dev/mki/user-guide/cluster-access`,
	Args: validateEnvironmentFlag,
	Run:  runMkiGenKibanaYaml,
}

var MkiRunLocalKibanaCmd = &cobra.Command{
	Use:   "run-local-kibana",
	Short: "Command to generata a Kibana yaml config file and run a local Kibana instance connected to a serverless project.",
	Long: `Command to generata a Kibana yaml config file and run a local Kibana instance connected to a serverless project.
The Kibana YAML generated allows to connect a local Kibana instance to the serverless project.

The following example connect to the qa MKI environment:

	oblt-cli services mki generate-kibana-yaml --cluster-name "keep_serverless-qa-oblt" --kibana-yaml-path "${PWD}/kibana.serverless.yaml" --kibana-src "${HOME}/src/kibana"

For more information about the MKI environment, please check https://docs.elastic.dev/mki/user-guide/cluster-access`,
	Args: validateEnvironmentFlag,
	Run:  runRunLocalKibana,
}

func init() {
	ServicesCmd.AddCommand(MkiCmd)
	MkiCmd.AddCommand(MkiLoginCmd)
	MkiCmd.AddCommand(MkiGenKibanaYamlCmd)
	MkiCmd.AddCommand(MkiRunLocalKibanaCmd)

	MkiLoginCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	MkiLoginCmd.Flags().Bool(config.VPNFlag, true, "Start the VPN connection, (default true)")
	MkiLoginCmd.Flags().String(config.ApiKeyFlag, "", "API key to use to login to the ESS Admin console.")

	cobra.MarkFlagRequired(MkiLoginCmd.Flags(), config.ClusterNameFlag)

	MkiGenKibanaYamlCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	MkiGenKibanaYamlCmd.Flags().String(config.KibanaYamlPathFlag, "", "Path to the Kibana yaml file to generate. (Required)")
	MkiGenKibanaYamlCmd.Flags().Bool(config.VPNFlag, true, "Start the VPN connection (default true)")
	MkiGenKibanaYamlCmd.Flags().Bool(config.MKILoginFlag, true, "Perform the Teleport login in the MKI environment. (default true)")
	MkiGenKibanaYamlCmd.Flags().String(config.ApiKeyFlag, "", "API key to use to login to the ESS Admin console.")

	cobra.MarkFlagRequired(MkiGenKibanaYamlCmd.Flags(), config.ClusterNameFlag)
	cobra.MarkFlagRequired(MkiGenKibanaYamlCmd.Flags(), config.KibanaYamlPathFlag)

	MkiRunLocalKibanaCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	MkiRunLocalKibanaCmd.Flags().String(config.KibanaYamlPathFlag, "", "Path to the Kibana yaml file to generate. (Required)")
	MkiRunLocalKibanaCmd.Flags().String(config.KibanaSrcPathFlag, "", "Path to the Kibana yaml file to generate. (Required)")
	MkiRunLocalKibanaCmd.Flags().Bool(config.VPNFlag, false, "Start the VPN connection (default true)")
	MkiRunLocalKibanaCmd.Flags().Bool(config.MKILoginFlag, true, "Perform the Teleport login in the MKI environment. (default true)")
	MkiRunLocalKibanaCmd.Flags().String(config.ApiKeyFlag, "", "API key to use to login to the ESS Admin console.")

	cobra.MarkFlagRequired(MkiRunLocalKibanaCmd.Flags(), config.ClusterNameFlag)
	cobra.MarkFlagRequired(MkiRunLocalKibanaCmd.Flags(), config.KibanaYamlPathFlag)
	cobra.MarkFlagRequired(MkiRunLocalKibanaCmd.Flags(), config.KibanaSrcPathFlag)
}

/*
runMkiLogin Login to the Elastic MKI environment,
*/
func runMkiLogin(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runMkiLogin", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	startVpn, _ := cmd.Flags().GetBool(config.VPNFlag)
	apiKey, _ := cmd.Flags().GetString(config.ApiKeyFlag)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)

	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	mkiClient := mki.NewMKIConfig(config, userConfig.GetDir(), obltTestEnvironments.GetPath(), startVpn, true, dryRun, apiKey)
	err = mkiClient.Login()
	apm.CobraCheckErr(err, tx, ctx)
}

// runMkiGenKibanaYaml generates a Kibana yaml config file to connet to a serverless project.
func runMkiGenKibanaYaml(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runMkiGenKibanaYaml", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	kibanaYamlPath, _ := cmd.Flags().GetString(config.KibanaYamlPathFlag)
	startVpn, _ := cmd.Flags().GetBool(config.VPNFlag)
	mkiLogin, _ := cmd.Flags().GetBool(config.MKILoginFlag)
	apiKey, _ := cmd.Flags().GetString(config.ApiKeyFlag)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	mkiClient := mki.NewMKIConfig(config, userConfig.GetDir(), obltTestEnvironments.GetPath(), startVpn, mkiLogin, dryRun, apiKey)
	err = mkiClient.GenerateKibanaConfig(kibanaYamlPath)
	apm.CobraCheckErr(err, tx, ctx)
}

func runRunLocalKibana(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runRunLocalKibana", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	kibanaYamlPath, _ := cmd.Flags().GetString(config.KibanaYamlPathFlag)
	kibanaSrcPath, _ := cmd.Flags().GetString(config.KibanaSrcPathFlag)
	startVpn, _ := cmd.Flags().GetBool(config.VPNFlag)
	mkiLogin, _ := cmd.Flags().GetBool(config.MKILoginFlag)
	apiKey, _ := cmd.Flags().GetString(config.ApiKeyFlag)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	mkiClient := mki.NewMKIConfig(config, userConfig.GetDir(), obltTestEnvironments.GetPath(), startVpn, mkiLogin, dryRun, apiKey)
	err = mkiClient.RunLocalKibana(kibanaYamlPath, kibanaSrcPath)
	apm.CobraCheckErr(err, tx, ctx)
}
