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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	hc "github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/healthcheck"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

var HealthcheckCommand = &cobra.Command{
	Use:   "healthcheck",
	Short: "prints cluster components health",
	Long: `prints cluster components health. \n
	Example: oblt-cli cluster healthcheck --cluster-name=[YOUR_CLUSTER_NAME]`,
	Run: runHealthcheck,
}

func init() {
	ClusterCmd.AddCommand(HealthcheckCommand)
	HealthcheckCommand.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	cobra.MarkFlagRequired(HealthcheckCommand.Flags(), config.ClusterNameFlag)
}

func runHealthcheck(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)

	labelClusterName := apm.Label{Key: "cluster-name", Value: clusterName}
	tx, ctx := apm.StartTransaction("runHealthcheck", "request", []apm.Label{labelVersion, labelClusterName}, userConfig)
	defer apm.Flush(tx)

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	clusterState, err := gcsm.ReadClusterStateSecret(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	healthchecks := []*hc.HealthcheckResponse{
		hc.ElasticsearchHC(clusterState),
		hc.KibanaHC(clusterState),
		hc.ApmHC(clusterState),
	}

	checklist := []*hc.HealthcheckResponse{}

	for _, hc := range healthchecks {
		if hc != nil {
			checklist = append(checklist, hc)
		}
	}

	for _, hcResponse := range checklist {
		logger.Infof(hc.PrettyPrintHC(hcResponse))
	}

	saveResults(hc.ResponsesToJson(checklist), outputFile)
}
