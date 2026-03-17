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
	"fmt"
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"gopkg.in/yaml.v3"
)

// YamlFile It stores the path and the content of a YAML file.
type YamlFile struct {
	Bytes []byte                      `json:"-"`
	Data  map[interface{}]interface{} `json:"-"`
	Owner string                      `json:"owner"`
	Path  string                      `json:"path"`
}

// Matches applies a set of filters to a Yaml file.
// The filters are applied in the order they are passed in.
func (f YamlFile) Matches(filters ...Filter) bool {
	// Make sure there are actually filters to be applied.
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		if !filter(f) {
			// at the moment a filter is not satisfied, return false
			return false
		}
	}

	// all filters satisfied
	return true
}

// Find It reads all YAML files in a folder, applying filters, and return a YamlFile struct for each one.
func Find(folder string, fileTypeFilterFn func(fileName string) bool, filters ...Filter) (list []YamlFile, err error) {
	err = filepath.Walk(folder, func(path string, file os.FileInfo, err error) error {
		if err == nil && fileTypeFilterFn(file.Name()) {
			logger.Debugf("Checking file %s", path)
			filePath := path
			owner := filepath.Base(filepath.Dir(filePath))

			yamlFile, err := ReadYamlOj(filePath, owner)
			if err == nil && len(filters) == 0 || yamlFile.Matches(filters...) {
				list = append(list, yamlFile)
			}
		}
		return err
	})

	if len(list) == 0 {
		err = fmt.Errorf("files not found at: %s", folder)
	}
	return list, err
}

// ReadYamlFile It reads a YML file an return a Mapa with their key values.
func ReadYamlFile(filePath string) (data map[interface{}]interface{}, bytes []byte, err error) {
	logger.Debugf("Reading file %s", filePath)
	yamlfile, err := os.ReadFile(filePath)
	if err == nil {
		data = make(map[interface{}]interface{})
		err = yaml.Unmarshal(yamlfile, &data)
	}
	return data, yamlfile, err
}

// ReadYamlOj It reads a YML file an return a YamlFile struct.
func ReadYamlOj(filePath, owner string) (yamlFile YamlFile, err error) {
	d, b, err := ReadYamlFile(filePath)
	if err == nil {
		yamlFile = YamlFile{
			Data: d, Owner: owner, Path: filePath, Bytes: b,
		}
	}
	return yamlFile, err
}
