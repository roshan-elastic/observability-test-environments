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

// ReleasesCmd represents the releases command
var ReleasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Command to list the releases.",
	Long:  `Command to list the releases in the Unified Release.`,
	Run:   runReleaseList,
}

func init() {
	UnifiedReleaseCmd.AddCommand(ReleasesCmd)
}

/*
It will query the artifacts api to gather all the current
releases or build candidates in the Unified Release.
*/
func runReleaseList(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, _ := apm.StartTransaction("runReleaseList", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	releases, _ := releases.GetReleases()

	var results = make(map[string]interface{})
	var items []interface{}

	data := [][]string{}
	for _, release := range releases {
		releaseVersion := release.Version
		releaseBuild := release.BuildId
		releaseDate := release.ReleaseDate
		releaseQuery := release.Query
		items = append(items, map[string]interface{}{"version": releaseVersion, "buildId": releaseBuild, "creationDate": releaseDate, "api": releaseQuery})

		data = append(data, []string{releaseVersion, releaseBuild, releaseDate, releaseQuery})
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Version", "BuildId", "Date", "API"})
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	results["items"] = items
	saveResults(results, outputFile)
}
