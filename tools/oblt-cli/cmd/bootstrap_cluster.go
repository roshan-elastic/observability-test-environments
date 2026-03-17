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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
	"github.com/spf13/cobra"
)

// BootstrapClusterCmd represents the cluster bootstrap command
var BootstrapClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Command to bootstrap an ESS cluster.",
	Long:  "Command to bootstrap an ESS cluster using known recipes.",
	Run:   runBootstrapCluster,
}

func init() {
	BootstrapCmd.AddCommand(BootstrapClusterCmd)

	BootstrapClusterCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to bootstrap. (Required)")
	BootstrapClusterCmd.Flags().String(config.RecipesElasticsearchFlag, "", "Apply only the selected recipes to Elasticsearch.[\"test\", \"create-users\"]")
	BootstrapClusterCmd.Flags().String(config.RecipesKibanaFlag, "", "Apply only the selected recipes to Kibana.[\"test\"]")
	BootstrapClusterCmd.Flags().String(config.RecipesApmFlag, "", "Apply only the selected recipes to APM.[\"test\"]")
	BootstrapClusterCmd.Flags().String(config.RecipesFleetFlag, "", "Apply only the selected recipes to Fleet.[\"test\"]")
	BootstrapClusterCmd.Flags().Bool(config.IgnoreCertificatesFlag, false, "Disable TLS certificate verification.")
	BootstrapClusterCmd.Flags().String(config.BootstrapFolderFlag, "", "Full path to the bootstarp folder root. (Optional, default: REPO/bootstrap)")
	BootstrapClusterCmd.Flags().String(config.ParametersFlag, "", "Parameters values defined in JSON '{ \"var_name1\": \"value\",\"var_name2\": \"value\"}'. (Required)")

	cobra.MarkFlagRequired(BootstrapClusterCmd.Flags(), config.ClusterNameFlag)
}

func runBootstrapCluster(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runBootstrapCluster", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	recipesEsJson, _ := cmd.Flags().GetString(config.RecipesElasticsearchFlag)
	recipesKbnJson, _ := cmd.Flags().GetString(config.RecipesKibanaFlag)
	recipesApmJson, _ := cmd.Flags().GetString(config.RecipesApmFlag)
	recipesFleetJson, _ := cmd.Flags().GetString(config.RecipesFleetFlag)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	ignoreCertificates, _ := cmd.Flags().GetBool(config.IgnoreCertificatesFlag)
	bootstrapFolder, _ := cmd.Flags().GetString(config.BootstrapFolderFlag)
	parameters, err := loadParametersMap(cmd)
	config.CheckErr(err)

	authType := http.AuthTypeUser

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	esCred, err := gcsm.ReadEsSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	params := map[string]string{"username": esCred.Username, "password": esCred.Password}

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	if bootstrapFolder == "" {
		bootstrapFolder = obltTestEnvironments.GetBootstrapRecipesDir()
	}

	// All recipes
	var recipes = bootstrap.RecipesStruct{
		BootstrapType:      bootstrap.TypeElasticsearch,
		RecipesJson:        recipesEsJson,
		AuthType:           authType,
		AuthParams:         params,
		DryRun:             dryRun,
		Url:                esCred.Url,
		IgnoreCertificates: ignoreCertificates,
		BootstarpFolder:    bootstrapFolder,
		TemplateParams:     parameters,
	}
	resultsEs := recipes.Apply(obltTestEnvironments)

	kbnCreds, err := gcsm.ReadKibanaSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	params = map[string]string{http.ParamUsername: kbnCreds.Username, http.ParamPassword: kbnCreds.Password}
	recipes = bootstrap.RecipesStruct{
		BootstrapType:      bootstrap.TypeKibana,
		RecipesJson:        recipesApmJson,
		AuthType:           authType,
		AuthParams:         params,
		DryRun:             dryRun,
		Url:                kbnCreds.Url,
		IgnoreCertificates: ignoreCertificates,
		BootstarpFolder:    bootstrapFolder,
	}
	resultsKbn := recipes.Apply(obltTestEnvironments)

	apmCreds, err := gcsm.ReadApmSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	params = map[string]string{http.ParamToken: apmCreds.Token}
	recipes = bootstrap.RecipesStruct{
		BootstrapType:      bootstrap.TypeApm,
		RecipesJson:        recipesKbnJson,
		AuthType:           http.AuthTypeToken,
		AuthParams:         params,
		DryRun:             dryRun,
		Url:                apmCreds.Url,
		IgnoreCertificates: ignoreCertificates,
		BootstarpFolder:    bootstrapFolder,
	}
	resultsApm := recipes.Apply(obltTestEnvironments)

	fleetCreds, err := gcsm.ReadFleetSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	params = map[string]string{http.ParamApiKey: fleetCreds.ApiToken}
	recipes = bootstrap.RecipesStruct{
		BootstrapType:      bootstrap.TypeFleet,
		RecipesJson:        recipesFleetJson,
		AuthType:           http.AuthTypeApiKey,
		AuthParams:         params,
		DryRun:             dryRun,
		Url:                fleetCreds.Url,
		IgnoreCertificates: ignoreCertificates,
		BootstarpFolder:    bootstrapFolder,
	}
	resultsFleet := recipes.Apply(obltTestEnvironments)

	results := resultsEs
	for k, v := range resultsKbn {
		results[k] = v
	}

	for k, v := range resultsApm {
		results[k] = v
	}

	for k, v := range resultsFleet {
		results[k] = v
	}
	saveResults(results, outputFile)
}
