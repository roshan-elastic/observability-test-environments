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
	"net/http"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestHandleCloudAccount(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop() // Ensure the server is stopped after the test
	s.Start()

	// Create a Slack client using the mock server's URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))
	github.MockedHTTPClient = github.NewMockedHTTPClient()

	interactions := slack.InteractionCallback{
		User:      slack.User{ID: "noAdmin"},
		TriggerID: "trigger123",
		View: slack.View{
			State: &slack.ViewState{
				Values: map[string]map[string]slack.BlockAction{
					"email": {
						"email": slack.BlockAction{
							Value: "no-elastic@example.com",
						},
					},
				},
			},
		},
	}
	// Assuming HandleCloudAccount takes a slack.Client and returns an error
	err := HandleCloudAccount(slackClient, interactions)
	// Use assert to check if error is nil (which means success)
	assert.NoError(t, err, "HandleCloudAccount should not return an error")

	// Test with an elastic email
	interactions.View.State.Values["email"]["email"] = slack.BlockAction{
		Value: "foo@elastic.co",
	}
	err = HandleCloudAccount(slackClient, interactions)
	assert.NoError(t, err, "HandleCloudAccount should not return an error")

	// Test with an elasticsearch email
	interactions.View.State.Values["email"]["email"] = slack.BlockAction{
		Value: "foo@elasticsearch.com",
	}
	err = HandleCloudAccount(slackClient, interactions)
	assert.NoError(t, err, "HandleCloudAccount should not return an error")
}

func TestShowCloudAccountForm(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop() // Ensure the server is stopped after the test
	s.Start()

	s.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	// Create a Slack client using the mock server's URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

	// Assuming ShowCloudAccountForm takes a slack.Client and returns an error
	err := ShowCloudAccountForm(slackClient, "trigger123", "user123")
	// Use assert to check if error is nil (which means success)
	assert.NoError(t, err, "ShowCloudAccountForm should not return an error")
}
