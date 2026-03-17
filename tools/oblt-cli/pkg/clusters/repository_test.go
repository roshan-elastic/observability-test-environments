// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http:// www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NormalizeElasticStackVersion(t *testing.T) {
	type testData struct {
		version                    string
		isRelease                  bool
		expectedStackVersion       string
		expectedDockerImageVersion string
	}

	testDatas := []testData{
		{
			version:                    "8.0.0",
			isRelease:                  true,
			expectedStackVersion:       "8.0.0",
			expectedDockerImageVersion: "8.0.0",
		},
		{
			version:                    "8.0.0-SNAPSHOT",
			isRelease:                  false,
			expectedStackVersion:       "8.0.0-SNAPSHOT",
			expectedDockerImageVersion: "8.0.0-SNAPSHOT",
		},
		{
			version:                    "8.0.0",
			isRelease:                  false,
			expectedStackVersion:       "8.0.0",
			expectedDockerImageVersion: "8.0.0",
		},
		{
			version:                    "8.0.0-aaaaaaaa",
			isRelease:                  false,
			expectedStackVersion:       "8.0.0-SNAPSHOT",
			expectedDockerImageVersion: "8.0.0-aaaaaaaa",
		},
		{
			version:                    "8.0.0-aaaaaaaa",
			isRelease:                  true,
			expectedStackVersion:       "8.0.0",
			expectedDockerImageVersion: "8.0.0-aaaaaaaa",
		},
	}

	for _, td := range testDatas {
		stackVersion, dockerImageVersion := normalizeElasticStackVersion(td.version, td.isRelease)
		assert.Equal(t, td.expectedStackVersion, stackVersion, fmt.Sprintf("Wrong Stack version %s - %t", td.version, td.isRelease))
		assert.Equal(t, td.expectedDockerImageVersion, dockerImageVersion, fmt.Sprintf("Wrong Docker image version %s - %t", td.version, td.isRelease))
	}
}

