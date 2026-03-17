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
//
// This file contains the functions to validate a YAML configuration file against a JSON schema.

package schema

import (
	"errors"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/api"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
)

// Validate validates a YAML configuration file against a JSON schema
func Validate(filename string) (err error) {
	var clusterConfig *api.ClusterConfig
	if clusterConfig, err = parseConfigFile(filename); err == nil {
		err = errors.Join(
			checkClusterName(clusterConfig),
			checkStackVersion(clusterConfig),
			checkDockerImages(clusterConfig),
		)
	}
	return err
}

// parseConfigFile reads a YAML configuration file and returns a ClusterConfig object
func parseConfigFile(filename string) (clusterConfig *api.ClusterConfig, err error) {
	var clusterConfigYaml files.YamlFile
	if clusterConfigYaml, err = files.ReadYamlOj(filename, "nobody"); err == nil {
		clusterConfig, err = clusters.NewClusterConfig(clusterConfigYaml)
	}
	return clusterConfig, err
}

// checkClusterName validates the cluster_name field
func checkClusterName(clusterConfig *api.ClusterConfig) (err error) {
	if clusterConfig.ClusterName == nil {
		err = errors.New("cluster_name is required")
	} else {
		err = config.ValidateClusterName(*clusterConfig.ClusterName)
	}
	return err
}

// checkStackVersion validates the stack.version field
func checkStackVersion(clusterConfig *api.ClusterConfig) (err error) {
	if clusterConfig.Stack != nil && clusterConfig.Stack.Version != nil {
		err = config.ValidateSemVer(*clusterConfig.Stack.Version)
	}
	return err
}

// checkDockerImages validates the Docker images fields
func checkDockerImages(clusterConfig *api.ClusterConfig) (err error) {
	if clusterConfig.Stack != nil {
		err = errors.Join(checkDockerImagesEss(clusterConfig), checkDockerImagesEck(clusterConfig))
	}
	return err
}

// checkDockerImagesEss validates the Docker images fields for the Elastic Stack ESS type
func checkDockerImagesEss(clusterConfig *api.ClusterConfig) (err error) {
	ess := clusterConfig.Stack.Ess
	if ess != nil {
		if ess.Elasticsearch != nil && ess.Elasticsearch.Image != nil {
			err = config.ValidateDockerImage(*ess.Elasticsearch.Image)
		}
		if ess.Kibana != nil && ess.Kibana.Image != nil {
			err = errors.Join(err, config.ValidateDockerImage(*ess.Kibana.Image))
		}
		if ess.Integrations != nil && ess.Integrations.Image != nil {
			err = errors.Join(err, config.ValidateDockerImage(*ess.Integrations.Image))
		}
		if ess.Profiling != nil && ess.Profiling.Image != nil {
			err = errors.Join(err, config.ValidateDockerImage(*ess.Profiling.Image))
		}
	}
	return err
}

// checkDockerImagesEck validates the Docker images fields for the Elastic Stack ECK type
func checkDockerImagesEck(clusterConfig *api.ClusterConfig) (err error) {
	eck := clusterConfig.Stack.Eck
	if eck != nil {
		if eck.Elasticsearch != nil && eck.Elasticsearch.Image != nil {
			err = config.ValidateDockerImage(*eck.Elasticsearch.Image)
		}
		if eck.Kibana != nil && eck.Kibana.Image != nil {
			err = errors.Join(err, config.ValidateDockerImage(*eck.Kibana.Image))
		}
		if eck.Agent != nil && eck.Agent.Image != nil {
			err = errors.Join(err, config.ValidateDockerImage(*eck.Agent.Image))
		}
	}
	return err
}
