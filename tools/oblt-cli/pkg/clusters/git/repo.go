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

// Package git contains the commands to manage the changes in the oblt test environments repository.
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"golang.org/x/exp/maps"
)

const (
	sshPrefix  = "git@github.com:"
	httpPrefix = "https://github.com/"

	// ModeSSH use SSH to make the git operation
	ModeSSH = "ssh"
	// ModeHTTP use HTTP to make the git operation
	ModeHTTP = "http"

	// DefaultBranch the default branch of the repository
	DefaultBranch = "main"
	// DefaultRemote the default remote of the repository
	DefaultRemote = "origin"
	// DefaultOwner the default owner of the repository
	DefaultOwner = "elastic"
	// ObltStatusRepo the name of the repository containing the status of the observability system
	ObltStatusRepo = "observability-system-status"
	// ObltRepoName the name of the repository containing the test environments
	ObltRepoName = "observability-test-environments"

	// Filter those clusters that have been created
	FilterCreateCluster = "A\tenvironments/users/"

	// Filter those clusters that have been updated
	FilterUpdateCluster = "M\tenvironments/users/"

	// Filter those clusters that have been renamed
	FilterRenameCluster = "R\tenvironments/users/"

	// Filter those clusters that have been destroyed
	FilterDestroyCluster = "D\tenvironments/users/"

	// GitHub Label for the creation of a cluster
	LabelCreate = "cluster:create"

	// GitHub Label for the update of a cluster
	LabelUpdate = "cluster:update"

	// GitHub Label for the destroy of a cluster
	LabelDestroy = "cluster:destroy"

	// Prefix for the GitHub Label related to the cluster management
	LabelClusterPrefix = "cluster"
)

// Repository represents a git repository
type Repository struct {
	// AlwaysCleanUp indicates if the repository should be cleaned up before the execution of the command
	AlwaysCleanUp bool
	// Branch of the repository
	Branch string
	// dryRun indicates if the git operations should be executed or not
	dryRun bool
	// Name of the repository
	Name string
	// Owner of the repository
	Owner string
	// path to the checkout repository
	path string // it must not be initialised nor modified externally
	// folder containing the repository folder
	parentPath string // it must not be initialised nor modified externally
	// SSHMode indicates if the git operations should be executed using SSH or HTTP
	SSHMode bool
	// Git wrapper
	git GitService
}

// NewRepository returns a new repository instance
func NewRepository(folder, owner, name, branch string, dryRun, alwaysCleanup, httpMode bool) *Repository {
	return &Repository{
		AlwaysCleanUp: alwaysCleanup,
		Name:          name,
		dryRun:        dryRun,
		Owner:         owner,
		Branch:        branch,
		path:          folder,
		parentPath:    filepath.Dir(folder),
		SSHMode:       !httpMode,
		git:           NewGitWrapper(),
	}
}

// WithRepositoryMock returns a new repository instance with a mock git wrapper
func WithRepositoryMock(repo *Repository) *Repository {
	tmpFolder := os.TempDir()
	return &Repository{
		AlwaysCleanUp: repo.AlwaysCleanUp,
		Name:          repo.Name,
		dryRun:        false,
		Owner:         repo.Owner,
		Branch:        repo.Branch,
		path:          filepath.Join(tmpFolder, repo.Name),
		parentPath:    tmpFolder,
		SSHMode:       repo.SSHMode,
		git:           NewGitWrapperMock(),
	}
}

// chdir changes the current working directory to the given path
func (g *Repository) chdir(path string) error {
	logger.Debugf("Changing to dir: %s", path)
	return os.Chdir(path)
}

// chdirRepo changes the current working directory to the repository path
func (g *Repository) chdirRepo() error {
	return g.chdir(g.path)
}

// Sync Sync a project repository, cloning or pulling if the repository already exists
func (g *Repository) Sync() (opErr error) {
	_, err := os.Stat(g.path)
	if os.IsNotExist(err) {
		// directory does not exist: clone and chdir
		if opErr = g.chdir(g.parentPath); opErr == nil {
			if g.path, opErr = g.git.Clone(g.Name, g.String(), g.Branch, g.parentPath, true); opErr == nil {
				opErr = g.chdirRepo()
			}
		}
	} else {
		// directory exists: pull
		if opErr = g.chdirRepo(); opErr == nil {
			// remove all changes
			opErr = g.resetChanges(opErr)
		}
	}
	return opErr
}

// resetChanges resets the changes in the repository
// if the AlwaysCleanUp is enabled
func (g *Repository) resetChanges(opErr error) error {
	if g.AlwaysCleanUp {
		if opErr = g.git.Checkout(g.Branch); opErr == nil {
			opErr = g.git.Reset()
		}
	}
	return opErr
}

