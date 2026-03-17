package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/google/go-github/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"golang.org/x/oauth2"
)

var MockedHTTPClient *http.Client

type ghIssueMock struct {
	Title string
	Body  string
}

type ghPullRequestMock struct {
	Title string
	Body  string
}

func NewMockedHTTPClient() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesByOwnerByRepo,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				decoder := json.NewDecoder(r.Body)
				var issue ghIssueMock
				err := decoder.Decode(&issue)
				if err != nil {
					panic(err)
				}
				w.Write(mock.MustMarshal(github.Issue{
					Title:   github.String(issue.Title),
					Body:    github.String(issue.Body),
					Number:  github.Int(100),
					State:   github.String("open"),
					HTMLURL: github.String("https://github.com/elastic/observability-robots/issues/100"),
				}))
			}),
		),
		mock.WithRequestMatchHandler(
			mock.PostReposPullsByOwnerByRepo,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				decoder := json.NewDecoder(r.Body)
				var pullRequest ghPullRequestMock
				err := decoder.Decode(&pullRequest)
				if err != nil {
					panic(err)
				}
				w.Write(mock.MustMarshal(github.PullRequest{
					Title:   github.String(pullRequest.Title),
					Body:    github.String(pullRequest.Body),
					Number:  github.Int(101),
					State:   github.String("open"),
					HTMLURL: github.String("https://github.com/elastic/observability-robots/pulls/101"),
				}))
			}),
		),
		mock.WithRequestMatchHandler(
			mock.GetReposCollaboratorsPermissionByOwnerByRepoByUsername,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pathConponents := strings.Split(r.URL.Path, "/")
				username := pathConponents[len(pathConponents)-2]
				if username == "obltmachine" {
					w.Write(mock.MustMarshal(github.RepositoryPermissionLevel{
						Permission: github.String("admin"),
					}))
				} else {
					w.Write(mock.MustMarshal(github.RepositoryPermissionLevel{
						Permission: github.String("read"),
					}))
				}
			}),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			github.PullRequest{
				Title:   github.String("Mock pull request"),
				Body:    github.String("Mock pull request body"),
				Number:  github.Int(101),
				State:   github.String(randomState()),
				HTMLURL: github.String("https://github.com/elastic/observability-robots/pulls/101"),
			},
		),
		mock.WithRequestMatch(
			mock.PatchReposPullsByOwnerByRepoByPullNumber,
			github.PullRequest{
				Title:   github.String("Mock pull request"),
				Body:    github.String("Mock pull request body"),
				Number:  github.Int(101),
				State:   github.String("closed"),
				HTMLURL: github.String("https://github.com/elastic/observability-robots/pulls/101"),
			},
		),
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepo,
			[]github.PullRequest{
				{
					Title:   github.String("Mock pull request"),
					Body:    github.String("Mock pull request body"),
					Number:  github.Int(101),
					State:   github.String("open"),
					HTMLURL: github.String("https://github.com/elastic/observability-robots/pulls/101"),
				},
			},
		),
		mock.WithRequestMatch(
			mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
			[]github.Label{
				{
					Name: github.String("label1"),
				},
				{
					Name: github.String("label2"),
				},
			},
		),
		mock.WithRequestMatch(
			mock.GetReposContentsByOwnerByRepoByPath,
			[]github.RepositoryContent{
				{
					Name:    github.String("8.13.properties"),
					Type:    github.String("file"),
					Size:    github.Int(100),
					Path:    github.String("cd/release/versions/8.13.properties"),
					Content: github.String("version=8.13.3\nreleaseDate=2024-05-02\nbuildId=8.13.3-xxxxxxxx"),
				},
				{
					Name:    github.String("8.14.properties"),
					Type:    github.String("file"),
					Size:    github.Int(100),
					Path:    github.String("cd/release/versions/8.14.properties"),
					Content: github.String("version=8.14.3\nreleaseDate=2024-05-02\nbuildId=8.14.3-xxxxxxxx"),
				},
			},
			github.RepositoryContent{
				Name:    github.String("8.13.properties"),
				Type:    github.String("file"),
				Size:    github.Int(100),
				Path:    github.String("cd/release/versions/8.13.properties"),
				Content: github.String("version=8.13.3\nreleaseDate=2024-05-02\nbuildId=8.13.3-xxxxxxxx"),
			},
			github.RepositoryContent{
				Name:    github.String("8.14.properties"),
				Type:    github.String("file"),
				Size:    github.Int(100),
				Path:    github.String("cd/release/versions/8.14.properties"),
				Content: github.String("version=8.14.3\nreleaseDate=2024-05-02\nbuildId=8.14.3-xxxxxxxx"),
			},
		),
	)
}

func randomState() string {
	states := []string{"open", "closed", "merged"}
	return states[time.Now().Unix()%3]
}

// CreateIssue creates a new issue
func CreateIssue(title, body, owner, repo string, labels []string) (issue *github.Issue, err error) {
	client, ctx := GetClient()

	input := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	issue, _, err = client.Issues.Create(ctx, owner, repo, input)
	return issue, err
}

// CreateIssueWithLabels creates a new issue with some labels
func CreateIssueWithLabels(title string, body string, labels []string, owner string, repo string) (issue *github.Issue, err error) {
	client, ctx := GetClient()

	input := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	issue, _, err = client.Issues.Create(ctx, owner, repo, input)
	return issue, err
}

