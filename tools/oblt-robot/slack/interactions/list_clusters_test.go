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
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

func TestListClusters(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop() // Ensure the server is stopped after the test
	s.Start()

	// Create a Slack client using the mock server URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

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
			{
				Owner: "user321",
				Path:  "path/to/cluster2.yml",
				// Missing "cluster_name" key
			},
		},
	}

	// Assuming ListClusters takes a slack.Client and returns an error
	err := ListClusters(slackClient, "trigger123", "user123", true)
	assert.NoError(t, err, "ListClusters should not return an error")

	err = ListClusters(slackClient, "trigger123", "user123", false)
	assert.NoError(t, err, "ListClusters should not return an error")
}

func TestListClustersNoClusters(t *testing.T) {
	// Initialize the mock Slack server
	s := slacktest.NewTestServer()
	defer s.Stop() // Ensure the server is stopped after the test
	s.Start()

	// Create a Slack client using the mock server URL
	slackClient := slack.New("dummy-token", slack.OptionAPIURL(s.GetAPIURL()))

	Clusters = SafeArrayFiles{
		files: []files.YamlFile{},
	}

	// Assuming ListClusters takes a slack.Client and returns an error
	err := ListClusters(slackClient, "trigger123", "user123", true)
	assert.NoError(t, err, "ListClusters should not return an error")

	err = ListClusters(slackClient, "trigger123", "user123", false)
	assert.NoError(t, err, "ListClusters should not return an error")
}