func Test_ObltTestEnvironments(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	currentUsername := fakeConfigInfo.CurrentUsername
	slackChannel := fakeConfigInfo.SlackChannel
	tmp := fakeConfigInfo.TempDir
	gitRepo := fakeConfigInfo.Obltrepo.gitRepo
	environmentsPath := fakeConfigInfo.EnvironmentsPath
	usersFolderPath := fakeConfigInfo.UsersFolderPath
	obltRepo := fakeConfigInfo.Obltrepo

	assert.True(t, gitRepo.AlwaysCleanUp)

	t.Run("Environment paths", func(t *testing.T) {
		assert.Equal(t, environmentsPath, obltRepo.getEnvironments())
		assert.Equal(t, usersFolderPath, obltRepo.GetUsersDir())
		assert.Equal(t, filepath.Join(usersFolderPath, currentUsername), obltRepo.GetCurrentUserEnvironments())
		assert.Equal(t, filepath.Join(usersFolderPath, "ccsTemplate.yml"), obltRepo.GetTemplate("ccsTemplate.yml"))
		assert.Equal(t, filepath.Join(usersFolderPath, currentUsername, "cluster-config.yml"), obltRepo.GetUserClusterConfig("cluster-config.yml"))
	})

	t.Run("Create CCS cluster: empty versions", func(t *testing.T) {
		remoteClusterName := "dev-oblt"

		clusterConfig, err := obltRepo.FindClusterConfig(remoteClusterName)
		assert.Nil(t, err)
		stack := clusterConfig.Data["stack"].(map[string]interface{})
		ess := stack["ess"].(map[string]interface{})
		esImage := ess["elasticsearch"].(map[string]interface{})["image"].(string)
		esVersion := stack["version"].(string)

		ccsCluster := &CCSCluster{
			TemplateName:      CCSTemplateName,
			ClusterNamePrefix: "fooprefix",
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			RemoteClusterName: remoteClusterName,
		}

		parametersMap, err := ccsCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var clusterParams CustomCluster
		var templateParams CustomTemplateParams
		var ccsTempateParams CCSTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &clusterParams)
		json.Unmarshal(jsonData, &templateParams)
		json.Unmarshal(jsonData, &ccsTempateParams)

		assert.Equal(t, ccsTempateParams.ElasticsearchDockerImage, esImage)
		assert.Equal(t, ccsTempateParams.StackVersion, esVersion)
	})

	t.Run("Create CCS cluster", func(t *testing.T) {
		remoteClusterName := "dev-oblt"
		repo := "elastic/acme"
		commit := "abcd1234"
		issue := "1"
		pullRequest := "PR-1"
		commentId := "1234556789"
		devConfigFile, _ := obltRepo.FindClusterConfig("dev-oblt")

		ccsCluster := &CCSCluster{
			TemplateName:      CCSTemplateName,
			ClusterNamePrefix: "fooprefix",
			ClusterNameSuffix: "barsuffix",
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			GitHubRepository:  repo,
			GitHubCommit:      commit,
			GitHubIssue:       issue,
			GitHubPullRequest: pullRequest,
			GitHubCommentId:   commentId,
			RemoteClusterName: remoteClusterName,
		}

		parametersMap, err := ccsCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var clusterParams CustomCluster
		var templateParams CustomTemplateParams
		var ccsTempateParams CCSTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &clusterParams)
		json.Unmarshal(jsonData, &templateParams)
		json.Unmarshal(jsonData, &ccsTempateParams)

		clusterName := templateParams.ClusterName

		assert.True(t, strings.HasPrefix(clusterName, "fooprefix-dev-oblt"))
		assert.True(t, strings.HasSuffix(clusterName, "barsuffix"))

		assert.Equal(t, clusterParams.Username, currentUsername)
		assert.Equal(t, clusterParams.SlackChannel, slackChannel)
		assert.Equal(t, clusterParams.TemplatePath, filepath.Join(usersFolderPath, "ccs.yml.tmpl"))
		assert.Equal(t, ccsTempateParams.RemoteClusterName, remoteClusterName)
		assert.Equal(t, ccsTempateParams.RemoteClusterConfigFile, devConfigFile.Path)
		assert.Equal(t, templateParams.ClusterConfigFile, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"))
		assert.NotEqual(t, "", templateParams.CommitSha)
		assert.NotEqual(t, "", templateParams.CommitMessage)
		assert.True(t, strings.HasPrefix(templateParams.CommitURL, obltTestEnvironmentsCommitURL))
		assert.True(t, templateParams.GitOps)
		assert.Equal(t, templateParams.GitHubCommit, commit)
		assert.Equal(t, templateParams.GitHubIssue, issue)
		assert.Equal(t, templateParams.GitHubPullRequest, pullRequest)
		assert.Equal(t, templateParams.GitHubRepository, repo)
		assert.Equal(t, templateParams.GitHubCommentId, commentId)
	})

	t.Run("Create Serverless cluster", func(t *testing.T) {
		repo := "elastic/acme"
		commit := "abcd1234"
		issue := "1"
		pullRequest := "PR-1"
		commentId := "1234556789"
		projectType := "security"
		target := "production"

		serverlessCluster := &ServerlessCluster{
			TemplateName:      ServerlessTemplateName,
			ClusterNamePrefix: "fooprefix",
			ClusterNameSuffix: "barsuffix",
			ProjectType:       projectType,
			Target:            target,
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			GitHubRepository:  repo,
			GitHubCommit:      commit,
			GitHubIssue:       issue,
			GitHubPullRequest: pullRequest,
			GitHubCommentId:   commentId,
		}

		parametersMap, err := serverlessCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var clusterParams CustomCluster
		var templateParams CustomTemplateParams
		var serverlessTempateParams ServerlessTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &clusterParams)
		json.Unmarshal(jsonData, &templateParams)
		json.Unmarshal(jsonData, &serverlessTempateParams)

		clusterName := templateParams.ClusterName

		assert.True(t, strings.HasPrefix(clusterName, "fooprefix"))
		assert.True(t, strings.HasSuffix(clusterName, "barsuffix"))

		assert.Equal(t, clusterParams.Username, currentUsername)
		assert.Equal(t, clusterParams.SlackChannel, slackChannel)
		assert.Equal(t, clusterParams.TemplatePath, filepath.Join(usersFolderPath, "serverless.yml.tmpl"))
		assert.Equal(t, templateParams.ClusterConfigFile, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"))
		assert.Equal(t, serverlessTempateParams.ProjectType, projectType)
		assert.Equal(t, serverlessTempateParams.Target, target)
		assert.NotEqual(t, "", templateParams.CommitSha)
		assert.NotEqual(t, "", templateParams.CommitMessage)
		assert.True(t, strings.HasPrefix(templateParams.CommitURL, obltTestEnvironmentsCommitURL))
		assert.True(t, templateParams.GitOps)
		assert.Equal(t, templateParams.GitHubCommit, commit)
		assert.Equal(t, templateParams.GitHubIssue, issue)
		assert.Equal(t, templateParams.GitHubPullRequest, pullRequest)
		assert.Equal(t, templateParams.GitHubRepository, repo)
		assert.Equal(t, templateParams.GitHubCommentId, commentId)
	})

	t.Run("Create custom cluster with template file path", func(t *testing.T) {
		templateFilePath := filepath.Join(usersFolderPath, "ess.yml.tmpl")
		obj := &ESSTemplateParams{
			StackVersion: "7.15.0",
		}
		parameters, _ := json.Marshal(obj)

		customCluster := &CustomCluster{
			ClusterNamePrefix: "fooprefix",
			ClusterNameSuffix: "barsuffix",
			TemplatePath:      templateFilePath,
			Parameters:        string(parameters),
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
		}

		parametersMap, err := customCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var clusterParams CustomCluster
		var templateParams CustomTemplateParams
		var essTempateParams ESSTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &clusterParams)
		json.Unmarshal(jsonData, &templateParams)
		json.Unmarshal(jsonData, &essTempateParams)

		clusterName := templateParams.ClusterName

		assert.True(t, strings.HasPrefix(clusterName, "fooprefix"))
		assert.True(t, strings.HasSuffix(clusterName, "barsuffix"))
		assert.Equal(t, templateParams.Username, currentUsername)
		assert.Equal(t, templateParams.SlackChannel, slackChannel)
		assert.Equal(t, essTempateParams.StackVersion, "7.15.0")
		assert.Equal(t, clusterParams.TemplatePath, filepath.Join(usersFolderPath, "ess.yml.tmpl"))
		assert.Equal(t, templateParams.ClusterConfigFile, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"))
		assert.NotEqual(t, templateParams.CommitSha, "")
		assert.NotEqual(t, templateParams.CommitMessage, "")
		assert.True(t, strings.HasPrefix(templateParams.CommitURL, obltTestEnvironmentsCommitURL))
		assert.False(t, templateParams.GitOps)
	})

	t.Run("Create custom cluster with template file path with GitOps", func(t *testing.T) {
		commit := "abcd1234"
		templateFilePath := filepath.Join(usersFolderPath, "ess.yml.tmpl")
		obj := &ESSTemplateParams{
			StackVersion: "7.15.0",
		}
		parameters, _ := json.Marshal(obj)
		customCluster := &CustomCluster{
			ClusterNamePrefix: "fooprefix",
			ClusterNameSuffix: "barsuffix",
			TemplatePath:      templateFilePath,
			Parameters:        string(parameters),
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			GitHubCommit:      commit,
		}

		parametersMap, err := customCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var templateParams CustomTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &templateParams)

		clusterName := templateParams.ClusterName

		assert.True(t, strings.HasPrefix(clusterName, "fooprefix"))
		assert.True(t, strings.HasSuffix(clusterName, "barsuffix"))
		assert.True(t, templateParams.GitOps)
		assert.Equal(t, commit, templateParams.GitHubCommit)
	})

	t.Run("Create custom cluster with template name", func(t *testing.T) {
		obj := &ESSTemplateParams{
			StackVersion:             "8.3.0-SNAPSHOT",
			ElasticsearchDockerImage: "docker.elastic.co/cloud-ci/elasticsearch:8.3.0-00f4f855",
		}
		parameters, _ := json.Marshal(obj)
		customCluster := &CustomCluster{
			ClusterNamePrefix: "fooprefix",
			TemplateName:      ESSTemplateName,
			Parameters:        string(parameters),
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
		}

		parametersMap, err := customCluster.Create()
		assert.Nil(t, err)
		assert.NotNil(t, parametersMap)

		var clusterParams CustomCluster
		var templateParams CustomTemplateParams
		var essTempateParams ESSTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &clusterParams)
		json.Unmarshal(jsonData, &templateParams)
		json.Unmarshal(jsonData, &essTempateParams)

		clusterName := templateParams.ClusterName

		assert.True(t, strings.HasPrefix(clusterName, "fooprefix"))
		assert.Equal(t, currentUsername, templateParams.Username)
		assert.Equal(t, slackChannel, templateParams.SlackChannel)
		assert.Equal(t, "8.3.0-SNAPSHOT", essTempateParams.StackVersion)
		assert.Equal(t, "docker.elastic.co/cloud-ci/elasticsearch:8.3.0-00f4f855", essTempateParams.ElasticsearchDockerImage)
		assert.Equal(t, clusterParams.TemplatePath, filepath.Join(usersFolderPath, "ess.yml.tmpl"))
		assert.Equal(t, templateParams.ClusterConfigFile, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"))
	})

	t.Run("Destroy CCS cluster", func(t *testing.T) {
		ccsCluster := &CCSCluster{
			TemplateName:      CCSTemplateName,
			ClusterNamePrefix: "fooprefix",
			ClusterNameSuffix: "barsuffix",
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			RemoteClusterName: "dev-oblt",
		}
		createParams, err := ccsCluster.Create()
		assert.Nil(t, err)

		clusterName := createParams["ClusterName"].(string)

		parametersMap, err := obltRepo.DestroyCluster(clusterName, true)
		assert.Nil(t, err)

		var templateParams CustomTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &templateParams)

		assert.NotEqual(t, templateParams.CommitSha, "")
		assert.NotEqual(t, templateParams.CommitMessage, "")
		assert.True(t, strings.HasPrefix(templateParams.CommitURL, obltTestEnvironmentsCommitURL))
		assert.True(t, strings.HasSuffix(parametersMap[clusterName].(string), clusterName+".yml"))

		_, err = obltRepo.FindClusterConfig(clusterName)
		assert.NotNil(t, err)
	})

	t.Run("Destroy CCS cluster that does not exist", func(t *testing.T) {
		_, err := obltRepo.DestroyCluster("fooprefix-dev-oblt-barsuffix-ccs", true)
		assert.NotNil(t, err)
	})

	t.Run("Find cluster config", func(t *testing.T) {
		devConfig, err := obltRepo.FindClusterConfig("dev-oblt")
		assert.Nil(t, err)
		assert.NotEmpty(t, devConfig.Path)

		edgeConfig, err := obltRepo.FindClusterConfig("edge-oblt")
		assert.Nil(t, err)
		assert.NotEmpty(t, edgeConfig.Path)

		releaseConfig, err := obltRepo.FindClusterConfig("release-oblt")
		assert.Nil(t, err)
		assert.NotEmpty(t, releaseConfig.Path)
	})

	t.Run("Find template", func(t *testing.T) {
		ccsTemplate, err := obltRepo.FindTemplate("ccs")
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(usersFolderPath, "ccs.yml.tmpl"), ccsTemplate.Path)

		deployKibanaTemplate, err := obltRepo.FindTemplate("deploy-kibana")
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(usersFolderPath, "deploy-kibana.yml.tmpl"), deployKibanaTemplate.Path)

		essTemplate, err := obltRepo.FindTemplate("ess")
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(usersFolderPath, "ess.yml.tmpl"), essTemplate.Path)
	})

	t.Run("Find current user cluster config", func(t *testing.T) {
		ccsCluster := &CCSCluster{
			TemplateName:      CCSTemplateName,
			ClusterNamePrefix: "fooprefix",
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			RemoteClusterName: "dev-oblt",
		}
		parametersMap, err := ccsCluster.Create()
		assert.Nil(t, err)

		var templateParams CustomTemplateParams
		jsonData, _ := json.Marshal(parametersMap)
		json.Unmarshal(jsonData, &templateParams)

		clusterName := templateParams.ClusterName

		userClusterConfig, err := obltRepo.FindCurrentUserClusterConfig(clusterName)
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"), userClusterConfig.Path)
	})

	t.Run("List templates", func(t *testing.T) {
		templates := obltRepo.ListTemplates()
		assert.True(t, len(templates) > 0)

		for template := range templates {
			assert.True(t, strings.HasSuffix(templates[template].Path, ".tmpl"))
		}
	})

	t.Run("List Bootstrap recipes", func(t *testing.T) {
		recipes := obltRepo.ListBootstrapRecipes("elasticsearch")
		assert.True(t, len(recipes) >= 1)
		es_test := false
		for _, recipe := range recipes {
			if filepath.Join(tmp, gitRepo.Name, obltRepo.GetBootstrapRecipesRelativeDir(), "elasticsearch", "test.yml") == recipe.Path {
				es_test = true
			}
		}
		assert.True(t, es_test)

		recipes = obltRepo.ListBootstrapRecipes("kibana")
		assert.True(t, len(recipes) >= 1)
		kbn_test := false
		for _, recipe := range recipes {
			if filepath.Join(tmp, gitRepo.Name, obltRepo.GetBootstrapRecipesRelativeDir(), "kibana", "test.yml") == recipe.Path {
				kbn_test = true
			}
		}
		assert.True(t, kbn_test)
	})

	t.Run("Load Bootstrap recipes from JSON", func(t *testing.T) {
		recipes := obltRepo.LoadRecipesFromJson("elasticsearch", "[\"test\"]")
		assert.Equal(t, 1, len(recipes))
		assert.Equal(t, filepath.Join(tmp, gitRepo.Name, obltRepo.GetBootstrapRecipesRelativeDir(), "elasticsearch", "test.yml"), recipes[0].Path)

		recipes = obltRepo.LoadRecipesFromJson("kibana", "[\"test\"]")
		assert.Equal(t, 1, len(recipes))
		assert.Equal(t, filepath.Join(tmp, gitRepo.Name, obltRepo.GetBootstrapRecipesRelativeDir(), "kibana", "test.yml"), recipes[0].Path)
	})

	t.Run("List clusters", func(t *testing.T) {
		// TO-DO refactor the tests so that each t.Run has an instance of its own git Repository
		// counting the number of clusters before adding a new one, which will prevent us from other tests writing in the repo
		initialClusters := obltRepo.ListClusters(false)
		initialClustersCount := len(initialClusters)

		// adding a new cluster
		ccsCluster := &CCSCluster{
			TemplateName:      CCSTemplateName,
			ClusterNamePrefix: "fooprefix",
			Username:          currentUsername,
			SlackChannel:      slackChannel,
			ObltRepo:          obltRepo,
			RemoteClusterName: "dev-oblt",
		}
		_, err := ccsCluster.Create()
		assert.Nil(t, err)

		// listing the clusters again to verify that one was added
		clusters := obltRepo.ListClusters(false)
		assert.Equal(t, initialClustersCount+1, len(clusters))
	})
}

