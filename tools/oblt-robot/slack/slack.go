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
	"context"
	"log"
	"os"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/interactions"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// DryRun enable/disable test mode.
var DryRun bool

// SlackSocketMode will be blocking and ingesting new Websocket messages on a channel at socketClient.Events
func StartSocketMode() {
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go mainLoop(ctx, client, socketClient)

	go updateClusterDataLoop()

	socketClient.Run()
}

// Create a for loop that selects either the context cancellation or the events incoming
// in case context cancel is called exit the goroutine
func mainLoop(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-socketClient.Events:
			processEvents(event, socketClient, client)
		}
	}
}

// We have a new Events, let's type switch the event
// Add more use cases here if you want to listen to other events.
func processEvents(event socketmode.Event, socketClient *socketmode.Client, client *slack.Client) {
	switch event.Type {
	case socketmode.EventTypeEventsAPI:
		processEventsAPI(event, socketClient, client)
	case socketmode.EventTypeSlashCommand:
		processSlashCommands(event, socketClient, client)
	case socketmode.EventTypeInteractive:
		processInteractive(event, socketClient, client)
	}
}

// handleSlashCommand will take care of the command
// Handle Interactions with commands
func processInteractive(event socketmode.Event, socketClient *socketmode.Client, client *slack.Client) bool {
	interaction, ok := event.Data.(slack.InteractionCallback)
	if !ok {
		log.Printf("Could not type cast the message to a Interaction callback: %v\n", interaction)
		return false
	}

	socketClient.Ack(*event.Request, nil)

	err := handleInteractionEvent(interaction, client)
	if err != nil {
		log.Printf("[ERROR] %v -> %v\n", interaction, err)
	}
	return true
}

// Handle Slash Commands
// type cast to the correct event type, this time a SlashEvent
// Dont forget to acknowledge the request and send the payload
// The payload is the response
func processSlashCommands(event socketmode.Event, socketClient *socketmode.Client, client *slack.Client) bool {
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
		return false
	}

	socketClient.Ack(*event.Request, nil)

	err := handleSlashCommand(command, client)
	if err != nil {
		log.Printf("[ERROR] %v -> %v\n", command, err)
	}
	return true
}

// handle EventAPI events
// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
// We need to send an Acknowledge to the slack server
// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
// Replace with actual err handling
func processEventsAPI(event socketmode.Event, socketClient *socketmode.Client, client *slack.Client) bool {
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
		return false
	}

	socketClient.Ack(*event.Request)

	err := handleEventMessage(eventsAPIEvent, client)
	if err != nil {

		log.Printf("[ERROR] %v -> %v\n", eventsAPIEvent, err)
	}
	return true
}

func updateClusterDataLoop() {
	for {
		interactions.UpdateClustersData()
		time.Sleep(1 * time.Minute)
	}
}
