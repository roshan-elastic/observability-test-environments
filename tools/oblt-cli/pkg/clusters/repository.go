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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/git"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/github"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

// maxPrefixLength the max length for prefixes when naming a cluster
const maxPrefixLength = 20

// environmentsPath Returns the absolute path to the oblt environments folder.
const environmentsPath = "environments"

// EnvFilesFolder relative path to the environments files in the configuration folder.
var EnvFilesFolder = filepath.Join(git.ObltRepoName, UsersFolder)

// UsersFolder relative path to the Users environments folder.
var UsersFolder = filepath.Join(environmentsPath, "users")

// obltTestEnvironmentsCommitURL url of the commits for the OBLT Test Environments repository
const obltTestEnvironmentsCommitURL = "https://github.com/elastic/observability-test-environments/commit"

// ObltEnvironmentsRepository struct representing the operations on the oblt-test-environments repository
type ObltEnvironmentsRepository struct {
	gitRepo *git.Repository
	config  config.ObltConfiguration
}

// Repository interface for the oblt-test-environments repository
type Repository interface {
	GetPath() string
	GetVersion() (repoVersion semver.Version, err error)
}

// NewObltTestEnvironments return the repository for the observability test environments.
func NewObltTestEnvironments(userConfig config.ObltConfiguration, dryRun bool) *ObltEnvironmentsRepository {
	return &ObltEnvironmentsRepository{
		gitRepo: git.NewRepository(path.Join(userConfig.GetDir(), git.ObltRepoName), git.DefaultOwner, git.ObltRepoName, git.DefaultBranch, dryRun, true, userConfig.GitHttpMode),
		config:  userConfig,
	}
}

// NewObltSystemStatusFromViper return the repository for the system status project.
func NewObltTestEnvironmentsFromViper(dryRun bool) *ObltEnvironmentsRepository {
	return NewObltTestEnvironments(config.NewObltConfigFromViper(), dryRun)
}

// Bootstrap initialises the oblt test environments repository, using dryRun and git options
func BootstrapRepository(dryRun bool) (*ObltEnvironmentsRepository, error) {
	obltEnv := NewObltTestEnvironmentsFromViper(dryRun)
	err := obltEnv.Repo(dryRun).Sync()
	return obltEnv, err
}

// Bootstrap initialises the oblt test environments repository, using dryRun and git options
func BootstrapRepositoryUserConfig(userConfig config.ObltConfiguration, dryRun bool) (*ObltEnvironmentsRepository, error) {
	obltEnv := NewObltTestEnvironments(userConfig, dryRun)
	err := obltEnv.Repo(dryRun).Sync()
	return obltEnv, err
}

// OwnerRepo returns the Github username/repo pair
func (r *ObltEnvironmentsRepository) OwnerRepo() string {
	return fmt.Sprintf("%s/%s", r.gitRepo.Owner, r.gitRepo.Name)
}

// Repo returns the Git repository for the oblt-test-environments repository
func (r *ObltEnvironmentsRepository) Repo(dryRun bool) *git.Repository {
	r.gitRepo.WithDryRun(dryRun)

	return r.gitRepo
}

// DestroyCluster removes the configuration file for the cluster, pushing the changes to the Git repository
func (r *ObltEnvironmentsRepository) DestroyCluster(clusterName string, all bool) (parametersMap map[string]interface{}, err error) {
	logger.Infof("Destroying cluster : %s", clusterName)
	var file files.YamlFile
	var sha string
	if all {
		file, err = r.FindClusterConfig(clusterName)
	} else {
		file, err = r.FindCurrentUserClusterConfig(clusterName)
	}

	if r.IsGoldenClusterFile(file) {
		err = fmt.Errorf("forbidden to destroy the following clusters: %s", clusterName)
	}
	if err == nil {
		logger.Infof("Removing cluster file : %s", file.Path)
		os.Remove(file.Path)

		commitMessage := "oblt-cli(" + r.config.Username + "): Destroy cluster " + clusterName
		sha, err = r.gitRepo.CommitAndPush(commitMessage)

		parametersMap = make(map[string]interface{})
		parametersMap[clusterName] = file.Path
		parametersMap["CommitSha"] = sha
		parametersMap["CommitMessage"] = commitMessage
		parametersMap["CommitURL"] = fmt.Sprintf("%s/%s", obltTestEnvironmentsCommitURL, sha)
	}

	return parametersMap, err
}

