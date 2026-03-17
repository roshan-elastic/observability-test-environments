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
// This file contains the functions to access Google Cloud Secret Manager.
package gcp

import (
	"path"
)

// SecretsManagerMock is a mock implementation of the SecretsManager interface for testing purposes.
type SecretsManagerMock struct {
	store map[string]string
}

// NewSecretsManagerMock returns a new instance of SecretsManagerMock.
func NewSecretsManagerMock() *SecretsManagerMock {
	return &SecretsManagerMock{
		store: make(map[string]string),
	}
}

// GetSecret returns the value of the secret.
func (s *SecretsManagerMock) GetSecret(secretPath string) (value string, err error) {
	return s.store[secretPath], nil
}

// CreateSecret creates a new secret with the given value.
func (s *SecretsManagerMock) CreateSecret(secretPath string, value interface{}) (err error) {
	data, err := marshallInterface(value)
	if err == nil {
		s.store[secretPath] = data
	}
	return err
}

// UpdateSecret updates the value of the secret.
func (s *SecretsManagerMock) UpdateSecret(secretPath, value string) (err error) {
	s.store[secretPath] = value
	return nil
}

// DeleteSecret deletes the secret.
func (s *SecretsManagerMock) DeleteSecret(secretPath string) (err error) {
	delete(s.store, secretPath)
	return nil
}

// ListSecrets returns the list of secrets.
// The filter parameter is not used in this mock implementation.
func (s *SecretsManagerMock) ListSecrets(secretPath string, filter string) (secretsPaths []string, err error) {
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys, nil
}

// Close closes the client.
func (s *SecretsManagerMock) Close() (err error) {
	return nil
}

// GetSecretParent returns the parent of the secret.
func (s *SecretsManagerMock) GetSecretParent(secretPath string) (parent string) {
	return path.Dir(path.Dir(secretPath))
}

// GetSecretName returns the name of the secret.
func (s *SecretsManagerMock) GetSecretName(secretPath string) (name string) {
	return path.Base(secretPath)
}
