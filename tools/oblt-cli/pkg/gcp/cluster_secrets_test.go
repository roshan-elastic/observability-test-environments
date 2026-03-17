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

func TestClusterStateSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadClusterStateSecret(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadEsSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadEsSecret(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret.EsSecret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadKibanaSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadKibanaSecret(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret.KibanaSecret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadApmSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadApmSecret(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret.ApmSecret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadFleetSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadFleetSecret(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret.FleetSecret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadESSDeployment(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := createClusterStateSecret(clusterSecrets)

	storedValue, err := clusterSecrets.ReadESSDeployment(secret.Name)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret.ESSDeployment, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadESSCredentials(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := NewESSCredentials()
	environment := createRandomName()
	secretPath := clusterSecrets.getESSCredentialsPath(environment)
	clusterSecrets.client.CreateSecret(secretPath, secret)

	storedValue, err := clusterSecrets.ReadESSCredentials(environment)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secretPath))
}

func TestReadServerlessCredentials(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := NewServerlessCredentials()
	environment := createRandomName()
	secretPath := clusterSecrets.getServerlessCredentialsPath(environment)
	clusterSecrets.client.CreateSecret(secretPath, secret)

	storedValue, err := clusterSecrets.ReadServerlessCredentials(environment)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secretPath))
}

func TestListClusterSecrets(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secret := NewClusterStateSecret()
	secret.Name = clusterName
	clusterSecrets.CreateClusterStateSecret(secret)
	secretPath := clusterSecrets.getClusterStateSecretPath(clusterName)
	SecretName := clusterSecrets.client.GetSecretName(secretPath)

	secretPaths, err := clusterSecrets.ListClusterSecrets(clusterName, true)
	assert.NoError(t, err, "failed to list the secrets: %v", err)
	assert.True(t, len(secretPaths) > 0, "the list is empty")
	assert.Contains(t, secretPaths, secretPath, "the secret is not in the list")
	secretNames, err := clusterSecrets.ListClusterSecrets(clusterName, false)
	assert.NoError(t, err, "failed to list the secrets: %v", err)
	assert.True(t, len(secretNames) > 0, "the list is empty")
	assert.Contains(t, secretNames, SecretName, "the secret is not in the list")

	clusterSecrets.client.DeleteSecret(clusterSecrets.getClusterStateSecretPath(secret.Name))
}

func TestReadLicense(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	secret := "{\"license\": \"test\"}"
	licanseType := createRandomName()
	secretPath := clusterSecrets.getLicenseSecretPath(licanseType)
	clusterSecrets.client.CreateSecret(secretPath, secret)

	storedValue, err := clusterSecrets.ReadLicense(licanseType)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")
}

// createRandomName creates a random cluster name.
func createRandomName() string {
	return "oblt-cli-test-" + config.Seed(5)
}

// createClusterStateSecret creates a new cluster state secret and stores it in the client.
func createClusterStateSecret(clusterSecrets *ClusterSecrets) ClusterStateSecret {
	secret := NewClusterStateSecret()
	secret.Name = createRandomName()
	clusterSecrets.CreateClusterStateSecret(secret)
	return secret
}

// createClusterSecretsMock creates a new ClusterSecrets instance using a mock client.
func createClusterSecretsMock() *ClusterSecrets {
	client := NewSecretsManagerMock()
	auth := NewAuthMock()
	return NewClusterSecretsWithClient(client, &auth)
}

// createClusterSecrets creates a new ClusterSecrets instance using the environment variables.
// If the GOOGLE_APPLICATION_CREDENTIALS is set, it will create a new client using the environment variables.
// Otherwise, it will create a mock client.
func createClusterSecrets(t *testing.T) *ClusterSecrets {
	var err error
	clusterSecrets := createClusterSecretsMock()
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		clusterSecrets, err = NewClusterSecrets()
		assert.NoError(t, err, "failed to create the client: %v", err)
	}
	return clusterSecrets
}

func TestReadDeployInfoSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secret := "test-secret"
	secretPath := clusterSecrets.getDeployInfoSecretPath(clusterName)
	clusterSecrets.client.CreateSecret(secretPath, secret)

	storedValue, err := clusterSecrets.ReadDeployInfoSecret(clusterName)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(secretPath)
}

func TestReadKibanaYamlSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secretValue := "test-secret-value"
	secretPath := clusterSecrets.getKibanaYamlSecretPath(clusterName)
	clusterSecrets.client.CreateSecret(secretPath, secretValue)

	storedValue, err := clusterSecrets.ReadKibanaYamlSecret(clusterName)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secretValue, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(secretPath)
}

func TestReadCredentialsSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secretValue := "test-secret-value"
	secretPath := clusterSecrets.getCredentialsSecretPath(clusterName)
	clusterSecrets.client.CreateSecret(secretPath, secretValue)

	storedValue, err := clusterSecrets.ReadCredentialsSecret(clusterName)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secretValue, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(secretPath)
}

func TestReadEnvSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secretValue := "test-secret-value"
	secretPath := clusterSecrets.getEnvSecretPath(clusterName)
	clusterSecrets.client.CreateSecret(secretPath, secretValue)

	storedValue, err := clusterSecrets.ReadEnvSecret(clusterName)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secretValue, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(secretPath)
}

func TestReadActivateSecret(t *testing.T) {
	clusterSecrets := createClusterSecrets(t)
	defer clusterSecrets.Close()
	clusterName := createRandomName()
	secret := "test-secret"
	secretPath := clusterSecrets.getActivateSecretPath(clusterName)
	clusterSecrets.client.CreateSecret(secretPath, secret)

	storedValue, err := clusterSecrets.ReadActivateSecret(clusterName)
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")

	clusterSecrets.client.DeleteSecret(secretPath)
}
