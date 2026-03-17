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
package clusters

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/viper"
)

const (
	itemsSeparator = "-"
)

// FakeConfigInfo is a struct that contains the information needed to create a fake config
type FakeConfigInfo struct {
	CurrentUsername  string
	SlackChannel     string
	TempDir          string
	Obltrepo         *ObltEnvironmentsRepository
	EnvironmentsPath string
	UsersFolderPath  string
}

// FakeConfigViper creates a fake config file and returns the information needed to use it
func FakeConfigViper(tempDir string) FakeConfigInfo {
	currentUsername := "the-user"
	slackChannel := "#0123456"

	// configuration is needed because the repo is cloned under config's workspace
	viper.SetConfigFile(filepath.Join(tempDir, "config.yml"))
	viper.Set(config.UsernameFlag, currentUsername)
	viper.Set(config.SlackChannelFlag, slackChannel)
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	viper.WriteConfig()

	// clone is needed to access all methods, using dry run
	obltRepo := NewObltTestEnvironmentsFromViper(true)
	gitRepo := obltRepo.Repo(true)
	gitRepo.Sync()
	environmentsPath := filepath.Join(tempDir, gitRepo.Name, "environments")
	usersFolderPath := filepath.Join(environmentsPath, "users")

	return FakeConfigInfo{
		CurrentUsername:  currentUsername,
		SlackChannel:     slackChannel,
		TempDir:          tempDir,
		EnvironmentsPath: environmentsPath,
		UsersFolderPath:  usersFolderPath,
		Obltrepo:         obltRepo,
	}
}

// FakeConfig creates a fake config file and returns the information needed to use it
func FakeConfig(tempDir string) FakeConfigInfo {
	currentUsername := "the-user"
	slackChannel := "#0123456"
	userConfig := config.NewObltConfig(
		filepath.Join(tempDir, "config.yml"),
		currentUsername,
		slackChannel,
		true,
		true,
	)
	obltRepo := NewObltTestEnvironments(userConfig, true)
	gitRepo := obltRepo.Repo(true)
	gitRepo.Sync()
	environmentsPath := filepath.Join(tempDir, gitRepo.Name, "environments")
	usersFolderPath := filepath.Join(environmentsPath, "users")
	return FakeConfigInfo{
		CurrentUsername:  currentUsername,
		SlackChannel:     slackChannel,
		TempDir:          tempDir,
		EnvironmentsPath: environmentsPath,
		UsersFolderPath:  usersFolderPath,
		Obltrepo:         obltRepo,
	}
}

// buildEsDockerImageString build the string representation for the Elasticsearch docker image
func buildEsDockerImageString(dockerImageVersion string) string {
	return fmt.Sprintf("docker.elastic.co/observability-ci/elasticsearch-cloud-ess:%s", dockerImageVersion)
}

// buildKibanaDockerImageString build the string representation for the Kibana docker image
func buildKibanaDockerImageString(dockerImageVersion string) string {
	return fmt.Sprintf("docker.elastic.co/observability-ci/kibana-cloud:%s", dockerImageVersion)
}

// buildElasticAgentDockerImageString build the string representation for the Elastic Agent docker image
func buildElasticAgentDockerImageString(dockerImageVersion string) string {
	return fmt.Sprintf("docker.elastic.co/observability-ci/elastic-agent-cloud:%s", dockerImageVersion)
}

// nowISOFormat return the current timestamp in ISO-8601 format
func nowISOFormat() (date string) {
	return time.Now().Format("2006-01-02T15:04:05-0700")
}

// newClusterName returns a new cluster name, which is a combination of the prefix, a string and the suffix,
// in this particular order:
// - the prefix will always contain the first 20 characters.
func newClusterName(prefix string, suffix string, name string) (string, error) {
	var err, err1 error
	clusterName := name
	if prefix != "" {
		err = config.ValidateNames(prefix)
	}
	if suffix != "" {
		err1 = config.ValidateNames(suffix)
	}
	if err == nil && err1 == nil {
		if len(prefix) > maxPrefixLength {
			prefix = prefix[:maxPrefixLength]
		}
		items := []string{prefix, clusterName, suffix}
		clusterName = joinIgnoreEmpty(items, itemsSeparator)
	}
	return clusterName, errors.Join(err, err1)
}

// NormalizeElasticStackVersion returns the normalized version of the Elastic Stack
func normalizeElasticStackVersion(version string, isRelease bool) (stackVersion string, dockerImageVersion string) {
	vesionSuffix := "-SNAPSHOT"
	dockerImageVersion = string(version)
	stackVersion = strings.Split(version, "-")[0]
	if !strings.Contains(dockerImageVersion, "-") && !isRelease {
		dockerImageVersion = string(stackVersion)
	}
	if strings.Contains(dockerImageVersion, "-") && !isRelease {
		stackVersion = fmt.Sprintf("%s%s", stackVersion, vesionSuffix)
	}
	return stackVersion, dockerImageVersion
}

// joinIgnoreEmpty joins the items using the separator, ignoring the empty ones
func joinIgnoreEmpty(items []string, separator string) string {
	var parts []string
	for _, s := range items {
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, separator)
}
