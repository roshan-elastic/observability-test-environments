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
	apm2 "go.elastic.co/apm/v2"
)

var ShowStatusAction = modals.Action{
	Name: "Show Clusters status",
	Desc: fmt.Sprintf("Show the status messages for the clusters with the `%s` command", ShowStatusCommand()),
	Tldr: fmt.Sprintf("Show the status for the Oblt Clusters with the `%s` command", ShowStatusCommand()),
}

// ShowStatusCommand what's the command
func ShowStatusCommand() string {
	return "/show-status"
}

// ShowStatuses handles the interaction to show the statuses
func ShowStatuses(client *slack.Client, triggerID string, channelID string) error {
	tx, ctx := apm.StartTransactionForm("ShowStatuses", "request", triggerID, channelID)
	defer tx.End()
	repository, err := OnMemoryUserConf(triggerID, triggerID, true)

	var messages []clusters.Message
	if err == nil {
		messages = clusters.ShowBanner(repository.GetPath())
	}

	err = sendShowStatusResponse(messages, err, client, channelID)

	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	_ = os.Remove(config.ForUser(triggerID))
	return nil
}

// sendShowStatusResponse sends the response to the user
func sendShowStatusResponse(messages []clusters.Message, err error, client *slack.Client, channelID string) error {
	if err == nil {
		headerText := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("We found *%d messages*", len(messages)), false, false)

		headerSection := slack.NewSectionBlock(headerText, nil, nil)

		divSection := slack.NewDividerBlock()

		list := ""
		for _, message := range messages {
			list += fmt.Sprintf("- `%s`: %s\n", message.Level, message.Text)
		}

		messagesInfo := slack.NewTextBlockObject(slack.MarkdownType, list, false, false)
		messagesSection := slack.NewSectionBlock(messagesInfo, nil, nil)

		msg := slack.NewBlockMessage(
			headerSection,
			divSection,
			messagesSection,
		)

		_, _, err = client.PostMessage(
			channelID,
			slack.MsgOptionText(msg.Text, false),
			slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
		)
	}
	return err
}
