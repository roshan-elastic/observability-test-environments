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

// Package artifacts this package contains the interactions with Elastic's artifacts-api
package artifacts

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
)

var artifactsApiUrl = "https://artifacts-api.elastic.co/v1"

type manifests struct {
	LastUpdateTime         string `json:"last-update-time"`
	SecondsSinceLastUpdate int    `json:"seconds-since-last-update"`
}

type apiResponse struct {
	Versions  []string  `json:"versions"`
	Aliases   []string  `json:"aliases"`
	Manifests manifests `json:"manifests"`
}

type buildsResponse struct {
	ErrorMessage string    `json:"error-message"`
	Builds       []string  `json:"builds"`
	Manifests    manifests `json:"manifests"`
}

func (br *buildsResponse) Error() string {
	return br.ErrorMessage
}

// GetBuilds returns a list of all available builds for a given version
func GetBuilds(version string) ([]string, error) {
	apiCall := fmt.Sprintf("versions/%s/builds?x-elastic-no-kpi=true", version)

	request := http.HttpRequest{
		Url:                getApiUrl(apiCall),
		Method:             "GET",
		IgnoreCertificates: true,
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
	}

	bodyBytes, err := request.DoHttpRequestProcessBody()
	if err != nil {
		return []string{}, err
	}

	var response buildsResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return []string{}, err
	}

	return response.Builds, nil
}

// GetVersions returns a list of all available versions
func GetVersions() ([]string, error) {
	req := http.HttpRequest{
		Url:                getApiUrl("versions?x-elastic-no-kpi=true"),
		Method:             "GET",
		IgnoreCertificates: true,
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
	}

	responseBytes, err := req.DoHttpRequestProcessBody()
	if err != nil {
		return []string{}, err
	}

	var response apiResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return []string{}, err
	}

	return response.Versions, nil
}

// getApiUrl returns the full URL for the artifacts API
func getApiUrl(api string) string {
	return fmt.Sprintf("%s/%s", artifactsApiUrl, api)
}
