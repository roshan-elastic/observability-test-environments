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
	"fmt"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// TemplatesCmd represents the template command
var TemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Command to list the cluster configuration templates.",
	Long:  `Command to list the cluster configuration templates and their descriptions.`,
	Run:   runTemplate,
}

func init() {
	ClusterCmd.AddCommand(TemplatesCmd)
}

/*
It will checkout the observability test environments repository,
and check all the files in the environments folder looking for templates.
*/
func runTemplate(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runTemplate", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	fs := obltTestEnvironments.ListTemplates()
	var results = make(map[string]interface{})
	var items []interface{}
	for _, file := range fs {
		templateName := file.Data["template_name"].(string)
		templateFilePath := file.Path
		templateDescription := file.Data["template_description"].(string)
		items = append(items, map[string]string{"name": templateName, "path": templateFilePath, "description": templateDescription})
		fmt.Printf("%s %s\n%s\n", logger.InfoColor.Sprint(fmt.Sprintf("📦 [%s] ->", templateName)), logger.WarnColor.Sprint(templateFilePath), templateDescription)
	}
	results["items"] = items
	saveResults(results, outputFile)
}
