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
package clusters

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomCluster_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	stackVersion := "8.8.0"
	paramsIn := CCSTemplateParams{
		ElasticsearchDockerImage: buildEsDockerImageString(stackVersion),
		RemoteClusterConfigFile:  "path/to/the/file",
		RemoteClusterName:        "edge-oblt",
		StackVersion:             stackVersion,
	}

	jsonStr, _ := json.Marshal(paramsIn)

	// Create a new CCSCluster instance with some test data
	custom := &CustomCluster{
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
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var ccsTempateParams CCSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &ccsTempateParams)

	esDockerImage := strings.Split(buildEsDockerImageString(ccsTempateParams.StackVersion), "-")[0]

	assert.Equal(t, "edge-oblt", ccsTempateParams.RemoteClusterName)
	assert.Equal(t, CCSTemplateName, clusterParams.TemplateName)
	assert.Equal(t, custom.GitHubCommentId, templateParams.GitHubCommentId)
	assert.Equal(t, custom.GitHubCommit, templateParams.GitHubCommit)
	assert.Equal(t, custom.GitHubIssue, templateParams.GitHubIssue)
	assert.Equal(t, custom.GitHubPullRequest, templateParams.GitHubPullRequest)
	assert.Equal(t, custom.GitHubRepository, templateParams.GitHubRepository)
	assert.True(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", ccsTempateParams.RemoteClusterConfigFile)
	assert.True(t, strings.HasPrefix(ccsTempateParams.ElasticsearchDockerImage, esDockerImage))
	assert.NotEqual(t, "", ccsTempateParams.StackVersion)
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestCustomCluster_with_clusterName_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	stackVersion := "8.8.0"
	paramsIn := CCSTemplateParams{
		ElasticsearchDockerImage: buildEsDockerImageString(stackVersion),
		RemoteClusterConfigFile:  "path/to/the/file",
		RemoteClusterName:        "edge-oblt",
		StackVersion:             stackVersion,
	}

	jsonStr, _ := json.Marshal(paramsIn)

	// Create a new CCSCluster instance with some test data
	custom := &CustomCluster{
		ClusterName:  "test-foo-bar",
		ObltRepo:     fakeConfigInfo.Obltrepo,
		SlackChannel: fakeConfigInfo.SlackChannel,
		Username:     fakeConfigInfo.CurrentUsername,
		Parameters:   string(jsonStr),
		TemplateName: "ccs",
	}

	// Call the Create() method and check the result
	params, err := custom.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var ccsTempateParams CCSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &ccsTempateParams)

	assert.Equal(t, "edge-oblt", ccsTempateParams.RemoteClusterName)
	assert.Equal(t, CCSTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.Equal(t, custom.ClusterName, templateParams.ClusterName)
}

func TestCustomCluster_with_long_clusterName_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new CCSCluster instance with a long cluster-name
	custom := &CustomCluster{
		ClusterName:  "kEQUz5wcRMTGq9saU3ecDhW2pvYNapFuupuKMzSABndzcmRCpQp4y6sLRPz7iB9ur",
		ObltRepo:     fakeConfigInfo.Obltrepo,
		SlackChannel: fakeConfigInfo.SlackChannel,
		Username:     fakeConfigInfo.CurrentUsername,
		Parameters:   string("{}"),
		TemplateName: "ccs",
	}

	// Call the Create() method and check the result
	params, err := custom.Create()
	assert.Error(t, err)
	assert.Nil(t, params)
}
