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
package modals

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/questions"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/releases"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/box"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/helper"
	"github.com/slack-go/slack"
)

// Action maps a incoming command or mention
type Action struct {
	Name string // Name is the name of the action.
	Tldr string // Tldr is a short description of what the action does. (80 chars)
	Desc string // Desc is a longer description of what the action does.
}

func AWSOnboardingModal(title, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
	}{
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/aws-onboarding-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func AWSResetModal(title, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
	}{
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/aws-reset-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func BugReportModal(title, user, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
		User  string
	}{
		title,
		user,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/bug-report-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func CIOnboardingModal(title, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
	}{
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/ci-onboarding-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func CreateClusterModal(templates []string, title string, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Templates []string
		Title     string
	}{
		templates,
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/create-cluster-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func CreateServerlessClusterModal(environments []string, projects []string, title string, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Environments []string
		Projects     []string
		Title        string
	}{
		environments,
		projects,
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/create-serverless-cluster-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func CloudOnboardingModal(title, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
	}{
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/cloud-onboarding-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

func SelectClusterModal(clusters []string, title string, block string, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Block    string
		Clusters []string
		Title    string
	}{
		block,
		clusters,
		title,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/select-cluster-modal.json.tmpl")), data)

	return createModalViewRequest(str, err, actionID)
}

// RenderGeneralHelp provides the general help command for any event channel
func RenderGeneralHelp(commands []Action, mentions []Action) ([]slack.Block, error) {

	data := struct {
		Commands []Action
		Mentions []Action
	}{
		commands,
		mentions,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/general-help.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderResetAWSAccount provides the AWS reset account output
func RenderResetAWSAccount(user string, email string, isElastic bool, success bool) ([]slack.Block, error) {

	data := struct {
		User    string
		Email   string
		Success bool
		Elastic bool
	}{
		user,
		email,
		success,
		isElastic,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/aws_account_reset.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderAwsOnboarding provides the AWS Onboarding output
func RenderAwsOnboarding(user string, email string, isElastic bool, success bool) ([]slack.Block, error) {

	data := struct {
		User    string
		Email   string
		Success bool
		Elastic bool
	}{
		user,
		email,
		success,
		isElastic,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/aws_account_onboarding.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderBugReport provides the Bug report output
func RenderBugReport(user string, issue string, err error, success bool) ([]slack.Block, error) {
	data := struct {
		User    string
		Issue   string
		Success bool
		Error   string
	}{
		user,
		issue,
		success,
		fmt.Sprintf("%v", err),
	}

	str, err := helper.TemplateStringToString(string(box.Get("/bug_report.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderCiOnboarding provides CI Onboarding output
func RenderCiOnboarding(user string, isAdmin bool, success bool) ([]slack.Block, error) {

	data := struct {
		User    string
		Success bool
		IsAdmin bool
	}{
		user,
		success,
		isAdmin,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/ci_onboarding.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderCloudOnboarding provides the Cloud Onboarding output
func RenderCloudOnboarding(user string, email string, isElastic bool, success bool) ([]slack.Block, error) {

	data := struct {
		User    string
		Email   string
		Success bool
		Elastic bool
	}{
		user,
		email,
		success,
		isElastic,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/cloud_account_onboarding.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderBranches provides the Active Branches output
func RenderBranches(versions []releases.Versions) ([]slack.Block, error) {

	data := struct {
		Versions []releases.Versions
	}{
		versions,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/active_branches.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderFAQ provides the current answers output
func RenderFAQ(answers []questions.Answer) ([]slack.Block, error) {
	data := struct {
		Answers []questions.Answer
	}{
		answers,
	}
	str, err := helper.TemplateStringToString(string(box.Get("/faq.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderHello provides the Hello output
func RenderHello(user string) ([]slack.Block, error) {

	data := struct {
		User string
	}{
		user,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/hello.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderReleases provides the current releases output
func RenderReleases(builds []releases.Builds) ([]slack.Block, error) {
	data := struct {
		Releases []releases.Builds
	}{
		builds,
	}
	str, err := helper.TemplateStringToString(string(box.Get("/releases.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderSnapshots provides the current snapshots output
func RenderSnapshots(snapshots []releases.Builds) ([]slack.Block, error) {
	data := struct {
		Snapshots []releases.Builds
	}{
		snapshots,
	}
	str, err := helper.TemplateStringToString(string(box.Get("/snapshots.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderUnknown provides the Unknown output
func RenderUnknown(user string) ([]slack.Block, error) {

	data := struct {
		User string
	}{
		user,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/unknown.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderVersion provides the Version output
func RenderVersion(version string, build string) ([]slack.Block, error) {

	data := struct {
		Version string
		Build   string
	}{
		version,
		build,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/version.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

// RenderMarkdown renders a simple markdown message
func RenderMarkdown(header, markdown string) ([]slack.Block, []slack.Attachment, error) {
	chunks := helper.Chunks(markdown, 2500)
	chunksLen := len(chunks)
	data := struct {
		Header   string
		Markdown []string
		Chunks   int
	}{
		helper.JsonEscape(header),
		[]string{},
		chunksLen,
	}

	for _, chunk := range chunks {
		data.Markdown = append(data.Markdown, helper.JsonEscape(chunk))
	}

	str, err := helper.TemplateStringToString(string(box.Get("/markdown-message.json.tmpl")), data)
	return createBlockSetRequest(str, err)
}

func RenderSecretsSelection(header string, actionPrefix string, secrets []string) ([]slack.Block, []slack.Attachment, error) {
	data := struct {
		Header       string
		ActionPrefix string
		Secrets      []string
	}{
		header,
		actionPrefix,
		secrets,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/cluster_secret_select.json.tmpl")), data)
	return createBlockSetRequest(str, err)
}

func createBlockSetRequest(str string, errStr error) ([]slack.Block, []slack.Attachment, error) {
	if errStr != nil {
		return nil, nil, errStr
	}

	view := slack.Msg{}

	err := json.Unmarshal([]byte(str), &view)
	if err != nil {
		return nil, nil, err
	}

	return view.Blocks.BlockSet, view.Attachments, nil
}

func createModalViewRequest(str string, errStr error, actionID string) (slack.ModalViewRequest, error) {
	if errStr != nil {
		return slack.ModalViewRequest{}, errStr
	}

	view := slack.ModalViewRequest{
		PrivateMetadata: actionID,
	}

	err := json.Unmarshal([]byte(str), &view)
	if err != nil {
		return slack.ModalViewRequest{}, err
	}

	return view, nil
}

func RenderConfigure(user string, githubUser string) ([]slack.Block, error) {
	data := struct {
		User       string
		GithubUser string
	}{
		user,
		githubUser,
	}

	str, err := helper.TemplateStringToString(string(box.Get("/configure.json.tmpl")), data)
	block, _, err := createBlockSetRequest(str, err)
	return block, err
}

func RenderConfigureForm(actionID string) (slack.ModalViewRequest, error) {
	data := struct{}{}
	str, err := helper.TemplateStringToString(string(box.Get("/configure-modal.json.tmpl")), data)
	return createModalViewRequest(str, err, actionID)
}
