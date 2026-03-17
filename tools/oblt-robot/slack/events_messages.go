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
	"errors"

	"github.com/elastic/observability-test-environments/tools/oblt-robot/slack/mentions"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// handleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {
	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	if mentions.IsHelloMention(event.Text) {
		return mentions.Hello(client, event.Channel, *user)
	} else if mentions.IsFAQMention(event.Text) {
		return mentions.FAQ(event.Text, client, event.Channel, *user)
	} else if mentions.IsHelpMention(event.Text) {
		return mentions.Help(client, event.Channel, *user)
	} else if mentions.IsVersionMention(event.Text) {
		return mentions.Version(client, event.Channel, *user)
	} else if mentions.IsActiveBranchesMention(event.Text) {
		return mentions.ActiveBraches(client, event.Channel, *user)
	} else if mentions.IsReleaseMention(event.Text) {
		return mentions.Releases(client, event.Channel, *user)
	} else if mentions.IsSnapshotMention(event.Text) {
		return mentions.Snapshots(client, event.Channel, *user)
	}
	return mentions.Unknown(client, event.Channel, *user)
}
