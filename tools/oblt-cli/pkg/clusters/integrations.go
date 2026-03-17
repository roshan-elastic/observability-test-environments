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

const IntegrationsTemplateName = "ess-integrations-v4"

// IntegrationCluster represents the structure of the integration cluster
type IntegrationCluster struct {
	Branch            string                      `json:"Branch,omitempty"`
	ClusterName       string                      `json:"ClusterName,omitempty"`
	ClusterNamePrefix string                      `json:"ClusterNamePrefix,omitempty"`
	ClusterNameSuffix string                      `json:"ClusterNameSuffix,omitempty"`
	Integration       string                      `json:"Integration,omitempty"`
	IsRelease         bool                        `json:"IsRelease"`
	Repository        string                      `json:"Repository,omitempty"`
	StackVersion      string                      `json:"StackVersion,omitempty"`
	ObltRepo          *ObltEnvironmentsRepository `json:"-"`
	SlackChannel      string                      `json:"SlackChannel,omitempty"`
	TemplateName      string                      `json:"TemplateName,omitempty"`
	Username          string                      `json:"Username,omitempty"`
}

type IntegrationClusterTemplate struct {
	Branch                   string `json:"Branch,omitempty"`
	ClusterName              string `json:"ClusterName,omitempty"`
	ElasticAgentDockerImage  string `json:"ElasticAgentDockerImage,omitempty"`
	ElasticsearchDockerImage string `json:"ElasticsearchDockerImage,omitempty"`
	Integration              string `json:"Integration,omitempty"`
	KibanaDockerImage        string `json:"KibanaDockerImage,omitempty"`
	Repository               string `json:"Repository,omitempty"`
	StackVersion             string `json:"StackVersion,omitempty"`
}

// IntegrationService represents the structure of the integration service
type IntegartionServiceTemplate struct {
	Integration string `json:"Integration,omitempty"`
	Repository  string `json:"Repository,omitempty"`
	Branch      string `json:"Branch,omitempty"`
}

// CreateESSIntegrationsCluster creates a cluster configuration file for the template ess-integrations.
func (r *IntegrationCluster) Create() (parameters map[string]interface{}, err error) {
	var jsonData string
	if err = r.validate(); err == nil {
		if jsonData, err = r.packClusterTemplateParams(); err == nil {
			parameters, err = r.createCluster(jsonData)
		}
	}
	return parameters, err
}

// CreateESSIntegrationsService creates a cluster configuration file for the template ess-integrations deploying only a service.
func (r *IntegrationCluster) CreateService() (parameters map[string]interface{}, err error) {
	var jsonData string
	if err = r.validate(); err == nil {
		if jsonData, err = r.packServiceTemplateParams(); err == nil {
			parameters, err = r.createCluster(jsonData)
		}
	}
	return parameters, err
}

// createCluster creates a cluster configuration file for the template ess-integrations.
func (r *IntegrationCluster) createCluster(json string) (map[string]interface{}, error) {
	var customCluster = &CustomCluster{
		ClusterName:       r.ClusterName,
		ClusterNamePrefix: r.ClusterNamePrefix,
		ClusterNameSuffix: r.ClusterNameSuffix,
		TemplateName:      IntegrationsTemplateName,
		Parameters:        string(json),
		ObltRepo:          r.ObltRepo,
		SlackChannel:      r.SlackChannel,
		Username:          r.Username,
	}
	return customCluster.Create()
}

// packServiceTemplateParams packs the parameters for the template ess-integrations.
func (r *IntegrationCluster) packServiceTemplateParams() (jsonContent string, err error) {
	if err == nil {
		templateParams := IntegartionServiceTemplate{
			Integration: r.Integration,
			Repository:  r.Repository,
			Branch:      r.Branch,
		}

		var jsonContentb []byte
		jsonContentb, err = json.Marshal(templateParams)
		if err == nil {
			jsonContent = string(jsonContentb)
		}
	}
	return jsonContent, err
}

// packClusterTemplateParams packs the parameters for the template ess-integrations.
func (r *IntegrationCluster) packClusterTemplateParams() (jsonContent string, err error) {
	if err == nil {
		stackVersion, dockerImageVersion := normalizeElasticStackVersion(r.StackVersion, r.IsRelease)

		templateParams := IntegrationClusterTemplate{
			StackVersion:             stackVersion,
			ElasticsearchDockerImage: buildEsDockerImageString(dockerImageVersion),
			KibanaDockerImage:        buildKibanaDockerImageString(dockerImageVersion),
			ElasticAgentDockerImage:  buildElasticAgentDockerImageString(dockerImageVersion),
			Integration:              r.Integration,
			Repository:               r.Repository,
			Branch:                   r.Branch,
		}

		var jsonContentb []byte
		jsonContentb, err = json.Marshal(templateParams)
		jsonContent = string(jsonContentb)
	}
	return jsonContent, err
}

func (r *IntegrationCluster) validate() (err error) {
	return errors.Join(
		config.ValidateClusterName(r.ClusterName),
	)
}
