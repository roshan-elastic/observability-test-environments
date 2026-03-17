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

// inspired by https://github.com/christophgysin/avpnc/blob/master/avpnc
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

// RunCommandCombinedOutput runs a command and show the output.
func RunCommandCombinedOutput(cmd string, dryRun bool, args ...string) (err error) {
	if !dryRun {
		connectVPN := exec.Command(cmd, args...)
		connectVPN.Env = os.Environ()
		var output []byte
		output, err = connectVPN.CombinedOutput()
		logger.Infof("out: %s", string(output))
	} else {
		logger.Infof("dry-run: %s %s", cmd, args)
	}
	return err
}

// RunBashScript runs a bash script and show the output.
func RunBashScript(script string, dryRun bool) (err error) {
	return RunCommandCombinedOutput("bash", dryRun, "-c", script)
}

// StartCommand starts a command and return the control to the main thread.
func StartCommand(name string, dryRun bool, args ...string) (command *exec.Cmd, err error) {
	if !dryRun {
		command = exec.Command(name, args...)
		command.Env = os.Environ()
		err = command.Start()
	} else {
		logger.Infof("dry-run: %s %s", name, args)
	}
	return command, err
}

// OpenBrowser opens a browser with the given URL.
func OpenBrowser(url string, dryRun bool) (err error) {
	osName := runtime.GOOS
	if dryRun && os.Getenv("GOOS") != "" {
		osName = os.Getenv("GOOS")
	}
	switch osName {
	case "linux":
		_, err = StartCommand("xdg-open", dryRun, url)
	case "windows":
		_, err = StartCommand("rundll32", dryRun, "url.dll,FileProtocolHandler", url)
	case "darwin":
		_, err = StartCommand("open", dryRun, url)
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
