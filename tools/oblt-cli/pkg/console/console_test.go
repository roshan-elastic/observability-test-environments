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

package console

import (
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestPrintYamlFiles(t *testing.T) {
	yamlfile := []byte(`name: data
description: some description`)
	data := make(map[interface{}]interface{})
	err := yaml.Unmarshal(yamlfile, &data)
	assert.Nil(t, err)

	title := "Test"
	items := []files.YamlFile{
		files.YamlFile{
			Owner: "foo",
			Bytes: []byte{1, 2, 3, 5},
			Data:  data,
			Path:  "/tmp/foo",
		},
		files.YamlFile{
			Owner: "bar",
			Bytes: []byte{1, 2, 3, 5},
			Data:  data,
			Path:  "/tmp/bar",
		},
	}
	t.Run("Test Print Yaml Files", func(t *testing.T) {
		data := PrintYamlFiles(title, items)
		assert.True(t, len(data) == 2)
		assert.True(t, len(data[0]) == 4)
		assert.True(t, len(data[1]) == 4)
	})
}
