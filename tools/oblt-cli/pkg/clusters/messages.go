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
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// Message struct representing a message in the banner. The content of the banner must follow a well-known format,
// containing a list of message elements. Each message should contain two fields: 'msg' and 'level'.
// An example of this struct:
// - msg: edge-oblt is under maintenance. It would probably not respond properly
// level: INFO
// - msg: dev-oblt cluster is down
// level: WARM
// - msg: A new version of the oblt-cli tool has been released. Please upgrade it to get the latest features
// level: INFO
type Message struct {
	Level string `yaml:"level"`
	Text  string `yaml:"msg"`
}

// coloriseLevel returns the level using a color scheme: info in Cyan, warning in yellow, and error in red
func coloriseLevel(level string) string {
	upperCaseLevel := strings.ToUpper(level)
	switch upperCaseLevel {
	case "INFO":
		return logger.InfoColor.Sprint(level)
	case "WARN", "WARNING":
		return logger.WarnColor.Sprint(level)
	case "ERROR":
		return logger.ErrColor.Sprint(level)
	}

	return level
}

// ShowBanner It reads the "banners.yml" file in the current directory and returns an array of Message
func ShowBanner(repositoryDir string) []Message {
	bannersPath := filepath.Join(repositoryDir, "banners.yml")
	logger.Debugf("Checking banners file %s", bannersPath)

	bannersContent, err := os.ReadFile(bannersPath)
	if err != nil {
		logger.Debugf("No banners found: %s", err)
		return []Message{}
	}

	statuses := []Message{}
	err = yaml.Unmarshal(bannersContent, &statuses)
	cobra.CheckErr(err)

	if len(statuses) > 0 {
		data := [][]string{}
		for _, message := range statuses {
			data = append(data, []string{coloriseLevel(message.Level), message.Text})
		}

		table := tablewriter.NewWriter(os.Stderr)
		table.SetHeader([]string{"Level", "Message"})
		table.SetAutoWrapText(false)
		table.AppendBulk(data)
		table.Render()
	}

	return statuses
}
