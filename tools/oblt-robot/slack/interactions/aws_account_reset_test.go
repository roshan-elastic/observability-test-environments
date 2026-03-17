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

func TestHandleResetAWSAccount(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	slackClient := slack.New("TEST_TOKEN", slack.OptionAPIURL(s.GetAPIURL()))
	github.MockedHTTPClient = github.NewMockedHTTPClient()

	tests := []struct {
		name         string
		client       *slack.Client
		wantErr      bool
		interactions slack.InteractionCallback
	}{
		{
			name:    "Reset AWS Account with no elastic",
			client:  slackClient,
			wantErr: false,
			interactions: slack.InteractionCallback{
				User:      slack.User{ID: "user123"},
				TriggerID: "trigger123",
				View: slack.View{
					State: &slack.ViewState{
						Values: map[string]map[string]slack.BlockAction{
							"email": {
								"email": slack.BlockAction{
									Value: "foo@example.com",
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Reset AWS Account with no elastic",
			client:  slackClient,
			wantErr: false,
			interactions: slack.InteractionCallback{
				User:      slack.User{ID: "user123"},
				TriggerID: "trigger123",
				View: slack.View{
					State: &slack.ViewState{
						Values: map[string]map[string]slack.BlockAction{
							"email": {
								"email": slack.BlockAction{
									Value: "foo@elastic.co",
								},
							},
						},
					},
				},
			},
		},
		// Add more test cases as necessary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleResetAWSAccount(tt.client, tt.interactions)
			if tt.wantErr {
				assert.Error(t, err, "HandleResetAWSAccount() should return an error")
			} else {
				assert.NoError(t, err, "HandleResetAWSAccount() should not return an error")
			}
		})
	}
}

func TestShowResetAWSAccountForm(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	s.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	slackClient := slack.New("TEST_TOKEN", slack.OptionAPIURL(s.GetAPIURL()))

	tests := []struct {
		name      string
		client    *slack.Client
		wantErr   bool
		triggerID string
		user      string
	}{
		{
			name:      "Show Reset AWS Account Form with valid client",
			client:    slackClient,
			wantErr:   false,
			triggerID: "trigger123",
			user:      "user123",
		},
		// Add more test cases as necessary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ShowResetAWSAccountForm(tt.client, tt.triggerID, tt.user)
			if tt.wantErr {
				assert.Error(t, err, "ShowResetAWSAccountForm() should return an error")
			} else {
				assert.NoError(t, err, "ShowResetAWSAccountForm() should not return an error")
			}
		})
	}
}
