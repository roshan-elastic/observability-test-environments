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

func TestIsFAQMention(t *testing.T) {
	// Test case 1: Text contains "faq"
	text1 := "I have a faq question"
	result := IsFAQMention(text1)
	assert.True(t, result, "Expected IsFAQMention to return true")

	// Test case 2: Text contains "know"
	text2 := "Do you know the answer?"
	result = IsFAQMention(text2)
	assert.True(t, result, "Expected IsFAQMention to return true")

	// Test case 3: Text does not contain "faq" or "know"
	text3 := "This is a regular message"
	result = IsFAQMention(text3)
	assert.False(t, result, "Expected IsFAQMention to return false")
}

func TestFAQ(t *testing.T) {
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()
	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	channelID := "channelID"
	user := slack.User{ID: "foo"}

	// Set up the test case
	question := "What is the answer to the FAQ?"

	// Call the FAQ function
	err := FAQ(question, client, channelID, user)

	// Assert that no error occurred
	assert.NoError(t, err, "Expected no error")
}
