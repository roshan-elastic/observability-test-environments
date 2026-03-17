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
package helper

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"path"
	"strings"
	"text/template"

	// Additional functions for Go Templates https://masterminds.github.io/sprig/
	"github.com/Masterminds/sprig/v3"

	"github.com/pkg/errors"
)

// TemplateToString takes a template file and a data interface and returns a string to print.
func TemplateToString(filePath string, data interface{}) (temp string, err error) {
	if data == nil {
		data = struct{}{}
	}

	var b bytes.Buffer

	t, err := template.New(path.Base(filePath)).Funcs(sprig.FuncMap()).ParseFiles(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template ")
	}

	err = t.Execute(&b, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return b.String(), nil
}

// TemplateStringToString takes a template content and a data interface and returns a string to print.
func TemplateStringToString(content string, data interface{}) (temp string, err error) {
	if data == nil {
		data = struct{}{}
	}

	var b bytes.Buffer

	t, err := template.New("my-template").Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	err = t.Execute(&b, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return b.String(), nil
}

// RenderTemplate receives a file for a template and returns it as bytes to convert
func RenderTemplate(fileSys fs.FS, file string, args interface{}) (bytes.Buffer, error) {
	var buff bytes.Buffer

	// read the block-kit definition as a go template
	tmpl, err := template.New(file).Funcs(sprig.FuncMap()).ParseFS(fileSys, file)
	if err != nil {
		return buff, err
	}

	// we render the view
	err = tmpl.Execute(&buff, args)
	if err != nil {
		return buff, err
	}

	return buff, nil
}

// JsonEscape escapes the string to be used in a JSON
func JsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

// Chunks will split a string into chunks of the given size
func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = []string{}
	currentStart := 0
	currentEnd := chunkSize
	for currentStart < len(s) {
		chunkInProgress := s[currentStart:currentEnd]
		// Try to chunk on a new line
		if lastNewLine := strings.LastIndex(chunkInProgress, "\n"); lastNewLine != -1 {
			currentEnd = currentStart + lastNewLine + 1
			chunkInProgress = s[currentStart:currentEnd]
		}
		chunks = append(chunks, chunkInProgress)
		currentStart = currentEnd
		if currentEnd += chunkSize; currentEnd > len(s) {
			currentEnd = len(s)
		}
	}
	chunks = append(chunks, s[currentStart:])
	return removeEmptyStrings(chunks)
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
