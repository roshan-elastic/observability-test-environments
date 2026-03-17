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
	"os"

	"github.com/blang/semver"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/troubleshoot"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/prompt"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// UpdateCLICmd represents the update CLI command
var UpdateCLICmd = &cobra.Command{
	Use:   "update",
	Short: "Command to self-update the tool.",
	Long:  `Command to self-update the tool.`,
	Run:   runUpdateCLI,
}

func init() {
	RootCmd.AddCommand(UpdateCLICmd)

	UpdateCLICmd.Flags().String(config.GithubTokenFlag, "", "Github token used to look for new releases.")
}

func runUpdateCLI(cmd *cobra.Command, args []string) {
	logger.Debugf("cli.self-update")

	// TODO: process.args contain sensitive data when instrumenting this command with apm.StartTransaction
	//       as long as that's the case we won't send traces for this particular command.

	token, _ := cmd.Flags().GetString(config.GithubTokenFlag)
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	} else {
		if err := prompt.ValidateGithubToken(token); err != nil {
			troubleshoot.CobraCheckErrWithError("Github token failed", err)
		}
	}
	if token == "" {
		// no token was passed as flag: prompt the user for a valid Github token
		token = prompt.GithubToken("Github token used to download the new release")
	}

	// to query the releases, the self-update library creates an HTTP request using GITHUB_TOKEN
	os.Setenv("GITHUB_TOKEN", token)

	if err := selfUpdate(currentVersion); err != nil {
		troubleshoot.CobraCheckErrWithError("Binary update failed.", err)
	}
}

func selfUpdate(current semver.Version) (err error) {
	obltRepo := clusters.NewObltTestEnvironmentsFromViper(true)
	latest, err := selfupdate.UpdateSelf(current, obltRepo.OwnerRepo())

	if err == nil && latest.AssetID == 0 {
		err = fmt.Errorf("an error occurred while looking for a new release, check yout GitHub token meet the requirements")
	}

	if err == nil {
		if latest.Version.Equals(currentVersion) {
			// latest version is the same as current version. It means current binary is up to date.
			err = fmt.Errorf("current binary is in the latest version: %s", currentVersion)
		} else {
			logger.Infof("Successfully updated to version %s", latest.Version)
			logger.Infof("Release note: %s", latest.ReleaseNotes)
		}
	}

	return err
}
