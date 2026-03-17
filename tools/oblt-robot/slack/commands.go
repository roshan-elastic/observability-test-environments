// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package slack

import (
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/interactions"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

// handleSlashCommand will take a slash command and route to the appropriate function
func handleSlashCommand(command slack.SlashCommand, client *slack.Client) error {
	client.Debugf("Slash command received: %+v", command)

	userID := command.UserID

	viper.SetConfigFile(config.ForUser(command.TriggerID))
	viper.Set(config.SlackChannelFlag, "@"+userID)
	viper.Set(config.UsernameFlag, interactions.SlackIDToUsername(userID))
	viper.Set(config.GitHttpModeFlag, true)

	logger.Debugf("Writing configuration file %s", viper.ConfigFileUsed())
	err := os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0700)
	if err != nil {
		return err
	}

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	// We need to switch depending on the command
	switch command.Command {
	case interactions.AddFAQEntryCommand():
		return interactions.ShowAddFAQEntryForm(client, command.TriggerID, command.UserName)
	case interactions.AWSAccountCommand():
		return interactions.ShowAWSAccountForm(client, command.TriggerID, command.UserName)
	case interactions.AWSResetAccountCommand():
		return interactions.ShowResetAWSAccountForm(client, command.TriggerID, command.UserName)
	case interactions.BugReportCommand():
		return interactions.ShowBugReportForm(client, command.TriggerID, command.UserName)
	case interactions.CIOnboardingCommand():
		return interactions.ShowCIOnboardingForm(client, command.TriggerID, command.UserName)
	case interactions.CloudAccountCommand():
		return interactions.ShowCloudAccountForm(client, command.TriggerID, command.UserName)
	case interactions.CreateCCSClusterCommand():
		return interactions.ShowCCSClusterForm(client, command.TriggerID, command.UserName)
	case interactions.CreateServerlessClusterCommand():
		return interactions.ShowServerlessClusterForm(client, command.TriggerID, command.UserName)
	case interactions.ClusterSecretCommand():
		return interactions.ShowSecretClusterForm(client, command.TriggerID, command.UserID)
	case interactions.DestroyClusterCommand():
		return interactions.ShowDestroyClusterForm(client, command.TriggerID, command.UserID)
	case interactions.ListClustersCommand():
		return interactions.ListClusters(client, command.TriggerID, command.UserID, true)
	case interactions.ListTemplatesCommand():
		return interactions.ListTemplates(client, command.TriggerID, command.UserID)
	case interactions.MyClustersCommand():
		return interactions.ListClusters(client, command.TriggerID, command.UserID, false)
	case interactions.ShowStatusCommand():
		return interactions.ShowStatuses(client, command.TriggerID, command.UserID)
	case interactions.ConfigurationCommand():
		return interactions.ShowConfigurationForm(client, command.TriggerID, command.UserID)
	}

	return nil
}
