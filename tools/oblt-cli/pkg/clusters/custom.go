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
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/template"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

// CustomCluster is the struct that contains the parameters for the Custom cluster
type CustomCluster struct {
	ClusterName       string                      `json:"ClusterName,omitempty"`
	ClusterNamePrefix string                      `json:"ClusterNamePrefix,omitempty"`
	ClusterNameSuffix string                      `json:"ClusterNameSuffix,omitempty"`
	GitHubCommentId   string                      `json:"GitHubCommentId,omitempty"`
	GitHubCommit      string                      `json:"GitHubCommit,omitempty"`
	GitHubIssue       string                      `json:"GitHubIssue,omitempty"`
	GitHubPullRequest string                      `json:"GitHubPullRequest,omitempty"`
	GitHubRepository  string                      `json:"GitHubRepository,omitempty"`
	ObltRepo          *ObltEnvironmentsRepository `json:"-"`
	Parameters        string                      `json:"Parameters,omitempty"`
	SlackChannel      string                      `json:"SlackChannel,omitempty"`
	TemplateName      string                      `json:"TemplateName,omitempty"`
	TemplatePath      string                      `json:"TemplatePath,omitempty"`
	Username          string                      `json:"Username,omitempty"`
}

// CustomTemplateParams is the struct that contains the parameters for the Custom template
type CustomTemplateParams struct {
	ClusterConfigFile string `json:"ClusterConfigFile,omitempty"`
	ClusterName       string `json:"ClusterName,omitempty"`
	CommitMessage     string `json:"CommitMessage,omitempty"`
	CommitSha         string `json:"CommitSha,omitempty"`
	CommitURL         string `json:"CommitURL,omitempty"`
	Date              string `json:"Date,omitempty"`
	GitHubCommentId   string `json:"GitHubCommentId,omitempty"`
	GitHubCommit      string `json:"GitHubCommit,omitempty"`
	GitHubIssue       string `json:"GitHubIssue,omitempty"`
	GitHubPullRequest string `json:"GitHubPullRequest,omitempty"`
	GitHubRepository  string `json:"GitHubRepository,omitempty"`
	GitOps            bool   `json:"GitOps"`
	SlackChannel      string `json:"SlackChannel,omitempty"`
	Username          string `json:"Username,omitempty"`
	TemplatePath      string `json:"TemplatePath,omitempty"`
	TemplateName      string `json:"TemplateName,omitempty"`
}

// Create creates the custom cluster
func (r *CustomCluster) Create() (parameters map[string]interface{}, err error) {
	var clusterName, clusterCfg, json string
	if err = r.validate(); err == nil {
		if err = r.loadTemplateInfo(r.ObltRepo.getTemplateFileByNameOrPath(r.TemplateName, r.TemplatePath)); err == nil {
			if clusterName, clusterCfg, err = r.createClusterFile(); err == nil {
				if json, err = r.packTemplateParams(clusterName, clusterCfg); err == nil {
					parameters, err = r.createCluster(clusterName, clusterCfg, json)
				}
			}
		}
	}
	if err == nil {
		logger.Infof("For more details about the workflow check https://elastic.github.io/observability-test-environments/tools/oblt-framework/ci/#oblt-cli-workflow")
		logger.Infof("Cluster creation is requested : %s", clusterName)
		logger.Infof("Cluster creation will start at : https://github.com/elastic/observability-test-environments/pulls?q=is:pr+%s+label:oblt-cli-cluster+is:open", clusterName)
		logger.Infof("You will receive updates about the creation process in Slack")
	} else {
		logger.LogError("Error creating cluster", err)
	}
	return parameters, err
}

// loadTemplateInfo loads the template information from the info returned by ObltEnvironmentsRepository.getTemplateFileByNameOrPath
func (r *CustomCluster) loadTemplateInfo(templateFile files.YamlFile, err error) error {
	if err == nil {
		logger.Infof("Using template : %s", templateFile.Path)
		r.TemplateName = templateFile.Data["template_name"].(string)
		r.TemplatePath = templateFile.Path
	}
	return err
}