// CommitAndPush It adds all changes in the current folder,then commit and push the changes.
func (g *Repository) CommitAndPush(menssage string) (sha string, err error) {
	if err = g.chdirRepo(); err == nil {
		pwd, _ := os.Getwd()
		logger.Debugf("Add changes at: %s", pwd)
		g.git.Add(".")
		if g.git.HasChanges() {
			logger.Debugf("Commit changes")
			sha, err = g.git.Commit(menssage)
			logger.Debugf("dry.run: %t", g.dryRun)
			if !g.dryRun {
				if err == nil {
					g.git.Push()
				}
			}
		}
	}

	return sha, err
}

// WithDryRun Defines the dry-run mode for all operations for this repository
func (g *Repository) WithDryRun(dryRun bool) *Repository {
	logger.Debugf("dry.run: %t", dryRun)
	g.dryRun = dryRun
	return g
}

// FullName returns the full name of the repository in the form owner/name
func (g *Repository) FullName() string {
	return fmt.Sprintf("%s/%s", g.Owner, g.Name)
}

// Path returns the path to the repository
func (g *Repository) Path() string {
	return g.path
}

// String returns the http URL to the repository
func (g *Repository) String() string {
	var prefix = httpPrefix
	if g.SSHMode {
		prefix = sshPrefix
	}

	return fmt.Sprintf("%s%s/%s.git", prefix, g.Owner, g.Name)
}

// IsDryRun returns true if the dry-run mode is enabled
func (g *Repository) IsDryRun() bool {
	return g.dryRun
}

// LabelsClusterType returns the labels for the type of changes
func LabelsClusterType(changes []string) []string {
	labels := make(map[string]bool)
	for _, change := range changes {
		maps.Copy(labels, getLabelsGivenAChange(change))
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}

	return keys
}

// getLabelsGivenAChange returns the labels for the given change
func getLabelsGivenAChange(change string) map[string]bool {
	labels := make(map[string]bool)
	// Ensure we only monitor files that are ending with .yml
	// and contains environments/users/
	if !strings.HasSuffix(change, ".yml") {
		return labels
	}
	filename := strings.Split(change, "\t")[1]
	if strings.HasPrefix(change, FilterCreateCluster) {
		labels[LabelCreate] = true
		maps.Copy(labels, getLabelsGivenACluster(filename))
	}
	if strings.HasPrefix(change, FilterUpdateCluster) {
		labels[LabelUpdate] = true
		maps.Copy(labels, getLabelsGivenACluster(filename))
	}
	if strings.HasPrefix(change, FilterRenameCluster) {
		labels[LabelUpdate] = true
		lines := strings.Split(change, "\t")
		if len(lines) > 2 {
			maps.Copy(labels, getLabelsGivenACluster(lines[2]))
		}
	}
	if strings.HasPrefix(change, FilterDestroyCluster) {
		labels[LabelDestroy] = true
	}
	return labels
}

// getLabelsGivenACluster returns the labels for the given cluster file
func getLabelsGivenACluster(filename string) map[string]bool {
	labels := make(map[string]bool)
	yamlFile, err := files.ReadYamlOj(filename, "NA")

	if err != nil {
		return labels
	}

	// template_name: ccs - cluster:template-ccs
	if yamlFile.Data["template_name"] != nil {
		templateName := yamlFile.Data["template_name"].(string)
		if templateName != "" {
			labels[fmt.Sprintf("%s:template-%s", LabelClusterPrefix, templateName)] = true
		}
	}

	// stack.mode: "serverless" - cluster:stack-mode-serverless
	// stack.template: "observability" - cluster:stack-template-observability
	// stack.target: "qa" - cluster:stack-target-qa
	stack, _ := yamlFile.Data["stack"].(map[string]interface{})
	mode, _ := stack["mode"].(string)
	template, _ := stack["template"].(string)
	target, _ := stack["target"].(string)

	if mode != "" {
		labels[fmt.Sprintf("%s:stack-mode-%s", LabelClusterPrefix, mode)] = true
	}
	if template != "" {
		labels[fmt.Sprintf("%s:stack-template-%s", LabelClusterPrefix, template)] = true
	}
	if target != "" {
		labels[fmt.Sprintf("%s:stack-target-%s", LabelClusterPrefix, target)] = true
	}
	return labels
}

// GetCurrentBranch returns the current branch of the repository
func (g *Repository) GetCurrentBranch() (branch string, err error) {
	err = g.chdirRepo()
	if err == nil {
		branch, err = g.git.GetCurrentBranch()
	}
	return branch, err
}

// GetChanges returns the changes between the base and the current branch
func (g *Repository) GetChanges(base string) (changes []string, err error) {
	err = g.chdirRepo()
	if err == nil {
		changes, err = g.git.GetChanges(base)
	}
	return changes, err
}
