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

package mentions

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
)

// Assuming Hello function is defined in the same package

func TestHelloWithAssert(t *testing.T) {
	// Initialize the Slack test server
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	// Create a new Slack client with the test server's URL
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))

	channelID := "channelID"
	user := slack.User{ID: "foo"}

	// Call the Hello function with the test server client, channelID, and user
	err := Hello(client, channelID, user)

	// Use assert.NoError to check if an error was not returned
	assert.NoError(t, err, "Hello function should not return an error")
}
func TestIsHelloMention(t *testing.T) {
	// Test case 1: Text contains "hello"
	text1 := "Hello, world!"
	assert.True(t, IsHelloMention(text1), "Expected true for text1")

	// Test case 2: Text contains "HELLO"
	text2 := "Say HELLO to everyone"
	assert.True(t, IsHelloMention(text2), "Expected true for text2")

	// Test case 3: Text does not contain "hello"
	text3 := "Hi there!"
	assert.False(t, IsHelloMention(text3), "Expected false for text3")
}
