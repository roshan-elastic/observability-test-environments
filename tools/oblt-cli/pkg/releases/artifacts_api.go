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

// Package releases It contains the functions to interact with the Unified Releases.
package releases

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/curl"
)

const (
	artifactsUrl = "https://artifacts-api.elastic.co/v1/versions"
)

// Builds defines a build defined in the Artifacts API
type Builds struct {
	Branch      string
	BuildId     string
	Version     string
	ReleaseDate string
	Query       string
}

func getLatestBuildVersion(version string) (Builds, error) {
	url := fmt.Sprintf("%s/%s/builds/latest", artifactsUrl, version)
	r := curl.HTTPRequest{
		URL: url,
	}
	response, err := curl.Get(r)
	if err != nil {
		return Builds{}, fmt.Errorf("failed to get the latest version for %s in the artifacts-api: %w", version, err)
	}
	jsonParsed, err := gabs.ParseJSON([]byte(response))
	if err != nil {
		return Builds{}, fmt.Errorf("failed to decode the latest version response: %w", err)
	}
	buildObject := jsonParsed.Path("build")
	build := Builds{
		Branch:      buildObject.Path("branch").Data().(string),
		ReleaseDate: buildObject.Path("end_time").Data().(string),
		Version:     version,
		BuildId:     buildObject.Path("build_id").Data().(string),
		Query:       url,
	}
	return build, nil
}

func getAliases() ([]string, error) {
	r := curl.HTTPRequest{
		URL: artifactsUrl,
	}

	response, err := curl.Get(r)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions in the artifacts-api: %w", err)
	}
	jsonParsed, err := gabs.ParseJSON([]byte(response))
	if err != nil {
		return nil, fmt.Errorf("failed to decode the aliases response: %w", err)
	}
	aliasesObject := jsonParsed.Path("aliases")
	aliases, _ := aliasesObject.Children()

	var ret []string
	for _, child := range aliases {
		ret = append(ret, child.Data().(string))
	}

	return ret, nil
}

func isSnapshot(version string) bool {
	return strings.Contains(strings.ToLower(version), "snapshot")
}

func isRelease(version string) bool {
	return !isSnapshot(version)
}
