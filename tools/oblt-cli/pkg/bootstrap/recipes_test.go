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

package bootstrap_test

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/bootstrap"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testsStruct struct {
	name     string
	fields   bootstrap.RecipesStruct
	validate func(results map[string]map[string]bool, output string)
}

func runTest(t *testing.T, tt testsStruct) {
	currentUsername := "the-user"
	slackChannel := "#0123456"
	tmp := t.TempDir()

	// configuration is needed because the repo is cloned under config's workspace
	viper.SetConfigFile(filepath.Join(tmp, "config.yml"))
	viper.Set(config.UsernameFlag, currentUsername)
	viper.Set(config.SlackChannelFlag, slackChannel)
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	viper.WriteConfig()

	// clone is needed to access all methods, using dry run
	obltTestEnvironments, err := clusters.BootstrapRepository(true)
	assert.NoError(t, err)

	t.Run(tt.name, func(t *testing.T) {
		var buf bytes.Buffer
		logger.Info.SetOutput(&buf)
		logger.Debug.SetOutput(&buf)
		logger.Error.SetOutput(&buf)
		logger.Warn.SetOutput(&buf)
		defer func() {
			logger.Info.SetOutput(os.Stdout)
			logger.Debug.SetOutput(os.Stdout)
			logger.Error.SetOutput(os.Stderr)
			logger.Warn.SetOutput(os.Stdout)
		}()
		r := &bootstrap.RecipesStruct{
			BootstrapType:      tt.fields.BootstrapType,
			RecipesJson:        tt.fields.RecipesJson,
			AuthType:           tt.fields.AuthType,
			AuthParams:         tt.fields.AuthParams,
			Url:                tt.fields.Url,
			DryRun:             tt.fields.DryRun,
			IgnoreCertificates: tt.fields.IgnoreCertificates,
			BootstarpFolder:    tt.fields.BootstarpFolder,
			TemplateParams:     tt.fields.TemplateParams,
		}
		results := r.Apply(obltTestEnvironments)
		output := buf.String()
		t.Log(output)
		tt.validate(results, output)
		assert.True(t, len(results) > 0)
		for _, recipes := range results {
			for _, result := range recipes {
				assert.True(t, result)
			}
		}
	})
}

func initConfiguration(t *testing.T) {
	logger.Verbose = true
	tmpDir := t.TempDir()
	parentDir := "foo"
	parentPath := filepath.Join(tmpDir, parentDir)
	fallbackCfgFile := filepath.Join(parentPath, "config.yml")
	config.Initialise(fallbackCfgFile)
}

// GetBootstrapRecipesDir returns the path to the test bootstrap recipes folder
func GetBootstrapRecipesDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return filepath.Join(basepath, "..", "..", "..", "..", "ansible", "ansible_collections", "oblt", "framework", "roles", "common", "files", "deployments", "bootstrap-test")
}

// GetBootstrapRecipesDir returns the path to the test bootstrap recipes folder
func GetBootstrapTemplatesRecipesDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return filepath.Join(basepath, "..", "..", "..", "..", "ansible", "ansible_collections", "oblt", "framework", "roles", "common", "files", "deployments", "bootstrap-test-template")
}

