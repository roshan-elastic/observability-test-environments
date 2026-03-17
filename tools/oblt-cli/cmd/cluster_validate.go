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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/schema"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// ClusterValidateCmd List the recipes available.
var ClusterValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the a cluster configuration file.",
	Long:  "Validate the a cluster configuration file.",
	Run:   runClusterValidate,
}

func init() {
	ClusterCmd.AddCommand(ClusterValidateCmd)
	ClusterValidateCmd.Flags().String(config.ConfigFileFlag, "", "Configuration file to use. (Required)")

	cobra.MarkFlagRequired(ClusterValidateCmd.Flags(), config.ConfigFileFlag)
}

func runClusterValidate(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString(config.ConfigFileFlag)

	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, _ := apm.StartTransaction("runClusterValidate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	err = schema.Validate(configFile)
	config.CheckErr(err)

	logger.Infof("Configuration file is valid.")
}
