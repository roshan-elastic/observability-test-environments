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

package cmd

import (
	"strings"

	"github.com/blang/semver"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/box"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// currentVersion version of the tool, read from the static files
var currentVersion semver.Version

// labelVersion APM label version of the tool, to be used by the other commands
var labelVersion apm.Label

func init() {
	v := strings.ReplaceAll(string(box.Get("/.version")), "\n", "")
	currentVersion = semver.MustParse(v)

	labelVersion = apm.Label{Key: "version", Value: currentVersion.String()}
}

// UpdateCLICmd represents the update CLI command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Command to show tool's version.",
	Long:  "Command to show tool's version.",
	Run:   runShowVersion,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

func runShowVersion(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runShowVersion", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	// the version is show in the init function from the root command so no operation need here.

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	if err == nil {
		version, err := obltTestEnvironments.GetVersion()
		apm.CobraCheckErr(err, tx, ctx)
		err = verifyCompatibility(version)
		apm.CobraCheckErr(err, tx, ctx)
	}
}
