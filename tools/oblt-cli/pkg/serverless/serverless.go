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

package serverless

import (
	"context"
	"fmt"
	"net/http"

	spec "github.com/elastic/observability-test-environments/tools/project-api/pkg/api/user/v1"
)

// ClusterInfo returns the cluster info for the given cluster ID.
func ClusterInfo(id spec.ProjectID, endpoint, apiKey string) (info *spec.GetElasticsearchProjectResponse, err error) {
	reqOptions := spec.RequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", apiKey))
		req.Header.Add("accept", "application/json")
		return nil
	})
	client, err := spec.NewClientWithResponses(endpoint)
	if err == nil {
		info, err = client.GetElasticsearchProjectWithResponse(context.TODO(), id, reqOptions)
	}
	return info, err
}
