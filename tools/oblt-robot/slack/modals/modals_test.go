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
package modals

import (
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/questions"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/releases"
)

func TestRenderGeneralHelp(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {

		mention := Action{
			Name: "mention",
			Desc: "Mention me",
			Tldr: "Display mention",
		}

		command := Action{
			Name: "command",
			Desc: "Command me",
			Tldr: "Display command",
		}

		_, err := RenderGeneralHelp([]Action{command}, []Action{mention})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderVersion(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {

		_, err := RenderVersion("1.0.0", "2022-07-20")
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderResetAWSAccount(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderResetAWSAccount("foo", "foo@elastic.co", true, true)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderAwsOnboarding(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderAwsOnboarding("foo", "foo@elastic.co", true, true)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderBugReport(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderBugReport("foo", "foo@elastic.co", nil, true)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderCloudOnboarding(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderCloudOnboarding("foo", "foo@elastic.co", true, true)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderBranches(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		version := releases.Versions{
			Branch:      "7.17",
			Version:     "7.17.15",
			ReleaseDate: "2022-07-24",
		}
		_, err := RenderBranches([]releases.Versions{version})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
	t.Run("if no versions should return an object without any error", func(t *testing.T) {
		_, err := RenderBranches([]releases.Versions{})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderReleases(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		build := releases.Builds{
			Branch:      "7.17",
			Version:     "7.17.15",
			ReleaseDate: "2022-07-24",
		}
		_, err := RenderReleases([]releases.Builds{build})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
	t.Run("if no versions should return an object without any error", func(t *testing.T) {
		_, err := RenderReleases([]releases.Builds{})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderSnapshots(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		build := releases.Builds{
			Branch:      "7.17",
			Version:     "7.17.15",
			ReleaseDate: "2022-07-24",
		}
		_, err := RenderSnapshots([]releases.Builds{build})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
	t.Run("if no versions should return an object without any error", func(t *testing.T) {
		_, err := RenderSnapshots([]releases.Builds{})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderAnswers(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		answer := questions.Answer{
			Title:  "my-title",
			Answer: "my-answer",
		}
		_, err := RenderFAQ([]questions.Answer{answer})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
	t.Run("if no versions should return an object without any error", func(t *testing.T) {
		_, err := RenderFAQ([]questions.Answer{})
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderHello(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderHello("foo")
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderUnknown(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := RenderUnknown("foo")
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestCreateServerlessClusterModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		_, err := CreateServerlessClusterModal([]string{"qa", "staging", "production"}, []string{"oblt", "acme", "es"}, "foo", "bar")
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderMarkdown(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		content := "foo-content"
		_, _, err := RenderMarkdown(header, content)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderSecretsSelection(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		content := []string{"foo-content", "bar-content"}
		_, _, err := RenderSecretsSelection(header, actionPrefix, content)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestAWSOnboardingModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		_, err := AWSOnboardingModal(header, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestAWSResetModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		_, err := AWSResetModal(header, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestBugReportModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		user := "foo"
		_, err := BugReportModal(header, user, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestCloudOnboardingModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		_, err := CloudOnboardingModal(header, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestCreateClusterModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		templates := []string{"foo", "bar"}
		actionPrefix := "foo-action"
		_, err := CreateClusterModal(templates, header, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestSelectClusterModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		clusters := []string{"foo", "bar"}
		block := "foo"
		actioID := "bar"
		_, err := SelectClusterModal(clusters, header, block, actioID)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}

func TestRenderCiOnboarding(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		user := "foo"
		isAdmin := true
		success := true
		_, err := RenderCiOnboarding(user, isAdmin, success)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}

		isAdmin = false
		success = false
		_, err = RenderCiOnboarding(user, isAdmin, success)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}

		isAdmin = true
		success = false
		_, err = RenderCiOnboarding(user, isAdmin, success)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}

		isAdmin = false
		success = true
		_, err = RenderCiOnboarding(user, isAdmin, success)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}

	})
}

func TestCIOnboardingModal(t *testing.T) {
	t.Run("should return an object without any error", func(t *testing.T) {
		header := "## foo-header"
		actionPrefix := "foo-action"
		_, err := CIOnboardingModal(header, actionPrefix)
		if err != nil {
			t.Errorf("Error actual = %v, and Expected = nil.", err)
		}
	})
}
