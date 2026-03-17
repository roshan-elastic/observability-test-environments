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

package prompt

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/artifacts"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// EssConfig it will ask for the configuration of the ESS cluster, such as the cluster prefix/suffix,
// the stack version and the docker image to be used for elasticsearch, kibana and elastic-agent.
func EssConfig(versions []string) (prefix string, suffix string, stackVersion string, dockerImageVersion string) {
	// the questions to ask
	qSelect := selectStackVersion(versions)
	err := survey.AskOne(qSelect, &stackVersion)
	cobra.CheckErr(err)

	qSelect = selectBuild(stackVersion)
	err = survey.AskOne(qSelect, &dockerImageVersion)
	cobra.CheckErr(err)

	answers := struct {
		Prefix string `survey:"prefix"`
		Suffix string `survey:"suffix"`
	}{}

	qSurvey := []*survey.Question{
		selectPrefix(),
		selectSuffix(),
	}

	// perform the questions
	err = survey.Ask(qSurvey, &answers)
	cobra.CheckErr(err)

	prefix = answers.Prefix
	suffix = answers.Suffix

	return
}

// selectBuild prompt for the specific build to be used for the selected version.
func selectBuild(version string) *survey.Select {
	builds, err := artifacts.GetBuilds(version)
	cobra.CheckErr(err)

	sanitizedBuilds := []string{
		"--- Do not use a custom Docker image", // adding a keyword to detect empty values
	}
	sanitizedBuilds = append(sanitizedBuilds, builds...)
	return &survey.Select{
		Message: "Which docker image build do you want to use?",
		Options: sanitizedBuilds,
	}
}

// selectPrefix prompt for a prefix.
func selectPrefix() *survey.Question {
	return &survey.Question{
		Name:   "prefix",
		Prompt: &survey.Input{Message: "Introduce a prefix for the cluster (Optional)"},
	}
}

// selectStackVersion prompt for the stack version to be used.
func selectStackVersion(versions []string) *survey.Select {
	return &survey.Select{
		Message: "Which stack version do you want to use?",
		Options: versions,
	}
}

// selectSuffix prompt for a suffix.
func selectSuffix() *survey.Question {
	return &survey.Question{
		Name:   "suffix",
		Prompt: &survey.Input{Message: "Introduce a suffix for the cluster (Optional)"},
	}
}
