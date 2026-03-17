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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuilds(t *testing.T) {
	version := "1.0.0"

	// Create a mock server to handle the HTTP request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/versions/"+version+"/builds", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Create a sample response
		builds := []string{"build1", "build2"}
		response := buildsResponse{
			Builds: builds,
		}

		// Marshal the response into JSON
		responseBytes, err := json.Marshal(response)
		assert.NoError(t, err)

		// Write the response to the client
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
	}))

	// Replace the URL in the GetBuilds function with the mock server URL
	oldURL := artifactsApiUrl
	artifactsApiUrl = mockServer.URL + "/v1"

	// Restore the original URL after the test
	defer func() {
		artifactsApiUrl = oldURL
	}()

	// Call the GetBuilds function
	builds, err := GetBuilds(version)

	// Assert that there are no errors
	assert.NoError(t, err)

	// Assert that the builds match the expected values
	expectedBuilds := []string{"build1", "build2"}
	assert.Equal(t, expectedBuilds, builds)
}
func TestGetVersions(t *testing.T) {
	// Create a mock server to handle the HTTP request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/versions", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Create a sample response
		versions := []string{"1.0.0", "2.0.0"}
		response := apiResponse{
			Versions: versions,
		}

		// Marshal the response into JSON
		responseBytes, err := json.Marshal(response)
		assert.NoError(t, err)

		// Write the response to the client
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
	}))

	// Replace the URL in the GetVersions function with the mock server URL
	oldURL := artifactsApiUrl
	artifactsApiUrl = mockServer.URL + "/v1"

	// Restore the original URL after the test
	defer func() {
		artifactsApiUrl = oldURL
	}()

	// Call the GetVersions function
	versions, err := GetVersions()

	// Assert that there are no errors
	assert.NoError(t, err)

	// Assert that the versions match the expected values
	expectedVersions := []string{"1.0.0", "2.0.0"}
	assert.Equal(t, expectedVersions, versions)
}
