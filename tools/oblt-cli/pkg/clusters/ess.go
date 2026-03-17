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

const ESSTemplateName = "ess"

// ESSCluster represents the structure of the ess cluster
type ESSCluster struct {
	ClusterName       string                      `json:"ClusterName,omitempty"`
	ClusterNamePrefix string                      `json:"ClusterNamePrefix,omitempty"`
	ClusterNameSuffix string                      `json:"ClusterNameSuffix,omitempty"`
	IsRelease         bool                        `json:"IsRelease"`
	ObltRepo          *ObltEnvironmentsRepository `json:"-"`
	StackVersion      string                      `json:"StackVersion,omitempty"`
	SlackChannel      string                      `json:"SlackChannel,omitempty"`
	TemplateName      string                      `json:"TemplateName,omitempty"`
	Username          string                      `json:"Username,omitempty"`
}

type ESSTemplateParams struct {
	ElasticAgentDockerImage  string `json:"ElasticAgentDockerImage,omitempty"`
	ElasticsearchDockerImage string `json:"ElasticsearchDockerImage,omitempty"`
	KibanaDockerImage        string `json:"KibanaDockerImage,omitempty"`
	StackVersion             string `json:"StackVersion,omitempty"`
}

// CreateESSCluster creates the configuration file for the ESS cluster.
// check the template file path, reading the YAML file represented by that path
func (r *ESSCluster) Create() (parameters map[string]interface{}, err error) {
	var jsonData string
	if err = r.validate(); err == nil {
		if jsonData, err = r.packTemplateParams(); err == nil {
			parameters, err = r.createCluster(jsonData)
		}
	}
	return parameters, err
}

// CreateESSCluster creates the configuration file for the ESS cluster.
func (r *ESSCluster) createCluster(json string) (map[string]interface{}, error) {
	var customCluster = &CustomCluster{
		ClusterName:       r.ClusterName,
		ClusterNamePrefix: r.ClusterNamePrefix,
		ClusterNameSuffix: r.ClusterNameSuffix,
		TemplateName:      ESSTemplateName,
		Parameters:        string(json),
		ObltRepo:          r.ObltRepo,
		SlackChannel:      r.SlackChannel,
		Username:          r.Username,
	}
	return customCluster.Create()
}

// packClusterTemplateParams packs the cluster template parameters into a JSON string.
func (r *ESSCluster) packTemplateParams() (jsonContent string, err error) {
	stackVersion, dockerImageVersion := normalizeElasticStackVersion(r.StackVersion, r.IsRelease)

	templateParams := ESSTemplateParams{
		StackVersion:             stackVersion,
		ElasticsearchDockerImage: buildEsDockerImageString(dockerImageVersion),
		KibanaDockerImage:        buildKibanaDockerImageString(dockerImageVersion),
		ElasticAgentDockerImage:  buildElasticAgentDockerImageString(dockerImageVersion),
	}

	var jsonContentb []byte
	jsonContentb, err = json.Marshal(templateParams)
	if err == nil {
		jsonContent = string(jsonContentb)
	}
	return jsonContent, err
}

func (r *ESSCluster) validate() (err error) {
	return errors.Join(
		config.ValidateClusterName(r.ClusterName),
		config.ValidateSemVer(r.StackVersion),
	)
}
