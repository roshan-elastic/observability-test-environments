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
	"context"
	"os"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestLoginAppRoleModeError(t *testing.T) {
	ctx := context.Background()
	vaultContainer := fixtures.CreateVaultContainer(ctx, t)
	ip, _ := vaultContainer.Host(ctx)
	port, _ := vaultContainer.MappedPort(ctx, "8200")
	VaultAddr := "http://" + ip + ":" + port.Port()
	defer fixtures.TerminateContainer(t, ctx, vaultContainer)

	roleId, secretId := fixtures.CreateAppRole(ctx, t, vaultContainer)

	os.Setenv("VAULT_ROLE_ID", roleId.RoleID)
	os.Setenv("VAULT_SECRET_ID", "foo")
	os.Setenv("VAULT_TOKEN", "")
	os.Setenv("HOME", t.TempDir())
	vault := Vault{DryRun: false, VaultAddr: VaultAddr}
	err := vault.NewClient()
	assert.Error(t, err)

	os.Setenv("VAULT_SECRET_ID", secretId.SecretID)
	err = vault.NewClient()
	assert.NoError(t, err)
}

func TestLoginTokenModeError(t *testing.T) {
	ctx := context.Background()
	vaultContainer := fixtures.CreateVaultContainer(ctx, t)
	ip, _ := vaultContainer.Host(ctx)
	port, _ := vaultContainer.MappedPort(ctx, "8200")
	VaultAddr := "http://" + ip + ":" + port.Port()
	defer fixtures.TerminateContainer(t, ctx, vaultContainer)

	os.Setenv("VAULT_ROLE_ID", "")
	os.Setenv("VAULT_SECRET_ID", "")
	os.Setenv("VAULT_TOKEN", "foo")
	os.Setenv("HOME", t.TempDir())
	vault := Vault{DryRun: false, VaultAddr: VaultAddr}
	err := vault.NewClient()
	assert.NoError(t, err)
	_, err = vault.readSecretRaw(fixtures.SecretPath)
	assert.Error(t, err)
}

func TestLoginTokenFileModeError(t *testing.T) {
	ctx := context.Background()
	vaultContainer := fixtures.CreateVaultContainer(ctx, t)
	ip, _ := vaultContainer.Host(ctx)
	port, _ := vaultContainer.MappedPort(ctx, "8200")
	VaultAddr := "http://" + ip + ":" + port.Port()
	defer fixtures.TerminateContainer(t, ctx, vaultContainer)

	os.Setenv("VAULT_TOKEN", "")
	os.Setenv("VAULT_ROLE_ID", "")
	os.Setenv("VAULT_SECRET_ID", "")
	os.Setenv("HOME", t.TempDir())
	vault := Vault{DryRun: false, VaultAddr: VaultAddr}
	err := vault.NewClient()
	assert.Error(t, err)
}

func TestLoginTokenMode(t *testing.T) {
	ctx := context.Background()
	vaultContainer := fixtures.CreateVaultContainer(ctx, t)
	ip, _ := vaultContainer.Host(ctx)
	port, _ := vaultContainer.MappedPort(ctx, "8200")
	VaultAddr := "http://" + ip + ":" + port.Port()
	defer fixtures.TerminateContainer(t, ctx, vaultContainer)

	os.Setenv("HOME", t.TempDir())
	os.Setenv("VAULT_ROLE_ID", "")
	os.Setenv("VAULT_SECRET_ID", "")
	os.Setenv("VAULT_TOKEN", fixtures.VaultTokenTest)
	vault := Vault{DryRun: false, VaultAddr: VaultAddr}
	err := vault.NewClient()
	assert.NoError(t, err)

	data, err := vault.readSecretRaw(fixtures.SecretPath)
	assert.NoError(t, err)
	assert.Equal(t, "password123", data["top_secret"])
}