// createClusterFile creates the cluster configuration file, returns the cluster name and the cluster configuration file path
func (r *CustomCluster) createClusterFile() (clusterName, clusterCfg string, err error) {
	clusterName = r.ClusterName

	if r.ClusterName == "" {
		prefix := joinIgnoreEmpty([]string{r.ClusterNamePrefix, r.TemplateName}, itemsSeparator)
		clusterName, err = newClusterName(prefix, r.ClusterNameSuffix, config.Seed(5))
	}

	if err == nil {
		clusterCfgFilename := clusterName + ".yml"
		clusterCfg = r.ObltRepo.GetUserClusterConfig(clusterCfgFilename)
		err = os.MkdirAll(filepath.Dir(clusterCfg), 0700)
	}
	return clusterName, clusterCfg, err
}

// createCluster creates the cluster configuration file
func (r *CustomCluster) createCluster(clusterName, clusterCfg, jsonStr string) (map[string]interface{}, error) {
	parametersMap := make(map[string]interface{})
	var sha string
	err := json.Unmarshal([]byte(jsonStr), &parametersMap)
	if err == nil {
		err = json.Unmarshal([]byte(r.Parameters), &parametersMap)
	}
	if err == nil {
		logger.Infof("Writing cluster configuration file %s", clusterCfg)
		err = template.Parse(parametersMap, r.TemplatePath, clusterCfg)
		if err == nil && template.CheckAllVariablesDefined(clusterCfg) {
			commitMessage := "oblt-cli(" + r.Username + "): Create " + clusterName + " cluster [template=" + r.TemplateName + "]"
			sha, err = r.ObltRepo.gitRepo.CommitAndPush(commitMessage)

			parametersMap["CommitSha"] = sha
			parametersMap["CommitMessage"] = commitMessage
			parametersMap["CommitURL"] = fmt.Sprintf("https://github.com/elastic/observability-test-environments/commit/%s", sha)
		} else {
			err1 := fmt.Errorf("there are parameters undefined, please check the template %s -> %s", r.TemplateName, r.TemplatePath)
			err = errors.Join(err, err1)
		}
	}
	return parametersMap, err
}

// packTemplateParams packs the template parameters in a JSON string
func (r *CustomCluster) packTemplateParams(clusterName, clusterCfg string) (jsonContent string, err error) {
	var jsonContentb []byte

	if err == nil {
		templateParams := CustomTemplateParams{
			ClusterName:       clusterName,
			ClusterConfigFile: clusterCfg,
			TemplatePath:      r.TemplatePath,
			TemplateName:      r.TemplateName,
			Date:              nowISOFormat(),
			GitHubCommentId:   r.GitHubCommentId,
			GitHubCommit:      r.GitHubCommit,
			GitHubIssue:       r.GitHubIssue,
			GitHubPullRequest: r.GitHubPullRequest,
			GitHubRepository:  r.GitHubRepository,
			GitOps:            len(r.GitHubCommentId) > 0 || len(r.GitHubCommit) > 0 || len(r.GitHubIssue) > 0 || len(r.GitHubPullRequest) > 0 || len(r.GitHubRepository) > 0,
			SlackChannel:      r.SlackChannel,
			Username:          r.Username,
		}

		jsonContentb, err = json.Marshal(templateParams)
		if err == nil {
			jsonContent = string(jsonContentb)
		}
	}
	return jsonContent, err
}

// validate validates the parameters
func (r *CustomCluster) validate() (err error) {
	return errors.Join(
		config.ValidateClusterName(r.ClusterName),
		config.ValidatePrefix(r.ClusterNamePrefix),
		config.ValidateSuffix(r.ClusterNameSuffix),
		config.ValidateSlackChannel(r.SlackChannel),
	)
}
