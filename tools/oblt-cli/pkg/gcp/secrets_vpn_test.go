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

// Package gcp provides the implementation to interact with Google Cloud Platform.
package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadVPNSvcCertificates(t *testing.T) {
	vpnSecrets := CreateVPNSecretsMock()
	defer vpnSecrets.Close()

	storedValue, err := vpnSecrets.ReadVPNSvcCertificates()
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.NotNil(t, storedValue, "the secret is not the expected one")
	assert.NotNil(t, storedValue.CaCert, "the secret is not the expected one")
	assert.NotNil(t, storedValue.CaKey, "the secret is not the expected one")
	assert.NotNil(t, storedValue.ServiceCert, "the secret is not the expected one")
	assert.NotNil(t, storedValue.ServiceKey, "the secret is not the expected one")
}

func TestReadVPNConfig(t *testing.T) {
	vpnSecrets := CreateVPNSecretsMock()
	defer vpnSecrets.Close()
	secret := defaultVpnConfig

	storedValue, err := vpnSecrets.ReadVPNConfig("pro")
	assert.NoError(t, err, "failed to get the secret: %v", err)
	assert.Equal(t, secret, storedValue, "the secret is not the expected one")
}
