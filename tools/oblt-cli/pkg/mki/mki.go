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

package mki

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/api"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/vpn"
)

const (
	scriptHeader = `#!/usr/bin/env bash
. "${OBLT_MKI_SCRIPT}"
. "${OBLT_KIBANA_SCRIPT}"
`
	scriptPlaceholder = "#!/usr/bin/env bash\nfunction install::vpn(){\ntrue\n};\nfunction elastic::vpn(){\ntrue\n};"
	mkiLib            = ".ci/scripts/lib/mki.sh"
	kibanaLib         = ".ci/scripts/lib/kibana.sh"
	loginScipt        = `
		mki::login ${OBLT_MKI_ENVIRONMENT}
	`
	loginK8sScript = `
		log::info "Getting MKI project details"
		log::info "CLUSTER_NAME: ${CLUSTER_NAME}"
		log::info "PROJECT_ID: ${PROJECT_ID}"
		log::info "ENVIRONMENT: ${OBLT_MKI_ENVIRONMENT}"
		log::info "CONSOLE: ${CONSOLE}"
		log::info "PROJECT_TYPE: ${PROJECT_TYPE}"

		PROJECT_JSON=$(mki::get-project "${CONSOLE}" "${PROJECT_TYPE}" "${CONSOLE_API_KEY}" "${PROJECT_ID}")
		KIBANA_K8S_CLUSTER=$(echo "${PROJECT_JSON}"|jq -r ".clusters.kibana")

		log::info "CLUSTER_NAME: ${CLUSTER_NAME}"
		log::info "PROJECT_ID: ${PROJECT_ID}"
		log::info "ENVIRONMENT: ${OBLT_MKI_ENVIRONMENT}"
		log::info "KIBANA_K8S_CLUSTER: ${KIBANA_K8S_CLUSTER}"

		mki::login-k8s "${KIBANA_K8S_CLUSTER}"
		k8s::set-current-context-namespace "${PROJECT_ID_K8S}"
		`
)

type MKIConfig struct {
	clusterConfig *api.ClusterConfig
	environment   string
	Script        string
	repoPath      string
	obltCliHome   string
	dryRun        bool
	startVpn      bool
	mkiLogin      bool
	apiKey        string
}

func NewMKIConfig(clusterConfig files.YamlFile, obltCliHome, repoPath string, startVpn, mkiLogin, dryRun bool, apiKey string) *MKIConfig {
	config, err := clusters.NewClusterConfig(clusterConfig)
	var environment string
	if err == nil {
		environment = *config.Stack.Target
	} else {
		environment = "qa"
	}
	return &MKIConfig{
		clusterConfig: config,
		environment:   environment,
		Script:        "",
		repoPath:      repoPath,
		obltCliHome:   obltCliHome,
		dryRun:        dryRun,
		startVpn:      startVpn,
		mkiLogin:      mkiLogin,
		apiKey:        apiKey,
	}
}

// Login to the MKI environment and authenticate in the k8s cluster used by the project
func (c *MKIConfig) Login() (err error) {
	scriptToRun := loginK8sScript
	if c.mkiLogin {
		scriptToRun = fmt.Sprintf(`
			%s
			%s
			`, loginScipt, loginK8sScript)
	}
	c.Script = fmt.Sprintf(`
		%s
		set -eo pipefail
		%s
		# replace the current shell
		bash
		`, scriptHeader, scriptToRun)
	err = c.runInsideVPN()
	return err
}

// setBaseEnv set the base environment variables to be used by the MKI scripts
func (c *MKIConfig) setBaseEnv() (err error) {
	clusterName := *c.clusterConfig.ClusterName
	if gcsm, err := gcp.NewClusterSecrets(); err == nil {
		clusterState, err := gcsm.ReadClusterStateSecret(clusterName)
		var serverelessCredentials gcp.ServerlessCredentials
		if err == nil {
			serverelessCredentials, err = gcsm.ReadServerlessCredentials(c.environment)
		}
		if err == nil {
			os.Setenv("OBLT_MKI_SCRIPT", filepath.Join(c.repoPath, mkiLib))
			os.Setenv("OBLT_KIBANA_SCRIPT", filepath.Join(c.repoPath, kibanaLib))
			os.Setenv("OBLT_CLI_HOME", c.obltCliHome)
			os.Setenv("CLUSTER_NAME", clusterName)
			os.Setenv("OBLT_MKI_ENVIRONMENT", c.environment)
			os.Setenv("PROJECT_ID", clusterState.ServerlessDeploymentId)
			os.Setenv("PROJECT_ID_K8S", fmt.Sprintf("project-%s", clusterState.ServerlessDeploymentId))
			os.Setenv("CONSOLE", serverelessCredentials.Console)
			if c.apiKey == "" {
				os.Setenv("CONSOLE_API_KEY", serverelessCredentials.ConsoleApiKey)
			} else {
				os.Setenv("CONSOLE_API_KEY", c.apiKey)
			}
			os.Setenv("ELASTICSEARCH_HOST", clusterState.EsSecret.Url)
			os.Setenv("ELASTICSEARCH_PASSWORD", clusterState.EsSecret.Password)
			os.Setenv("ELASTICSEARCH_USERNAME", clusterState.EsSecret.Username)
			os.Setenv("PROJECT_TYPE", clusterState.StackTemplate)
			if logger.Verbose {
				os.Setenv("DEBUG", "true")
			}
		}
	}

	return err
}

// GenerateKibanaConfig generates a Kibana yaml config file to connet to a serverless project.
func (c *MKIConfig) GenerateKibanaConfig(kibanaYamlPath string) (err error) {
	os.Setenv("KIBANA_FILE", kibanaYamlPath)
	c.Script = fmt.Sprintf(`
		%s
		set -eo pipefail
		%s
		mki::generate-kibana-yml "${PROJECT_ID_K8S}" "${KIBANA_FILE}"
		log::info "Kibana yaml config file generated at ${KIBANA_FILE}"
		`, scriptHeader, loginScipt)
	err = c.runInsideVPN()
	return err
}

// RunLocalKibana runs a local Kibana instance using the Kibana yaml config file generated by GenerateKibanaConfig
func (c *MKIConfig) RunLocalKibana(kibanaYamlPath, KibanaSource string) (err error) {
	os.Setenv("KIBANA_FOLDER", KibanaSource)
	os.Setenv("KIBANA_FILE", kibanaYamlPath)
	c.Script = fmt.Sprintf(`
		%s
		set -eo pipefail
		%s
		mki::generate-kibana-yml "${PROJECT_ID_K8S}" "${KIBANA_FILE}"
		log::info "Kibana yaml config file generated at ${KIBANA_FILE}"
		log::info "Running Kibana from ${KIBANA_FOLDER} with config ${KIBANA_FILE}"
		kibana::run-dev "${KIBANA_FOLDER}" "${KIBANA_FILE}"
		`, scriptHeader, loginScipt)
	err = c.runInsideVPN()
	return err
}

// runInsideVPN runs the script inside a VPN connection
func (c *MKIConfig) runInsideVPN() (err error) {
	if c.startVpn {
		vpnClient := vpn.NewVPNConfig(c.obltCliHome, c.repoPath, c.dryRun)
		err = vpnClient.Connect()
	}
	if err == nil {
		err = c.RunScript()
	}
	return err
}

// RunScript runs the script
func (c *MKIConfig) RunScript() (err error) {
	if err = c.setBaseEnv(); err == nil {
		err = syscall.Exec("/bin/bash", []string{"-l", "-c", c.Script}, os.Environ())
	}
	return err
}
