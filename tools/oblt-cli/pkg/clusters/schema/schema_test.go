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
package schema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/stretchr/testify/assert"
)

func createFakeClusterConfig(t *testing.T) string {
	fakeConfigInfo := clusters.FakeConfig(t.TempDir())
	stackVersion := "8.8.0"
	paramsIn := clusters.CCSTemplateParams{
		ElasticsearchDockerImage: fmt.Sprintf("docker.elastic.co/es-image:%s", stackVersion),
		RemoteClusterConfigFile:  "path/to/the/file",
		RemoteClusterName:        "edge-oblt",
		StackVersion:             stackVersion,
	}

	jsonStr, _ := json.Marshal(paramsIn)

	// Create a new CCSCluster instance with some test data
	custom := &clusters.CustomCluster{
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
		GitHubCommentId:   "1234",
		GitHubCommit:      "aaaa",
		GitHubIssue:       "56789",
		GitHubPullRequest: "89",
		GitHubRepository:  "owner/repo",
		Parameters:        string(jsonStr),
		TemplateName:      "ccs",
	}

	// Call the Create() method and check the result
	params, err := custom.Create()
	assert.NoError(t, err)

	var templateParams clusters.CustomTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &templateParams)
	return templateParams.ClusterConfigFile
}

func TestValidate(t *testing.T) {
	clusterConfigFile := createFakeClusterConfig(t)
	err := Validate(clusterConfigFile)
	assert.NoError(t, err)
}
