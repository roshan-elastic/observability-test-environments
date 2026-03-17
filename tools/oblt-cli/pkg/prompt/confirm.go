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
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// Confirm asks a question which answer is yes or no, if yes it returns, otherwise it aborts.
func Confirm(message string) (result string, err error) {
	// the answers will be written to this struct
	answers := struct {
		Answer string `survey:"answer"`
	}{}

	// the questions to ask
	var qs = []*survey.Question{
		{
			Name:   "answer",
			Prompt: &survey.Input{Message: message + " (yes/no)"},
			Validate: func(val interface{}) error {
				return validateYes(val.(string))
			},
			Transform: survey.Title,
		},
	}

	// perform the questions
	err = survey.Ask(qs, &answers)
	cobra.CheckErr(err)

	result = answers.Answer

	if strings.TrimSpace(strings.ToLower(result)) == "yes" {
		return result, nil
	}

	logger.Fatal("the operation was aborted")
	return result, nil //never executed
}

// validateYes It checks if the answer is yes or no.
func validateYes(value string) (err error) {
	if strings.TrimSpace(strings.ToLower(value)) == "yes" || strings.TrimSpace(strings.ToLower(value)) == "no" {
		return nil
	}
	return errors.New("the answer must be 'yes' or 'no'")
}
