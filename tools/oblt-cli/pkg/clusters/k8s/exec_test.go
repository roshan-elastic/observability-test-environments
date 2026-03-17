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
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestK8sShellExecScript(t *testing.T) {
	// create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "test-k8s-shell")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// create a temporary file for the cluster config
	clusterConfigFile, err := os.CreateTemp(tmpDir, "cluster-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(clusterConfigFile.Name())

	// write the cluster config to the file
	_, err = clusterConfigFile.WriteString(`cluster_name: test-oblt`)
	assert.NoError(t, err)

	// create a K8sShell instance
	script := "echo 'hello world'"
	yamlFile, err := files.ReadYamlOj(clusterConfigFile.Name(), "NA")
	assert.NoError(t, err)
	k := NewK8sShell(
		yamlFile,
		script,
		tmpDir,
		true,
	)

	// check that the script is executable
	err = k.ExecScript()
	assert.NoError(t, err)

	// check that the script file was created
	scriptPath := k.getActivateScriptPath()
	assert.FileExists(t, scriptPath)

	// check that the script path is correct
	expectedPath := filepath.Join(tmpDir, "test-oblt-activate.sh")
	assert.Equal(t, expectedPath, scriptPath)

	// check that the script content is correct
	scriptContent, err := os.ReadFile(scriptPath)
	assert.NoError(t, err)
	assert.Equal(t, scriptPlaceholder, string(scriptContent))
}
