// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package interactions

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/box"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/helper"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
)

var AWSResetAccountAction = modals.Action{
	Name: "Reset your AWS account",
	Desc: fmt.Sprintf("Display a form to reset an AWS account in the `elastic-observability` with the `%s` command", AWSResetAccountCommand()),
	Tldr: fmt.Sprintf("Reset an AWS account with the `%s` command", AWSResetAccountCommand()),
}

// AWSResetAccountCommand what's the command
func AWSResetAccountCommand() string {
	return "/oblt-reset-aws"
}

// resetAWSAccountPostMessage posts message before an AWS account is reset
func resetAWSAccountPostMessage(client *slack.Client, interaction slack.InteractionCallback, email string) error {
	user := interaction.User
	isElastic := strings.Contains(email, "@elastic.co")
	success := false
	if isElastic {
		success = true

		user := interaction.User

		// Generate data from the template
		var data AWSIssueData
		data.User = user.Name
		if user.RealName != "" {
			data.User = user.RealName
		}
		data.Email = email
		body, errInTemplate := helper.TemplateStringToString(string(box.Get("/aws-account-reset.issue")), data)
		// See https://github.com/elastic/observability-robots/pull/1681
		// This automation uses GitHub labels so the GitHub actions can filter what GitHub issues should be filtered
		labels := []string{"aws-account-reset-automation"}
		issue, errGitHubIssue := github.CreateIssueWithLabels(fmt.Sprintf("[AWS account]: reset %s", email), body, labels, "elastic", "observability-robots")
		// If the issue is nil, we need to create a fake issue to display the URL
		url := "http://invalid-issue-report.example.com"
		if issue != nil {
			url = *issue.HTMLURL
		}
		success = errInTemplate == nil && errGitHubIssue == nil
		fmt.Printf("Issue has been created: %v\n", url)
	}
	blocks, _ := modals.RenderResetAWSAccount(user.ID, email, isElastic, success)
	_, _, err := client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return nil
}

// HandleResetAWSAccount handles the interaction with the select AWS for reset an existing account
func HandleResetAWSAccount(client *slack.Client, interaction slack.InteractionCallback) error {
	var email string

	viewState := interaction.View.State

	// pick up all values from the aws-onboarding form
	for k, v := range viewState.Values {
		switch k {
		case "email":
			email = v["email"].Value
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}

	return resetAWSAccountPostMessage(client, interaction, email)
}

// ShowResetAWSAccountForm creates a ModalViewRequest to gather the email to be used
func ShowResetAWSAccountForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowResetAWSAccountForm", "request", triggerID, user)
	defer tx.End()

	modalRequest, err := modals.AWSResetModal(AWSResetAccountAction.Name, triggerID)
	if err != nil {
		return err
	}

	_, err = client.OpenView(triggerID, modalRequest)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}