func Test_Wipeup(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	currentUsername := fakeConfigInfo.CurrentUsername
	slackChannel := fakeConfigInfo.SlackChannel
	obltRepo := fakeConfigInfo.Obltrepo

	ccsCluster := &CCSCluster{
		TemplateName:      CCSTemplateName,
		ClusterNamePrefix: "fooprefix",
		ClusterNameSuffix: "barsuffix",
		Username:          currentUsername,
		SlackChannel:      slackChannel,
		ObltRepo:          obltRepo,
		RemoteClusterName: "dev-oblt",
	}
	createParams00, err := ccsCluster.Create()
	assert.Nil(t, err)
	createParams01, err := ccsCluster.Create()
	assert.Nil(t, err)

	results, err := obltRepo.Wipeup()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(results))

	for _, result := range results {
		assert.True(t, result[createParams00["ClusterName"].(string)] != nil || result[createParams01["ClusterName"].(string)] != nil)
	}
}

func Test_NewClusterName(t *testing.T) {
	t.Run("(With Prefix, suffix and random name", func(t *testing.T) {
		clusterName, err := newClusterName("prefix", "suffix", "random")
		assert.Nil(t, err)
		assert.True(t, strings.HasPrefix(clusterName, "prefix"))
		assert.True(t, strings.Contains(clusterName, "random"))
		assert.True(t, strings.HasSuffix(clusterName, "suffix"))
	})
	t.Run("(With Prefix only", func(t *testing.T) {
		clusterName, err := newClusterName("prefix", "", "")
		assert.Nil(t, err)
		assert.True(t, strings.HasPrefix(clusterName, "prefix"))
		assert.False(t, strings.Contains(clusterName, "random"))
		assert.False(t, strings.HasSuffix(clusterName, "suffix"))
	})
	t.Run("(With Suffix only", func(t *testing.T) {
		clusterName, err := newClusterName("", "suffix", "")
		assert.Nil(t, err)
		assert.False(t, strings.HasPrefix(clusterName, "prefix"))
		assert.False(t, strings.Contains(clusterName, "random"))
		assert.True(t, strings.HasSuffix(clusterName, "suffix"))
	})
}

