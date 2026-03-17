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
package slack

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/interactions"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/slack-go/slack/socketmode"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationCommand(t *testing.T) {
	// Initialize the mock Slack server
	// Create a Slack client using the mock server URL
	// Option to set a custom logger
	server, slackClient, socketClient := newSlackServerConnection()
	defer server.Stop()
	server.Start()

	event := interactions.NewTestEventSampleConfiguration()
	assert.True(t, processInteractive(event, socketClient, slackClient))

	// TODO investigate if we can test the whole bot as connected to the server
	// https://github.com/slack-go/slack/blob/master/slacktest/server_test.go
	// go mainLoop(context.Background(), slackClient, socketClient)
	// slackClient.SendMessage("dummy-channel", slack.MsgOptionText("/oblt-configure", false))
	// time.Sleep(100 * time.Second)
}

// newSlackServerConnection creates a new slack server, a Slack Client, and Socket Client connection
func newSlackServerConnection() (*slacktest.Server, *slack.Client, *socketmode.Client) {
	server := slacktest.NewTestServer()

	server.Handle("/views.open", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	slackClient := slack.New("dummy-token", slack.OptionAPIURL(server.GetAPIURL()))

	socketClient := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),

		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return server, slackClient, socketClient
}
