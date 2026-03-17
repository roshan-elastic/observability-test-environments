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

var ListClustersAction = modals.Action{
	Name: "List clusters",
	Desc: fmt.Sprintf("List all clusters with the `%s` command", ListClustersCommand()),
	Tldr: fmt.Sprintf("List all clusters with the `%s` command", ListClustersCommand()),
}

var MyClustersAction = modals.Action{
	Name: "My clusters",
	Desc: fmt.Sprintf("List your clusters with the `%s` command", MyClustersCommand()),
	Tldr: fmt.Sprintf("List your clusters with the `%s` command", MyClustersCommand()),
}

// ListClustersCommand what's the command
func ListClustersCommand() string {
	return "/list-clusters"
}

// MyClustersCommand what's the command
func MyClustersCommand() string {
	return "/my-clusters"
}

// ListClusters handles the interaction to list the clusters
func ListClusters(client *slack.Client, triggerID string, slackID string, all bool) error {
	tx, ctx := apm.StartTransactionForm("ListClusters", "request", triggerID, slackID)
	defer tx.End()

	Clusters.Lock()
	files := make([]files.YamlFile, len(Clusters.files))
	copy(files, Clusters.files)
	Clusters.Unlock()
	if !all {
		files = FilterClustersByUser(slackID)
	}

	headerText := slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("We found *%d Clusters*", len(files)), false, false)

	// Create the header section
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Shared Divider
	divSection := slack.NewDividerBlock()

	list := ""
	for i, file := range files {
		if file.Data["cluster_name"] != nil && i < 40 {
			list += fmt.Sprintf("- *%s* [%s]\n", file.Data["cluster_name"].(string), file.Owner)
		}
	}

	if list == "" {
		list = "No clusters found"
	}

	clustersInfo := slack.NewTextBlockObject(slack.MarkdownType, list, false, false)
	clustersSection := slack.NewSectionBlock(clustersInfo, nil, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		clustersSection,
	)

	_, _, err := client.PostMessage(
		slackID,
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
