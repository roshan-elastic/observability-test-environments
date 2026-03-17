package github_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWaitPullRequest(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN is not set")
	}

	currentUsername := "the-user"
	slackChannel := "#0123456"
	tmp := t.TempDir()

	// configuration is needed because the repo is cloned under config's workspace
	viper.SetConfigFile(filepath.Join(tmp, "config.yml"))
	viper.Set(config.UsernameFlag, currentUsername)
	viper.Set(config.SlackChannelFlag, slackChannel)
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	viper.WriteConfig()

	obltTestEnvironments, err := clusters.BootstrapRepository(false)
	assert.Nil(t, err, "Error should be nil")
	choreFile := filepath.Join(obltTestEnvironments.GetCurrentUserEnvironments(), "chore.txt")
	os.MkdirAll(filepath.Dir(choreFile), 0700)
	os.Create(choreFile)
	repo := obltTestEnvironments.Repo(false)
	repo.CommitAndPush("test: test PR")

	currentBranch, err := repo.GetCurrentBranch()
	assert.Nil(t, err, "Error should be nil")
	pr, _, err0 := github.CreatePullRequest("test: Test PR foobar", "test body", "elastic", "observability-test-environments", repo.Branch, currentBranch)
	assert.Nil(t, err0, "Error should be nil")

	pr1, err1 := github.FindPullRequest("Test PR foobar", "elastic", "observability-test-environments")
	assert.Nil(t, err1, "Error should be nil")
	assert.Equal(t, pr.GetNumber(), pr1.GetNumber(), "PR numbers should be equal")

	state, _, err3 := github.WaitPullRequestStateRetry(pr.GetNumber(), "elastic", "observability-test-environments", "open", 1*time.Minute)
	assert.Nil(t, err3, "Error should be nil")
	assert.Equal(t, "open", state, "State should be open")

	state, _, err4 := github.WaitPullRequestStateRetry(pr.GetNumber(), "elastic", "observability-test-environments", "merged", 1*time.Minute)
	github.ClosePullRequest(pr.GetNumber(), "elastic", "observability-test-environments")
	assert.NotNil(t, err4, "Error should not be nil")
	assert.Equal(t, "open", state, "State should be open")

	err5 := github.DeleteRef("elastic", "observability-test-environments", "heads/"+pr.GetHead().GetRef())
	t.Logf("PR Head: %s", pr.GetHead().GetRef())
	assert.Nil(t, err5, "Error should be nil")
}

func TestCreatePullRequestWithLabels(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN is not set")
	}

	currentUsername := "the-user"
	slackChannel := "#0123456"
	tmp := t.TempDir()

	// configuration is needed because the repo is cloned under config's workspace
	viper.SetConfigFile(filepath.Join(tmp, "config.yml"))
	viper.Set(config.UsernameFlag, currentUsername)
	viper.Set(config.SlackChannelFlag, slackChannel)

	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	viper.WriteConfig()

	obltTestEnvironments, err := clusters.BootstrapRepository(false)
	assert.Nil(t, err, "Error should be nil")
	choreFile := filepath.Join(obltTestEnvironments.GetCurrentUserEnvironments(), "chore-test.txt")
	os.MkdirAll(filepath.Dir(choreFile), 0700)
	os.Create(choreFile)
	repo := obltTestEnvironments.Repo(false)
	repo.CommitAndPush("test: test bar PR")

	currentBranch, err := repo.GetCurrentBranch()
	assert.Nil(t, err, "Error should be nil")
	pr, _, err := github.CreatePullRequestWithLabels("test: Test PR bar with labels", "test body", "elastic", "observability-test-environments", repo.Branch, currentBranch, []string{"skip-changelog"})
	assert.Nil(t, err, "Error should be nil")

	pr1, err1 := github.FindPullRequest("Test PR bar with labels", "elastic", "observability-test-environments")
	assert.Nil(t, err1, "Error should be nil")
	assert.Equal(t, pr.GetNumber(), pr1.GetNumber(), "PR numbers should be equal")

	err4 := github.DeleteRef("elastic", "observability-test-environments", "heads/"+pr.GetHead().GetRef())
	t.Logf("PR Head: %s", pr.GetHead().GetRef())
	assert.Nil(t, err4, "Error should be nil")
}

func TestGhMock(t *testing.T) {
	github.MockedHTTPClient = github.NewMockedHTTPClient()
	issue, err := github.CreateIssue("test", "test", "elastic", "observability-test-environments", []string{"test"})
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, issue, "Issue should not be nil")
	assert.Equal(t, issue.GetTitle(), "test", "Title should be test")

	issue, err = github.CreateIssueWithLabels("test", "test", []string{"test"}, "elastic", "observability-test-environments")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, issue, "Issue should not be nil")
	assert.Equal(t, issue.GetTitle(), "test", "Title should be test")

	pullrequest, _, err := github.CreatePullRequest("test", "test", "elastic", "observability-test-environments", "main", "test")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, pullrequest, "Pull request should not be nil")
	assert.Equal(t, pullrequest.GetTitle(), "test", "Title should be test")

	pullrequest, _, err = github.CreatePullRequestWithLabels("test", "test", "elastic", "observability-test-environments", "main", "test", []string{"test"})
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, pullrequest, "Pull request should not be nil")
	assert.Equal(t, pullrequest.GetTitle(), "test", "Title should be test")
}
