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
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveResults(t *testing.T) {
	type testData struct {
		Name string
		Age  int
	}
	data := testData{Name: "John", Age: 30}
	filePath := "test.json"
	err := SaveResults(data, filePath)
	assert.NoError(t, err, "SaveResults failed with error")
	defer os.Remove(filePath)

	file, err := os.Open(filePath)
	assert.NoError(t, err, "Failed to open file")
	defer file.Close()

	var result testData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&result)
	assert.NoError(t, err, "Failed to decode file")
	assert.Equal(t, data.Name, result.Name)
	assert.Equal(t, data.Age, result.Age)
}

func TestSaveResultsRaw(t *testing.T) {
	rawFilePath := "test.txt"
	rawData := "Hello, world!"
	err := SaveResultsRaw(rawData, rawFilePath)
	assert.NoError(t, err, "SaveResultsRaw failed with error")
	defer os.Remove(rawFilePath)

	rawFile, err := os.Open(rawFilePath)
	assert.NoError(t, err, "Failed to open file")
	defer rawFile.Close()

	rawResult, err := io.ReadAll(rawFile)
	assert.NoError(t, err, "Failed to read file")
	assert.Equal(t, rawData, string(rawResult))
}
