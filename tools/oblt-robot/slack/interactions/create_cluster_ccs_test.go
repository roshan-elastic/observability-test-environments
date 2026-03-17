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

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateCCSCluster(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop()
	s.Start()

	// Create a Slack client using the mock server URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

	// Define test cases for all switch cases in HandleCreateCCSCluster
	tests := []struct {
		name        string
		setup       func() // Setup function if needed for specific test case
		wantErr     bool
		interaction slack.InteractionCallback
	}{
		{
			name:    "CCS Cluster from template",
			wantErr: false,
			interaction: slack.InteractionCallback{
				User:      slack.User{ID: "dry-run"},
				TriggerID: "trigger123",
				View: slack.View{
					State: &slack.ViewState{
						Values: map[string]map[string]slack.BlockAction{
							"create-ccs-cluster-from-template": {
								"select-from-template": slack.BlockAction{
									SelectedOption: slack.OptionBlockObject{
										Value: "edge-oblt",
									},
								},
							},
							"create-ccs-cluster-name-prefix": {
								"cluster-name-prefix": slack.BlockAction{
									Value: "foo",
								},
							},
							"create-ccs-cluster-name-suffix": {
								"cluster-name-suffix": slack.BlockAction{
									Value: "bar",
								},
							},
							"ignore": {
								"ignore": slack.BlockAction{
									Value: "baz",
								},
							},
						},
					},
				},
			},
		},
		// Add more cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assuming HandleCreateCCSCluster takes a slack.Client and returns an error
			err := HandleCreateCCSCluster(slackClient, tt.interaction)
			assert.NoError(t, err, "HandleCreateCCSCluster should not return an error")
		})
	}
}

func TestShowCCSClusterForm(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop()
	s.Start()

	s.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	// Create a Slack client using the mock server's URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

	// Assuming ShowCCSClusterForm takes a slack.Client and returns an error
	err := ShowCCSClusterForm(slackClient, "trigger123", "user123")
	assert.NoError(t, err, "ShowCCSClusterForm should not return an error")
}
