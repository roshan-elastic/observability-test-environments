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

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/releases"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// BranchesCmd represents the branches command
var BranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "Command to list the active branches.",
	Long:  `Command to list the active branches in the Unified Release.`,
	Run:   runBranches,
}

func init() {
	UnifiedReleaseCmd.AddCommand(BranchesCmd)
}

/*
It will query the unified release.
*/
func runBranches(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, _ := apm.StartTransaction("runBranches", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		apm.ReportError(tx, "please run with the GITHUB_TOKEN=YourToken")
	}
	os.Setenv("GITHUB_PASSWORD", token)

	versions, _ := releases.GetVersions()

	var results = make(map[string]interface{})
	var items []interface{}

	data := [][]string{}

	for _, entry := range versions {
		branch := entry.Branch
		version := entry.Version
		releaseDate := entry.ReleaseDate
		items = append(items, map[string]interface{}{"branch": branch, "version": version, "releaseDate": releaseDate})

		data = append(data, []string{branch, version, releaseDate})
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Branch", "Version", "Release Date"})
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	results["items"] = items
	saveResults(results, outputFile)
}
