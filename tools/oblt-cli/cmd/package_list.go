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
	"os"
	"syscall"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// PackagesListCmd represents the packages list command
var PackagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Command to list the integrations available.",
	Long: `Command to list the integrations available.
	With this command you can list the integrations available to install in your cluster.
	Alternatively you can save the list in a JSON file with the flag --output-file.`,
	Run: runPackagesList,
}

func init() {
	PackagesCmd.AddCommand(PackagesListCmd)
}

// TODO implement it in Go
// runList shows the list of integrations available
func runPackagesList(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	tx, ctx := apm.StartTransaction("runPackagesList", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	cfgDir := userConfig.GetDir()
	os.Setenv("OBLT_CLI_HOME", cfgDir)
	script := `
		cd ${OBLT_CLI_HOME}
		if [ ! -d "integrations" ]; then
			git clone https://github.com/elastic/integrations.git
		fi
		cd integrations
		git pull
		ls packages
		exit 0
	`
	err = syscall.Exec("/bin/bash", []string{"-l", "-c", script}, os.Environ())
	apm.CobraCheckErr(err, tx, ctx)
}
