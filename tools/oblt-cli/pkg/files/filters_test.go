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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFilesFilter(t *testing.T) {
	type testData struct {
		path     string
		expected bool
	}

	testDatas := []testData{
		{
			path:     "a.yml",
			expected: true,
		},
		{
			path:     "a.yml.tmpl",
			expected: false,
		},
		{
			path:     "a.yaml",
			expected: true,
		},
		{
			path:     "a.yaml.tmpl",
			expected: false,
		},
	}

	for _, td := range testDatas {
		assert.Equal(t, td.expected, ConfigFilesFilter(td.path))
	}
}

func TestTemplatesFilter(t *testing.T) {
	type testData struct {
		path     string
		expected bool
	}

	testDatas := []testData{
		{
			path:     "a.yml",
			expected: false,
		},
		{
			path:     "a.yml.tmpl",
			expected: true,
		},
		{
			path:     "a.yaml",
			expected: false,
		},
		{
			path:     "a.yaml.tmpl",
			expected: true,
		},
	}

	for _, td := range testDatas {
		assert.Equal(t, td.expected, TemplatesFilter(td.path))
	}
}

func TestHasSuffix(t *testing.T) {
	type testData struct {
		path     string
		suffixes []string
	}

	testDatas := []testData{
		{
			path:     "a.yml",
			suffixes: []string{".yml"},
		},
		{
			path:     "a.yml.tmpl",
			suffixes: []string{".yml.tmpl"},
		},
		{
			path:     "a.yaml",
			suffixes: []string{".yaml"},
		},
		{
			path:     "a.yaml.tmpl",
			suffixes: []string{".yaml.tmpl"},
		},
	}

	for _, td := range testDatas {
		assert.True(t, hasSuffix(td.path, td.suffixes))
	}

	assert.False(t, hasSuffix("a.md", []string{".yml"}))
}
