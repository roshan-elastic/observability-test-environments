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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	osExec "os/exec"
	"path/filepath"
	"time"

	"github.com/blang/semver"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
	apm_v2 "go.elastic.co/apm/v2"
)

const (
	// isoDateFormat flag name.
	isoDateFormat = "2006-01-02T15:04:05-0700"
)

// Exec it executes a command show the output, it exits on a failure.
func Exec(cmd string) {
	shellCmd := osExec.Command(cmd)
	stdout, err := shellCmd.Output()
	config.CheckErr(err)
	logger.Infof("%s: %s", cmd, stdout)
}

// CopyFile it copies the content of a file to another file, if there is an error the program exits.
func CopyFile(srcPath string, dstPath string) {
	logger.Debugf("copy file %s to %s", srcPath, dstPath)
	input, err := os.ReadFile(srcPath)
	config.CheckErr(err)
	err = os.WriteFile(dstPath, input, 0644)
	config.CheckErr(err)
}

// nowISOFormat return the current timestamp in ISO-8601 format
func nowISOFormat() (date string) {
	return time.Now().Format(isoDateFormat)
}

// checkOpperationLock check if it is possible to operate over the cluster.
// to avoid issues there is a block of 30 seconds between operations.
func checkOpperationLock(userConfig config.ObltConfiguration) {
	if dryRun || os.Getenv("CI") == "true" {
		return
	}
	cfgDir := userConfig.GetDir()
	lockFile := filepath.Join(cfgDir, "lock")
	now := time.Now()
	lastOperation, _ := os.ReadFile(lockFile)
	if len(lastOperation) > 0 {
		timeLastOp, err := time.Parse(isoDateFormat, string(lastOperation))
		if err == nil && now.Before(timeLastOp.Add(time.Second*30)) {
			log.Fatalf("operations are locked until %s", timeLastOp.Add(time.Second*30).Format(isoDateFormat))
		}
	}
	err := os.WriteFile(lockFile, []byte(nowISOFormat()), 0644)
	config.CheckErr(err)
}

// openBrowser opens an URL using the system browser
// TODO: commented for future use right now is not used
// func openBrowser(url string) {
// 	logger.Debugf("browser.open: %s", url)

// 	var err error

// 	switch runtime.GOOS {
// 	case "linux":
// 		err = osExec.Command("xdg-open", url).Start()
// 	case "windows":
// 		err = osExec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 	case "darwin":
// 		err = osExec.Command("open", url).Start()
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}

// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// }

// saveResults It saves the reults of the operation in a JSON file.
func saveResults(data interface{}, filePath string) {
	err := files.SaveResults(data, filePath)
	config.CheckErr(err)
}

// saveResultsRaw It saves the reults of the operation in a Raw file.
func saveResultsRaw(data string, filePath string) {
	err := files.SaveResultsRaw(data, filePath)
	config.CheckErr(err)
}

// loadParameters loads the parameters values from the parameters flag or the parameters file flag
func loadParameters(cmd *cobra.Command) (currentParameters string, err error) {
	parameters, _ := cmd.Flags().GetString(config.ParametersFlag)
	parametersFile, _ := cmd.Flags().GetString(config.ParametersFileFlag)
	currentParameters = "{}"

	if len(parametersFile) > 0 && len(parameters) > 0 {
		err = fmt.Errorf("%s and %s cannot be set at the same time", config.ParametersFlag, config.ParametersFileFlag)
	} else {
		if len(parametersFile) > 0 {
			var fileBytes []byte
			fileBytes, err = os.ReadFile(parametersFile)
			if err == nil {
				currentParameters = string(fileBytes)
			}
		}
		if len(parameters) > 0 {
			currentParameters = parameters
		}
	}

	return currentParameters, err
}

// loadParametersMap loads the parameters values from the parameters flag or the parameters file flag in to a map[string]interface{}
func loadParametersMap(cmd *cobra.Command) (parameters map[string]interface{}, err error) {
	parametersJson, _ := loadParameters(cmd)
	err = json.Unmarshal([]byte(parametersJson), &parameters)
	return parameters, err
}

// returns the first not empty string of certain strings
func firstNotEmpty(values ...string) (ret string) {
	for _, value := range values {
		if value != "" {
			ret = value
			break
		}
	}
	return ret
}

// showBanner shows the banner if it is not disabled
func showBanner(obltRepo clusters.Repository) (msgs []clusters.Message) {
	if !disableBanner {
		msgs = clusters.ShowBanner(obltRepo.GetPath())
	}
	return msgs
}

// verifyCompatibility checks if the repository is compatible with the current version of the tool
func verifyCompatibility(repoVersion semver.Version) (err error) {
	if repoVersion.GT(currentVersion) {
		logger.Warnf("a new version of the oblt-cli tool has been released. %s > %s", repoVersion, currentVersion)
		logger.Warnf("please update oblt-cli https://elastic.github.io/observability-test-environments/tools/oblt-cli/update/")
	}
	if repoVersion.Major != currentVersion.Major {
		err = fmt.Errorf("the version of the repository is not compatible with the current version of the tool")
	}
	return err
}

// initialChecks performs the initial checks before running a command
func initialChecks(obltRepo clusters.Repository, tx *apm_v2.Transaction, ctx context.Context) {
	showBanner(obltRepo)
	version, err := obltRepo.GetVersion()
	apm.CobraCheckErr(err, tx, ctx)
	err = verifyCompatibility(version)
	apm.CobraCheckErr(err, tx, ctx)
}
