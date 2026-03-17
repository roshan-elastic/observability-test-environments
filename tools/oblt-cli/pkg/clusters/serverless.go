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
)

const (
	ServerlessTemplateName = "serverless"
)

// ServerlessCluster is the struct that contains the parameters for the Serverless cluster
type ServerlessCluster struct {
	ClusterNamePrefix        string                      `json:"ClusterNamePrefix,omitempty"`
	ClusterNameSuffix        string                      `json:"ClusterNameSuffix,omitempty"`
	GitHubCommentId          string                      `json:"GitHubCommentId,omitempty"`
	GitHubCommit             string                      `json:"GitHubCommit,omitempty"`
	GitHubIssue              string                      `json:"GitHubIssue,omitempty"`
	GitHubPullRequest        string                      `json:"GitHubPullRequest,omitempty"`
	GitHubRepository         string                      `json:"GitHubRepository,omitempty"`
	ObltRepo                 *ObltEnvironmentsRepository `json:"-"`
	ProjectType              string                      `json:"ProjectType,omitempty"`
	SlackChannel             string                      `json:"SlackChannel,omitempty"`
	Target                   string                      `json:"Target,omitempty"`
	TemplateName             string                      `json:"TemplateName,omitempty"`
	Username                 string                      `json:"Username,omitempty"`
	ElasticsearchDockerImage string                      `json:"ElasticsearchDockerImage,omitempty"`
	KibanaDockerImage        string                      `json:"KibanaDockerImage,omitempty"`
	FleetDockerImage         string                      `json:"FleetDockerImage,omitempty"`
}

// ServerlessTemplateParams is the struct that contains the parameters for the Serverless template
type ServerlessTemplateParams struct {
	ProjectType              string `json:"ProjectType,omitempty"`
	Target                   string `json:"Target,omitempty"`
	ElasticsearchDockerImage string `json:"ElasticsearchDockerImage,omitempty"`
	KibanaDockerImage        string `json:"KibanaDockerImage,omitempty"`
	FleetDockerImage         string `json:"FleetDockerImage,omitempty"`
}

// Create creates a new Serverless cluster
func (r *ServerlessCluster) Create() (parameters map[string]interface{}, err error) {
	var json string

	if err = r.validate(); err == nil {
		json, err = r.packTemplateParams()
	}
	if err == nil {
		parameters, err = r.createCluster(json)
	}
	return parameters, err
}

// createCluster creates the cluster configuration file
func (r *ServerlessCluster) createCluster(jsonStr string) (map[string]interface{}, error) {
	var customCluster = &CustomCluster{
		ClusterNamePrefix: r.ClusterNamePrefix,
		ClusterNameSuffix: r.ClusterNameSuffix,
		TemplateName:      ServerlessTemplateName,
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
func (r *ServerlessCluster) packTemplateParams() (jsonContent string, err error) {
	var jsonContentb []byte

	templateParams := ServerlessTemplateParams{
		ProjectType:              r.ProjectType,
		Target:                   r.Target,
		ElasticsearchDockerImage: r.ElasticsearchDockerImage,
		KibanaDockerImage:        r.KibanaDockerImage,
		FleetDockerImage:         r.FleetDockerImage,
	}

	jsonContentb, err = json.Marshal(templateParams)
	if err == nil {
		jsonContent = string(jsonContentb)
	}

	return jsonContent, err
}

// validate validates the parameters
func (r *ServerlessCluster) validate() (err error) {
	return errors.Join(
		config.ValidateProjectType(r.ProjectType),
		config.ValidateTarget(r.Target),
		config.ValidateDockerImage(r.ElasticsearchDockerImage),
		config.ValidateDockerImage(r.KibanaDockerImage),
		config.ValidateDockerImage(r.FleetDockerImage),
	)
}
