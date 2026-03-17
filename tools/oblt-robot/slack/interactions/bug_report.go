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
	"errors"
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

const (
	org           = "elastic"
	repo          = "observability-robots-playground"
	issueTemplate = "/bug-report.issue"
)

type bugReportData struct {
	Title            string
	ExpectedBehavior string
	CurrentBehavior  string
	Suggestions      string
	Context          string
	User             string
	labels           []string
}

var BugReportAction = modals.Action{
	Name: "Submit a Bug Report",
	Desc: fmt.Sprintf("Display a form to submit a Bug with the `%s` command", BugReportCommand()),
	Tldr: fmt.Sprintf("Submit a Bug with the `%s` command", BugReportCommand()),
}

// BugReportCommand what's the command
func BugReportCommand() string {
	return "/oblt-bug-report"
}

// createBugReportPostMessage posts message before the Bug Report is created
func createBugReportPostMessage(client *slack.Client, interaction slack.InteractionCallback, data bugReportData) error {
	user := interaction.User

	// Generate data from the template
	data.User = user.Name
	if user.RealName != "" {
		data.User = user.RealName
	}
	body, errInTemplate := helper.TemplateStringToString(string(box.Get(issueTemplate)), data)
	issue, errGitHubIssue := github.CreateIssue(data.Title, body, org, repo, data.labels)

	// If the issue is nil, we need to create a fake issue to display the URL
	url := "http://invalid-issue-report.example.com"
	if issue != nil {
		url = *issue.HTMLURL
	}
	success := errInTemplate == nil && errGitHubIssue == nil
	blocks, _ := modals.RenderBugReport(user.ID, url, errors.Join(errInTemplate, errGitHubIssue), success)
	_, _, errSlackMsg := client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...))
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return errors.Join(errInTemplate, errGitHubIssue, errSlackMsg)
}

// HandleBugReport handles the interaction with the form for creating a Bug report
func HandleBugReport(client *slack.Client, interaction slack.InteractionCallback) error {
	data := bugReportData{}

	viewState := interaction.View.State

	for k, v := range viewState.Values {
		switch k {
		case "issue_title":
			data.Title = v["title"].Value
		case "issue_expectedBehavior":
			data.ExpectedBehavior = v["expectedBehavior"].Value
		case "issue_currentBehavior":
			data.CurrentBehavior = v["currentBehavior"].Value
		case "issue_suggestions":
			data.Suggestions = v["suggestions"].Value
		case "issue_context":
			data.Context = v["context"].Value
		case "issue_labels":
			data.labels = strings.Split(v["labels"].Value, ",")
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}

	return createBugReportPostMessage(client, interaction, data)
}

// ShowBugReportForm creates a ModalViewRequest to gather the details to report a bug
func ShowBugReportForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowBugReportForm", "request", triggerID, user)
	defer tx.End()

	modalRequest, err := modals.BugReportModal(BugReportAction.Name, user, triggerID)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	_, err = client.OpenView(triggerID, modalRequest)
	if err != nil {
		apm2.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}
