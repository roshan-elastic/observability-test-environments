// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interactions

import (
	"net/http"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestAddFAQEntryPostMessage(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	interaction := slack.InteractionCallback{
		User:      slack.User{ID: "foo"},
		TriggerID: "bar",
	}
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	tests := []struct {
		name        string
		client      *slack.Client
		interaction slack.InteractionCallback
		message     string
		wantErr     bool
		data        faqData
	}{
		{
			name:        "Valid parameters",
			client:      client,
			interaction: interaction,
			message:     "This is a test FAQ entry",
			wantErr:     false,
			data: faqData{
				question: "What is the meaning of life?",
				answer:   "42",
				url:      "https://example.com",
				user:     "foo",
			},
		},
		// Define additional test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addFAQEntryPostMessage(tt.client, tt.interaction, tt.data)
			if tt.wantErr {
				assert.Error(t, err, "addFAQEntryPostMessage() should return an error")
			} else {
				assert.NoError(t, err, "addFAQEntryPostMessage() should not return an error")
			}
		})
	}
}

func TestHandleAddFAQEntry(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	interaction := slack.InteractionCallback{
		User: slack.User{
			ID:       "foo",
			Name:     "John",
			RealName: "John Doe",
		},
		View: slack.View{
			State: &slack.ViewState{
				Values: map[string]map[string]slack.BlockAction{
					"question": {
						"question": slack.BlockAction{
							Value: "What is the meaning of life?",
						},
					},
					"answer": {
						"answer": slack.BlockAction{
							Value: "42",
						},
					},
					"url": {
						"url": slack.BlockAction{
							Value: "https://example.com",
						},
					},
				},
			},
		},
	}
	err := HandleAddFAQEntry(client, interaction)
	assert.NoError(t, err, "HandleAddFAQEntry() should not return an error")
}

func TestShowAddFAQEntryForm(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	// handle the views.open request
	s.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	// Define your test cases
	tests := []struct {
		name      string
		client    *slack.Client
		wantErr   bool
		triggerID string
		user      string
	}{
		{
			name:      "Valid request",
			client:    client,
			wantErr:   false,
			triggerID: "bar",
			user:      "foo",
		},
		// Add more test cases as necessary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assuming ShowAddFAQEntryForm returns an error
			err := ShowAddFAQEntryForm(tt.client, tt.triggerID, tt.user)
			if tt.wantErr {
				assert.Error(t, err, "ShowAddFAQEntryForm() should return an error")
			} else {
				assert.NoError(t, err, "ShowAddFAQEntryForm() should not return an error")
			}
		})
	}
}
