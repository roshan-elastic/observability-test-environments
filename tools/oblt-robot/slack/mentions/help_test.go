// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

// Assuming Help function is defined in yourpackage and operates similarly to the Hello function

func TestHelp(t *testing.T) {
	// Initialize the Slack test server
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	// Create a new Slack client with the test server's URL
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	channelID := "channelID"
	user := slack.User{ID: "foo"}

	// Call the Help function with the test server client, channelID, and user
	err := Help(client, channelID, user)

	// Use assert.NoError to check if an error was not returned
	assert.NoError(t, err, "Help function should not return an error")
}

func TestIsHelpMention(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "is help message",
			message: "help",
			want:    true,
		},
		{
			name:    "is help message with prefix",
			message: "/help",
			want:    true,
		},
		{
			name:    "is not help message",
			message: "hello",
			want:    false,
		},
		{
			name:    "empty message",
			message: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHelpMention(tt.message)
			assert.Equal(t, tt.want, got, "IsHelpMention() = %v, want %v", got, tt.want)
		})
	}
}
