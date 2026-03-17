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

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
)

var ListTemplatesAction = modals.Action{
	Name: "List templates",
	Desc: fmt.Sprintf("List all templates with the `%s` command", ListTemplatesCommand()),
	Tldr: fmt.Sprintf("List all templates with the `%s` command", ListTemplatesCommand()),
}

// ListTemplatesCommand what's the command
func ListTemplatesCommand() string {
	return "/list-templates"
}

// ListTemplates handles the interaction to list the templates
func ListTemplates(client *slack.Client, triggerID string, channelID string) error {
	tx, ctx := apm.StartTransactionForm("ListTemplates", "request", triggerID, channelID)
	defer tx.End()

	Templates.Lock()
	files := make([]files.YamlFile, len(Clusters.files))
	copy(files, Templates.files)
	Templates.Unlock()

	headerText := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("We found *%d templates*", len(files)), false, false)

	// Create the header section
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Shared Divider
	divSection := slack.NewDividerBlock()

	list := ""
	for i, tpl := range files {
		if tpl.Data["template_name"] != nil && i < 40 {
			list += fmt.Sprintf("- *%s*\n", tpl.Data["template_name"].(string))
		}
	}

	templatesInfo := slack.NewTextBlockObject(slack.MarkdownType, list, false, false)
	templatesSection := slack.NewSectionBlock(templatesInfo, nil, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		templatesSection,
	)

	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
	)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	_ = os.Remove(config.ForUser(triggerID))
	return nil
}
