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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/bootstrap"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
	"github.com/spf13/cobra"
)

// BootstrapKibanaCmd represents the kibana bootstrap subcommand
var BootstrapKibanaCmd = &cobra.Command{
	Use:   "kibana",
	Short: "Command to bootstrap Kibana.",
	Long:  "Command to bootstrap Kibana using known recipes.",
	Run:   runBootstrapKbn,
}

func init() {
	BootstrapCmd.AddCommand(BootstrapKibanaCmd)

	BootstrapKibanaCmd.Flags().String(config.UrlFlag, "", "URL of the Kibana service to bootstrap.(Required)")
	BootstrapKibanaCmd.Flags().String(config.UsernameFlag, "", "Username to authenticate. (Not needed if you provide "+config.ApiKeyFlag+")")
	BootstrapKibanaCmd.Flags().String(config.PasswordFlag, "", "Password to authenticate. (Not neededif you provide "+config.ApiKeyFlag+")")
	BootstrapKibanaCmd.Flags().String(config.ApiKeyFlag, "", "API Key to authenticate. (Not needed if you provide "+config.UsernameFlag+")")
	BootstrapKibanaCmd.Flags().String(config.RecipesFlag, "", "Apply only the selected recipes.([\"ml-logs-ui-categories\"])")
	BootstrapKibanaCmd.Flags().Bool(config.IgnoreCertificatesFlag, false, "Disable TLS certificate verification.")
	BootstrapKibanaCmd.Flags().String(config.BootstrapFolderFlag, "", "Full path to the bootstarp folder root. (Optional, default: REPO/bootstrap)")
	BootstrapKibanaCmd.Flags().String(config.ParametersFlag, "", "Parameters values defined in JSON '{ \"var_name1\": \"value\",\"var_name2\": \"value\"}'. (Required)")

	cobra.MarkFlagRequired(BootstrapKibanaCmd.Flags(), config.UrlFlag)
}

func runBootstrapKbn(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runBootstrapEs", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	recipesJson, _ := cmd.Flags().GetString(config.RecipesFlag)
	kibanaUrl, _ := cmd.Flags().GetString(config.UrlFlag)
	username, _ := cmd.Flags().GetString(config.UsernameFlag)
	password, _ := cmd.Flags().GetString(config.PasswordFlag)
	apiKey, _ := cmd.Flags().GetString(config.ApiKeyFlag)
	ignoreCertificates, _ := cmd.Flags().GetBool(config.IgnoreCertificatesFlag)
	bootstrapFolder, _ := cmd.Flags().GetString(config.BootstrapFolderFlag)
	parameters, err := loadParametersMap(cmd)
	config.CheckErr(err)

	authType, params, err := http.ChooseAuthType(username, password, apiKey, "")
	apm.CobraCheckErr(err, tx, ctx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)
	if bootstrapFolder == "" {
		bootstrapFolder = obltTestEnvironments.GetBootstrapRecipesDir()
	}
	var recipes = &bootstrap.RecipesStruct{
		BootstrapType:      bootstrap.TypeKibana,
		RecipesJson:        recipesJson,
		AuthType:           authType,
		AuthParams:         params,
		DryRun:             dryRun,
		Url:                kibanaUrl,
		IgnoreCertificates: ignoreCertificates,
		BootstarpFolder:    bootstrapFolder,
		TemplateParams:     parameters,
	}
	results := recipes.Apply(obltTestEnvironments)
	saveResults(results, outputFile)
}
