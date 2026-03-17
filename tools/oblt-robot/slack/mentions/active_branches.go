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

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/releases"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
)

var ActiveBranchesAction = modals.Action{
	Name: "active branches",
	Desc: "Mention me in the right channel including the `branches` word. I will let you know what are the active branches in the Unified Release process.",
	Tldr: "Ask me (including `branches` word) and I'll list the active branches",
}

// IsActiveBranchesMention whether it's a version mention
func IsActiveBranchesMention(text string) bool {
	return strings.Contains(strings.ToLower(text), "branches")
}

// ActiveBraches handles the interaction to show the current active branches in the Unified Release
func ActiveBraches(client *slack.Client, channelID string, user slack.User) error {

	versions, _ := releases.GetVersions()

	blocks, _ := modals.RenderBranches(versions)
	_, _, err := client.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}
