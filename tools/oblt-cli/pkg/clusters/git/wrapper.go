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

package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	gitAdd "github.com/ldez/go-git-cmd-wrapper/v2/add"
	gitCheckout "github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	gitClone "github.com/ldez/go-git-cmd-wrapper/v2/clone"
	gitCommit "github.com/ldez/go-git-cmd-wrapper/v2/commit"
	gitConfig "github.com/ldez/go-git-cmd-wrapper/v2/config"
	"github.com/ldez/go-git-cmd-wrapper/v2/fetch"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/pull"
	"github.com/ldez/go-git-cmd-wrapper/v2/reset"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	gitTypes "github.com/ldez/go-git-cmd-wrapper/v2/types"
)

type GitService interface {
	Clone(name, repository, branch, basePath string, fast bool) (path string, err error)
	Add(path string) error
	Commit(message string) (string, error)
	Push() error
	Pull(fast bool) error
	GetChanges(basePath string) ([]string, error)
	GetLastCommitMessage() (string, error)
	GetCurrentBranch() (string, error)
	NewBranch(name string) error
	RevParse() (string, error)
	PullBranch(branch string) (err error)
	Reset() (err error)
	Checkout(branch string) (err error)
	GitUser() (string, error)
	GitEmail() (string, error)
	Status() (string, error)
	HasChanges() bool
}

// cmdExecutorMock is a mock executor for testing purposes
func cmdExecutorMock(_ context.Context, name string, _ bool, args ...string) (string, error) {
	cmd := fmt.Sprintf("%s %s", name, strings.Join(args, " "))
	logger.Debugf("cmdExecutorMock: %s", cmd)
	return fmt.Sprintln(cmd), nil
}

// GitWrapper represents a git wrapper
type GitWrapper struct {
	CmdExecutor gitTypes.Executor
}

// NewGitWrapper creates a new GitWrapper
func NewGitWrapper() *GitWrapper {
	return &GitWrapper{}
}

// NewGitWrapperMock creates a new GitWrapper with a mock executor
func NewGitWrapperMock() *GitWrapper {
	return &GitWrapper{
		CmdExecutor: cmdExecutorMock,
	}
}

// WithGitWrapperMock sets the mock executor to the given GitWrapper
func WithGitWrapperMock(wrapper *GitWrapper) *GitWrapper {
	wrapper.CmdExecutor = cmdExecutorMock
	return wrapper
}

// CmdExecutor Allow to override the Git command call (useful for testing purpose).
func (w *GitWrapper) cmdExecutor(executor gitTypes.Executor) gitTypes.Option {
	if executor == nil {
		return func(g *gitTypes.Cmd) {}
	}
	return func(g *gitTypes.Cmd) {
		g.Executor = executor
	}
}

// GetChanges retrieves the changes between the current and the given base ref
func (w *GitWrapper) GetChanges(base string) ([]string, error) {
	out, err := git.Raw("diff", func(g *gitTypes.Cmd) {
		g.AddOptions("--name-status")
		g.AddOptions(base)
	},
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor))
	linesAll := strings.Split(out, "\n")
	var lines []string
	for _, line := range linesAll {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, err
}

// clone clones the git repository represented by the struct
func (w *GitWrapper) Clone(name, repository, branch, basePath string, fast bool) (path string, err error) {
	logger.Debugf("git.clone.project: %s", repository)

	options := []gitTypes.Option{
		gitClone.Repository(repository),
		gitClone.Branch(branch),
		gitClone.Config("push.default", "current"),
		gitClone.Config("pull.default", "current"),
		gitClone.Directory(name),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor),
	}
	if fast {
		options = append(options, gitClone.Depth("1"))
		options = append(options, gitClone.SingleBranch)
	}
	if err = os.Chdir(basePath); err == nil {
		_, err = logger.LogError(git.Clone(options...))
		if err == nil {
			path = filepath.Join(basePath, name)
		}
	}
	return path, err
}

