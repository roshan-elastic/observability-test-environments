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
// This file contain the access to VPN secrets in GCP.
package gcp

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/certs"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"gopkg.in/yaml.v3"
)

const (
	defaultVpnConfig = `#AviatrixController https://localhost
#ClientLocalPort 5005
`
)

// VpnSecrets struct to manage VPN related secrets
type VpnSecrets struct {
	client SecretsManager
	auth   Auth
}

// VPNSvcCert struct to store the callback server certificated
type VPNSvcCert struct {
	CaCert      string `json:"ca_cert" yaml:"ca_cert"`
	CaKey       string `json:"ca_key" yaml:"ca_key"`
	ServiceCert string `json:"service_cert" yaml:"service_cert"`
	ServiceKey  string `json:"service_key" yaml:"service_key"`
}

func NewVPNSvcCert() VPNSvcCert {
	return VPNSvcCert{
		CaCert:      "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----\n",
		CaKey:       "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
		ServiceCert: "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----\n",
		ServiceKey:  "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
	}
}

// NewVpnSecrets returns a new instance of VpnSecrets.
func NewVpnSecrets() (vpnSecrets *VpnSecrets, err error) {
	return NewVpnSecretsWithProjectId(default_project_id)
}

// NewVpnSecretsWithProjectId returns a new instance of VpnSecrets with the given project ID.
func NewVpnSecretsWithProjectId(projectId string) (vpnSecrets *VpnSecrets, err error) {
	auth := NewGcpAuthWithProjectId(projectId)
	auth.AuthenticateInteractive()
	gcsmClient, err := NewClientWithOpts(auth.Options...)
	if err != nil {
		logger.Errorf("failed to create secretmanager client: %v", err)
	} else {
		vpnSecrets = NewVpnSecretsWithClient(gcsmClient, &auth)
	}
	return vpnSecrets, err
}

// NewVpnSecretsWithClient returns a new instance of VpnSecrets with the given client.
func NewVpnSecretsWithClient(client SecretsManager, auth Auth) (vpnSecrets *VpnSecrets) {
	return &VpnSecrets{
		client: client,
		auth:   auth,
	}
}

// SecretVpnUnmarsahler interfaz to help unmarshall all VPN related secrets
type SecretVpnUnmarsahler interface {
	VPNSvcCert
}

// unmarshallVpnSecret read an unmarshall VPN related secrets
func unmarshallVpnSecret[T SecretVpnUnmarsahler](client SecretsManager, secretsPath string, value *T) (err error) {
	data, err := client.GetSecret(secretsPath)
	if err == nil {
		err = yaml.Unmarshal([]byte(data), value)
	}
	if err != nil {
		logger.Errorf("failed to read the secret %s: %v", secretsPath, err)
	}
	return err
}

// Close closes the secrets manager client.
func (v *VpnSecrets) Close() error {
	return v.client.Close()
}

// getProjectPath returns the project path.
func (v *VpnSecrets) getProjectPath() string {
	return "projects/" + v.auth.GetProjectId()
}

// getSecretsRootPath returns the root path for the secrets.
func (v *VpnSecrets) getSecretsRootPath() string {
	return v.getProjectPath() + "/secrets"
}

// getVpnConfigPath returns the path to the VPN configuration file for the given environment.
func (v *VpnSecrets) getVpnConfigPath(environment string) string {
	return v.getSecretsRootPath() + "/elastic-vpn-" + environment
}

// getVpnSvcCertPath returns the path to the VPN callback server certificates.
func (v *VpnSecrets) getVpnSvcCertPath() string {
	return v.getSecretsRootPath() + "/oblt-cli-vpn-certificates"
}

// ReadVPNSvcCertificates read the callback server certificates
func (v *VpnSecrets) ReadVPNSvcCertificates() (value VPNSvcCert, err error) {
	err = unmarshallVpnSecret(v.client, v.getVpnSvcCertPath(), &value)
	return value, err
}

// ReadVPNConfig read the VPN configuration file for the given environment
func (v *VpnSecrets) ReadVPNConfig(env string) (data string, err error) {
	return v.client.GetSecret(v.getVpnConfigPath(env))
}

// CreateVPNSecretsMock creates a new ClusterSecrets instance using a mock client.
func CreateVPNSecretsMock() (vnpSevrets *VpnSecrets) {
	key, pemKey, _ := certs.CreateKey()
	pemCert, _ := certs.CreateCertificate(key)
	client := NewSecretsManagerMock()
	auth := NewAuthMock()
	mockCerts := NewVPNSvcCert()
	mockCerts.CaCert = pemCert
	mockCerts.CaKey = pemKey
	mockCerts.ServiceCert = pemCert
	mockCerts.ServiceKey = pemKey
	vnpSevrets = NewVpnSecretsWithClient(client, &auth)
	client.CreateSecret(vnpSevrets.getVpnSvcCertPath(), mockCerts)
	client.CreateSecret(vnpSevrets.getVpnConfigPath("pro"), defaultVpnConfig)
	client.CreateSecret(vnpSevrets.getVpnConfigPath("staging"), defaultVpnConfig)
	return vnpSevrets
}
