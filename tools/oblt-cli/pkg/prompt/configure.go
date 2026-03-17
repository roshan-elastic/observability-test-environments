// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package prompt It contains the functions to interact with the user in the console.
package prompt

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var SlackChannel string
var User string

// Configure it will ask for a few configuration settings in the terminal.
func Configure() {
	// the answers will be written to this struct
	answers := struct {
		SlackChannel string `survey:"slackChannel"`
		Username     string `survey:"username"`
	}{}

	// the questions to ask
	var qs = []*survey.Question{
		slackChannel(),
		username(),
	}

	// perform the questions
	err := survey.Ask(qs, &answers)
	cobra.CheckErr(err)

	SlackChannel = answers.SlackChannel
	User = answers.Username

	if len(SlackChannel) == 0 {
		logger.Fatal("No Slack channel/member ID set.")
	}

	if len(User) == 0 {
		logger.Fatal("No User set.")
	}
}

// slackChannel prompt for a slack channel.
func slackChannel() *survey.Question {
	return &survey.Question{
		Name:   "slackChannel",
		Prompt: &survey.Input{Message: "Slack channel(#myChannel)/member ID(@F1A2DFG)"},
		Validate: func(val interface{}) error {
			return config.ValidateSlackChannel(val.(string))
		},
	}
}

// Username prompt for a username.
func username() *survey.Question {
	return &survey.Question{
		Name:   "username",
		Prompt: &survey.Input{Message: "User name ([a-z0-9-.])"},
		Validate: func(val interface{}) error {
			return config.ValidateUsername(val.(string))
		},
		Transform: survey.ToLower,
	}
}
