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
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	git "github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/git"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apmLib "go.elastic.co/apm/v2"
)

// interactive if true, it will ask for the configuration settings in a interactive way.
var interactive bool

// ConfigureCmd represents the configure command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "It configures the oblt-cli.",
	Long: `It configures the Slack member ID and other parameters needed, also default values to use.
  The configuration is saved in a file in ` + config.DefaultFile(),
	Run: runConfigure,
}

func init() {
	RootCmd.AddCommand(ConfigureCmd)

	ConfigureCmd.Flags().String(config.SlackChannelFlag, "", "Slack member ID.")
	ConfigureCmd.Flags().String(config.UsernameFlag, "", "GitHub user to tag deployments [a-z0-9].")
	ConfigureCmd.Flags().Bool(config.GitHttpModeFlag, false, "If true oblt-cli uses HTTP to checkout the code and the GITHUB_TOKEN environment variable to authenticate.")
	ConfigureCmd.Flags().BoolVar(&interactive, config.InteractiveFlag, false, "If true, it will ask for the configuration settings in a interactive way.")

	viper.BindPFlag(config.SlackChannelFlag, ConfigureCmd.Flags().Lookup(config.SlackChannelFlag))
	viper.BindPFlag(config.UsernameFlag, ConfigureCmd.Flags().Lookup(config.UsernameFlag))
	viper.BindPFlag(config.GitHttpModeFlag, ConfigureCmd.Flags().Lookup(config.GitHttpModeFlag))
}

func runConfigure(cmd *cobra.Command, args []string) {
	tx, ctx := apm.StartTransaction("runConfigure", "request", []apm.Label{labelVersion}, config.NewObltConfig("", "no-set", "no-set", false, false))
	defer apm.Flush(tx)
	saveConfig = true

	if interactive {
		prompt.Configure()
		viper.Set(config.SlackChannelFlag, prompt.SlackChannel)
		viper.Set(config.UsernameFlag, prompt.User)
	}

	_, err := newObltConfig(cmd)
	apm.CobraCheckErr(err, tx, ctx)

	createUserConfigRepo(config.NewObltConfigFromViper(), tx, ctx)
}

// gitMode it is the Git mode selected (SSH/HTTP)
func gitMode(isHTTPMode bool) (ret string) {
	if isHTTPMode {
		return git.ModeHTTP
	}

	return git.ModeSSH
}

// validateConfig It checks that all configuration settings are correct.
func validateConfig(userConfig config.ObltConfiguration) (err error) {
	logger.Debugf("SlackChannel: '%s'", userConfig.SlackChannel)
	logger.Debugf("User: '%s'", userConfig.Username)
	logger.Debugf("Git mode: '%s'", gitMode(userConfig.GitHttpMode))

	err = errors.Join(
		config.ValidateUsername(userConfig.Username),
		config.ValidateSlackChannel(userConfig.SlackChannel),
	)
	return err
}

// newObltConfig It creates a new ObltConfig.
func newObltConfig(cmd *cobra.Command) (userConfig config.ObltConfiguration, err error) {
	userConfig = config.NewObltConfigFromViper()
	username, _ := cmd.Flags().GetString(config.UsernameFlag)
	slackChannel, _ := cmd.Flags().GetString(config.SlackChannelFlag)
	gitHttpMode, _ := cmd.Flags().GetBool(config.GitHttpModeFlag)

	userConfig.Username = firstNotEmpty(username, userConfig.Username)
	userConfig.SlackChannel = firstNotEmpty(slackChannel, userConfig.SlackChannel)
	userConfig.GitHttpMode = gitHttpMode || userConfig.GitHttpMode

	if err = validateConfig(userConfig); err == nil {
		viper.Set(config.SlackChannelFlag, userConfig.SlackChannel)
		viper.Set(config.UsernameFlag, userConfig.Username)
		viper.Set(config.GitHttpModeFlag, userConfig.GitHttpMode)
		err = writeConfiguration()
	}

	return userConfig, err
}

// writeConfiguration It writes the configuration file.
func writeConfiguration() (err error) {
	if saveConfig {
		logger.Infof("Writing configuration file %s", viper.ConfigFileUsed())
		if err = os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0700); err == nil {
			err = viper.WriteConfig()
		}
	}
	return err
}

func createUserConfigRepo(userConfig config.ObltConfiguration, tx *apmLib.Transaction, ctx context.Context) (err error) {
	if !ciMode {
		obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
		apm.CobraCheckErr(err, tx, ctx)

		customCluster := &clusters.CustomCluster{
			TemplateName: "oblt-user",
			ClusterName:  userConfig.Username,
			Username:     userConfig.Username,
			SlackChannel: userConfig.SlackChannel,
			ObltRepo:     obltTestEnvironments,
			Parameters:   "{}",
		}

		_, err = customCluster.Create()
		apm.CobraCheckErr(err, tx, ctx)
	}
	return err
}
