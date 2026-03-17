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
	"fmt"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/interactions"
	"github.com/slack-go/slack"
)

func handleInteractionEvent(interaction slack.InteractionCallback, client *slack.Client) error {
	// Switch depending on the Type
	switch interaction.Type {
	case slack.InteractionTypeViewSubmission:

		view := interaction.View
		// select the view by form title
		title := view.Title.Text
		switch title {
		case interactions.AddFAQEntryAction.Name:
			interactions.HandleAddFAQEntry(client, interaction)
		case interactions.AWSAccountAction.Name:
			interactions.HandleAWSAccount(client, interaction)
		case interactions.AWSResetAccountAction.Name:
			interactions.HandleResetAWSAccount(client, interaction)
		case interactions.CIOnboardingAction.Name:
			interactions.HandleCIOnboarding(client, interaction)
		case interactions.CloudAccountAction.Name:
			interactions.HandleCloudAccount(client, interaction)
		case interactions.CreateCCSClusterAction.Name:
			interactions.HandleCreateCCSCluster(client, interaction)
		case interactions.CreateServerlessClusterAction.Name:
			interactions.HandleCreateServerlessCluster(client, interaction)
		case interactions.BugReportAction.Name:
			interactions.HandleBugReport(client, interaction)
		case interactions.DestroyClusterAction.Name:
			interactions.HandleSelectCluster(client, interaction)
		case interactions.ClusterSecretAction.Name:
			interactions.HandleSelectCluster(client, interaction)
		case interactions.ConfigurationAction.Name:
			interactions.HandleConfiguration(client, interaction)
		default:
			fmt.Printf("form interaction not supported: %v\n", title)
		}
	case slack.InteractionTypeBlockActions:
		actionCallback := interaction.ActionCallback
		messageTS := interaction.Message.Timestamp // use timestamp of the message to respond in a thread

		for _, blockAction := range actionCallback.BlockActions {
			if strings.HasPrefix(blockAction.ActionID, interactions.ReadClusterSecretPrefix) {
				gcsm, err := gcp.NewClusterSecrets()
				if err == nil {
					_, _, err = interactions.PostSecretMessage(client, interaction, blockAction.ActionID, messageTS, gcsm)
				}
				return err
			}
		}

	default:
		return fmt.Errorf("interaction type not found: %s", interaction.Type)
	}

	return nil
}
