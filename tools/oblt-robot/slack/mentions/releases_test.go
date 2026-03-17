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

// Assuming the Releases function sends a message about the latest releases
func TestReleases(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	channelID := "channelID"
	user := slack.User{ID: "foo"}

	err := Releases(client, channelID, user)
	assert.NoError(t, err, "Releases function should not return an error")
}

// Assuming IsReleaseMention checks if a message is asking for releases
func TestIsReleaseMention(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "is release mention",
			message: "releases",
			want:    true,
		},
		{
			name:    "is release mention with prefix",
			message: "/releases",
			want:    true,
		},
		{
			name:    "is not release mention",
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
			got := IsReleaseMention(tt.message)
			assert.Equal(t, tt.want, got, "IsReleaseMention(%v) = %v, want %v", tt.message, got, tt.want)
		})
	}
}