func Test_UpdateCluster(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	currentUsername := fakeConfigInfo.CurrentUsername
	slackChannel := fakeConfigInfo.SlackChannel
	obltRepo := fakeConfigInfo.Obltrepo
	usersFolderPath := fakeConfigInfo.UsersFolderPath

	ccsCluster := &CCSCluster{
		TemplateName:      CCSTemplateName,
		ClusterNamePrefix: "fooprefix",
		ClusterNameSuffix: "barsuffix",
		Username:          currentUsername,
		SlackChannel:      slackChannel,
		ObltRepo:          obltRepo,
		RemoteClusterName: "dev-oblt",
	}
	createParams00, err := ccsCluster.Create()
	assert.Nil(t, err)

	parameters := make(map[string]interface{})
	parameters["foo"] = "bar"
	parameters["bar"] = make(map[string]interface{})
	parameters["bar"].(map[string]interface{})["foo"] = "baz"
	parameters["stack"] = make(map[string]interface{})
	parameters["stack"].(map[string]interface{})["version"] = "foo"

	clusterName := createParams00["ClusterName"].(string)
	data, err := json.Marshal(parameters)
	assert.Nil(t, err, "Error marshalling parameters")
	reaultClusterConfig, err := obltRepo.UpdateCluster(clusterName, string(data))
	assert.Nil(t, err, "Error updating cluster")
	assert.Equal(t, "bar", reaultClusterConfig["foo"])
	assert.Equal(t, "baz", reaultClusterConfig["bar"].(map[string]interface{})["foo"], "error creating a new key")
	assert.Equal(t, "foo", reaultClusterConfig["stack"].(map[string]interface{})["version"], "error updating a key")
	assert.Equal(t, clusterName, reaultClusterConfig["cluster_name"])

	clusterConfig, path, err := obltRepo.ReadClusterConfig(clusterName)
	assert.Nil(t, err, "Error reading cluster config")
	assert.Equal(t, "baz", clusterConfig["bar"].(map[string]interface{})["foo"], "error creating a new key")
	assert.Equal(t, "foo", clusterConfig["stack"].(map[string]interface{})["version"], "error updating a key")
	assert.Equal(t, filepath.Join(usersFolderPath, currentUsername, clusterName+".yml"), path)
	assert.Equal(t, clusterName, clusterConfig["cluster_name"])
}

func Test_UpdateClusterChangeName(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	currentUsername := fakeConfigInfo.CurrentUsername
	slackChannel := fakeConfigInfo.SlackChannel
	obltRepo := fakeConfigInfo.Obltrepo

	ccsCluster := &CCSCluster{
		TemplateName:      CCSTemplateName,
		ClusterNamePrefix: "fooprefix",
		ClusterNameSuffix: "barsuffix",
		Username:          currentUsername,
		SlackChannel:      slackChannel,
		ObltRepo:          obltRepo,
		RemoteClusterName: "dev-oblt",
	}
	createParams00, err := ccsCluster.Create()
	assert.Nil(t, err)

	parameters := make(map[string]interface{})
	parameters["cluster_name"] = "bar"

	clusterName := createParams00["ClusterName"].(string)
	data, err := json.Marshal(parameters)
	assert.Nil(t, err, "Error marshalling parameters")
	_, err = obltRepo.UpdateCluster(clusterName, string(data))
	assert.Error(t, err, "The name of the cluster cannot be changed")
}

func Test_GetVersion(t *testing.T) {
	fakeConfigInfo := FakeConfig(t.TempDir())
	obltRepo := fakeConfigInfo.Obltrepo
	version, err := obltRepo.GetVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
}