// IsGoldenClusterFile returns true if the file is a golden cluster file
func (r *ObltEnvironmentsRepository) IsGoldenClusterFile(file files.YamlFile) bool {
	return file.Data["golden_cluster"] != nil && file.Data["golden_cluster"].(bool)
}

// IsGoldenCluster returns true if the cluster is a golden cluster
func (r *ObltEnvironmentsRepository) IsGoldenCluster(clusterName string) (ret bool) {
	ret = false
	goldenClusters := r.FindGoldenClustersNames()
	for _, goldenCluster := range goldenClusters {
		if strings.EqualFold(clusterName, goldenCluster) {
			ret = true
			break
		}
	}
	return ret
}

// FindClusterConfig returns the YAML configuration for a cluster
func (r *ObltEnvironmentsRepository) FindClusterConfig(clusterName string) (file files.YamlFile, err error) {
	logger.Debugf("Reading cluster config %s", clusterName)
	files, err := files.Find(r.getEnvironments(), files.ConfigFilesFilter, files.FilterByYamlPath("cluster_name", clusterName))
	if len(files) > 0 {
		return files[0], err
	}

	return file, fmt.Errorf("cluster config for 'cluster_name=%s' not found", clusterName)
}

// FindTemplate returns the YAML configuration for a template
func (r *ObltEnvironmentsRepository) FindTemplate(templateName string) (file files.YamlFile, err error) {
	logger.Debugf("Searching template : %s", templateName)
	files, _ := files.Find(r.GetUsersDir(), files.TemplatesFilter, files.FilterByYamlPath("template_name", templateName))
	if len(files) > 0 {
		file = files[0]
	} else {
		err = fmt.Errorf("template for 'template_name=%s' not found", templateName)
	}
	return file, err
}

// FindCurrentUserClusterConfig returns the YAML configuration for a cluster of the current user
func (r *ObltEnvironmentsRepository) FindCurrentUserClusterConfig(clusterName string) (file files.YamlFile, err error) {
	logger.Debugf("Reading current user cluster config %s", clusterName)
	files, _ := files.Find(r.GetCurrentUserEnvironments(), files.ConfigFilesFilter, files.FilterByYamlPath("cluster_name", clusterName))
	if len(files) > 0 {
		file = files[0]
	} else {
		err = fmt.Errorf("current user's cluster config for 'cluster_name=%s' not found", clusterName)
	}
	return file, err
}

// FindCurrentUserClusterConfig returns the YAML configuration for a cluster of the current user
func (r *ObltEnvironmentsRepository) FindClustersByRemote(all bool, remoteCluster string) ([]files.YamlFile, error) {
	clustersDir := r.GetCurrentUserEnvironments()
	if all {
		clustersDir = r.getEnvironments()
	}

	logger.Debugf("Cluster Configurations Available:")
	return files.Find(clustersDir, files.ConfigFilesFilter, files.FilterByYamlPath("elasticsearch.ccs_remote_cluster", remoteCluster))
}

// FindGoldenClusters returns the YAML configuration for the golden clusters
func (r *ObltEnvironmentsRepository) FindGoldenClusters() ([]files.YamlFile, error) {
	logger.Debugf("Searching for Golden Clusters")
	return r.FindClusterByYamlPath(map[string]string{"golden_cluster": "true"})
}

// FindClusterByYamlPath returns the YAML configuration for the clusters that match the yamlPath and yamlValue
func (r *ObltEnvironmentsRepository) FindClusterByYamlPath(filters map[string]string) ([]files.YamlFile, error) {
	var filtersFunc []files.Filter
	for key, value := range filters {
		logger.Debugf("Searching for clusters with %s=%s", key, value)
		filtersFunc = append(filtersFunc, files.FilterByYamlPath(key, value))
	}
	files, err := files.Find(r.getEnvironments(), files.ConfigFilesFilter, filtersFunc...)
	return files, err
}

