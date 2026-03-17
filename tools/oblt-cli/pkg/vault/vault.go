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

// Package vault It contains the functions to interact with Vault.
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/troubleshoot"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	vault "github.com/hashicorp/vault/api"
	yaml "gopkg.in/yaml.v3"
)

const (
	VaultAddr      = "https://secrets.elastic.co:8200"
	VaultTokenFile = ".vault-token"
	SecretBase     = "secret/observability-team/ci/test-clusters"
)

type Vault struct {
	client    *vault.Client
	VaultAddr string
	DryRun    bool
}

// NewClient initializes a new Vault client.
func (v *Vault) NewClient() error {
	logger.Debugf("Vault NewClient")
	logger.Warnf("Vault client usage is deprecated. Please change to configure Google Secret Manager.")
	config := vault.DefaultConfig()
	if v.VaultAddr == "" {
		config.Address = VaultAddr
	} else {
		config.Address = v.VaultAddr
	}

	client, err := vault.NewClient(config)
	if err == nil {
		v.client = client
		err = v.Login()
	}
	return err
}

// readSecretRaw retrieve a secret from a Vault service.
func (v *Vault) readSecretRaw(secretPath string) (data map[string]interface{}, err error) {
	logger.Debugf("readSecretRaw: %s", secretPath)

	err = v.NewClient()
	if err == nil {
		var secret *vault.Secret
		secret, err = v.client.Logical().Read(secretPath)

		if secret == nil && err == nil {
			return nil, fmt.Errorf("error reading the secret '%s': it could not exist", secretPath)
		} else if err == nil {
			data = secret.Data
		}
	}
	if err != nil {
		troubleshoot.Warn("You might have misconfigured your Vault environment.")
	}
	return data, err
}

// readSecret retrieve a secret from a Vault service, and get a field.
func (v *Vault) ReadSecret(secretPath string, field string) (value interface{}, err error) {
	logger.Debugf("ReadSecret: %s field: %s", secretPath, field)
	var data map[string]interface{}
	if data, err = v.readSecretRaw(secretPath); err == nil && data != nil {
		if data[field] == nil {
			if valueBytes, err := yaml.Marshal(data); err == nil {
				value = string(valueBytes)
			}
		} else {
			value = data[field].(string)
		}
	}
	return value, err
}

// Login login to Vault.
func (v *Vault) Login() (err error) {
	home, _ := os.UserHomeDir()
	tokenFilePath := filepath.Join(home, VaultTokenFile)
	roleID := os.Getenv("VAULT_ROLE_ID")
	secretID := os.Getenv("VAULT_SECRET_ID")
	token := os.Getenv("VAULT_TOKEN")

	if roleID != "" && secretID != "" {
		err = v.LoginAppRoleMode(roleID, secretID)
	} else if token != "" {
		err = v.LoginTokenMode(token)
	} else if _, err = os.Stat(tokenFilePath); err == nil {
		err = v.LoginTokenFileMode(tokenFilePath)
	} else {
		err = fmt.Errorf("no auth method provided")
	}
	return err
}

// Login login to Vault using the AppRole auth method.
func (v *Vault) LoginAppRoleMode(roleID, secretID string) (err error) {
	resp, err := v.client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	})
	if err == nil {
		logger.Debugf("Login: policies %v", resp.Auth.Policies)
		v.client.SetToken(resp.Auth.ClientToken)
	} else {
		err = fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	return err
}

// Login login to Vault using the Token auth method using the env var VAULT_TOKEN.
func (v *Vault) LoginTokenMode(token string) (err error) {
	v.client.SetToken(token)
	return err
}

// Login login to Vault using the Token auth method using the  ~/.vault-token file.
func (v *Vault) LoginTokenFileMode(filePath string) (err error) {
	token, err := os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("unable to read token file: %w", err)
	}
	v.LoginTokenMode(string(token))
	return err
}

// WriteSecret write a secret to Vault.
func (v *Vault) WriteSecret(secretPath string, data map[string]interface{}) (secret *vault.Secret, err error) {
	logger.Debugf("WriteSecret: %s", secretPath)
	err = v.NewClient()
	if err == nil {
		secret, err = v.client.Logical().Write(secretPath, data)
	}
	return secret, err
}