// add it add the pathSpec files/folder to the git index.
func (w *GitWrapper) Add(pathSpec string) (err error) {
	logger.Debugf("git.add")
	_, err = logger.LogError(git.Add(
		gitAdd.PathSpec(pathSpec),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	return err
}

// commit it commits the current changes.
func (w *GitWrapper) Commit(msg string) (sha string, err error) {
	logger.Debugf("git.commit")
	_, err = logger.LogError(git.Commit(
		gitCommit.Message(msg),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	if err == nil {
		sha, err = logger.LogError(w.RevParse())
	}
	return sha, err
}

// push it pushes the current pending commits.
func (w *GitWrapper) Push() (err error) {
	logger.Debugf("git.push")
	w.createTemporalBranch()
	_, err = logger.LogError(git.Push(
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	return err
}

// createTemporalBranch creates a temporal branch from the current SHA.
func (w *GitWrapper) createTemporalBranch() error {
	logger.Debugf("createTemporalBranch")
	newBranch, err := w.newTempBranch()
	if err == nil {
		err = w.NewBranch(newBranch)
	}
	return err
}

// newTempBranch generate a temporal branch name from the current SHA.
func (w *GitWrapper) newTempBranch() (newBranch string, err error) {
	if sh, err := w.RevParse(); err == nil {
		newBranch = fmt.Sprintf("oblt-cli/%s", sh)
	}
	return newBranch, err
}

// revParse retrieves last commit ID
func (w *GitWrapper) RevParse() (string, error) {
	logger.Debugf("git.rev-parse")
	sha, err := git.RevParse(
		revparse.Args("HEAD"),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor))
	return strings.ReplaceAll(sha, "\n", ""), err
}

// GetLastCommitMessage retrieves last commit message
func (w *GitWrapper) GetLastCommitMessage() (string, error) {
	out, err := git.Raw("log",
		func(g *gitTypes.Cmd) {
			g.AddOptions("-1")
			g.AddOptions("--pretty=%B")
		},
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor))
	// replace the last newline character
	return strings.ReplaceAll(out, "\n", ";"), err
}

// GetCurrentBranch retrieves the current branch name
func (w *GitWrapper) GetCurrentBranch() (string, error) {
	out, err := git.RevParse(
		revparse.AbbrevRef(""),
		revparse.Args("HEAD"),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor))
	return strings.ReplaceAll(out, "\n", ""), err
}

// pull moves to the repository path and pulls from the configured upstream
func (w GitWrapper) Pull(fast bool) (err error) {
	logger.Debugf("git.pull")

	options := []gitTypes.Option{
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor),
	}
	if fast {
		options = append(options, pull.Depth("1"))
		options = append(options, pull.Ff)
	}
	_, err = logger.LogError(git.Pull(options...))
	return err
}

func (w *GitWrapper) PullBranch(branch string) (err error) {
	logger.Debugf("git.pull.branch")
	_, err = logger.LogError(git.Pull(
		originRef,
		func(g *gitTypes.Cmd) {
			g.AddOptions(branch)
		},
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	return err
}

// newBranch creates a new branch from the current branch.
func (w *GitWrapper) NewBranch(branch string) (err error) {
	logger.Debugf("Create new branch: %s", branch)
	_, err = logger.LogError(git.Checkout(
		gitCheckout.NewBranch(branch),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	return err
}

// Checkout get the branch/sha from the repository.
func (w *GitWrapper) Checkout(branch string) (err error) {
	logger.Debugf("git.checkout")
	_, err = logger.LogError(git.Checkout(
		gitCheckout.Branch(branch),
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor)))
	return err
}

// Reset it resets the current branch to the last commit.
func (w *GitWrapper) Reset() (err error) {
	logger.Debugf("git.reset")
	if _, err = logger.LogError(git.Fetch(fetch.NoTags, fetch.Remote(DefaultRemote), fetch.RefSpec(DefaultBranch))); err == nil {
		_, err = logger.LogError(git.Reset(
			reset.Hard,
			mainRef,
			git.Debugger(logger.Verbose),
			w.cmdExecutor(w.CmdExecutor)))
	}
	return err
}

// mainRef is a helper function to add the ref origin/main as git argument.
func mainRef(g *gitTypes.Cmd) {
	g.AddOptions(DefaultRemote + "/" + DefaultBranch)
}

// originRef is a helper function to add the ref origin as git argument.
func originRef(g *gitTypes.Cmd) {
	g.AddOptions(DefaultRemote)
}

// GitUser retrieves the git user
func (w *GitWrapper) GitUser() (string, error) {
	out, err := git.Config(gitConfig.Global, gitConfig.Get("user.name", ""), w.cmdExecutor(w.CmdExecutor))
	return strings.ReplaceAll(out, "\n", ""), err
}

// GitEmail retrieves the git email
func (w *GitWrapper) GitEmail() (string, error) {
	out, err := git.Config(gitConfig.Global, gitConfig.Get("user.email", ""), w.cmdExecutor(w.CmdExecutor))
	return strings.ReplaceAll(out, "\n", ""), err
}

// Status retrieves the git status
func (w *GitWrapper) Status() (string, error) {
	out, err := git.Status(
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor))
	return out, err
}

// hasChanges checks if there are changes in the repository
func (w *GitWrapper) HasChanges() bool {
	options := []gitTypes.Option{
		git.Debugger(logger.Verbose),
		w.cmdExecutor(w.CmdExecutor),
		func(g *gitTypes.Cmd) {
			g.AddOptions("--exit-code")
			g.AddOptions("--quiet")
			g.AddOptions("HEAD")
		},
	}
	_, err := git.Raw("diff", options...)
	return err != nil
}