// FindGoldenClustersNames returns the names of the golden clusters
func (r *ObltEnvironmentsRepository) FindGoldenClustersNames() []string {
	var goldenClustersNames []string
	goldenClusters, _ := r.FindGoldenClusters()
	for _, goldenCluster := range goldenClusters {
		goldenClustersNames = append(goldenClustersNames, goldenCluster.Data["cluster_name"].(string))
	}
	return goldenClustersNames
}

// GetCurrentUserEnvironments returns the absolute path to the current user's environments folder.
func (r *ObltEnvironmentsRepository) GetCurrentUserEnvironments() string {
	return filepath.Join(r.GetUsersDir(), r.config.Username)
}

// getEnvironments returns the absolute path to the environments folder.
func (r *ObltEnvironmentsRepository) getEnvironments() string {
	return filepath.Join(r.gitRepo.Path(), environmentsPath)
}

// GetTemplate retrieves the path to the CCS template file
func (r *ObltEnvironmentsRepository) GetTemplate(ccsTemplateFile string) string {
	return filepath.Join(r.GetUsersDir(), ccsTemplateFile)
}

// GetUserClusterConfig retrieves the path to the cluster config file for a user
func (r *ObltEnvironmentsRepository) GetUserClusterConfig(clusterCfgFilename string) string {
	return filepath.Join(r.GetCurrentUserEnvironments(), clusterCfgFilename)
}

// GetUsersDir retrieves the path to the users dir
func (r *ObltEnvironmentsRepository) GetUsersDir() string {
	return filepath.Join(r.getEnvironments(), "users")
}

// GetBootstrapRecipesDir retrieves the path to the Bootstrap recipes dir
func (r *ObltEnvironmentsRepository) GetBootstrapRecipesDir() string {
	return filepath.Join(r.gitRepo.Path(), r.GetBootstrapRecipesRelativeDir())
}

// GetBootstrapRecipesRelativeDir retrieves the relative path to the Bootstrap recipes dir
func (r *ObltEnvironmentsRepository) GetBootstrapRecipesRelativeDir() string {
	return filepath.Join("ansible", "ansible_collections", "oblt", "framework", "roles", "common", "files", "deployments", "bootstrap")
}

// GetPath retrieves the path to the repository
func (r *ObltEnvironmentsRepository) GetPath() string {
	return r.gitRepo.Path()
}

// ListClusters retrieves the list of clusters, being possible retrieve all clusters or current user's ones
func (r *ObltEnvironmentsRepository) ListClusters(all bool) []files.YamlFile {
	clustersDir := r.GetCurrentUserEnvironments()
	if all {
		clustersDir = r.getEnvironments()
	}

	logger.Debugf("Cluster Configurations Available:")
	clusters, _ := files.Find(clustersDir, files.ConfigFilesFilter)
	if !all {
		clusterByUsername, _ := r.FindClusterByYamlPath(map[string]string{"oblt_username": r.config.Username})
		clusters = append(clusters, clusterByUsername...)
		// check for Slack Channel if it is a user
		if len(r.config.SlackChannel) > 0 && strings.HasPrefix(r.config.SlackChannel, "@") {
			clustersBySlack, _ := r.FindClusterByYamlPath(map[string]string{"slack_channel": r.config.SlackChannel})
			clusters = append(clusters, clustersBySlack...)
		}
		data := make(map[string]files.YamlFile)
		for _, item := range clusters {
			data[item.Data["cluster_name"].(string)] = item
		}
		clusters = maps.Values(data)
	}

	if !all {
		clusters = filterUserConfigFiles(clusters)
	}
	return clusters
}

func filterUserConfigFiles(clusters []files.YamlFile) []files.YamlFile {
	var filteredClusters []files.YamlFile
	for _, cluster := range clusters {
		if obltManaged, ok := cluster.Data["oblt_managed"].(bool); !ok || obltManaged {
			filteredClusters = append(filteredClusters, cluster)
		}
	}
	clusters = filteredClusters
	return clusters
}

