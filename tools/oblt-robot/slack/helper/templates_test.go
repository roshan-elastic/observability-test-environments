// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package helper

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateToString(t *testing.T) {
	// Setup: Create a temporary template file
	tempFile, err := os.CreateTemp("", "template-*.txt")
	assert.NoError(t, err, "Failed to create temporary template file")

	tempFileContent := "Hello, {{ .Name }}!"
	_, err = tempFile.WriteString(tempFileContent)
	assert.NoError(t, err, "Failed to create temporary template file")
	tempFile.Close()

	defer os.Remove(tempFile.Name()) // Clean up after the test

	// Test case
	tests := []struct {
		name     string
		filePath string
		data     interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid template rendering",
			filePath: tempFile.Name(),
			data: struct {
				Name string
			}{
				Name: "World",
			},
			want:    "Hello, World!",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemplateToString(tt.filePath, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTemplateStringToString(t *testing.T) {
	// Test case
	tests := []struct {
		name    string
		content string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "Valid template rendering",
			content: "Hello, {{ .Name }}!",
			data: struct {
				Name string
			}{
				Name: "World",
			},
			want:    "Hello, World!",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemplateStringToString(tt.content, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	// Setup: Create a temporary template file
	pattern := "template-*.txt"
	tempFile, err := os.CreateTemp("", pattern)
	assert.NoError(t, err, "Failed to create temporary template file")

	tempFileContent := "Hello, {{ .Name }}!"
	_, err = tempFile.WriteString(tempFileContent)
	assert.NoError(t, err, "Failed to create temporary template file")
	tempFile.Close()

	defer os.Remove(tempFile.Name()) // Clean up after the test

	fileName := path.Base(tempFile.Name())
	dir := path.Dir(tempFile.Name())
	fileSys := fs.FS(os.DirFS(dir))
	args := struct {
		Name string
	}{
		Name: "World",
	}

	// Test case
	want := "Hello, World!"
	got, err := RenderTemplate(fileSys, fileName, args)
	assert.NoError(t, err, "Failed to render template")
	assert.Equal(t, want, got.String())
}

func TestJsonEscape(t *testing.T) {
	// Test case
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String without special characters",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "String with special characters",
			input:    `{"name":"John","age":30}`,
			expected: `{\"name\":\"John\",\"age\":30}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JsonEscape(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestChunks(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		input      string
		chunkSize  int
		wantChunks []string
	}{
		{
			name:       "Empty string",
			input:      "",
			chunkSize:  3,
			wantChunks: nil,
		},
		{
			name:       "String with length less than chunk size",
			input:      "Hello",
			chunkSize:  10,
			wantChunks: []string{"Hello"},
		},
		{
			name:       "String with length equal to chunk size",
			input:      "World",
			chunkSize:  5,
			wantChunks: []string{"World"},
		},
		{
			name:       "String with length greater than chunk size",
			input:      "Lorem ipsum dolor sit amet",
			chunkSize:  7,
			wantChunks: []string{"Lorem i", "psum do", "lor sit", " amet"},
		},
		{
			name:       "String with length greater than chunk size and with new line character",
			input:      "Lorem ipsum dolor sit amet\n consectetur\n adipiscing\n elit\n sed\n do\n eiusmod\n tempor\n incididunt\n ut\n labore\n et\n dolore\n magna\n aliqua",
			chunkSize:  10,
			wantChunks: []string{"Lorem ipsu", "m dolor si", "t amet\n", " consectet", "ur\n", " adipiscin", "g\n elit\n", " sed\n do\n", " eiusmod\n", " tempor\n", " incididun", "t\n ut\n", " labore\n", " et\n", " dolore\n", " magna\n", " aliqua"},
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChunks := Chunks(tt.input, tt.chunkSize)
			assert.Equal(t, tt.wantChunks, gotChunks)
		})
	}
}
