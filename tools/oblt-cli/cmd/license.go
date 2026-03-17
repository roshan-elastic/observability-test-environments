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
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// LicenseCmd List the recipes available.
var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Command to push a license to an Elasticsearch.",
	Long:  "Command to push a license to an Elasticsearch.",
	Run:   runLicense,
}

func init() {
	RootCmd.AddCommand(LicenseCmd)
	LicenseCmd.Flags().String(config.TypeFlag, "", "Type of license to push [release, dev, orchestration, orchestration-dev]. (Required)")
	LicenseCmd.Flags().String(config.UrlFlag, "", "URL of the Elastisearch service to bootstrap. (Required)")
	LicenseCmd.Flags().String(config.UsernameFlag, "", "Username to authenticate. (Not needed if you provide "+config.ApiKeyFlag+")")
	LicenseCmd.Flags().String(config.PasswordFlag, "", "Password to authenticate. (Not needed if you provide "+config.ApiKeyFlag+")")
	LicenseCmd.Flags().String(config.ApiKeyFlag, "", "API Key to authenticate. (Not needed if you provide "+config.UsernameFlag+")")
	LicenseCmd.Flags().Bool(config.IgnoreCertificatesFlag, false, "Disable TLS certificate verification.")

	cobra.MarkFlagRequired(LicenseCmd.Flags(), config.TypeFlag)
	cobra.MarkFlagRequired(LicenseCmd.Flags(), config.UrlFlag)
}

func runLicense(cmd *cobra.Command, args []string) {
	licenseType, _ := cmd.Flags().GetString(config.TypeFlag)
	esUrl, _ := cmd.Flags().GetString(config.UrlFlag)
	username, _ := cmd.Flags().GetString(config.UsernameFlag)
	password, _ := cmd.Flags().GetString(config.PasswordFlag)
	apiKey, _ := cmd.Flags().GetString(config.ApiKeyFlag)
	ignoreCertificates, _ := cmd.Flags().GetBool(config.IgnoreCertificatesFlag)

	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	runLicenseCommon(userConfig, "runLicense", licenseType, username, password, apiKey, esUrl, ignoreCertificates)
}

// runLicenseCommon is a common function to push the license.
func runLicenseCommon(userConfig config.ObltConfiguration, apmRequest, licenseType string, username string, password string, apiKey string, esUrl string, ignoreCertificates bool) {
	tx, ctx := apm.StartTransaction(apmRequest, "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	err := config.ValidateLicenseType(licenseType)
	apm.CobraCheckErr(err, tx, ctx)

	authType, params, err := http.ChooseAuthType(username, password, apiKey, "")
	apm.CobraCheckErr(err, tx, ctx)

	gcsm, err := gcp.NewClusterSecrets()
	apm.CobraCheckErr(err, tx, ctx)
	license, err := gcsm.ReadLicense(licenseType)
	apm.CobraCheckErr(err, tx, ctx)

	request := http.HttpRequest{
		Url:                esUrl + "/_license",
		IgnoreCertificates: ignoreCertificates,
		Method:             "PUT",
		AuthType:           authType,
		Headers: map[string]interface{}{
			"Content-Type":         "application/json",
			"X-Management-Request": true,
		},
		Body:       license,
		DryRun:     dryRun,
		AuthParams: params,
	}
	if !dryRun {
		_, err = request.DoHttpRequest()
		apm.CobraCheckErr(err, tx, ctx)
	} else {
		logger.Infof("Dry run mode, no request sent")
		logger.Debugf("Request: %v", request)
	}
	logger.Infof("License %s pushed successfully", licenseType)
}
