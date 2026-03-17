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
package files

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

// SaveResults It saves the reults of the operation in a JSON file.
func SaveResults(data interface{}, filePath string) (err error) {
	var jsonStr []byte
	jsonStr, err = json.MarshalIndent(data, "", "  ")
	if err == nil {
		err = writeFile(filePath, jsonStr)
	}
	return err
}

// SaveResultsRaw It saves the reults of the operation in a Raw file.
func SaveResultsRaw(data string, filePath string) (err error) {
	return writeFile(filePath, []byte(data))
}

// writeFile It writes the data in the file.
// If the file name is "-" or "stdout" it writes the data in the stdout.
func writeFile(filePath string, data []byte) (err error) {
	if len(filePath) > 0 {
		if filePath != "-" && filePath != "stdout" {
			var absolutePath string
			if absolutePath, err = filepath.Abs(filePath); err == nil {
				logger.Infof("Writing output file %s", absolutePath)
				err = os.WriteFile(filePath, data, 0644)
			}
		} else {
			logger.Infof("Writing output file stdout")
			_, err = os.Stdout.Write(data)
		}
	}
	return err
}
