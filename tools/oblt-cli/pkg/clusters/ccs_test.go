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

func TestCCSCluster_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ccs := &CCSCluster{
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		RemoteClusterName: "edge-oblt",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
		GitHubCommentId:   "1234",
		GitHubCommit:      "aaaa",
		GitHubIssue:       "56789",
		GitHubPullRequest: "89",
		GitHubRepository:  "owner/repo",
	}

	// Call the Create() method and check the result
	params, err := ccs.Create()
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

	assert.Equal(t, ccs.RemoteClusterName, ccsTempateParams.RemoteClusterName)
	assert.Equal(t, CCSTemplateName, clusterParams.TemplateName)
	assert.Equal(t, ccs.GitHubCommentId, templateParams.GitHubCommentId)
	assert.Equal(t, ccs.GitHubCommit, templateParams.GitHubCommit)
	assert.Equal(t, ccs.GitHubIssue, templateParams.GitHubIssue)
	assert.Equal(t, ccs.GitHubPullRequest, templateParams.GitHubPullRequest)
	assert.Equal(t, ccs.GitHubRepository, templateParams.GitHubRepository)
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

func TestCCSCluster_with_clusterName_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new CCSCluster instance with some test data
	ccs := &CCSCluster{
		ClusterName:       "test",
		RemoteClusterName: "edge-oblt",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ccs.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var ccsTempateParams CCSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &ccsTempateParams)

	assert.Equal(t, ccs.RemoteClusterName, ccsTempateParams.RemoteClusterName)
	assert.Equal(t, CCSTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.Equal(t, ccs.ClusterName, templateParams.ClusterName)
}

func TestCCSCluster_with_long_clusterName_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new CCSCluster instance with a long cluster-name
	ccs := &CCSCluster{
		ClusterName:       "kEQUz5wcRMTGq9saU3ecDhW2pvYNapFuupuKMzSABndzcmRCpQp4y6sLRPz7iB9ur",
		RemoteClusterName: "edge-oblt",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ccs.Create()
	assert.Error(t, err)
	assert.Nil(t, params)
}
