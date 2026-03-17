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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/spf13/cobra"
)

// ClusterLicenseCmd List the recipes available.
var ClusterLicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Command to push a license to Cluster.",
	Long:  "Command to push a license to Cluster.",
	Run:   runClusterLicense,
}

func init() {
	ClusterCmd.AddCommand(ClusterLicenseCmd)
	ClusterLicenseCmd.Flags().String(config.TypeFlag, "", "Type of license to push [release, dev, orchestration, orchestration-dev]. (Required)")
	ClusterLicenseCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	ClusterLicenseCmd.Flags().String(config.EnvironmentFlag, "", "Environment to use [production, staging]. (Default production)")

	cobra.MarkFlagRequired(ClusterLicenseCmd.Flags(), config.ClusterNameFlag)

	cobra.MarkFlagRequired(ClusterLicenseCmd.Flags(), config.TypeFlag)
}

func runClusterLicense(cmd *cobra.Command, args []string) {
	licenseType, _ := cmd.Flags().GetString(config.TypeFlag)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	environment, _ := cmd.Flags().GetString(config.EnvironmentFlag)

	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runClusterLicense", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	essCred, err := gcsm.ReadESSCredentials(environment)
	apm.CobraCheckErr(err, tx, ctx)

	essDeployment, err := gcsm.ReadESSDeployment(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	url := fmt.Sprintf("%s/api/v1/deployments/%s/elasticsearch/main-elasticsearch/proxy", essCred.Url, essDeployment.Id)
	runLicenseCommon(userConfig, "runLicense", licenseType, "", "", essCred.ApiKey, url, false)
}
