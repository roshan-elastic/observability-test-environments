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

// Package gcp contains the functions related to the Google Cloud Platform.
package gcp

import (
	"os"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	seed := config.Seed(5)
	secretPath := "projects/8560181848/secrets/oblt-cli-test-" + seed

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Skipping test that requires GCP credentials")
	}

	gcsm, err := NewClient()
	assert.NoError(t, err, "Error should be nil")
	assert.NotNil(t, gcsm, "SecretsManager should not be nil")
	defer gcsm.Close()
	secretName := gcsm.GetSecretName(secretPath)

	err = gcsm.CreateSecret(secretPath, "foo")
	assert.NoError(t, err, "Error should be nil")

	secretValue, err := gcsm.GetSecret(secretPath)
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, "foo", secretValue, "Secret value should match")

	listSecrets, err := gcsm.ListSecrets("projects/8560181848", "name:"+secretName)
	assert.NoError(t, err, "Error should be nil")
	assert.Contains(t, listSecrets, secretPath, "Secret path should be in the list")

	err = gcsm.DeleteSecret(secretPath)
	assert.NoError(t, err, "Error should be nil")

	_, err = gcsm.GetSecret(secretPath)
	assert.Error(t, err, "Error should not be nil")
}
