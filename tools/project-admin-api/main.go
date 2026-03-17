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
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	spec "github.com/elastic/observability-test-environments/tools/project-admin-api/pkg/api/admin"
)

// main is a simple example of how to use the generated client code
func main() {
	apiKey := os.Getenv("API_KEY")
	client, err := spec.NewClientWithResponses("https://global.qa.cld.elstc.co")
	reqOptions := spec.RequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", apiKey))
		req.Header.Add("accept", "application/json")
		return nil
	})
	params := spec.AdminListElasticsearchProjectsParams{}
	if err == nil {
		client.AdminListElasticsearchProjects(context.Background(), &params, reqOptions)
	}
}
