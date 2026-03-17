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

// Package releases It contains the functions to interact with the Unified Releases.
package releases

import (
	"bufio"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	gh "github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"

	"github.com/google/go-github/github"
)

const (
	org        = "elastic"
	repo       = "infra"
	folderPath = "cd/release/versions"
)

// Versions defines a version file defined in the Unified Release process within the infra repository
type Versions struct {
	Branch      string
	Filename    string
	Version     string
	ReleaseDate string
}

// IsUserAdmin whether the given GitHub user has admin privileges in the given GitHub repository
func IsUserAdmin(repoName string, user string) (bool, error) {
	client, ctx := gh.GetClient()

	rpl, _, err := client.Repositories.GetPermissionLevel(ctx, org, repoName, user)
	if err != nil {
		return false, fmt.Errorf("failed to get the permissions in the %s repo: %w", repoName, err)
	}

	if *rpl.Permission == "admin" {
		return true, nil
	}
	return false, nil
}

// GetVersions searches for all the active branches in the Unified Release
// and parse them
func GetVersions() ([]Versions, error) {
	client, ctx := gh.GetClient()

	// Get the versions directory from the infra repository
	_, dirContent, _, err := client.Repositories.GetContents(ctx, org, repo, folderPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s directory: %w", folderPath, err)
	}

	allProperties := make([]Versions, len(dirContent))
	for i, fileContent := range dirContent {
		versionFile, err := getVersionFile(ctx, client, fileContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version file: %w", err)
		}
		allProperties[i] = versionFile
	}

	return allProperties, nil
}

// getVersionFile parses a version file from a given GH filepath
// The filepath must be a valid .properties file
func getVersionFile(ctx context.Context, client *github.Client, fileContent *github.RepositoryContent) (Versions, error) {
	versions := Versions{
		Branch:      fileContent.GetName()[:len(fileContent.GetName())-len(filepath.Ext(fileContent.GetName()))],
		Filename:    fileContent.GetName(),
		Version:     "unknown",
		ReleaseDate: "unknown",
	}

	// Fetch a given file
	content, _, _, err := client.Repositories.GetContents(ctx, org, repo, fileContent.GetPath(), nil)
	if err != nil {
		return versions, err
	}

	decodedContent, err := content.GetContent()
	if err != nil {
		return versions, err
	}

	// Parse properties file
	scanner := bufio.NewScanner(strings.NewReader(decodedContent))
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				switch key {
				case "branch":
					versions.Branch = value
				case "fileName":
					versions.Filename = value
				case "releaseDate":
					versions.ReleaseDate = value
				case "version":
					versions.Version = value
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
		return versions, err
	}

	return versions, nil
}
