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

package console

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/maps"
	"github.com/olekukonko/tablewriter"
)

// PrintYamlFiles print a list of YAML files on the console.
func PrintYamlFiles(title string, items []files.YamlFile) [][]string {
	data := [][]string{}
	titleColor := logger.WarnColor.Sprint(title)
	fmt.Printf("\n--%s--\n", logger.WarnColor.Sprint(titleColor))

	for _, item := range items {
		row := []string{}
		fileName := filepath.Base(item.Path)
		name := logger.InfoColor.Sprint(strings.TrimSuffix(fileName, path.Ext(fileName)))
		filePath := logger.WarnColor.Sprint(item.Path)

		// extract the object from the data
		obj := maps.ConvertToMap(item.Data)
		description := ""
		if obj["description"] != nil {
			description = obj["description"].(string)
		}
		row = append(row, name, titleColor, description, filePath)
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Name", "Type", "Description", "File path"})
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return data
}
