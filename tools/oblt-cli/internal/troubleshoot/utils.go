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
package troubleshoot

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	// url troubleshooting.
	url = "https://elastic.github.io/observability-test-environments/user-guide/troubleshooting/#oblt-cli for further details."
)

// Warn it prints a warning with the given message and the troubleshoot URL.
func Warn(message string) {
	logger.Warnf("%s", message)
	logger.Warnf("%s", url)
}

// Error it prints an error with the given message and the troubleshoot URL.
func Error(message string) {
	logger.Errorf("%s", message)
	logger.Errorf("%s", url)
}

// CobraCheckErrWithWarning if an error then prints a warning message
func CobraCheckErrWithWarning(message string, err interface{}) {
	if err != nil {
		Warn(message)
	}
	cobra.CheckErr(err)
}

// CobraCheckErrWithError if an error then prints an error message
func CobraCheckErrWithError(message string, err interface{}) {
	if err != nil {
		Error(message)
	}
	cobra.CheckErr(err)
}