// ListTemplates retrieves the list of existing templates
func (r *ObltEnvironmentsRepository) ListTemplates() []files.YamlFile {
	logger.Debugf("Cluster Templates Available:")
	templates, _ := files.Find(r.GetUsersDir(), files.TemplatesFilter)
	return templates
}

// ListBootstrapRecipes retrieves the list of existing Bootstrap recipes
func (r *ObltEnvironmentsRepository) ListBootstrapRecipes(category string) []files.YamlFile {
	return r.ListBootstrapRecipesInFolder(category, r.GetBootstrapRecipesDir())
}

// ListBootstrapRecipesInFolder retrieves the list of existing Bootstrap recipes from a specific folder
func (r *ObltEnvironmentsRepository) ListBootstrapRecipesInFolder(category string, bootstrapFolder string) []files.YamlFile {
	logger.Debugf("Cluster Bootstrap recipes Available at: %s", bootstrapFolder)
	recipes, _ := files.Find(filepath.Join(bootstrapFolder, category), files.ConfigFilesFilter)
	return recipes
}

// LoadRecipesFromJson load the recipes declared in the JSON ["recipe01", "recipe01", "recipe01"]
// where category01 is the name of the folder and recipe01 the name of the file without extension.
func (r *ObltEnvironmentsRepository) LoadRecipesFromJson(category string, recipesJson string) []files.YamlFile {
	recipes, _ := r.LoadRecipesFromJsonInFolder(category, recipesJson, r.GetBootstrapRecipesDir())
	return recipes
}

// LoadRecipesFromJsonInFolder load the recipes declared in the JSON ["recipe01", "recipe01", "recipe01"]
// where category01 is the name of the folder and recipe01 the name of the file without extension
// from a specific folder.
func (r *ObltEnvironmentsRepository) LoadRecipesFromJsonInFolder(category string, recipesJson string, bootstrapFolder string) ([]files.YamlFile, error) {
	var recipesList []files.YamlFile
	logger.Debugf("Using recipes %s", recipesJson)
	var recipes []string
	err := json.Unmarshal([]byte(recipesJson), &recipes)
	for _, recipe := range recipes {
		recipeFile := recipe + ".yml"
		filePath := filepath.Join(bootstrapFolder, category, recipeFile)
		logger.Debugf("Loading recipe %s", filePath)
		d, b, err := files.ReadYamlFile(filePath)
		if err == nil {
			yamlFile := files.YamlFile{
				Data: d, Owner: category, Path: filePath, Bytes: b,
			}
			recipesList = append(recipesList, yamlFile)
		}
	}
	return recipesList, err
}

// WaitForClusterCreation waits for the cluster to be created, it monitoring the pull request status until it is merged.
func (r *ObltEnvironmentsRepository) WaitForClusterCreation(clusterName string, timeout time.Duration) {
	repo := r.Repo(false)
	expectedState := "closed"
	owner := repo.Owner
	repoName := repo.Name

	pr, err := github.FindPullRequestRetry(clusterName, owner, repoName, timeout)
	cobra.CheckErr(err)

	_, _, err = github.WaitPullRequestStateRetry(pr.GetNumber(), owner, repoName, expectedState, timeout)
	cobra.CheckErr(err)
}

// getTemplateFileByNameOrPath returns the template file by name or path
func (r *ObltEnvironmentsRepository) getTemplateFileByNameOrPath(templateName string, templateFilePath string) (yamlFile files.YamlFile, err error) {
	var templateFile files.YamlFile
	if len(templateName) > 0 {
		yamlFile, err = r.FindTemplate(templateName)
		templateFile = yamlFile
	} else {
		logger.Infof("Reading template : %s", templateFilePath)
		d, b, err := files.ReadYamlFile(templateFilePath)
		if err == nil {
			yamlFile = files.YamlFile{
				Data:  d,
				Bytes: b,
				Owner: r.config.Username,
				Path:  templateFilePath,
			}
			templateFile = yamlFile
		}
	}
	return templateFile, err
}

