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
	"log"
	"os"
	"runtime"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// dryRun enable/disable test mode.
var dryRun bool

// experimental enable/disable experiemental features.
var experimental bool

// outputFile absolute path to a file to save the operation output results.
var outputFile string

// wait it waits N minutes for the operation to finish.
var wait int

// disableBanner disable the system status banner.
var disableBanner bool

// ciMode enable/disable CI mode.
var ciMode bool

// saveConfig force to save the configuration in the config file.
var saveConfig bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "oblt-cli",
	Short: "CLI tool to operate the oblt clusters.",
	Long: `CLI tool to operate the oblt clusters, it allows to create/update/destroy oblt clusters.

	For more information, please visit:
	* clusters docs at https://ela.st/oblt-clusters
	* oblt-cli docs at https://ela.st/oblt-cli
	* oblt-robot docs at https://ela.st/oblt-robot
	* Servreless docs at https://ela.st/oblt-serverless`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	logger.Debugf("Root Args: %s", os.Args)

	// Configure the APM instrumentation, hardcoded for now and using the alias
	os.Setenv("ELASTIC_APM_SERVICE_NAME", "oblt-cli")
	os.Setenv("ELASTIC_APM_ENVIRONMENT", "production")
	os.Setenv("ELASTIC_APM_SERVER_URL", "https://observability-ci.apm.us-west2.gcp.elastic-cloud.com")

	config.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&config.CfgFile, config.ConfigFlag, "", "Config file (default is "+config.DefaultFile()+")")
	RootCmd.PersistentFlags().BoolVar(&dryRun, config.DryRunFlag, false, "If true, the Git changes will not commit the changes made to the Git repository.")
	RootCmd.PersistentFlags().BoolVar(&disableBanner, config.DisableBannerFlag, false, "If true, the system status banner will not be shown.")
	RootCmd.PersistentFlags().BoolVar(&logger.Verbose, config.VerboseFlag, false, "If true, the logs output is more verbose.")
	RootCmd.PersistentFlags().BoolVar(&experimental, config.ExperimentalFlag, false, "If true, the experimental features will be available.")
	RootCmd.PersistentFlags().StringVar(&outputFile, config.OutputFileFlag, "", "It is the absolute path to a file to save the operation output results. If the file name is '-' or 'stdout' it writes the data in the stdout.")
	RootCmd.PersistentFlags().IntVar(&wait, config.WaitFlag, -1, "it waits N minutes for the operation to finish.")
	RootCmd.PersistentFlags().StringVar(&config.RepoBranch, config.BranchFlag, "main", "Change the default branch to checkout from observability-test-environments repository. This flag is for testing new features. (default is main)")

	RootCmd.PersistentFlags().String(config.SlackChannelFlag, "", "Slack member ID.")
	RootCmd.PersistentFlags().String(config.UsernameFlag, "", "Username to show in the deployments [a-z0-9] (any name identifies you in Elastic no matter, it is for tagging purposes e.g myuser).")
	RootCmd.PersistentFlags().Bool(config.GitHttpModeFlag, false, "If true oblt-cli uses HTTP to checkout the code and the GITHUB_TOKEN environment variable to authenticate.")
	RootCmd.PersistentFlags().BoolVar(&saveConfig, config.SaveConfigFlag, false, "If true, it forces to save the configuration in the config file.")
}

// initConfig reads in config file.
func initConfig() {
	initDefaultValues()
	config.Initialise(config.DefaultFile(), whenConfigured)
}

func initDefaultValues() {
	if logger.Verbose {
		// show the file name and line number in the logs
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if os.Getenv("CI") == "true" {
		logger.Infof("Running in CI environment")
		disableBanner = true
		ciMode = true
	}

	logger.Infof("oblt-cli version %s %s/%s\n", currentVersion, runtime.GOOS, runtime.GOARCH)
}

func whenConfigured() {
	logger.Debugf("Initialization after configured")
}
