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
	"sort"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
)

var CreateServerlessClusterAction = modals.Action{
	Name: "New Serverless cluster",
	Desc: fmt.Sprintf("Display a form to create a Serverless cluster with the `%s` command", CreateServerlessClusterCommand()),
	Tldr: fmt.Sprintf("Create a Serverless cluster with the `%s` command", CreateServerlessClusterCommand()),
}

// CreateServerlessClusterCommand what's the command
func CreateServerlessClusterCommand() string {
	return "/create-serverless-cluster"
}

// createServerlessClusterPostMessage posts message after a Serverless cluster is created
func createServerlessClusterPostMessage(client *slack.Client, interaction slack.InteractionCallback, clusterPrefix string, projectType string, target string, clusterSuffix string) error {
	user := interaction.User
	obltUser := SlackToUsername(user.ID, user.Name)
	dryRun := user.ID == "dry-run"

	labelUser := apm.Label{Key: "slack-user", Value: user.ID}
	labelUserOblt := apm.Label{Key: "oblt-user", Value: obltUser}
	labelTarget := apm.Label{Key: "target", Value: target}
	labelProjectType := apm.Label{Key: "project-type", Value: projectType}
	tx, ctx := apm.StartTransaction("createServerlessClusterPostMessage", "request", []apm.Label{labelUser, labelUserOblt, labelTarget, labelProjectType})
	defer tx.End()

	obltTestEnvironments, err := OnMemoryUserConf(user.ID, obltUser, dryRun)

	var params map[string]interface{}
	if err == nil {
		serverless := &clusters.ServerlessCluster{
			ClusterNamePrefix: clusterPrefix,
			ClusterNameSuffix: clusterSuffix,
			Target:            target,
			ProjectType:       projectType,
			ObltRepo:          obltTestEnvironments,
			SlackChannel:      "@" + user.ID,
			Username:          obltUser,
		}
		params, err = serverless.Create()

	}

	err = sendServerlessCreateResponse(user, target, err, params, client)

	if err != nil {
		logger.LogError("Error creating cluster", err)
		apm2.CaptureError(ctx, err).Send()
	}

	_ = os.Remove(config.ForUser(user.ID))
	return err
}

// sendRespose sends a response to the user after a Serverless cluster is created
func sendServerlessCreateResponse(user slack.User, target string, err error, params map[string]interface{}, client *slack.Client) error {
	// skip those keys including paths, as they will be relative to the deployed bot
	pretext := fmt.Sprintf("Hey <@%s>, your Serverless cluster using `%s` template is going to be created! GitHub will send you a DM with the configuration details. In the meantime, I'm adding here a summary of what you requested", user.ID, target)
	color := "#4af030"
	var text string
	if err != nil {
		color = "#ff0000"
		pretext = fmt.Sprintf("Hey <@%s>, the creation of the Serverless cluster using `%s` target failed.", user.ID, target)
		text = fmt.Sprintf("error: %s", err)
	} else {
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if strings.Contains(k, "File") {

				continue
			}

			text += fmt.Sprintf("- `%s`: %+v\n", k, params[k])
		}
	}

	attachment := slack.Attachment{
		Text:    text,
		Pretext: pretext,
		Color:   color,
		Fields:  []slack.AttachmentField{},
	}

	_, _, errPost := client.PostMessage(user.ID,
		slack.MsgOptionAttachments(attachment),
	)
	return errors.Join(err, errPost)
}

// HandleServerlessClusterHandle handles the interaction with the create Serverless cluster form
func HandleCreateServerlessCluster(client *slack.Client, interaction slack.InteractionCallback) error {
	var clusterPrefix string
	var clusterProject string
	var clusterTarget string
	var clusterSuffix string

	viewState := interaction.View.State

	// pick up all values from the create cluster form
	for k, v := range viewState.Values {
		switch k {
		case "create-serverless-cluster-from-target":
			clusterTarget = v["select-from-target"].SelectedOption.Value
		case "create-serverless-cluster-from-project":
			clusterProject = v["select-from-project"].SelectedOption.Value
		case "create-serverless-cluster-name-prefix":
			clusterPrefix = v["cluster-name-prefix"].Value
		case "create-serverless-cluster-name-suffix":
			clusterSuffix = v["cluster-name-suffix"].Value
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
		}
	}
	err := createServerlessClusterPostMessage(client, interaction, clusterPrefix, clusterProject, clusterTarget, clusterSuffix)
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return err
}

// ShowServerlessClusterForm creates a ModalViewRequest with a select box to present the available template clusters
func ShowServerlessClusterForm(client *slack.Client, triggerID string, user string) error {
	tx, ctx := apm.StartTransactionForm("ShowServerlessClusterForm", "request", triggerID, user)
	defer tx.End()

	serverlessList := ServerlessEnvironments()
	projectList := ServerlessProjects()
	modalRequest, err := modals.CreateServerlessClusterModal(serverlessList, projectList, CreateServerlessClusterAction.Name, triggerID)
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
