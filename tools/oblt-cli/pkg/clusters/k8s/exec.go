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
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
)

const (
	scriptHeader = `#!/usr/bin/env bash
. "${OBLT_ACTIVATE_SCRIPT}"
oblt-auth
oblt-context
`
	scriptPlaceholder = "#!/usr/bin/env bash\nfunction oblt-auth(){\ntrue\n};\nfunction oblt-context(){\ntrue\n};"
)

type K8sShell struct {
	clusterConfig files.YamlFile
	script        string
	cfgDir        string
	dryRun        bool
}

func NewK8sShell(clusterConfig files.YamlFile, script string, cfgDir string, dryRun bool) *K8sShell {
	return &K8sShell{clusterConfig: clusterConfig, script: script, cfgDir: cfgDir, dryRun: dryRun}
}

// ExecScript Export the environment variables and execute the script
func (k *K8sShell) ExecScript() (err error) {
	path, err := k.saveActivationScript()
	if err == nil {
		os.Setenv("OBLT_ACTIVATE_SCRIPT", path)

		err = syscall.Exec("/bin/bash", []string{"-l", "-c", scriptHeader + k.script}, os.Environ())
	}
	return err
}

// getActivateScriptPath return the path of the activation script
func (k *K8sShell) getActivateScriptPath() string {
	clusterName := k.clusterConfig.Data["cluster_name"].(string)
	return filepath.Join(k.cfgDir, clusterName+"-activate.sh")
}

// saveActivationScript read the activation secret and save the activation script in a file
func (k *K8sShell) saveActivationScript() (path string, err error) {
	path = k.getActivateScriptPath()
	secretValue, err := k.readActivateScript()
	if err == nil {
		files.SaveResultsRaw(fmt.Sprintf("%s", secretValue), path)
	}
	return path, err
}

// readActivateScript read the activation secret from GCSM, in dryRun mode returns a placeholder.
func (k *K8sShell) readActivateScript() (secretValue interface{}, err error) {
	clusterName := k.clusterConfig.Data["cluster_name"].(string)
	if k.dryRun {
		secretValue = scriptPlaceholder
	} else {
		var gcsm *gcp.ClusterSecrets
		if gcsm, err = gcp.NewClusterSecrets(); err == nil {
			secretValue, err = gcsm.ReadActivateSecret(clusterName)
		}
	}
	return secretValue, err
}
