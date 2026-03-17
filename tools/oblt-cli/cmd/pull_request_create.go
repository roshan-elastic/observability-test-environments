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
	git "github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/git"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// PullRequestCreateCmd represents the create command
var PullRequestCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Command to create a Pull Request.",
	Long:    `Command to create a Pull Request in the current directory.`,
	Example: `$ oblt-cli ci pull-request create --title "oblt-cli: update cluster" --body "chore" --labels "skip-changelog,test"`,
	Run:     pullRequestCreate,
}

func init() {
	PullRequestCmd.AddCommand(PullRequestCreateCmd)

	PullRequestCreateCmd.Flags().String(config.BaseFlag, "", "The branch into which you want your code merged. (default: main)")
	PullRequestCreateCmd.Flags().String(config.BodyFlag, "", "GitHub Pull Request body. (Required)")
	PullRequestCreateCmd.Flags().String(config.HeadFlag, "", "The branch that contains commits for your pull request (default: current branch)")
	PullRequestCreateCmd.Flags().String(config.TitleFlag, "", "GitHub Pull Request title. (Required)")
	PullRequestCreateCmd.Flags().StringSlice(config.LabelsFlag, []string{}, `GitHub labels to apply, ex. --labels=cluster-creation,new-cluster,ephemeral (Optional)`)

	cobra.MarkFlagRequired(PullRequestCreateCmd.Flags(), config.BodyFlag)
	cobra.MarkFlagRequired(PullRequestCreateCmd.Flags(), config.TitleFlag)
}

// pullRequestCreate is the function executed to create a pull request
func pullRequestCreate(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("pullRequestCreate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	base, _ := cmd.Flags().GetString(config.BaseFlag)
	body, _ := cmd.Flags().GetString(config.BodyFlag)
	head, _ := cmd.Flags().GetString(config.HeadFlag)
	title, _ := cmd.Flags().GetString(config.TitleFlag)
	labels, _ := cmd.Flags().GetStringSlice(config.LabelsFlag)

	if base == "" {
		base = git.DefaultBranch
	}
	dir, err := os.Getwd()
	apm.CobraCheckErr(err, tx, ctx)
	NewRepository := git.NewRepository(dir, git.DefaultOwner, git.ObltRepoName, base, dryRun, false, true)

	if head == "" {
		head, err = NewRepository.GetCurrentBranch()
		apm.CobraCheckErr(err, tx, ctx)
	}

	changes, _ := NewRepository.GetChanges(base)
	logger.Debugf("ChangeSet %+q", changes)
	newLabels := append(labels, git.LabelsClusterType(changes)...)
	logger.Infof("Inferred GitHub Labels %+q", newLabels)
	pullRequest, response, err := github.CreatePullRequestWithLabels(title, body, git.DefaultOwner, git.ObltRepoName, base, head, newLabels)
	logger.Debugf("Pull Request response %s", response.String())
	apm.CobraCheckErr(err, tx, ctx)
	if err == nil {
		logger.Infof("Pull Request has been created %s", pullRequest.GetHTMLURL())
	}
}
