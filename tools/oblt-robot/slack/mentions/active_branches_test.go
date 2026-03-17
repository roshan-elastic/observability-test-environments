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
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestActiveBraches(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	channelID := "channelID"
	user := slack.User{ID: "foo"}
	github.MockedHTTPClient = github.NewMockedHTTPClient()

	err := ActiveBraches(client, channelID, user)
	assert.NoError(t, err)
}
func TestIsActiveBranchesMention(t *testing.T) {
	text := "Please mention branches"
	result := IsActiveBranchesMention(text)
	assert.True(t, result, "Expected IsActiveBranchesMention to return true")

	text = "Please mention versions"
	result = IsActiveBranchesMention(text)
	assert.False(t, result, "Expected IsActiveBranchesMention to return false")
}
