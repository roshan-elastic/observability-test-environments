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

// Package gcp contains the functions related to the Google Cloud Platform.
// This file contains the functions to authenticate to GCP.
package gcp

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/cmd"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// googleOauthConfig is the OAuth2 configuration for Google
// for info avout scopes see https://developers.google.com/identity/protocols/oauth2/scope
var googleOauthConfig = &oauth2.Config{
	// from lib/googlecloudsdk/core/config.py
	ClientID:     "32555940559.apps.googleusercontent.com",
	ClientSecret: "ZmssLNjJy2998hD4CTg2ejr2",
	Endpoint:     google.Endpoint,
	Scopes: []string{"openid",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
		// "https://www.googleapis.com/auth/appengine.admin",
		// "https://www.googleapis.com/auth/sqlservice.login",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/accounts.reauth",
	},
	RedirectURL: "http://localhost:8000/auth/google/callback",
}

// use PKCE to protect against CSRF attacks
// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
var verifier = oauth2.GenerateVerifier()

const (
	// oauthGoogleUrlAPI is the URL to get the user info from Google
	oauthGoogleUrlAPI  = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	default_project_id = "elastic-observability"
)

// wg is used to wait for the callback server to finish
var wg sync.WaitGroup

// GcpAuth is the struct to authenticate to GCP
type GcpAuth struct {
	ProjectId string
	Options   []option.ClientOption
}

type Auth interface {
	GetProjectId() string
	Authenticate() (err error)
	AuthenticateInteractive() (err error)
}

// NewGcpAuth returns a new instance of GcpAuth
func NewGcpAuth() GcpAuth {
	return GcpAuth{ProjectId: default_project_id, Options: []option.ClientOption{}}
}

// NewGcpAuthWithProjectId returns a new instance of GcpAuth with the given project ID
func NewGcpAuthWithProjectId(projectId string) GcpAuth {
	return GcpAuth{ProjectId: projectId}
}

// GetProjectId returns the project ID
func (a *GcpAuth) GetProjectId() string {
	return a.ProjectId
}

// Authenticate authenticates to GCP
// If GOOGLE_APPLICATION_CREDENTIALS is set or there is a Application Default Credentials defined, it uses it to authenticate
// Otherwise, it starts an interactive authentication
// for more info see https://cloud.google.com/docs/authentication/application-default-credentials
func (a *GcpAuth) Authenticate() (err error) {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err == nil {
		a.Options = append(a.Options, option.WithCredentials(credentials))
	} else if os.Getenv("CI") == "" {
		a.AuthenticateInteractive()
	} else {
		logger.Errorf("failed to find default credentials: %v", err)
	}
	return err
}

// AuthenticateInteractive starts an interactive authentication
// It starts a callback server and opens the browser to authenticate
func (a *GcpAuth) AuthenticateInteractive() (err error) {
	logger.Debugf("Authenticating to GCP")
	wg.Add(1)
	go a.callbackServer()
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
	logger.Debugf("Visit the URL for the auth dialog: %v", url)
	err = cmd.OpenBrowser(url, false)
	if err != nil {
		logger.Errorf("failed to open the browser: %v", err)
	} else {
		wg.Wait()
	}
	return err
}

// CallbackServer starts a callback server to handle the authentication callback
func (a *GcpAuth) callbackServer() {
	server := &http.Server{
		Addr:    ":8000",
		Handler: a.newHandler(),
	}
	logger.Infof("Starting Google Authentication Callback HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Infof("%v", err)
	} else {
		logger.Infof("Server closed!")
	}
}

// oauthGoogleCallback handles the authentication callback
// It exchanges the code for a token and gets the user info
// It writes the user info to the response
func (a *GcpAuth) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	defer wg.Done()
	code := r.FormValue("code")
	ctx := context.Background()
	token, err := googleOauthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err == nil {
		a.Options = append(a.Options, option.WithTokenSource(googleOauthConfig.TokenSource(ctx, token)))
		testToken(token, w)
	} else {
		logger.Errorf("failed to exchange token: %s", err)
	}
}

// newHandler returns a new http.Handler to handle the authentication callback
func (a *GcpAuth) newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google/callback", a.oauthGoogleCallback)
	return mux
}

// testToken tests the token and writes the user info to the response.
func testToken(token *oauth2.Token, w http.ResponseWriter) (err error) {
	var contents []byte
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err == nil {
		defer response.Body.Close()
		contents, err = io.ReadAll(response.Body)
		if err == nil {
			logger.Debugf("Response: %s", string(contents))
			w.Write(contents)
		}
	}
	if err != nil {
		logger.Errorf("failed test the token: %v", err)
	}
	return err
}
