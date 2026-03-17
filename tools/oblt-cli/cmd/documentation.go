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
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// DocumentationCmd generates the Markdown documentation for oblt-cli.
var DocumentationCmd = &cobra.Command{
	Use:   "documentation",
	Short: "Command to generate the Markdown documentation for oblt-cli.",
	Long: `Command to generate the Markdown documentation for oblt-cli.
	It generates a Markdown document for each command available in oblt-cli with all the options available.`,
	Run: runDocumentation,
}

func init() {
	RootCmd.AddCommand(DocumentationCmd)

	DocumentationCmd.Flags().String(config.OutputFolderFlag, "", "Folder where the documentation will be generated. (Required)")
	cobra.MarkFlagRequired(DocumentationCmd.Flags(), config.OutputFolderFlag)
}

func runDocumentation(cmd *cobra.Command, args []string) {
	outputFolder, _ := cmd.Flags().GetString(config.OutputFolderFlag)
	err := os.MkdirAll(outputFolder, 0777)
	if err != nil {
		logger.Fatal(err)
	}
	err = doc.GenMarkdownTree(RootCmd, outputFolder)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Documentation generated in %s", outputFolder)
}
