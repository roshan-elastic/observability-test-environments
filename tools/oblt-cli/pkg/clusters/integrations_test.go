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

func TestIntegrationsCluster_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ess := &IntegrationCluster{
		Branch:            "main",
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		Integration:       "test",
		IsRelease:         true,
		ObltRepo:          fakeConfigInfo.Obltrepo,
		Repository:        "owner/repo",
		SlackChannel:      fakeConfigInfo.SlackChannel,
		StackVersion:      "8.8.0",
		TemplateName:      IntegrationsTemplateName,
		Username:          fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var iTempateParams IntegrationClusterTemplate

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &iTempateParams)

	esDockerImage := strings.Split(buildEsDockerImageString(iTempateParams.StackVersion), "-")[0]
	kibanaDockerImage := strings.Split(buildKibanaDockerImageString(iTempateParams.StackVersion), "-")[0]
	agentDockerImage := strings.Split(buildElasticAgentDockerImageString(iTempateParams.StackVersion), "-")[0]

	assert.Equal(t, IntegrationsTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", iTempateParams.ElasticsearchDockerImage)
	assert.NotEqual(t, "", iTempateParams.KibanaDockerImage)
	assert.NotEqual(t, "", iTempateParams.ElasticAgentDockerImage)
	assert.Equal(t, ess.StackVersion, iTempateParams.StackVersion)
	assert.True(t, strings.HasPrefix(iTempateParams.ElasticsearchDockerImage, esDockerImage))
	assert.True(t, strings.HasPrefix(iTempateParams.KibanaDockerImage, kibanaDockerImage))
	assert.True(t, strings.HasPrefix(iTempateParams.ElasticAgentDockerImage, agentDockerImage))
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestIntegrationsClusterService_Create(t *testing.T) {
	t.Skip("Skipping test as it is not implemented yet")
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new CCSCluster instance with some test data
	ess := &IntegrationCluster{
		Branch:            "main",
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		Integration:       "test",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		Repository:        "owner/repo",
		SlackChannel:      fakeConfigInfo.SlackChannel,
		TemplateName:      IntegrationsTemplateName,
		Username:          fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ess.CreateService()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var iTempateParams IntegrationClusterTemplate

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &iTempateParams)

	assert.Equal(t, IntegrationsTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.Equal(t, "", iTempateParams.ElasticsearchDockerImage)
	assert.Equal(t, "", iTempateParams.KibanaDockerImage)
	assert.Equal(t, "", iTempateParams.ElasticAgentDockerImage)
	assert.Equal(t, "", iTempateParams.StackVersion)
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestIntegrationsCluster_with_cluster_name_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ess := &IntegrationCluster{
		Branch:       "main",
		ClusterName:  "test",
		Integration:  "test",
		IsRelease:    true,
		ObltRepo:     fakeConfigInfo.Obltrepo,
		Repository:   "owner/repo",
		SlackChannel: fakeConfigInfo.SlackChannel,
		StackVersion: "8.8.0",
		TemplateName: IntegrationsTemplateName,
		Username:     fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var iTempateParams IntegrationClusterTemplate

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &iTempateParams)

	assert.Equal(t, IntegrationsTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.Equal(t, ess.ClusterName, templateParams.ClusterName)
}

func TestIntegrationsClusterService_with_long_cluster_name_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with long cluster-name
	ess := &IntegrationCluster{
		Branch:       "main",
		ClusterName:  "kEQUz5wcRMTGq9saU3ecDhW2pvYNapFuupuKMzSABndzcmRCpQp4y6sLRPz7iB9ur",
		ObltRepo:     fakeConfigInfo.Obltrepo,
		Repository:   "owner/repo",
		SlackChannel: fakeConfigInfo.SlackChannel,
		TemplateName: IntegrationsTemplateName,
		Username:     fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.Error(t, err)
	assert.Nil(t, params)
}
