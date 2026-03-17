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

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/releases"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/box"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/helper"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
)

type ciIssueData struct {
	Project string
	Type    string
	User    string
}

var CIOnboardingAction = modals.Action{
	Name: "Enroll your GitHub",
	Desc: fmt.Sprintf("Display a form to enroll a new GitHub project in the CI with the `%s` command", CIOnboardingCommand()),
	Tldr: fmt.Sprintf("Onboard in the CI your project with the `%s` command", CIOnboardingCommand()),
}

// CIOnboardingCommand what's the command
func CIOnboardingCommand() string {
	return "/oblt-onboarding-ci"
}

// createCIOnboardingPostMessage posts message before the CI account is created
func createCIOnboardingPostMessage(client *slack.Client, interaction slack.InteractionCallback, project string, projectType string) error {
	user := interaction.User
	success := false
	isAdmin, _ := releases.IsUserAdmin(project, "obltmachine")

	if isAdmin {
		user := interaction.User
		success = true

		// Generate data from the template
		var data ciIssueData
		data.User = user.Name
		if user.RealName != "" {
			data.User = user.RealName
		}
		data.Project = project
		data.Type = projectType
		body, errInTemplate := helper.TemplateStringToString(string(box.Get("/ci-onboarding.issue")), data)
		// See https://github.com/elastic/observability-robots/pull/1633
		// This automation uses GitHub labels so the GitHub actions can filter what GitHub issues should be filtered
		labels := []string{"ci-onboarding-automation"}
		issue, errGitHubIssue := github.CreateIssueWithLabels(fmt.Sprintf("[CI onboarding] elastic/%s request %s", project, data.User), body, labels, "elastic", "observability-robots")
		// If the issue is nil, we need to create a fake issue to display the URL
		url := "http://invalid-issue-report.example.com"
		if issue != nil {
			url = *issue.HTMLURL
		}
		success = errInTemplate == nil && errGitHubIssue == nil
		fmt.Printf("Issue has been created: %v\n", url)
	}

	blocks, _ := modals.RenderCiOnboarding(user.ID, isAdmin, success)
	_, _, err := client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return nil
}

// HandleCIOnboarding handles the interaction with the form for enrolling a project in the CI
func HandleCIOnboarding(client *slack.Client, interaction slack.InteractionCallback) error {
	var project string
	var projectType string

	viewState := interaction.View.State

	// pick up all values from the aws-onboarding form
	for k, v := range viewState.Values {
		switch k {
		case "project":
			project = v["project"].Value
		case "projectType":
			projectType = v["projectType"].SelectedOption.Value
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}

	return createCIOnboardingPostMessage(client, interaction, project, projectType)
}

// ShowCIOnboardingForm creates a ModalViewRequest to gather the details to onboard the project in the CI
func ShowCIOnboardingForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowCIOnboardingForm", "request", triggerID, user)
	defer tx.End()

	modalRequest, err := modals.CIOnboardingModal(CIOnboardingAction.Name, triggerID)
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
