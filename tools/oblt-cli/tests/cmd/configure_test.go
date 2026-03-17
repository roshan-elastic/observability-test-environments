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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCommandslackChannel(t *testing.T) {
	dir, err := os.MkdirTemp("", "oblt-cli")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	cfgFile := filepath.Join(dir, "config.yml")
	mode := ""
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		mode = "--git-http-mode"
	}

	_, err = executeCommand(t, cmd.RootCmd, "configure", "--slack-channel=#foo", "--username=bar", "--config="+cfgFile, mode, "--verbose", "--dry-run")
	assert.NoError(t, err)

	// empty configuration
	viper.Reset()
	slackChannel := viper.GetString(config.SlackChannelFlag)
	username := viper.GetString(config.UsernameFlag)
	assert.Equal(t, username, "")
	assert.Equal(t, slackChannel, "")

	// check written configuration
	viper.SetConfigFile(cfgFile)
	viper.ReadInConfig()
	slackChannel = viper.GetString(config.SlackChannelFlag)
	username = viper.GetString(config.UsernameFlag)
	assert.Equal(t, slackChannel, "#foo")
	assert.Equal(t, username, "bar")

	// check that we can overwrite the configuration
	_, err = executeCommand(t, cmd.RootCmd, "version", "--slack-channel=#foo1", "--username=bar1", "--config="+cfgFile, mode, "--verbose", "--save-config", "--dry-run")
	assert.NoError(t, err)
	slackChannel = viper.GetString(config.SlackChannelFlag)
	username = viper.GetString(config.UsernameFlag)
	assert.Equal(t, slackChannel, "#foo1")
	assert.Equal(t, username, "bar1")
	viper.Reset()

	// check written configuration
	viper.SetConfigFile(cfgFile)
	viper.ReadInConfig()
	slackChannel = viper.GetString(config.SlackChannelFlag)
	username = viper.GetString(config.UsernameFlag)
	assert.Equal(t, slackChannel, "#foo1")
	assert.Equal(t, username, "bar1")
}
