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

func TestHandleCIOnboarding(t *testing.T) {
	// Start a new Slack test server
	s := slacktest.NewTestServer()
	defer s.Stop()
	s.Start()

	// Create a Slack client using the test server
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))
	github.MockedHTTPClient = github.NewMockedHTTPClient()

	interaction := slack.InteractionCallback{
		User:      slack.User{ID: "noAdmin"},
		TriggerID: "trigger123",
		View: slack.View{
			State: &slack.ViewState{
				Values: map[string]map[string]slack.BlockAction{
					"project": {
						"project": slack.BlockAction{
							Value: "oblservability-test-environments",
						},
					},
					"type": {
						"type": slack.BlockAction{
							Value: "Go",
						},
					},
				},
			},
		},
	}

	// Assuming HandleCIOnboarding takes a slack.Client and returns an error
	err := HandleCIOnboarding(slackClient, interaction)
	assert.NoError(t, err, "HandleCIOnboarding should not return an error")
}

func TestShowCIOnboardingForm(t *testing.T) {
	// Start a new Slack test server
	s := slacktest.NewTestServer()
	defer s.Stop()
	s.Start()

	s.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	// Create a Slack client using the test server
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

	// Assuming ShowCIOnboardingForm takes a slack.Client and returns an error
	err := ShowCIOnboardingForm(slackClient, "trigger123", "user123")
	assert.NoError(t, err, "ShowCIOnboardingForm should not return an error")
}
