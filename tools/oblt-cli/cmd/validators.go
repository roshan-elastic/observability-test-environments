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
	"errors"
	"fmt"
	"path"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/spf13/cobra"
)

func validateCreateCustomFlags(cmd *cobra.Command, args []string) error {
	templateName, _ := cmd.Flags().GetString(config.TemplateNameFlag)
	templateFilePath, _ := cmd.Flags().GetString(config.TemplateFileFlag)

	if len(templateName) == 0 && len(templateFilePath) == 0 {
		return fmt.Errorf("required flag(s) \"%s\" or \"%s\" not set", config.TemplateNameFlag, config.TemplateFileFlag)
	}

	if len(templateName) > 0 && len(templateFilePath) > 0 {
		return fmt.Errorf("flags conflict, you have to define only one of \"%s\" or \"%s\"", config.TemplateNameFlag, config.TemplateFileFlag)
	}

	if len(templateFilePath) > 0 && !path.IsAbs(templateFilePath)   {
		return fmt.Errorf("\"%s\" requires an absolute path", config.TemplateFileFlag)
	}

	return nil
}

func validateDeprecatedCCSFlags(cmd *cobra.Command, args []string) error {
	deprecatedClusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	if deprecatedClusterName != "" {
		return fmt.Errorf("--%s flag is deprecated. Use --%s instead, or check 'oblt-cli cluster create ccs --help' from more details", config.ClusterNameFlag, config.RemoteClusterFlag)
	}

	return nil
}

// validateCIMinimumArguments checks that the minimum arguments are provided when running in CI mode
func validateCIMinimumArguments(cmd *cobra.Command, args []string) (err error) {
	clusterNamePrefix, _ := cmd.Flags().GetString(config.ClusterNamePrefixFlag)
	clusterNameSuffix, _ := cmd.Flags().GetString(config.ClusterNameSuffixFlag)
	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	if ciMode {
		if clusterName == "" && clusterNamePrefix == "" && clusterNameSuffix == "" {
			err = errors.New("must provide `--cluster-name-prefix` or `--cluster-name-prefix` or `--cluster-name` when running in the CI")
		}
		if clusterName != "" && (clusterNamePrefix != "" || clusterNameSuffix != "") {
			err = errors.New("`--cluster-name` is not supported with `--cluster-name-prefix` or `--cluster-name-suffix` when running in the CI")
		}

	} else {
		// Avoid using `--cluster-name` by default but only when running in the CI
		// That's the edge case we want to have when reusing deployments created automatically.
		if clusterName != "" {
			err = errors.New("`--cluster-name` is not supported when running in the CI")
		}
	}
	return err
}

// validateEnvironmentFlag checks that the environment flag is valid
func validateEnvironmentFlag(cmd *cobra.Command, args []string) error {
	environment, _ := cmd.Flags().GetString(config.EnvironmentFlag)
	return config.ValidateTarget(environment)
}
