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
	"encoding/json"
	"fmt"
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/questions"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/box"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/colors"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/helper"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"

	apm2 "go.elastic.co/apm/v2"
)

const (
	faqTemplate = "/add-faq-entry-modal.json.tmpl"
)

type faqData struct {
	question string
	answer   string
	url      string
	user     string
}

var AddFAQEntryAction = modals.Action{
	Name: "Submit a FAQ entry",
	Desc: fmt.Sprintf("Display a form to submit a FAQ entry with the `%s` command", AddFAQEntryCommand()),
	Tldr: fmt.Sprintf("Submit your knowdledge with the `%s` command", AddFAQEntryCommand()),
}

// AddFAQEntryCommand what's the command
func AddFAQEntryCommand() string {
	return "/oblt-faq-entry"
}

// addFAQEntryPostMessage posts message before the FAQ entry is created
func addFAQEntryPostMessage(client *slack.Client, interaction slack.InteractionCallback, data faqData) error {
	user := interaction.User
	var attachment slack.Attachment

	attachment = slack.Attachment{
		Text:    "Everything went smooth and the FAQ has been created. Thanks",
		Pretext: fmt.Sprintf("Hey <@%s>,", user.ID),
		Color:   colors.Green,
	}

	// Generate data from the template
	data.user = user.Name
	if user.RealName != "" {
		data.user = user.RealName
	}

	err := questions.AddAnswer(data.question, data.answer, data.url, data.user)
	if err != nil {
		attachment = slack.Attachment{
			Text:    fmt.Sprintf("something went wrong when creating the issue: %v", err),
			Pretext: fmt.Sprintf("Hey <@%s>,", user.ID),
			Color:   colors.Red,
		}
	}

	_, _, err = client.PostMessage(user.ID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("Error creating the issue: %v", err)
		return err
	}

	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return nil
}

// HandleAddFAQEntry handles the interaction with the form for creating a FAQ entry
func HandleAddFAQEntry(client *slack.Client, interaction slack.InteractionCallback) error {
	data := faqData{}

	viewState := interaction.View.State

	for k, v := range viewState.Values {
		switch k {
		case "question":
			data.question = v["question"].Value
		case "answer":
			data.answer = v["answer"].Value
		case "url":
			data.url = v["url"].Value
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}

	return addFAQEntryPostMessage(client, interaction, data)
}

// ShowAddFAQEntryForm creates a ModalViewRequest
func ShowAddFAQEntryForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowAddFAQEntryForm", "request", triggerID, user)
	defer tx.End()

	modalRequest, err := addFAQEntryModal(AddFAQEntryAction.Name, user, triggerID)
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

func addFAQEntryModal(title, user, actionID string) (slack.ModalViewRequest, error) {
	data := struct {
		Title string
		User  string
	}{
		title,
		user,
	}

	str, err := helper.TemplateStringToString(string(box.Get(faqTemplate)), data)
	if err != nil {
		return slack.ModalViewRequest{}, err
	}

	view := slack.ModalViewRequest{
		PrivateMetadata: actionID,
	}

	err = json.Unmarshal([]byte(str), &view)
	if err != nil {
		return slack.ModalViewRequest{}, err
	}

	return view, nil
}