// Wipeup destroys all clusters of a user
func (r *ObltEnvironmentsRepository) Wipeup() (ret []map[string]interface{}, errs error) {
	files := r.ListClusters(false)

	for _, file := range files {
		if file.Data["cluster_name"] != nil {
			// do not clean up if dry run is enabled, tests will fail
			if !r.gitRepo.IsDryRun() {
				r.gitRepo.AlwaysCleanUp = true
				r.gitRepo.Sync()
			}
			clusterName := file.Data["cluster_name"].(string)
			results, err := r.DestroyCluster(clusterName, false)
			ret = append(ret, results)
			errs = errors.Join(errs, err)
		}
	}
	return ret, errs
}

// UpdateCluster updates a cluster using some parameters
func (r *ObltEnvironmentsRepository) UpdateCluster(clusterName string, parameters string) (clusterConfig map[string]interface{}, err error) {
	logger.Debugf("Updating cluster config %s: %s", clusterName, parameters)
	var path string
	if clusterConfig, path, err = r.ReadClusterConfig(clusterName); err == nil {
		var clusterConfigNew map[string]interface{}
		if err = json.Unmarshal([]byte(parameters), &clusterConfigNew); err == nil {
			if clusterConfigNew["cluster_name"] != nil && clusterConfigNew["cluster_name"].(string) != clusterName {
				err = fmt.Errorf("cluster_name %s is different from %s", clusterConfigNew["cluster_name"].(string), clusterName)
			}
			if clusterConfig1, err := Merge(clusterConfigNew, clusterConfig); err == nil {
				if err = r.WriteClusterConfig(clusterConfig1, path); err == nil {
					commitMessage := "oblt-cli(" + r.config.Username + "): Updating cluster " + clusterName
					r.gitRepo.CommitAndPush(commitMessage)
				}
				clusterConfig = clusterConfig1.(map[string]interface{})
			}
		}
	}
	if err != nil {
		logger.Infof("Error updating cluster config %s: %s", clusterName, err)
	} else {
		logger.Infof("Cluster config updated %s: %s", clusterName, path)
	}
	return clusterConfig, err
}

// ReadClusterConfig reads the cluster configuration from a file, the clusters is searched by name.
func (r *ObltEnvironmentsRepository) ReadClusterConfig(clusterName string) (clusterConfig map[string]interface{}, path string, err error) {
	var clusterConfigFiles []files.YamlFile
	clusterConfig = make(map[string]interface{})
	if clusterConfigFiles, err = r.FindClusterByYamlPath(map[string]string{"cluster_name": clusterName}); err == nil && len(clusterConfigFiles) > 0 {
		clusterConfigFile := clusterConfigFiles[0]
		path = clusterConfigFile.Path
		logger.Debugf("Reading cluster config %s: %s", clusterName, path)
		err = yaml.Unmarshal(clusterConfigFile.Bytes, &clusterConfig)
	}
	return clusterConfig, path, err
}

// WriteClusterConfig writes the cluster configuration to a file.
func (r *ObltEnvironmentsRepository) WriteClusterConfig(clusterConfig interface{}, path string) (err error) {
	logger.Debugf("Writing cluster config %s", path)
	var data []byte
	if data, err = yaml.Marshal(clusterConfig); err == nil {
		err = os.WriteFile(path, data, 0644)
	}
	return err
}

// https://go.dev/play/p/8jlJUbEJKf
// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
func Merge(patch, org interface{}) (interface{}, error) {
	data1, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}
	var jpatch interface{}
	err = json.Unmarshal(data1, &jpatch)
	if err != nil {
		return nil, err
	}
	var jorg interface{}
	err = json.Unmarshal(data2, &jorg)
	if err != nil {
		return nil, err
	}
	return merge1(jpatch, jorg), nil
}

func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = merge1(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x1
}

// GetVersion returns the .ci/.version file content
func (r *ObltEnvironmentsRepository) GetVersion() (repoVersion semver.Version, err error) {
	versionFile := filepath.Join(r.GetPath(), ".ci", ".version")
	version, err := os.ReadFile(versionFile)
	if err == nil {
		repoVersion = semver.MustParse(strings.ReplaceAll(strings.ReplaceAll(string(version), "\n", ""), "\r", ""))
	}
	return repoVersion, err
}
