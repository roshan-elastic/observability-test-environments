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
package mentions

import (
	"fmt"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/interactions"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
)

var HelpAction = modals.Action{
	Name: "Help",
	Desc: "Mention me in the right channel including the `help` word. I will display this help message",
	Tldr: "Display help",
}

// IsHelpMention whether it's a help mention
func IsHelpMention(text string) bool {
	return strings.Contains(strings.ToLower(text), "help")
}

// Help handles the interaction to show the help commands
func Help(client *slack.Client, channelID string, user slack.User) error {

	mentions := []modals.Action{HelpAction, ActiveBranchesAction, FAQAction, HelloAction, ReleaseAction, SnapshotAction, VersionAction}

	commands := []modals.Action{interactions.AddFAQEntryAction, interactions.AWSAccountAction, interactions.CIOnboardingAction, interactions.ClusterSecretAction, interactions.CreateCCSClusterAction, interactions.CreateServerlessClusterAction, interactions.DestroyClusterAction, interactions.ListClustersAction, interactions.ListTemplatesAction, interactions.MyClustersAction, interactions.ShowStatusAction}

	blocks, _ := modals.RenderGeneralHelp(commands, mentions)
	_, _, err := client.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}
