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
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	defer func() {
		logger.SetOutput(os.Stderr)
	}()

	err := RunCommandCombinedOutput("echo", false, "hello", "world")
	assert.NoError(t, err)
	logs := buf.String()
	assert.Contains(t, logs, "out: hello world")
}

func TestRunCommandDryRun(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	defer func() {
		logger.SetOutput(os.Stderr)
	}()

	err := RunCommandCombinedOutput("echo", true, "hello", "world")
	assert.NoError(t, err)
	logs := buf.String()
	assert.Contains(t, logs, "dry-run: echo [hello world]")
}

func TestOpenBrowser(t *testing.T) {
	tests := []struct {
		name    string
		os      string
		url     string
		wantCmd string
		wantErr bool
	}{
		{
			name:    "Linux",
			os:      "linux",
			url:     "http://example.com",
			wantCmd: "xdg-open [http://example.com]",
			wantErr: false,
		},
		{
			name:    "Windows",
			os:      "windows",
			url:     "http://example.com",
			wantCmd: "rundll32 [url.dll,FileProtocolHandler http://example.com]",
			wantErr: false,
		},
		{
			name:    "Darwin",
			os:      "darwin",
			url:     "http://example.com",
			wantCmd: "open [http://example.com]",
			wantErr: false,
		},
		{
			name:    "Unsupported",
			os:      "unsupported",
			url:     "http://example.com",
			wantCmd: "",
			wantErr: true,
		},
	}

	var buf bytes.Buffer
	logger.SetOutput(&buf)
	defer func() {
		logger.SetOutput(os.Stderr)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the runtime.GOOS variable
			oldGOOS := runtime.GOOS
			defer func() { os.Setenv("GOOS", oldGOOS) }()
			os.Setenv("GOOS", tt.os)

			err := OpenBrowser(tt.url, true)

			logs := buf.String()
			assert.Contains(t, logs, tt.wantCmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			buf.Reset()
		})
	}
}
