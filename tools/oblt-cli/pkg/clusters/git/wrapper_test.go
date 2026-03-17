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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GitWrapper(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	repo := fmt.Sprintf("%s%s/%s.git", httpPrefix, "octocat", "Hello-World")
	defer os.Remove(tmpDir)

	git := NewGitWrapper()

	user, err := git.GitUser()
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	email, err := git.GitEmail()
	assert.NoError(t, err)
	assert.NotEmpty(t, email)

	repoPath, err := git.Clone("test-fast", repo, "master", tmpDir, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, repoPath)

	repoPath, err = git.Clone("test", repo, "master", tmpDir, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, repoPath)

	err = os.Chdir(repoPath)
	assert.NoError(t, err)

	branchName, err := git.GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "master", branchName)

	git.Checkout("553c2077f0edc3d5dc5d17262f6aa498e69d6f8e")
	commitMessage, err := git.GetLastCommitMessage()
	assert.NoError(t, err)
	assert.Equal(t, "first commit;;", commitMessage)

	err = git.PullBranch("master")
	assert.NoError(t, err)

	changes, err := git.GetChanges(repoPath)
	assert.NoError(t, err)
	assert.Empty(t, changes)

	ref, err := git.RevParse()
	assert.NoError(t, err)
	assert.NotNil(t, ref)

	err = git.Checkout("master")
	assert.NoError(t, err)

	ref, err = git.GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "master", ref)

	err = git.Pull(true)
	assert.NoError(t, err)

	err = git.NewBranch("test")
	assert.NoError(t, err)

	branchName, err = git.GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "test", branchName)

	err = os.WriteFile(filepath.Join(repoPath, "README"), []byte("test"), 0644)
	assert.NoError(t, err)

	changes, err = git.GetChanges(repoPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, changes)

	err = git.Add(repoPath)
	assert.NoError(t, err)

	commitsha, err := git.Commit("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, commitsha)

	WithGitWrapperMock(git)
	err = git.Push()
	assert.NoError(t, err)
}
