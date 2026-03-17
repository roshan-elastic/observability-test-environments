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

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/cmd"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestClusterCssCommand(t *testing.T) {
	t.Skip("Skipping testing since it failed with a non-existing cluster name")
	dir, err := os.MkdirTemp("", "oblt-cli")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	cfgFile := filepath.Join(dir, "config.yml")
	viper.SetConfigFile(cfgFile)
	viper.Set(config.UsernameFlag, "baruser")
	viper.Set(config.SlackChannelFlag, "#fooChannel")
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	_, err = executeCommand(t, cmd.RootCmd, "cluster", "create", "ccs", "--config="+cfgFile, "--remote-cluster=foocluster", "--stack-version=barVersion", "--dry-run")
	assert.NoError(t, err)

	username := viper.GetString(config.UsernameFlag)
	clusterName, _ := cmd.CcsCmd.Flags().GetString(config.RemoteClusterFlag)
	clusterNamePrefix, _ := cmd.CcsCmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterCfgFilename := clusterNamePrefix + "-" + clusterName + "-ccs.yml"
	cfgFile = viper.ConfigFileUsed()
	cfgDir := filepath.Dir(cfgFile)
	clusterConfigFile := filepath.Join(cfgDir, clusters.EnvFilesFolder, username, clusterCfgFilename)

	yamlfile, err1 := os.ReadFile(clusterConfigFile)
	assert.NoError(t, err1)
	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlfile, &data)
	assert.NoError(t, err)

	assert.Equal(t, data["cluster_name"].(string), "baruser-foocluster")
	assert.Equal(t, data["stack_version"].(string), "barVersion")
	assert.Equal(t, data["slack_channel"].(string), "#fooChannel")
}
