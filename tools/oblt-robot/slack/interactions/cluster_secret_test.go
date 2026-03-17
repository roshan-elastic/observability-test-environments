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

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

// https://github.com/lusis/slack-test
// https://github.com/slack-go/slack/tree/master/slacktest
func TestSecrets(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	interaction := slack.InteractionCallback{
		User:      slack.User{ID: "foo"},
		TriggerID: "bar",
	}
	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	gcsm := createClusterSecretsMock()
	for _, secret := range listOfSecrets {
		blocks, attachments, err := PostSecretMessage(api, interaction, "fuzzy__foobar__"+secret, "", gcsm)
		assert.NoError(t, err, secret)
		assert.True(t, attachments != nil || blocks != nil, secret)
	}

	_, _, err := PostSecretMessage(api, interaction, "fuzzy__foobar__foo", "", gcsm)
	assert.Error(t, err)
}

func TestSecretSelection(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	interaction := slack.InteractionCallback{
		User:      slack.User{ID: "foo"},
		TriggerID: "bar",
	}
	api := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	blocks, err := postSecretsMessage(api, interaction, "cluster-foo")
	assert.NoError(t, err)
	assert.NotNil(t, blocks)
}

// createClusterSecretsMock creates a new ClusterSecrets instance using a mock client.
func createClusterSecretsMock() *gcp.ClusterSecrets {
	client := gcp.NewSecretsManagerMock()
	auth := gcp.NewAuthMock()
	return gcp.NewClusterSecretsWithClient(client, &auth)
}

func TestHandleSelectCluster(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	slackClient := slack.New("TEST_TOKEN", slack.OptionAPIURL(s.GetAPIURL()))

	Clusters = SafeArrayFiles{
		files: []files.YamlFile{
			{
				Owner: "user123",
				Path:  "path/to/cluster1.yml",
				Data:  map[interface{}]interface{}{"cluster_name": "cluster1"},
			},
			{
				Owner: "user321",
				Path:  "path/to/cluster2.yml",
				Data:  map[interface{}]interface{}{"cluster_name": "cluster2"},
			},
		},
	}
	tests := []struct {
		name         string
		client       *slack.Client
		wantErr      bool
		interactions slack.InteractionCallback
	}{
		{
			name:    "Select cluster to destroy",
			client:  slackClient,
			wantErr: false,
			interactions: slack.InteractionCallback{
				User:      slack.User{ID: "user123"},
				TriggerID: "trigger123",
				View: slack.View{
					State: &slack.ViewState{
						Values: map[string]map[string]slack.BlockAction{
							"destroy-select-cluster": {
								"select-cluster": slack.BlockAction{
									SelectedOption: slack.OptionBlockObject{
										Value: "cluster1",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Select Cluster to retrieve secrets",
			client:  slackClient,
			wantErr: false,
			interactions: slack.InteractionCallback{
				User:      slack.User{ID: "user123"},
				TriggerID: "trigger123",
				View: slack.View{
					State: &slack.ViewState{
						Values: map[string]map[string]slack.BlockAction{
							"select-cluster-for-secrets": {
								"select-cluster": slack.BlockAction{
									SelectedOption: slack.OptionBlockObject{
										Value: "cluster1",
									},
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
			err := HandleSelectCluster(tt.client, tt.interactions)
			if tt.wantErr {
				assert.Error(t, err, "HandleSelectCluster() should return an error")
			} else {
				assert.NoError(t, err, "HandleSelectCluster() should not return an error")
			}
		})
	}
}

func TestShowSecretClusterForm(t *testing.T) {
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
	err := ShowSecretClusterForm(slackClient, "trigger123", "user123")
	assert.NoError(t, err, "ShowCIOnboardingForm should not return an error")
}
