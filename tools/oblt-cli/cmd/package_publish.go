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
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/k8s"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"

	"github.com/spf13/cobra"
)

// PackagesPublishCmd represents the packages publish command
var PackagesPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Command to publish the integration package to the cluster.",
	Long: `Command to publish the integration package to the cluster.
	The package will be build and published to the cluster.`,
	Run: runPackagesPublish,
}

func init() {
	PackagesCmd.AddCommand(PackagesPublishCmd)

	PackagesPublishCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	PackagesPublishCmd.Flags().String(config.PackageFolderFlag, "", "Absolute path to the source code of the package. (Required)")

	cobra.MarkFlagRequired(PackagesPublishCmd.Flags(), config.ClusterNameFlag)
	cobra.MarkFlagRequired(PackagesPublishCmd.Flags(), config.PackageFolderFlag)
}

// runPublish publish the integration package to the cluster
func runPackagesPublish(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	packageFolder, _ := cmd.Flags().GetString(config.PackageFolderFlag)

	tx, ctx := apm.StartTransaction("runPackagesPublish", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	esSecrets, err := gcsm.ReadEsSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)
	kbnSecrets, err := gcsm.ReadKibanaSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	os.Setenv("ELASTIC_PACKAGE_ELASTICSEARCH_HOST", esSecrets.Url)
	os.Setenv("ELASTIC_PACKAGE_ELASTICSEARCH_USERNAME", esSecrets.Username)
	os.Setenv("ELASTIC_PACKAGE_ELASTICSEARCH_PASSWORD", esSecrets.Password)
	os.Setenv("ELASTIC_PACKAGE_KIBANA_HOST", kbnSecrets.Url)
	os.Setenv("PACKAGE_FOLDER", packageFolder)

	script := `
	oblt-publish
	exit 0`
	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)
	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}
