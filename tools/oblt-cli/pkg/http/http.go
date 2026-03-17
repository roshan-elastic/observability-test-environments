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
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	// Type of Authentication
	AuthTypeApiKey string = "apiKey"
	AuthTypeUser   string = "user"
	AuthTypeToken  string = "token"

	// Params
	ParamUsername string = "username"
	ParamPassword string = "password"
	ParamApiKey   string = "apiKey"
	ParamToken    string = "token"
)

// HttpRequest Struct to hold a HTTP request.
type HttpRequest struct {
	Url                string
	IgnoreCertificates bool
	Method             string
	AuthType           string
	Headers            map[string]interface{}
	Body               string
	DryRun             bool
	AuthParams         map[string]string
}

// DoHttpRequest performs an HTTP request.
func (r *HttpRequest) DoHttpRequest() (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: r.IgnoreCertificates},
	}
	client := &http.Client{Transport: tr}
	var req *http.Request
	req, err = http.NewRequest(r.Method, r.Url, strings.NewReader(r.Body))
	cobra.CheckErr(err)
	r.addAuthentication(req)
	for header, value := range r.Headers {
		req.Header.Add(header, fmt.Sprintf("%v", value))
	}
	logger.Debugf("Request: %v", req)

	if !r.DryRun {
		resp, err = client.Do(req)
		logger.Debugf("Response: %v", resp)
	} else {
		resp = &http.Response{}
		resp.StatusCode = 200
		logger.Info.Print("Dry run, no request performed")
	}
	return resp, err
}

// DoHttpRequestProcessBody performs an HTTP request and returns the body.
// It returns the body only if the status code is in the range [200, 400).
func (r *HttpRequest) DoHttpRequestProcessBody() (resp []byte, err error) {
	var respHttp *http.Response
	var bodyBytes []byte
	if respHttp, err = r.DoHttpRequest(); err == nil {
		if bodyBytes, err = io.ReadAll(respHttp.Body); err == nil {
			// http.Status ==> [2xx, 4xx)
			if respHttp.StatusCode >= http.StatusOK && respHttp.StatusCode < http.StatusBadRequest {
				resp = bodyBytes
			}
		}
	}
	if err != nil {
		logger.Infof("could not read response body: %s", err)
	}
	return resp, err
}

// addAuthentication add the authentication headers selected with the AuthType value.
func (r *HttpRequest) addAuthentication(req *http.Request) {
	if r.AuthType == AuthTypeUser {
		logger.Debugf("Enable username authentication")
		req.SetBasicAuth(r.AuthParams[ParamUsername], r.AuthParams[ParamPassword])
	} else if r.AuthType == AuthTypeApiKey {
		logger.Debugf("Enable ApiKey authentication")
		value := fmt.Sprintf("ApiKey %s", r.AuthParams[ParamApiKey])
		req.Header.Add("Authorization", value)
	} else if r.AuthType == AuthTypeToken {
		logger.Debugf("Enable Token authentication")
		value := fmt.Sprintf("Bearer %s", r.AuthParams[ParamToken])
		req.Header.Add("Authorization", value)
	} else if r.AuthType == "" {
		logger.Debugf("No authentication")
	} else {
		log.Fatalf("unknown Authentication type: %s", r.AuthType)
	}
}

// ChooseAuthType choose the authentication type based on the parameters passed.
func ChooseAuthType(username string, password string, apiKey string, token string) (authType string, params map[string]string, err error) {
	if username != "" && password != "" && apiKey == "" {
		params = map[string]string{"username": username, "password": password}
		authType = AuthTypeUser
	} else if apiKey != "" && username == "" {
		params = map[string]string{"apiKey": apiKey}
		authType = AuthTypeApiKey
	} else if token != "" && username == "" {
		params = map[string]string{"token": token}
		authType = AuthTypeToken
	} else {
		err = fmt.Errorf("you should set and authentication method, and only one")
	}
	return authType, params, err
}
