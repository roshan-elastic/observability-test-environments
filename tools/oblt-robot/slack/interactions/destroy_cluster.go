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
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
)

var DestroyClusterAction = modals.Action{
	Name: "Destroy a cluster",
	Desc: fmt.Sprintf("Display a form to destroy a cluster with the `%s` command", DestroyClusterCommand()),
	Tldr: fmt.Sprintf("Destroy a cluster with the `%s` command", DestroyClusterCommand()),
}

// DestroyClusterCommand how to invoke the command from Slack with a slash command
func DestroyClusterCommand() string {
	return "/destroy-cluster"
}

// destroyClusterPostMessage posts message after a cluster is destroyed
func destroyClusterPostMessage(client *slack.Client, interaction slack.InteractionCallback, clusterName string) error {
	user := interaction.User
	obltUser := SlackToUsername(user.ID, user.Name)
	labelUser := apm.Label{Key: "slack-user", Value: user.ID}
	labelCluster := apm.Label{Key: "cluster-name", Value: clusterName}
	tx, ctx := apm.StartTransaction("destroyClusterPostMessage", "request", []apm.Label{labelUser, labelCluster})
	defer tx.End()

	obltTestEnvironments, err := OnMemoryUserConf(user.ID, obltUser, false)

	var params map[string]interface{}
	if err == nil {
		params, err = obltTestEnvironments.DestroyCluster(clusterName, true)
	}

	err = sendDestroyResponse(user, clusterName, err, params, client, ctx)
	if err != nil {
		fmt.Printf("Error sending the message for the destruction of the cluster: %s", err)
		apm2.CaptureError(ctx, err).Send()
	}
	_ = os.Remove(config.ForUser(user.ID))
	return err
}

// sendDestroyResponse sends a message to the user with the result of the destruction of the cluster
func sendDestroyResponse(user slack.User, clusterName string, err error, params map[string]interface{}, client *slack.Client, ctx context.Context) error {
	color := "#4af030"
	pretext := fmt.Sprintf("Hey <@%s>, the cluster `%s` is going to be destroyed! GitHub will send you a DM when the real destructions takes place.", user.ID, clusterName)
	var text string
	if err != nil {
		pretext = fmt.Sprintf("Hey <@%s>, the cluster `%s` **could not** be destroyed!", user.ID, clusterName)
		color = "#ff0000"
		text = fmt.Sprintf("error: %s", err)
	} else {
		text = fmt.Sprintf("You can find the commit destroying your cluster here: %s", params["CommitURL"])
	}

	attachment := slack.Attachment{
		Text:    text,
		Pretext: pretext,
		Color:   color,
		Fields:  []slack.AttachmentField{},
	}

	_, _, err = client.PostMessage(user.ID,
		slack.MsgOptionAttachments(attachment),
	)
	return err
}

// ShowDestroyClusterForm creates a ModalViewRequest with a select box to present the available clusters
func ShowDestroyClusterForm(client *slack.Client, triggerID string, slackID string) error {
	tx, ctx := apm.StartTransactionForm("ShowDestroyClusterForm", "request", triggerID, slackID)
	defer tx.End()

	var clusterList []string

	files := FilterClustersByUser(slackID)
	for _, cluster := range files {
		clusterList = append(clusterList, cluster.Data["cluster_name"].(string))
	}
	sort.Strings(clusterList)
	modalRequest, err := modals.SelectClusterModal(clusterList, DestroyClusterAction.Name, "destroy-select-cluster", triggerID)
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
