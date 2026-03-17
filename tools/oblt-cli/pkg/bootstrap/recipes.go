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

package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpLib "net/http"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/files"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/http"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/maps"
	"github.com/spf13/cobra"
)

const (
	// Type of recipe
	TypeElasticsearch string = "elasticsearch"
	TypeKibana        string = "kibana"
	TypeFleet         string = "fleet"
	TypeApm           string = "apm"
	TypeCluster       string = "cluster"
)

// Recipe Struct to hold a recipe.
type Recipe struct {
	Description  string                 `json:"description"`
	Api          string                 `json:"api"`
	Method       string                 `json:"method"`
	Headers      map[string]interface{} `json:"headers"`
	Body         string                 `json:"body"`
	ReturnCode   []int                  `json:"return_code"`
	IgnoreErrors bool                   `json:"ignore_errors"`
}

// RecipesStruct struct representing the bootstrap recipe.
type RecipesStruct struct {
	BootstrapType      string
	RecipesJson        string
	AuthType           string
	AuthParams         map[string]string
	Url                string
	DryRun             bool
	IgnoreCertificates bool
	BootstarpFolder    string
	TemplateParams     map[string]interface{}
}

// RecipesList struct to store a list of all recipes and types.
type RecipesList struct {
	Elasticsearch []files.YamlFile `json:"elasticsearch"`
	Kibana        []files.YamlFile `json:"kibana"`
}

// Apply the recipes configured to the type of bootstrap configured.
func (r *RecipesStruct) Apply(obltTestEnvironments *clusters.ObltEnvironmentsRepository) map[string]map[string]bool {
	recipesList, err := r.loadRecipes(obltTestEnvironments)
	if err != nil {
		logger.Infof("%v", err)
	}
	var results = make(map[string]map[string]bool)
	for _, recipe := range recipesList {
		if recipe.Owner == r.BootstrapType {
			logger.Debugf("Applying recipe : %s:%s", recipe.Owner, recipe.Path)
			err := r.doRequest(recipe.Data)
			if err != nil {
				logger.Debugf("Recipe Failed: %s:%s", recipe.Owner, recipe.Path)
			}
			if results[recipe.Owner] == nil {
				results[recipe.Owner] = make(map[string]bool)
			}
			recipeName := filepath.Base(recipe.Path)
			results[recipe.Owner][recipeName] = (err == nil)
		}
	}
	return results
}

// ListRecipes return the list of recipes available for each category.
func ListRecipes(obltTestEnvironments *clusters.ObltEnvironmentsRepository, bootstrapFolder string) map[string][]files.YamlFile {
	logger.Debugf("Load all recipes.")

	var results = make(map[string][]files.YamlFile)
	recipesElasticsearch := obltTestEnvironments.ListBootstrapRecipesInFolder(TypeElasticsearch, bootstrapFolder)
	recipesKibana := obltTestEnvironments.ListBootstrapRecipesInFolder(TypeKibana, bootstrapFolder)
	results[TypeElasticsearch] = recipesElasticsearch
	results[TypeKibana] = recipesKibana
	return results
}

// loadRecipes load the recipes os the category BootstrapType from the repository.
// if a list of recipes is passed in RecipesJson only those recipes are loaded.
func (r *RecipesStruct) loadRecipes(obltTestEnvironments *clusters.ObltEnvironmentsRepository) (recipesList []files.YamlFile, err error) {
	if r.RecipesJson != "" {
		logger.Debugf("Load list of recipes : %s:%s", r.BootstrapType, r.RecipesJson)
		recipesList, err = obltTestEnvironments.LoadRecipesFromJsonInFolder(r.BootstrapType, r.RecipesJson, r.BootstarpFolder)
	} else {
		logger.Debugf("Load all recipes.")
		recipesList = obltTestEnvironments.ListBootstrapRecipesInFolder(r.BootstrapType, r.BootstarpFolder)
	}
	return recipesList, err
}

// doRequest perform a HTTP request using the recipe passed as parameter.
func (r *RecipesStruct) doRequest(recipe map[interface{}]interface{}) (err error) {
	recipeObj := validateRecipe(recipe)
	var body string
	var resp *httpLib.Response
	logger.Debugf("Performing the recipe %v", recipe)
	if body, err = r.processBodyTemplate(recipeObj.Body); err == nil {
		request := http.HttpRequest{
			Url:                fmt.Sprintf("%s%s", r.Url, recipeObj.Api),
			IgnoreCertificates: r.IgnoreCertificates,
			Method:             recipeObj.Method,
			AuthType:           r.AuthType,
			Headers:            recipeObj.Headers,
			Body:               body,
			DryRun:             r.DryRun,
			AuthParams:         r.AuthParams,
		}
		logger.Debugf("Performing the request %v", request)
		resp, err = request.DoHttpRequest()
		if err == nil {
			logger.Debugf("Response: %v", resp)
			err = r.processResponse(resp, recipeObj)
		}
	}
	return err
}

func (r *RecipesStruct) processBodyTemplate(bodyTemplate string) (body string, err error) {
	var tmpl *template.Template
	var buf bytes.Buffer = bytes.Buffer{}
	tmpl = template.New("template").Funcs(sprig.FuncMap())
	if tmpl, err = tmpl.Parse(bodyTemplate); err == nil {
		err = tmpl.Execute(&buf, r.TemplateParams)
	}
	return buf.String(), err
}

// processResponse check if the response status code is in the list of expected status code.
func (r *RecipesStruct) processResponse(resp *httpLib.Response, recipeObj Recipe) (err error) {
	found := false
	for _, code := range recipeObj.ReturnCode {
		if code == resp.StatusCode {
			found = true
			break
		}
	}

	if !found {
		logger.Debugf("wrong status code %d != %v\n", recipeObj.ReturnCode, resp.StatusCode)
		if !recipeObj.IgnoreErrors {
			err = fmt.Errorf("wrong status code %d != %v", recipeObj.ReturnCode, resp.StatusCode)
		}
	}
	return err
}

// validateRecipe perform a conversion to a structure this ensure that types are correct.
func validateRecipe(recipe interface{}) (obj Recipe) {
	logger.Debugf("Validate Recipe.")
	jsonString, err := json.Marshal(maps.InterfaceToMapInterface(recipe))
	cobra.CheckErr(err)
	err = json.Unmarshal(jsonString, &obj)
	cobra.CheckErr(err)
	return obj
}
