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
package http

import (
	"encoding/base64"
	"fmt"
	"io"
	httpLib "net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDoHttpRequest(t *testing.T) {
	// Start a local HTTP server to mock the response
	server := httptest.NewServer(httpLib.HandlerFunc(func(w httpLib.ResponseWriter, r *httpLib.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Authorization"), "Basic "+base64.StdEncoding.EncodeToString([]byte("foo:bar")))
		fmt.Fprintln(w, "Hello, world!")
	}))
	defer server.Close()

	// Make a request to the local server
	url := server.URL
	request := HttpRequest{
		Url:                url,
		Method:             "GET",
		IgnoreCertificates: true,
		AuthType:           AuthTypeUser,
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body:   `{"foo": "bar"}`,
		DryRun: false,
		AuthParams: map[string]string{
			ParamUsername: "foo",
			ParamPassword: "bar",
		},
	}
	response, err := request.DoHttpRequest()
	assert.NilError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.NilError(t, err)
	// Verify the response body
	expectedBody := "Hello, world!\n"
	assert.Equal(t, string(body), expectedBody)

	bodyBytes, err := request.DoHttpRequestProcessBody()
	assert.NilError(t, err)
	assert.Equal(t, string(bodyBytes), expectedBody)
}

func TestDoHttpRequestApiKey(t *testing.T) {
	// Start a local HTTP server to mock the response
	server := httptest.NewServer(httpLib.HandlerFunc(func(w httpLib.ResponseWriter, r *httpLib.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Authorization"), "ApiKey foo")
		fmt.Fprintln(w, "Hello, world!")
	}))
	defer server.Close()

	// Make a request to the local server
	url := server.URL
	request := HttpRequest{
		Url:                url,
		Method:             "GET",
		IgnoreCertificates: true,
		AuthType:           AuthTypeApiKey,
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body:   `{"foo": "bar"}`,
		DryRun: false,
		AuthParams: map[string]string{
			ParamApiKey: "foo",
		},
	}
	response, err := request.DoHttpRequest()
	assert.NilError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.NilError(t, err)
	// Verify the response body
	expectedBody := "Hello, world!\n"
	assert.Equal(t, string(body), expectedBody)
}

func TestDoHttpRequestToken(t *testing.T) {
	// Start a local HTTP server to mock the response
	server := httptest.NewServer(httpLib.HandlerFunc(func(w httpLib.ResponseWriter, r *httpLib.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer foo")
		fmt.Fprintln(w, "Hello, world!")
	}))
	defer server.Close()

	// Make a request to the local server
	url := server.URL
	request := HttpRequest{
		Url:                url,
		Method:             "GET",
		IgnoreCertificates: true,
		AuthType:           AuthTypeToken,
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body:   `{"foo": "bar"}`,
		DryRun: false,
		AuthParams: map[string]string{
			ParamToken: "foo",
		},
	}
	response, err := request.DoHttpRequest()
	assert.NilError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.NilError(t, err)
	// Verify the response body
	expectedBody := "Hello, world!\n"
	assert.Equal(t, string(body), expectedBody)
}
