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
	"fmt"
	"runtime"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/box"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

// ShellCmd represents the shell command
var ShellCmd = &cobra.Command{
	Use:   "shell-functions",
	Short: "Command to prepare the shell with the set of functions.",
	Long:  `Command to prepare the shell with the set of common functions.`,
	Example: `$ oblt-cli ci shell-functions --output-file "$(pwd)/shell-functions.sh"
	$ eval "$(cat $(pwd)/shell-functions.sh)"`,
	Run: shell,
}

func init() {
	CICmd.AddCommand(ShellCmd)
}

func shell(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("shell", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	if runtime.GOOS == "windows" {
		apm.CobraCheckErr(fmt.Errorf("unsupported platform"), tx, ctx)
	}
	shellFunctions := string(box.Get("/shell-functions"))
	fmt.Printf("%s", shellFunctions)
	if len(outputFile) > 0 {
		saveResultsRaw(shellFunctions, outputFile)
	}
}