func TestRecipesStruct_ElasticsearchUserAuth(t *testing.T) {
	test := testsStruct{
		name: "Elasticsearch user authentication",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeElasticsearch,
			RecipesJson:   "[\"test\"]",
			AuthType:      http.AuthTypeUser,
			AuthParams: map[string]string{
				http.ParamUsername: "elastic",
				http.ParamPassword: "changeme",
			},
			Url:                "https://es.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Enable username authentication"))
			assert.True(t, strings.Contains(output, "Authorization:[Basic ZWxhc3RpYzpjaGFuZ2VtZQ==]"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/elasticsearch/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_KibanaUserAuth(t *testing.T) {
	test := testsStruct{
		name: "Kibana user authentication",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeKibana,
			RecipesJson:   "[\"test\"]",
			AuthType:      http.AuthTypeUser,
			AuthParams: map[string]string{
				http.ParamUsername: "elastic",
				http.ParamPassword: "changeme",
			},
			Url:                "https://es.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Enable username authentication"))
			assert.True(t, strings.Contains(output, "Authorization:[Basic ZWxhc3RpYzpjaGFuZ2VtZQ==]"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/kibana/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_ElasticsearchApiKeyAuth(t *testing.T) {
	test := testsStruct{
		name: "Elasticsearch apiKey authentication",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeElasticsearch,
			RecipesJson:   "[\"test\"]",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://es.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			assert.True(t, strings.Contains(output, "Authorization:[ApiKey foo]"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/elasticsearch/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_KibanaApiKeyAuth(t *testing.T) {
	test := testsStruct{
		name: "Kibana apiKey authentication",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeKibana,
			RecipesJson:   "[\"test\"]",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://kibana.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Authorization:[ApiKey foo]"))
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/kibana/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_ApmTokenyAuth(t *testing.T) {
	test := testsStruct{
		name: "Apm Token authentication",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeApm,
			RecipesJson:   "[\"test\"]",
			AuthType:      http.AuthTypeToken,
			AuthParams: map[string]string{
				http.ParamToken: "foo",
			},
			Url:                "https://kibana.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Authorization:[Bearer foo]"))
			assert.True(t, strings.Contains(output, "Enable Token authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/apm/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_ElasticsearchAllrecipes(t *testing.T) {
	test := testsStruct{
		name: "Elasticsearch all recipes",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeElasticsearch,
			RecipesJson:   "",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://es.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) >= 1)
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/elasticsearch/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_KibanaAllrecipes(t *testing.T) {
	test := testsStruct{
		name: "Kibana all recipes",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeKibana,
			RecipesJson:   "",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://kibana.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) >= 1)
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/kibana/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_TemplateRecipes(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "value00")
	test := testsStruct{
		name: "Template recipes",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeElasticsearch,
			RecipesJson:   "",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://es.example.com",
			DryRun:             true,
			IgnoreCertificates: false,
			BootstarpFolder:    GetBootstrapTemplatesRecipesDir(),
			TemplateParams: map[string]interface{}{
				"TestValue": "value01",
			},
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) >= 1)
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test-template/elasticsearch/test.yml`, output)
			assert.True(t, matched)
			assert.True(t, strings.Contains(output, "\"field00\": \"interpolate value00\""))
			assert.True(t, strings.Contains(output, "\"field01\": \"value01\""))
			assert.True(t, strings.Contains(output, "\"field00\": \"value00\""))
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_IgnoreCertificates(t *testing.T) {
	test := testsStruct{
		name: "Ignore certificates",
		fields: bootstrap.RecipesStruct{
			BootstrapType: bootstrap.TypeKibana,
			RecipesJson:   "",
			AuthType:      http.AuthTypeApiKey,
			AuthParams: map[string]string{
				http.ParamApiKey: "foo",
			},
			Url:                "https://kibana.example.com",
			DryRun:             true,
			IgnoreCertificates: true,
			BootstarpFolder:    GetBootstrapRecipesDir(),
		},
		validate: func(results map[string]map[string]bool, output string) {
			assert.True(t, len(results) == 1)
			assert.True(t, strings.Contains(output, "Enable ApiKey authentication"))
			matched, _ := regexp.MatchString(`Applying recipe.*bootstrap-test/kibana/test.yml`, output)
			assert.True(t, matched)
		},
	}
	initConfiguration(t)
	runTest(t, test)
}

func TestRecipesStruct_ListRecipes(t *testing.T) {
	currentUsername := "the-user"
	slackChannel := "#0123456"
	tmp := t.TempDir()

	// configuration is needed because the repo is cloned under config's workspace
	viper.SetConfigFile(filepath.Join(tmp, "config.yml"))
	viper.Set(config.UsernameFlag, currentUsername)
	viper.Set(config.SlackChannelFlag, slackChannel)
	// if CI environment variable is set, we use the https
	if os.Getenv("CI") != "" {
		viper.Set(config.GitHttpModeFlag, true)
		viper.Set(config.VerboseFlag, true)
	}
	viper.WriteConfig()

	// clone is needed to access all methods, using dry run
	obltTestEnvironments, err := clusters.BootstrapRepository(true)
	assert.NoError(t, err)

	t.Run("ListRecipes", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()
		initConfiguration(t)
		results := bootstrap.ListRecipes(obltTestEnvironments, GetBootstrapRecipesDir())
		assert.True(t, len(results) > 0)
	})
}
