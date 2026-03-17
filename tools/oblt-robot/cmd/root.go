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

// Package cmd this package contains the oblt-cli commands and flags.
package cmd

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&slack.DryRun, config.DryRunFlag, false, "If true, the Git changes will not commit the changes made to the Git repository.")
	RootCmd.PersistentFlags().BoolVar(&logger.Verbose, config.VerboseFlag, false, "If true, the logs output is more verbose.")
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "oblt-robot",
	Short: "Runs a Slack bot using websockets to interact with the CLI",
	Long:  `Runs a Slack bot using websockets to interact with the CLI, allowing running it on a server`,
	Run: func(cmd *cobra.Command, args []string) {
		slack.StartSocketMode()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}
