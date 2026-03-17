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
	"strings"

	yaml "github.com/goccy/go-yaml"
)

type Filter func(yamlFile YamlFile) bool

func FilterByYamlPath(key string, value string) Filter {
	return func(yamlFile YamlFile) bool {
		yamlPath, err := yaml.PathString("$." + key)
		if err != nil {
			return false
		}

		var yamlValue string
		if err = yamlPath.Read(strings.NewReader(string(yamlFile.Bytes)), &yamlValue); err != nil {
			return false
		}

		return (yamlValue == value)
	}
}

// ConfigFilesFilter returns a filter for configuration files
func ConfigFilesFilter(fileName string) bool {
	return hasSuffix(fileName, []string{".yml", ".yaml"})
}

// TemplatesFilter returns a filter for templates
func TemplatesFilter(fileName string) bool {
	return hasSuffix(fileName, []string{".yml.tmpl", ".yaml.tmpl"})
}

// hasSuffix checks a string for a range of suffixes.
func hasSuffix(value string, suffixes []string) (ret bool) {
	for _, suffix := range suffixes {
		if strings.HasSuffix(value, suffix) {
			return true
		}
	}
	return false
}
