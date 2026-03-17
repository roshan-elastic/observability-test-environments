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

// Package template It contains the functions to manipulate cluster configuration files.
package template

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

// Parse it uses the data passed to parse the CCS template to build the CCS cluster configuration file.
func Parse(data map[string]interface{}, templatePath, outputFilePath string) error {
	logger.Debugf("Template parameters: %v", data)
	var outfile *os.File
	os.MkdirAll(filepath.Dir(outputFilePath), 0777)
	outfile, err := os.Create(outputFilePath)
	if err == nil {
		defer outfile.Close()
		var tmpl *template.Template
		tmpl, err = template.New("template").Funcs(sprig.FuncMap()).ParseFiles(templatePath)
		if err == nil {
			err = tmpl.ExecuteTemplate(outfile, filepath.Base(templatePath), data)
		}
		if err == nil {
			err = CleanFile(outputFilePath)
		}
	}
	return err
}

func CleanFile(filePath string) (err error) {
	var inputFile []byte
	if inputFile, err = os.ReadFile(filePath); err == nil {
		removeComments := regexp.MustCompile(`(?m)^#.*\n?`)
		removeEmptyNewLines := regexp.MustCompile(`(?m)^\n`)
		removeEmptyCarrierReturnWin := regexp.MustCompile(`\r`)
		cleanedFile := removeComments.ReplaceAllString(string(inputFile), "")
		cleanedFile = removeEmptyCarrierReturnWin.ReplaceAllString(cleanedFile, "")
		cleanedFile = removeEmptyNewLines.ReplaceAllString(cleanedFile, "")
		err = os.WriteFile(filePath, []byte(cleanedFile), 0644)
	}
	return err
}

// CheckAllVariablesDefined Searches for the string '<no value>' in the file,
// if there is a variable undefined, the file will contain the string '<no value>'.
func CheckAllVariablesDefined(filePath string) (ret bool) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Errorf("Error reading file %s: %v", filePath, err)
	}
	return !strings.Contains(string(data), "<no value>") && err == nil
}
