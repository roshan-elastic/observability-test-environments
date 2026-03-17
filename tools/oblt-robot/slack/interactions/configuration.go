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
package interactions

import (
	"fmt"
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	apm2 "go.elastic.co/apm/v2"
)

// NewTestEventSampleConfiguration creates a test event for configuration
func NewTestEventSampleConfiguration() socketmode.Event {
	return socketmode.Event{
		Type: socketmode.EventTypeSlashCommand,
		Request: &socketmode.Request{
			EnvelopeID: "dummy",
		},
		Data: slack.InteractionCallback{
			User:      slack.User{ID: "dry-run"},
			TriggerID: "trigger123",
			Type:      slack.InteractionTypeViewSubmission,
			View: slack.View{
				Title: &slack.TextBlockObject{
					Type: slack.PlainTextType,
					Text: ConfigurationAction.Name,
				},
				State: &slack.ViewState{
					Values: map[string]map[string]slack.BlockAction{
						"oblt-configure": {
							"oblt-configure-github-username": slack.BlockAction{
								Value: "dry-run",
							},
						},
						"ignore": {
							"ignore": slack.BlockAction{
								Value: "baz",
							},
						},
					},
				},
			},
		},
	}
}

type ObltConfigurationData struct {
	Username string
	SlackID  string
}

var ConfigurationAction = modals.Action{
	Name: "Configure Oblt clusters.",
	Desc: fmt.Sprintf("Display a form to configure a user to start using oblt clusters `%s` command", ConfigurationCommand()),
	Tldr: fmt.Sprintf("Configure oblt clusters with the `%s` command", ConfigurationCommand()),
}

// ConfigurationCommand what's the command
func ConfigurationCommand() string {
	return "/oblt-configure"
}

// configurationPostMessage posts message after Configure Oblt clusters
func configurationPostMessage(client *slack.Client, interaction slack.InteractionCallback, githubUser string) (err error) {
	user := interaction.User
	var data ObltConfigurationData
	data.Username = githubUser
	data.SlackID = "@" + user.ID
	dryRun := user.ID == "dry-run"
	var obltTestEnvironments *clusters.ObltEnvironmentsRepository
	if obltTestEnvironments, err = OnMemoryUserConf(user.ID, data.Username, dryRun); err == nil {

		customCluster := &clusters.CustomCluster{
			TemplateName: "oblt-user",
			ClusterName:  data.Username,
			Username:     data.Username,
			SlackChannel: data.SlackID,
			ObltRepo:     obltTestEnvironments,
			Parameters:   "{}",
		}

		if _, err = customCluster.Create(); err == nil {
			blocks, _ := modals.RenderConfigure(user.ID, githubUser)
			_, _, err = client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...))
		}
	}

	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return nil
}

// HandleConfiguration handles the interaction with the configure form
func HandleConfiguration(client *slack.Client, interaction slack.InteractionCallback) error {
	var githubUsername string
	viewState := interaction.View.State

	// pick up all values from the create cluster form
	for k, v := range viewState.Values {
		switch k {
		case "oblt-configure":
			githubUsername = v["oblt-configure-github-username"].Value
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}
	err := configurationPostMessage(client, interaction, githubUsername)
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return err
}

// ShowConfigurationForm creates a ModalViewRequest requesting the configuration information
func ShowConfigurationForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowConfigurationForm", "request", triggerID, user)
	defer tx.End()

	modalRequest, err := modals.RenderConfigureForm(triggerID)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	_, err = client.OpenView(triggerID, modalRequest)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}
