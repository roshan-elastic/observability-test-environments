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
	"sort"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/modals"
	"github.com/slack-go/slack"
	apm2 "go.elastic.co/apm/v2"
	"gopkg.in/yaml.v3"
)

const (
	ReadClusterSecretPrefix = "readsecret__"
	credentialsSecret       = "credentials"
	deployInfoSecret        = "deploy-info"
	environmentSecret       = "environment-vars"
	kibanaConfigSecret      = "kibana-yaml"
	clusterStateSecret      = "cluster-state"
)

var listOfSecrets = []string{
	credentialsSecret,
	deployInfoSecret,
	environmentSecret,
	kibanaConfigSecret,
	clusterStateSecret,
}

var ClusterSecretAction = modals.Action{
	Name: "Select a cluster",
	Desc: fmt.Sprintf("Display a form to select a secret from a cluster with the `%s` command", ClusterSecretCommand()),
	Tldr: fmt.Sprintf("Display secret for a cluster with the `%s` command", ClusterSecretCommand()),
}

// ClusterSecretCommand what's the command
func ClusterSecretCommand() string {
	return "/cluster-secret"
}

// PostSecretMessage posts message with the content of the secret
func PostSecretMessage(client *slack.Client, interaction slack.InteractionCallback, secretActionID string, threadTS string, gcsm *gcp.ClusterSecrets) (blocks []slack.Block, attachments []slack.Attachment, err error) {
	user := interaction.User
	tokens := strings.Split(secretActionID, "__")
	clusterName := tokens[1]
	secret := tokens[2]
	var content string
	if content, err = ReadClusterSecret(gcsm, clusterName, secret); err == nil {
		header := fmt.Sprintf("Retrieving `%s` secret from *%s*", secret, clusterName)
		blocks, attachments, err = modals.RenderMarkdown(header, content)
		if err == nil {
			_, _, err = client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...), slack.MsgOptionAttachments(attachments...), slack.MsgOptionTS(threadTS))
		}
	}
	_ = os.Remove(config.ForUser(interaction.TriggerID))
	return blocks, attachments, err
}

// ReadClusterSecret reads the secret from the cluster, the secrets is identified by the secret name
func ReadClusterSecret(gcsm *gcp.ClusterSecrets, clusterName, secret string) (content string, err error) {
	switch secret {
	case credentialsSecret:
		content, err = gcsm.ReadCredentialsSecret(clusterName)
	case deployInfoSecret:
		content, err = gcsm.ReadDeployInfoSecret(clusterName)
	case environmentSecret:
		content, err = gcsm.ReadEnvSecret(clusterName)
		content = "```" + content + "```"
	case kibanaConfigSecret:
		content, err = gcsm.ReadKibanaYamlSecret(clusterName)
		content = "```" + content + "```"
	case clusterStateSecret:
		var clusterStateSecret gcp.ClusterStateSecret
		clusterStateSecret, err = gcsm.ReadClusterStateSecret(clusterName)
		if err != nil {
			content = fmt.Sprintf("Error reading the cluster state secret: %s", err)
		} else {
			var bytes []byte
			if bytes, err = yaml.Marshal(clusterStateSecret); err == nil {
				content = "```" + string(bytes) + "```"
			}
		}
	default:
		content = "secret %s not supported"
		fmt.Printf(content, secret)
		err = fmt.Errorf("%s", content)
	}
	if err != nil {
		content = fmt.Sprintf("Error reading the secret: %s", err)
	}
	return content, err
}

// postSecretsMessage posts message with the list of secrets
func postSecretsMessage(client *slack.Client, interaction slack.InteractionCallback, clusterName string) (blocks []slack.Block, err error) {
	user := interaction.User
	header := fmt.Sprintf("Hey there 👋 <@%s>. the secrets for the `%s` cluster are listed below. Please use them with caution", user.ID, clusterName)
	blocks, _, err = modals.RenderSecretsSelection(header, ReadClusterSecretPrefix+clusterName, listOfSecrets)
	if err == nil {
		_, _, err = client.PostMessage(user.ID, slack.MsgOptionBlocks(blocks...))
	}
	if err != nil {
		fmt.Printf("Error posting the secrets: %s", err)
	}

	return blocks, err
}

// HandleSelectCluster handles the interaction with the select cluster form
// this function test is in cluster_secret_test.go
func HandleSelectCluster(client *slack.Client, interaction slack.InteractionCallback) (err error) {
	viewState := interaction.View.State
	// pick up all values from the create cluster form
	for k, v := range viewState.Values {
		switch k {
		case "destroy-select-cluster":
			clusterName := v["select-cluster"].SelectedOption.Value
			err = destroyClusterPostMessage(client, interaction, clusterName)
			return err
		case "select-cluster-for-secrets":
			clusterName := v["select-cluster"].SelectedOption.Value
			_, err = postSecretsMessage(client, interaction, clusterName)
			return err
		default:
			fmt.Printf("viewSubmission interaction not supported: %v\n", k)
			err = fmt.Errorf("viewSubmission interaction not supported: %v", k)
		}
	}
	return err
}

// ShowSecretClusterForm creates a ModalViewRequest with a select box to present the available clusters
func ShowSecretClusterForm(client *slack.Client, triggerID string, slackID string) error {
	tx, ctx := apm.StartTransactionForm("ShowSecretClusterForm", "request", triggerID, slackID)
	defer tx.End()

	var clusterList []string

	files := FilterClustersByUser(slackID)
	for _, cluster := range files {
		clusterList = append(clusterList, cluster.Data["cluster_name"].(string))
	}

	goldenFiles := GoldenClusters()
	clusterList = append(clusterList, goldenFiles...)

	if len(clusterList) == 0 {
		clusterList = append(clusterList, "No clusters available")
	}

	sort.Strings(clusterList)

	modalRequest, err := modals.SelectClusterModal(clusterList, ClusterSecretAction.Name, "select-cluster-for-secrets", triggerID)
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
