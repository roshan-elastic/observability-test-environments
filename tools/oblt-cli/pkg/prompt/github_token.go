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
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// GithubToken it will ask for a github token
func GithubToken(message string) string {
	// the answers will be written to this struct
	answers := struct {
		GithubToken string `survey:"githubToken"`
	}{}

	// the questions to ask
	var qs = []*survey.Question{
		githubToken(message),
	}

	// perform the questions
	err := survey.Ask(qs, &answers)
	cobra.CheckErr(err)

	return answers.GithubToken
}

// githubToken prompt for a github token.
func githubToken(message string) *survey.Question {
	return &survey.Question{
		Name:   "githubToken",
		Prompt: &survey.Password{Message: message},
		Validate: func(res interface{}) error {
			return ValidateGithubToken(res.(string))
		},
	}
}

// IsGithubToken It checks if the GitHub token is a valid token.
func IsGithubToken(token string) bool {
	// ref. https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/
	// ref. https://github.blog/changelog/2022-10-18-introducing-fine-grained-personal-access-tokens/
	re := regexp.MustCompile(`((ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{36}|github_pat_[0-9a-zA-Z_]{82})`)
	return re.Match([]byte(token))
}

// ValidateGithubToken It checks if the GitHub token is a valid token.
func ValidateGithubToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("a Github token was not set and it's required to look for new releases on Github")
	}
	if !IsGithubToken(token) {
		return fmt.Errorf("this is not a Github token please read https://elastic.github.io/observability-test-environments/user-guide/github-token/")
	}
	return nil
}
