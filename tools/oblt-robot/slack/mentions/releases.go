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

var ReleaseAction = modals.Action{
	Name: "Current releases",
	Desc: "Mention me in the right channel including the `releases` word. I will let you know what are the current list of releases/build candidates in the Unified Release process.",
	Tldr: "Ask me (including `releases` word) and I'll list all the current releases",
}

// IsReleaseMention whether it's a releases mention
func IsReleaseMention(text string) bool {
	return strings.Contains(strings.ToLower(text), "releases")
}

// Releases handles the interaction to show the current releases
func Releases(client *slack.Client, channelID string, user slack.User) error {

	releases, _ := releases.GetReleases()

	blocks, _ := modals.RenderReleases(releases)
	_, _, err := client.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}
