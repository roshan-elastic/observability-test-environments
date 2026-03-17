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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerlessCluster_Create_Default(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new ServerlessCluster instance with some test data
	serverless := &ServerlessCluster{
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
	}

	// Call the Create() method and check the result
	params, err := serverless.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var serverlessTempateParams ServerlessTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &serverlessTempateParams)

	assert.Equal(t, "", serverlessTempateParams.ProjectType)
	assert.Equal(t, "", serverlessTempateParams.Target)
	assert.Equal(t, ServerlessTemplateName, clusterParams.TemplateName)
	assert.Equal(t, serverless.GitHubCommentId, templateParams.GitHubCommentId)
	assert.Equal(t, serverless.GitHubCommit, templateParams.GitHubCommit)
	assert.Equal(t, serverless.GitHubIssue, templateParams.GitHubIssue)
	assert.Equal(t, serverless.GitHubPullRequest, templateParams.GitHubPullRequest)
	assert.Equal(t, serverless.GitHubRepository, templateParams.GitHubRepository)
	assert.True(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestServerlessCluster_Create_With_Parameters(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new ServerlessCluster instance with some test data
	serverless := &ServerlessCluster{
		ClusterNamePrefix:        "test",
		ClusterNameSuffix:        "suffix",
		ObltRepo:                 fakeConfigInfo.Obltrepo,
		SlackChannel:             fakeConfigInfo.SlackChannel,
		Username:                 fakeConfigInfo.CurrentUsername,
		GitHubCommentId:          "1234",
		GitHubCommit:             "aaaa",
		GitHubIssue:              "56789",
		GitHubPullRequest:        "89",
		GitHubRepository:         "owner/repo",
		ProjectType:              "security",
		Target:                   "production",
		ElasticsearchDockerImage: "docker.elastic.co/elasticsearch/elasticsearch:7.14.0-SNAPSHOT",
		KibanaDockerImage:        "docker.elastic.co/kibana/kibana:7.14.0-SNAPSHOT",
		FleetDockerImage:         "docker.elastic.co/beats/elastic-agent:7.14.0-SNAPSHOT",
	}

	// Call the Create() method and check the result
	params, err := serverless.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var serverlessTempateParams ServerlessTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &serverlessTempateParams)

	assert.Equal(t, "security", serverlessTempateParams.ProjectType)
	assert.Equal(t, "production", serverlessTempateParams.Target)
	assert.Equal(t, "docker.elastic.co/elasticsearch/elasticsearch:7.14.0-SNAPSHOT", serverlessTempateParams.ElasticsearchDockerImage)
	assert.Equal(t, "docker.elastic.co/kibana/kibana:7.14.0-SNAPSHOT", serverlessTempateParams.KibanaDockerImage)
	assert.Equal(t, "docker.elastic.co/beats/elastic-agent:7.14.0-SNAPSHOT", serverlessTempateParams.FleetDockerImage)
	assert.Equal(t, ServerlessTemplateName, clusterParams.TemplateName)
	assert.Equal(t, serverless.GitHubCommentId, templateParams.GitHubCommentId)
	assert.Equal(t, serverless.GitHubCommit, templateParams.GitHubCommit)
	assert.Equal(t, serverless.GitHubIssue, templateParams.GitHubIssue)
	assert.Equal(t, serverless.GitHubPullRequest, templateParams.GitHubPullRequest)
	assert.Equal(t, serverless.GitHubRepository, templateParams.GitHubRepository)
	assert.True(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}
