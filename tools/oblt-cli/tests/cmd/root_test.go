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
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRootCommandNoArgs(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd)
	assert.NoError(t, err)
	assert.Contains(t, output, "CLI tool to operate the oblt clusters")
}

func TestRootCommandHelp(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "CLI tool to operate the oblt clusters")
}

func TestConfigureCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "configure", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "It configures the Slack member ID")
}

func TestCusterCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "cluster", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "Command to operate a cluster")
}

func TestClusterCreateCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "cluster", "create", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "Command to create a cluster.")
}

func TestClusterCreateCssCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "cluster", "create", "ccs", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "Command to create a Cross Cluster Search cluster.")
}

func TestClusterDestroyCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "cluster", "destroy", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "Command to destroy a cluster.")
}

func TestClusterUpdateCommand(t *testing.T) {
	output, err := executeCommand(t, cmd.RootCmd, "cluster", "update", "--help")
	assert.NoError(t, err)
	assert.Contains(t, output, "Command to update a cluster.")
}
