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
	"errors"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
)

const (
	CCSTemplateName = "ccs"
)

// CCSCluster is the struct that contains the parameters for the CCS cluster
type CCSCluster struct {
	ClusterName       string                      `json:"ClusterName,omitempty"`
	ClusterNamePrefix string                      `json:"ClusterNamePrefix,omitempty"`
	ClusterNameSuffix string                      `json:"ClusterNameSuffix,omitempty"`
	GitHubCommentId   string                      `json:"GitHubCommentId,omitempty"`
	GitHubCommit      string                      `json:"GitHubCommit,omitempty"`
	GitHubIssue       string                      `json:"GitHubIssue,omitempty"`
	GitHubPullRequest string                      `json:"GitHubPullRequest,omitempty"`
	GitHubRepository  string                      `json:"GitHubRepository,omitempty"`
	RemoteClusterName string                      `json:"RemoteClusterName,omitempty"`
	ObltRepo          *ObltEnvironmentsRepository `json:"-"`
	SlackChannel      string                      `json:"SlackChannel,omitempty"`
	TemplateName      string                      `json:"TemplateName,omitempty"`
	Username          string                      `json:"Username,omitempty"`
}

// CCSTemplateParams is the struct that contains the parameters for the CCS template
type CCSTemplateParams struct {
	ElasticsearchDockerImage string `json:"ElasticsearchDockerImage,omitempty"`
	RemoteClusterConfigFile  string `json:"RemoteClusterConfigFile,omitempty"`
	RemoteClusterName        string `json:"RemoteClusterName,omitempty"`
	StackVersion             string `json:"StackVersion,omitempty"`
}

// Create creates a new CCS cluster
func (r *CCSCluster) Create() (parameters map[string]interface{}, err error) {
	var json string
	if err = r.validate(); err == nil {
		if json, err = r.packTemplateParams(); err == nil {
			parameters, err = r.createCluster(json)
		}
	}
	return parameters, err
}

// createCluster creates the cluster configuration file
func (r *CCSCluster) createCluster(jsonStr string) (map[string]interface{}, error) {
	var customCluster = &CustomCluster{
		ClusterName:       r.ClusterName,
		ClusterNamePrefix: joinIgnoreEmpty([]string{r.ClusterNamePrefix, r.RemoteClusterName}, itemsSeparator),
		ClusterNameSuffix: r.ClusterNameSuffix,
		TemplateName:      CCSTemplateName,
		Parameters:        string(jsonStr),
		ObltRepo:          r.ObltRepo,
		SlackChannel:      r.SlackChannel,
		Username:          r.Username,
		GitHubCommentId:   r.GitHubCommentId,
		GitHubCommit:      r.GitHubCommit,
		GitHubIssue:       r.GitHubIssue,
		GitHubPullRequest: r.GitHubPullRequest,
		GitHubRepository:  r.GitHubRepository,
	}
	return customCluster.Create()
}

// packTemplateParams packs the template parameters in a JSON string
func (r *CCSCluster) packTemplateParams() (jsonContent string, err error) {
	var file files.YamlFile
	var jsonContentb []byte

	file, err = r.ObltRepo.FindClusterConfig(r.RemoteClusterName)
	if err == nil {
		stack, _ := file.Data["stack"].(map[string]interface{})
		ess, _ := stack["ess"].(map[string]interface{})
		elasticsearch, _ := ess["elasticsearch"].(map[string]interface{})
		stackVersion, _ := stack["version"].(string)
		dockerImageVersion, _ := elasticsearch["image"].(string)

		templateParams := CCSTemplateParams{
			ElasticsearchDockerImage: dockerImageVersion,
			RemoteClusterConfigFile:  file.Path,
			RemoteClusterName:        r.RemoteClusterName,
			StackVersion:             stackVersion,
		}

		jsonContentb, err = json.Marshal(templateParams)
		if err == nil {
			jsonContent = string(jsonContentb)
		}
	}
	return jsonContent, err
}

// validate validates the parameters
func (r *CCSCluster) validate() (err error) {
	return errors.Join(
		config.ValidateClusterName(r.ClusterName),
		config.ValidateNames(r.RemoteClusterName),
	)
}
