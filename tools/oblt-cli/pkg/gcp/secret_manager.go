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
	"context"
	"path"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

// SecretsManager is the interface to interact with the Google Cloud Secret Manager.
type SecretsManager interface {
	GetSecret(secretPath string) (value string, err error)
	CreateSecret(secretPath string, value interface{}) (err error)
	UpdateSecret(secretPath, value string) (err error)
	DeleteSecret(secretPath string) (err error)
	ListSecrets(secretPath string, filter string) (secretsPaths []string, err error)
	Close() (err error)
	GetSecretName(secretPath string) (name string)
	GetSecretParent(secretPath string) (parent string)
}

// GCPSecretsManager is the implementation of the SecretsManager interface for Google Cloud Secret Manager.
type GCPSecretsManager struct {
	client *secretmanager.Client
}

// NewClient returns a new instance of GCPSecretsManager.
func NewClient() (clientGCSM SecretsManager, err error) {
	return NewClientWithOpts([]option.ClientOption{}...)
}

func NewClientWithOpts(opts ...option.ClientOption) (clientGCSM SecretsManager, err error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		logger.Errorf("failed to create secretmanager client: %v", err)
	} else {
		clientGCSMImpl := GCPSecretsManager{}
		clientGCSMImpl.client = client
		clientGCSM = &clientGCSMImpl
	}
	return clientGCSM, err
}

// GetSecret returns the value of the secret at the given path. It takes the latest version.
func (gcsm *GCPSecretsManager) GetSecret(secretPath string) (value string, err error) {
	logger.Debugf("Reading %s secret:", secretPath)
	ctx := context.Background()
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath + "/versions/latest",
	}
	resp, err := gcsm.client.AccessSecretVersion(ctx, req)
	if err != nil {
		logger.Errorf("failed to access secret version: %v", err)
		value = ""
	} else {
		value = string(resp.Payload.Data)
	}
	return value, err
}

// createSecret creates a new secret with the given value.
func (gcsm *GCPSecretsManager) createSecret(secretPath, value string) (err error) {
	logger.Debugf("Writing %s secret:", secretPath)
	ctx := context.Background()
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   gcsm.GetSecretParent(secretPath),
		SecretId: gcsm.GetSecretName(secretPath),
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}
	_, err = gcsm.client.CreateSecret(ctx, req)
	if err != nil {
		logger.Errorf("failed to create secret: %v", err)
	} else {
		err = gcsm.UpdateSecret(secretPath, value)
	}
	return err
}

// UpdateSecret updates the value of the secret at the given path.
func (gcsm *GCPSecretsManager) UpdateSecret(secretPath, value string) (err error) {
	logger.Debugf("Updating %s secret:", secretPath)
	ctx := context.Background()
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretPath,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(value),
		},
	}
	_, err = gcsm.client.AddSecretVersion(ctx, req)
	if err != nil {
		logger.Errorf("failed to add secret version: %v", err)
	}
	return err
}

// DeleteSecret deletes the secret at the given path.
func (gcsm *GCPSecretsManager) DeleteSecret(secretPath string) (err error) {
	logger.Debugf("Deleting %s secret:", secretPath)
	ctx := context.Background()
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: secretPath,
	}
	err = gcsm.client.DeleteSecret(ctx, req)
	if err != nil {
		logger.Errorf("failed to delete secret: %v", err)
	}
	return err
}

// getSecretParent returns the parent of the secret path.
func (gcsm *GCPSecretsManager) GetSecretParent(secretPath string) (parent string) {
	// TODO it has to return "projects/PROJECT_ID"
	return path.Dir(path.Dir(secretPath))
}

// getSecretName returns the name of the secret from the secret path.
func (gcsm *GCPSecretsManager) GetSecretName(secretPath string) (name string) {
	return path.Base(secretPath)
}

// Close closes the client.
func (gcsm *GCPSecretsManager) Close() (err error) {
	err = gcsm.client.Close()
	if err != nil {
		logger.Errorf("failed to close secretmanager client: %v", err)
	}
	return err
}

// ListSecrets returns the list of secrets at the given path.
// It returns the first 100 secrets.
// for details about the filters see https://cloud.google.com/secret-manager/docs/filtering
func (gcsm *GCPSecretsManager) ListSecrets(secretPath, filter string) (secretsPaths []string, err error) {
	logger.Debugf("Listing %s secrets with filter %s", secretPath, filter)
	ctx := context.Background()
	secretsPaths = []string{}
	req := &secretmanagerpb.ListSecretsRequest{
		Parent:   secretPath,
		Filter:   filter,
		PageSize: 100,
	}
	it := gcsm.client.ListSecrets(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			logger.Errorf("failed to list secrets: %v", err)
			break
		} else {
			secretsPaths = append(secretsPaths, resp.Name)
		}
	}
	return secretsPaths, err
}

// CreateSecret creates a new secret with the given value.
// It takes the value as an interface and it writes it as a YAML string.
// if the interface is a string, it writes it directly.
func (gcsm *GCPSecretsManager) CreateSecret(secretPath string, value interface{}) (err error) {
	// in case of String we write it directly
	data, err := marshallInterface(value)
	if err == nil {
		gcsm.createSecret(secretPath, data)
	} else {
		logger.Errorf("failed to create secret: %v", err)
	}
	return err
}

// marshallInterface returns the string representation of the interface.
func marshallInterface(value interface{}) (data string, err error) {
	var dataBytes []byte
	data, ok := value.(string)
	if !ok {
		if dataBytes, err = yaml.Marshal(value); err == nil {
			data = string(dataBytes)
		}
	}
	return data, err
}
