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

// Package fixtures contains the functions to create and terminate the fixtures for the tests.
package fixtures

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/modules/vault"
)

const (
	VaultTokenTest    = "1234567890"
	VaultRoleIDTest   = "foo"
	VaultSecretIDTest = "bar"
	SecretPath        = "secret/testing1"
)

type RoleID struct {
	RoleID string `json:"role_id"`
}

type SecretID struct {
	SecretID string `json:"secret_id"`
}

// TerminateContainer terminates a Vault container
func TerminateContainer(t *testing.T, ctx context.Context, container *vault.VaultContainer) {
	assert.NoError(t, container.Terminate(ctx))
}

// CreateVaultContainer creates a new Vault container
func CreateVaultContainer(ctx context.Context, t *testing.T) *vault.VaultContainer {
	vaultContainer, err := vault.RunContainer(ctx, vault.WithToken(VaultTokenTest), vault.WithInitCommand(
		"auth enable approle",
		"secrets disable secret",
		"secrets enable -version=1 -path=secret kv",
		"write --force auth/approle/role/"+VaultRoleIDTest,
		"write "+SecretPath+" top_secret=password123",
	))
	assert.NoError(t, err)
	return vaultContainer
}

// CreateAppRole creates a new AppRole in Vault and returns the RoleID and SecretID
func CreateAppRole(ctx context.Context, t *testing.T, vaultContainer *vault.VaultContainer) (roleId RoleID, secretId SecretID) {
	cmds := []string{
		"vault", "read", "-format", "json", "-field", "data", "auth/approle/role/" + VaultRoleIDTest + "/role-id",
	}
	_, result, err := vaultContainer.Exec(ctx, cmds, tcexec.Multiplexed())
	assert.NoError(t, err)
	str, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(str, &roleId))

	cmds = []string{
		"vault", "write", "--force", "-format", "json", "-field", "data", "auth/approle/role/" + VaultRoleIDTest + "/secret-id",
	}
	_, result, err = vaultContainer.Exec(ctx, cmds, tcexec.Multiplexed())
	assert.NoError(t, err)
	str, err = io.ReadAll(result)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(str, &secretId))
	return roleId, secretId
}