// GetClient returns a new GitHub API client using the environment variable
// GITHUB_TOKEN or GITHUB_PASSWORD
func GetClient() (*github.Client, context.Context) {
	if MockedHTTPClient != nil {
		return github.NewClient(MockedHTTPClient), context.Background()
	}
	// Connect to the private GitHub repository
	ctx := context.Background()
	token := os.Getenv("GITHUB_PASSWORD")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		logger.Fatal("GITHUB_TOKEN or GITHUB_PASSWORD environment variable not set")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

// CreatePullRequest creates a new pull request
func CreatePullRequest(title, body, owner, repo, base, head string) (pr *github.PullRequest, response *github.Response, err error) {
	client, ctx := GetClient()

	input := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Base:  github.String(base),
		Head:  github.String(head),
	}

	pullRequest, response, err := client.PullRequests.Create(ctx, owner, repo, input)
	if err == nil {
		logger.Debugf("Created pull request: %s", pullRequest.GetHTMLURL())
	}
	return pullRequest, response, err
}

// CreatePullRequestWithLabels creates a new pull request with the given labels
func CreatePullRequestWithLabels(title, body, owner, repo, base, head string, labels []string) (pr *github.PullRequest, response *github.Response, err error) {
	var pullRequest *github.PullRequest

	pullRequest, response, err = CreatePullRequest(title, body, owner, repo, base, head)
	if err != nil {
		return nil, response, err
	}

	if len(labels) != 0 {
		client, ctx := GetClient()
		// If labels don't exist then they will be created.
		_, response, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, pullRequest.GetNumber(), labels)
		if err == nil {
			logger.Debugf("Edited pull request: %s", pullRequest.GetHTMLURL())
		}
	}
	return pullRequest, response, err
}

// WaitPullRequestState waits until the pull request has a given state
func WaitPullRequestState(prNumber int, owner, repo string, expectedState string) (currentState, mergeableState string, err error) {
	logger.Debugf("Waiting for pull request: %d", prNumber)
	client, ctx := GetClient()
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err == nil {
		currentState = pr.GetState()
		mergeableState = pr.GetMergeableState()
		logger.Debugf("Compare Pull request status: %s -> %s - %s", expectedState, currentState, mergeableState)
		if strings.Compare(currentState, expectedState) != 0 {
			err = fmt.Errorf("pull request '%d' has an incorrect state '%s' while it's expected to be '%s'", prNumber, currentState, expectedState)
		}
	}
	return currentState, mergeableState, err
}

// WaitPullRequestStateRetry waits until the pull request has a given state, retrying with backoff.
func WaitPullRequestStateRetry(prNumber int, owner, repo string, expectedState string, timeout time.Duration) (currentState, mergeableState string, err error) {
	logger.Debugf("Compare Pull request status: %s -> %s - %s timeout: %s", expectedState, currentState, mergeableState, timeout)
	_wait := func() (err error) {
		currentState, mergeableState, err = WaitPullRequestState(prNumber, owner, repo, expectedState)
		return err
	}
	err = Retry(_wait, timeout)
	// Let's help with a log message saying a timeout has been reached if err is not unset
	if err != nil {
		fmt.Printf("%s\n", logger.WarnColor.Sprint(fmt.Sprintf("timeout after '%s'", timeout.String())))
	}
	return currentState, mergeableState, err
}

// Retry retries an operation during a timeout.
func Retry(operation backoff.Operation, timeout time.Duration) (err error) {
	exp := backoff.NewExponentialBackOff()
	exp.InitialInterval = 5 * time.Second
	exp.RandomizationFactor = 0.5
	exp.Multiplier = 2.0
	exp.MaxInterval = 15 * time.Second
	exp.MaxElapsedTime = timeout
	err = backoff.Retry(operation, exp)
	return err
}

// ClosePullRequest closes a pull request
func ClosePullRequest(prNumber int, owner, repo string) (pr *github.PullRequest, err error) {
	logger.Debugf("Closing pull request: %d", prNumber)
	client, ctx := GetClient()
	pr, _, err = client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err == nil {
		pr.State = github.String("closed")
		pr, _, err = client.PullRequests.Edit(ctx, owner, repo, prNumber, pr)
	}
	return pr, err
}

// FindPullRequest finds a pull request by title
func FindPullRequest(title string, owner, repo string) (pr *github.PullRequest, err error) {
	logger.Debugf("Searching for pull request with title: %s", title)
	client, ctx := GetClient()
	prs, _, err := client.PullRequests.List(ctx, owner, repo, nil)
	if err == nil {
		for _, item := range prs {
			if strings.Contains(item.GetTitle(), title) {
				pr = item
				break
			}
		}
		if pr == nil {
			err = errors.New("pull request not found")
		}
	}
	return pr, err
}

// FindPullRequestRetry finds a pull request by title, retrying with backoff.
func FindPullRequestRetry(title string, owner, repo string, timeout time.Duration) (pr *github.PullRequest, err error) {
	logger.Debugf("Searching for pull request with title: %s timeout: %s", title, timeout)
	_wait := func() (err error) {
		pr, err = FindPullRequest(title, owner, repo)
		return err
	}
	err = Retry(_wait, timeout)
	return pr, err
}

// DeleteRef deletes a ref from a repository
func DeleteRef(owner, repo, ref string) (err error) {
	logger.Debugf("Deleting ref: %s from the repo %s/%s", ref, owner, repo)
	client, ctx := GetClient()
	_, err = client.Git.DeleteRef(ctx, owner, repo, ref)
	return err
}
