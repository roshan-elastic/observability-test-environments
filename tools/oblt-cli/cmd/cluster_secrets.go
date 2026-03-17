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

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/prompt"
)

// ClusterSecretsCmd represents the secrets command
var ClusterSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Command to read the secrets from a cluster.",
	Long:  `Command to read the secrets from a cluster. It will use Google Secrets Manager to read the secrets, therefore accessing a secret will depend on Google Secrets Manager.`,
	Run:   runSecrets,
}

// ClusterSecretInfoCmd represents the command to retrieve the deployment info
var ClusterSecretInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Command to read the deploy-info secret from a cluster.",
	Long:  `Command to read the deploy-info secret from a cluster. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.`,
	Run:   runReadSecretInfo,
}

// ClusterSecretKibanaCmd represents the secrets command to retrieve the kibana config file
var ClusterSecretKibanaCmd = &cobra.Command{
	Use:   "kibana-config",
	Short: "Command to read the kibana-yml secret from a cluster.",
	Long:  `Command to read the kibana-yml secret from a cluster, retrieving Kibana's configuration file. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.`,
	Run:   runReadKibanaYaml,
}

// ClusterSecretUsersCmd represents the secrets command to retrieve the credentials of the cluster
var ClusterSecretUsersCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Command to read the credentials secret from a cluster.",
	Long:  `Command to read the credentials secret from a cluster, retrieving user credentials and URLs to interact with the cluster. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.`,
	Run:   runReadCredentialsSecret,
}

var ClusterSecretEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Command to read the Elastic Stack environment variables secret from a cluster.",
	Long:  `Command to read the Elastic Stack environment variables secret from a cluster, Elastic Stack environment variables to interact with the cluster. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.`,
	Run:   runReadEnvSecret,
}

var ClusterSecretStateCmd = &cobra.Command{
	Use:   "cluster-state",
	Short: "Command to read the Cluster state YAML secret from a cluster.",
	Long:  `Command to read the Cluster state YAML secret from a cluster. The contains credentials and endpoints to interact with the cluster. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.`,
	Run:   runReadClusterStateSecret,
}

var ClusterSecretActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Command to read the Activate local environment script secret from a cluster.",
	Long:  `Command to read the local environment script secret from a cluster. It contains a script to initialice a shell to connect to the k8s cluster to operate with it. It will use Google Secrets Manager to read the secret, therefore accessing the secret will depend on Google Secrets Manager.This script is only created for cluster with k8s deployments. You can activate the environment by running the script ". ./activate", then you can use k8s commands.`,
	Run:   runReadActivateSecret,
}

func init() {
	subcommands := []*cobra.Command{
		ClusterSecretInfoCmd, ClusterSecretKibanaCmd, ClusterSecretUsersCmd, ClusterSecretEnvCmd, ClusterSecretStateCmd, ClusterSecretActivateCmd,
	}

	for _, subCmd := range subcommands {
		ClusterSecretsCmd.AddCommand(subCmd)

		subCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster name. (Required)")
		cobra.MarkFlagRequired(subCmd.Flags(), config.ClusterNameFlag)
	}

	ClusterCmd.AddCommand(ClusterSecretsCmd)
	ClusterSecretsCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster name. (Required)")
	cobra.MarkFlagRequired(ClusterSecretsCmd.Flags(), config.ClusterNameFlag)
}

// runSecrets is the main function for the secrets command
func runSecrets(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runSecrets", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secrets, err := gcsm.ListClusterSecrets(clusterName, false)
	apm.CobraCheckErr(err, tx, ctx)

	secret := prompt.ClusterSecrets(secrets)

	switch secret {
	case "deploy-info":
		runReadSecretInfo(cmd, args)
	case "kibana-yml":
		runReadKibanaYaml(cmd, args)
	case "credentials":
		runReadCredentialsSecret(cmd, args)
	case "env":
		runReadEnvSecret(cmd, args)
	case "cluster-state":
		runReadClusterStateSecret(cmd, args)
	case "activate":
		runReadActivateSecret(cmd, args)
	}
}

// runReadSecretInfo reads the deploy-info secret from the cluster
func runReadSecretInfo(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadSecretInfo", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadDeployInfoSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(secretValue)
}

// runReadKibanaYaml reads the kibana-yml secret from the cluster
func runReadKibanaYaml(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadKibanaYaml", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadKibanaYamlSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(secretValue)
}

// runReadCredentialsSecret reads the credentials secret from the cluster
func runReadCredentialsSecret(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadCredentialsSecret", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadCredentialsSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(secretValue)
}

// runReadEnvSecret reads the Elastic Stack environment variables secret from the cluster
func runReadEnvSecret(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadEnvSecret", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadEnvSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(secretValue)
}

// runReadClusterStateSecret reads the Cluster state YAML secret from the cluster
func runReadClusterStateSecret(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadClusterStateSecret", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadClusterStateSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	str, err := yaml.Marshal(secretValue)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(string(str))
}

// runReadActivateSecret reads the Activate local environment script secret from the cluster
func runReadActivateSecret(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runReadActivateSecret", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	secretValue, err := gcsm.ReadActivateSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	writeOutput(secretValue)
}

// writeOutput writes the secret value to the output file or to the console.
func writeOutput(secretValue string) {
	secretOut := fmt.Sprintf("%s\n", secretValue)
	if len(outputFile) == 0 {
		fmt.Printf("%s", secretOut)
	} else {
		saveResultsRaw(secretOut, outputFile)
	}
}
