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

package k8s

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Integrations struct {
	Name           string `yaml:"name"`
	NameNormalized string `yaml:"name_normalized"`
	Version        string `yaml:"version"`
}
type ClusterConfigIntegration struct {
	ClusterName  string       `yaml:"cluster_name"`
	Integrations Integrations `yaml:"integrations"`
}

// ExecPackagesScript Export the environment variables and execute the script
func (k *K8sShell) ExecPackagesScript() (err error) {
	err = k.setPackagesEnvVars()
	if err == nil {
		err = k.ExecScript()
	} else {
		err = errors.New("there is no integration configured, so no script executed")
	}
	return err
}

// setPackagesEnvVars set the environment variables needed by the packages shell
func (k *K8sShell) setPackagesEnvVars() (err error) {
	var clusterConfig ClusterConfigIntegration
	err = yaml.Unmarshal(k.clusterConfig.Bytes, &clusterConfig)
	if err == nil {
		os.Setenv("PACKAGE", clusterConfig.Integrations.Name)
		os.Setenv("PACKAGE_NORMALIZED", clusterConfig.Integrations.NameNormalized)
		os.Setenv("INTEGRATIONS_BRANCH", clusterConfig.Integrations.Version)
	}
	return err
}
