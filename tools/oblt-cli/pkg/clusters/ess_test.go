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

func TestESSCluster_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ess := &ESSCluster{
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
		StackVersion:      "8.8.0",
		IsRelease:         true,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var essTempateParams ESSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &essTempateParams)

	esDockerImage := strings.Split(buildEsDockerImageString(essTempateParams.StackVersion), "-")[0]
	kibanaDockerImage := strings.Split(buildKibanaDockerImageString(essTempateParams.StackVersion), "-")[0]
	agentDockerImage := strings.Split(buildElasticAgentDockerImageString(essTempateParams.StackVersion), "-")[0]

	assert.Equal(t, ESSTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", essTempateParams.ElasticsearchDockerImage)
	assert.NotEqual(t, "", essTempateParams.KibanaDockerImage)
	assert.NotEqual(t, "", essTempateParams.ElasticAgentDockerImage)
	assert.Equal(t, ess.StackVersion, essTempateParams.StackVersion)
	assert.True(t, strings.HasPrefix(essTempateParams.ElasticsearchDockerImage, esDockerImage))
	assert.True(t, strings.HasPrefix(essTempateParams.KibanaDockerImage, kibanaDockerImage))
	assert.True(t, strings.HasPrefix(essTempateParams.ElasticAgentDockerImage, agentDockerImage))
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestESSCluster_CreateSNAPSHOT(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ess := &ESSCluster{
		ClusterNamePrefix: "test",
		ClusterNameSuffix: "suffix",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
		StackVersion:      "8.8.0-SNAPSHOT",
		IsRelease:         false,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var essTempateParams ESSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &essTempateParams)

	esDockerImage := strings.Split(buildEsDockerImageString(essTempateParams.StackVersion), "-")[0]
	kibanaDockerImage := strings.Split(buildKibanaDockerImageString(essTempateParams.StackVersion), "-")[0]
	agentDockerImage := strings.Split(buildElasticAgentDockerImageString(essTempateParams.StackVersion), "-")[0]

	assert.Equal(t, ESSTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.CommitMessage)
	assert.NotEqual(t, "", templateParams.CommitSha)
	assert.NotEqual(t, "", templateParams.CommitURL)
	assert.NotEqual(t, "", essTempateParams.ElasticsearchDockerImage)
	assert.NotEqual(t, "", essTempateParams.KibanaDockerImage)
	assert.NotEqual(t, "", essTempateParams.ElasticAgentDockerImage)
	assert.Equal(t, ess.StackVersion, essTempateParams.StackVersion)
	assert.True(t, strings.HasPrefix(essTempateParams.ElasticsearchDockerImage, esDockerImage))
	assert.True(t, strings.HasPrefix(essTempateParams.KibanaDockerImage, kibanaDockerImage))
	assert.True(t, strings.HasPrefix(essTempateParams.ElasticAgentDockerImage, agentDockerImage))
	assert.NotEqual(t, "", templateParams.ClusterName)
	assert.NotEqual(t, "", templateParams.ClusterConfigFile)
	assert.NotEqual(t, "", templateParams.Date)
}

func TestESSCluster_with_cluster_name_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	// Create a new CCSCluster instance with some test data
	ess := &ESSCluster{
		ClusterName:       "test",
		ClusterNameSuffix: "suffix",
		ObltRepo:          fakeConfigInfo.Obltrepo,
		SlackChannel:      fakeConfigInfo.SlackChannel,
		Username:          fakeConfigInfo.CurrentUsername,
		StackVersion:      "8.8.0",
		IsRelease:         true,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.NoError(t, err)
	assert.NotNil(t, params)

	var clusterParams CustomCluster
	var templateParams CustomTemplateParams
	var essTempateParams ESSTemplateParams

	jsonData, _ := json.Marshal(params)
	json.Unmarshal(jsonData, &clusterParams)
	json.Unmarshal(jsonData, &templateParams)
	json.Unmarshal(jsonData, &essTempateParams)

	assert.Equal(t, ESSTemplateName, clusterParams.TemplateName)
	assert.False(t, templateParams.GitOps)
	assert.Equal(t, ess.ClusterName, templateParams.ClusterName)
}

func TestESSCluster_with_long_clusterName_Create(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())

	// Create a new CCSCluster instance with a long cluster-name
	ess := &ESSCluster{
		ClusterName:  "kEQUz5wcRMTGq9saU3ecDhW2pvYNapFuupuKMzSABndzcmRCpQp4y6sLRPz7iB9ur",
		ObltRepo:     fakeConfigInfo.Obltrepo,
		SlackChannel: fakeConfigInfo.SlackChannel,
		Username:     fakeConfigInfo.CurrentUsername,
	}

	// Call the Create() method and check the result
	params, err := ess.Create()
	assert.Error(t, err)
	assert.Nil(t, params)
}
