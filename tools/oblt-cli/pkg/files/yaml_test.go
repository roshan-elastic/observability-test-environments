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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatches(t *testing.T) {
	tmpDir := t.TempDir()
	ymlFilePath := filepath.Join(tmpDir, "file.yml")

	ymlContent := `
country: spain
city: toledo
store:
  book:
    - author: john
      price: 10
    - author: ken
      price: 12
  bicycle:
    color: red
    price: 19.95
`
	err := os.WriteFile(ymlFilePath, []byte(ymlContent), 0644)
	assert.Nil(t, err, "failed to write yaml file")

	d, b, err := ReadYamlFile(ymlFilePath)
	assert.Nil(t, err, "failed to read yaml file")
	file := YamlFile{
		Data: d, Owner: "foo", Path: ymlFilePath, Bytes: b,
	}

	t.Run("single match", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("country", "spain")))
		assert.True(t, file.Matches(FilterByYamlPath("city", "toledo")))
		assert.False(t, file.Matches(FilterByYamlPath("store", "1")))
		assert.False(t, file.Matches(FilterByYamlPath("non-existing-key", "3")))
	})

	t.Run("double match", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("country", "spain"), FilterByYamlPath("city", "toledo")))
		assert.True(t, file.Matches(FilterByYamlPath("city", "toledo"), FilterByYamlPath("country", "spain")))
		assert.False(t, file.Matches(FilterByYamlPath("two", "2"), FilterByYamlPath("one", "1111")))
		assert.False(t, file.Matches(FilterByYamlPath("one", "1"), FilterByYamlPath("two", "2222")))
	})

	t.Run("single match including nested", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("store.bicycle.color", "red")))
		assert.False(t, file.Matches(FilterByYamlPath("store.color", "red")))
	})

	t.Run("single match for arrays", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("store.book[0].author", "john")))
		assert.True(t, file.Matches(FilterByYamlPath("store.book[1].price", "12")))
	})

	t.Run("double match including nested", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("store.book[0].author", "john"), FilterByYamlPath("store.bicycle.color", "red")))
		assert.True(t, file.Matches(FilterByYamlPath("store.bicycle.color", "red"), FilterByYamlPath("store.book[0].author", "john")))
		assert.False(t, file.Matches(FilterByYamlPath("store.book[0].author", "john"), FilterByYamlPath("one", "1111")))
		assert.False(t, file.Matches(FilterByYamlPath("one", "1"), FilterByYamlPath("store.book[0].author", "john")))
	})

	t.Run("double match for arrays", func(t *testing.T) {
		assert.True(t, file.Matches(FilterByYamlPath("country", "spain"), FilterByYamlPath("store.bicycle.color", "red")))
		assert.True(t, file.Matches(FilterByYamlPath("store.bicycle.color", "red"), FilterByYamlPath("country", "spain")))
		assert.False(t, file.Matches(FilterByYamlPath("store.bicycle.color", "red"), FilterByYamlPath("one", "1111")))
		assert.False(t, file.Matches(FilterByYamlPath("one", "1"), FilterByYamlPath("store.bicycle.color", "red")))
	})
}
