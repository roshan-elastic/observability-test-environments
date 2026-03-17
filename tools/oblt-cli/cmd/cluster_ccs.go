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
	"time"

	"github.com/spf13/cobra"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
)

// CcsCmd represents the css command
var CcsCmd = &cobra.Command{
	Use:   "ccs",
	Short: "Command to create a CCS cluster.",
	Long:  `Command to create a Cross Cluster Search cluster. This cluster configure a oblt cluster as remote cluster to use CCS.`,
	Args:  validateDeprecatedCCSFlags,
	Run:   runCcs,
}

func init() {
	CreateCmd.AddCommand(CcsCmd)

	CcsCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster we want to create. This is only supported when running on the CI to help with redeploying clusters. (Optional)")
	CcsCmd.Flags().String(config.RemoteClusterFlag, "", "Oblt cluster to use (release-oblt, dev-oblt, edge-oblt, or a custom name) as a remote cluster for CCS, not a unique name for the cluster. (Required)")
	CcsCmd.Flags().String(config.ClusterNamePrefixFlag, "", "Prefix to be prepended to the randomised cluster name. (Optional)")
	CcsCmd.Flags().String(config.ClusterNameSuffixFlag, "", "Suffix to be appended to the randomised cluster name. If not present, a random seed will be used. This parameter ensures that the name of the cluster is unique, use it with caution. (Optional)")
	CcsCmd.Flags().String(config.RepoFlag, "", "The GitHub repository that requested this cluster, to help with the GitOps automation (Optional)")
	CcsCmd.Flags().String(config.CommitFlag, "", "The GitHub sha commit that requested this cluster, to help with the GitOps automation (Optional)")
	CcsCmd.Flags().String(config.CommentIdFlag, "", "The GitHub comment ID that requested this cluster., to help with the GitOps automation (Optional)")
	CcsCmd.Flags().String(config.IssueFlag, "", "The GitHub issue that requested this cluster, to help with the GitOps automation (Optional)")
	CcsCmd.Flags().String(config.PullRequestFlag, "", "The GitHub Pull Request that requested this cluster, to help with the GitOps automation (Optional)")

	cobra.MarkFlagRequired(CcsCmd.Flags(), config.RemoteClusterFlag)
}

/*
It will create a cluster configuration file based on the parameters passed.
Then this cluster configuration file is commit and pushed to the oblt test environments repo.
*/
func runCcs(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	remoteClusterName, _ := cmd.Flags().GetString(config.RemoteClusterFlag)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	clusterNamePrefix, _ := cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ := cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	repo, _ := cmd.Flags().GetString(config.RepoFlag)
	commit, _ := cmd.Flags().GetString(config.CommitFlag)
	commentId, _ := cmd.Flags().GetString(config.CommentIdFlag)
	issue, _ := cmd.Flags().GetString(config.IssueFlag)
	pullRequest, _ := cmd.Flags().GetString(config.PullRequestFlag)
	wait, _ := cmd.Flags().GetInt(config.WaitFlag)

	labelClusterTemplate := apm.Label{Key: "cluster.template", Value: remoteClusterName}
	tx, ctx := apm.StartTransaction("runCcs", "request", []apm.Label{labelVersion, labelClusterTemplate}, userConfig)
	defer apm.Flush(tx)

	err = validateCIMinimumArguments(cmd, args)
	apm.CobraCheckErr(err, tx, ctx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	ccsCluster := &clusters.CCSCluster{
		TemplateName:      clusters.CCSTemplateName,
		ClusterName:       clusterName,
		ClusterNamePrefix: clusterNamePrefix,
		ClusterNameSuffix: clusterNameSuffix,
		Username:          userConfig.Username,
		SlackChannel:      userConfig.SlackChannel,
		ObltRepo:          obltTestEnvironments,
		GitHubRepository:  repo,
		GitHubCommit:      commit,
		GitHubIssue:       issue,
		GitHubPullRequest: pullRequest,
		GitHubCommentId:   commentId,
		RemoteClusterName: remoteClusterName,
	}

	parametersMap, err := ccsCluster.Create()
	apm.CobraCheckErr(err, tx, ctx)

	if wait > 0 {
		obltTestEnvironments.WaitForClusterCreation(parametersMap["ClusterName"].(string), time.Duration(wait)*time.Minute)
	}

	saveResults(parametersMap, outputFile)
}
