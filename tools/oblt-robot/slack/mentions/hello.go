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

	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
)

var HelloAction = modals.Action{
	Name: "Hello",
	Desc: "Mention me in the right channel including the `hello` word. I will greet you",
	Tldr: "I'll greet you",
}

// IsHelloMention whether it's a hello mention
func IsHelloMention(text string) bool {
	return strings.Contains(strings.ToLower(text), "hello")
}

// Hello handles the interaction to show greetings
func Hello(client *slack.Client, channelID string, user slack.User) error {
	blocks, _ := modals.RenderHello(user.ID)
	_, _, err := client.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}
