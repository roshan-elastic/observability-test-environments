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

package vpn

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/stretchr/testify/assert"
)

func createOVPNFakeConfig(t *testing.T, tmpDir string) (tmpFilePath string) {
	var err error
	var tmpFile *os.File
	vpnPath := filepath.Join(tmpDir, ".oblt-cli", "vpn")
	os.MkdirAll(vpnPath, 0755)
	if tmpFile, err = os.Create(filepath.Join(vpnPath, "production.ovpn")); err == nil {
		io.WriteString(tmpFile, "#AviatrixController https://localhost\n#ClientLocalPort 5005\n")
		os.Chmod(tmpFile.Name(), 0777)
		tmpFilePath = tmpFile.Name()
	}
	assert.NoError(t, err)
	return tmpFilePath
}

func newVPNConfigFake(t *testing.T, tmpDir string) (vpnConfig VPNConfig) {
	os.Setenv("HOME", tmpDir)
	createOVPNFakeConfig(t, tmpDir)
	vpnConf := NewVPNConfig(tmpDir, tmpDir, true)
	return *vpnConf
}

func newAuthVPNFake(t *testing.T, tmpDir string) (authVPN *AuthVPN) {
	os.Setenv("HOME", tmpDir)
	configFile := createOVPNFakeConfig(t, tmpDir)
	auth := NewAuthVPN(configFile)
	auth.DryRun = true
	return auth
}

func CreateTemp(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "oblt-cli-test")
	assert.NoError(t, err)
	return tmpDir
}

func TestSaveCertificates(t *testing.T) {
	v := newAuthVPNFake(t, CreateTemp(t))
	v.gcsm = gcp.CreateVPNSecretsMock()
	v.saveCertificates()
	certs, _ := v.gcsm.ReadVPNSvcCertificates()

	// Read the files and compare with the expected content
	content, err := os.ReadFile(v.getSvcCertPath())
	assert.NoError(t, err)
	assert.Equal(t, certs.ServiceCert, string(content))

	content, err = os.ReadFile(v.getSvcKeyPath())
	assert.NoError(t, err)
	assert.Equal(t, certs.ServiceKey, string(content))

	content, err = os.ReadFile(v.getCACertPath())
	assert.NoError(t, err)
	assert.Equal(t, certs.CaCert, string(content))

	content, err = os.ReadFile(v.getCAKeyPath())
	assert.NoError(t, err)
	assert.Equal(t, certs.CaKey, string(content))
}

func TestRunServer(t *testing.T) {
	v := newAuthVPNFake(t, CreateTemp(t))
	v.gcsm = gcp.CreateVPNSecretsMock()
	v.saveCertificates()
	go v.RunServer()

	time.Sleep(time.Second * 10)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://localhost:5005")
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "SuccessAviatrix", string(body))
}

func TestOpenVPN(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	v := newAuthVPNFake(t, CreateTemp(t))
	v.DryRun = true
	v.OpenVPN()

	_, err := os.Stat(v.AuthFile)
	assert.NoError(t, err)
}

func TestWaitForCertImported(t *testing.T) {
	var waitGroup sync.WaitGroup
	v := newAuthVPNFake(t, CreateTemp(t))
	waitGroup.Add(1)
	go waitForCertImported(t, v, &waitGroup)
	time.Sleep(time.Second * 2)
	os.WriteFile(v.getConfigureMark(), []byte("ok"), 0777)
	time.Sleep(time.Second * 5)
}

func waitForCertImported(t *testing.T, v *AuthVPN, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	err := v.waitForCertImported()
	assert.NoError(t, err)
}

func TestWaitForConnection(t *testing.T) {
	var waitGroup sync.WaitGroup
	v := newAuthVPNFake(t, CreateTemp(t))
	waitGroup.Add(1)
	go waitForConnection(t, v, &waitGroup)
	time.Sleep(time.Second * 2)
	os.WriteFile(v.GetLogFile(), []byte("Initialization Sequence Completed"), 0777)
	waitGroup.Wait()
}

func TestWaitForConnectionTimeout(t *testing.T) {
	v := newAuthVPNFake(t, CreateTemp(t))
	err := v.WaitForConnection()
	assert.Error(t, err)
}

func waitForConnection(t *testing.T, v *AuthVPN, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	err := v.WaitForConnection()
	assert.NoError(t, err)
}

func TestWaitForCredentials(t *testing.T) {
	var waitGroup sync.WaitGroup
	v := newAuthVPNFake(t, CreateTemp(t))
	waitGroup.Add(1)
	go waitForCredentials(t, v, &waitGroup)
	time.Sleep(time.Second * 2)
	v.mutex.Lock()
	v.Token = "token"
	v.Email = "email"
	v.mutex.Unlock()
	waitGroup.Wait()
}

func waitForCredentials(t *testing.T, v *AuthVPN, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	err := v.WaitForCredentials()
	assert.NoError(t, err)
}

func TestWaitForCredentialsTimeout(t *testing.T) {
	v := newAuthVPNFake(t, CreateTemp(t))
	err := v.WaitForCredentials()
	assert.Error(t, err)
}

func TestVPN(t *testing.T) {
	t.Skip("Unresolved issue: panic: http: multiple registrations for /")
	tmpDir := CreateTemp(t)
	auth := newAuthVPNFake(t, tmpDir)
	auth.gcsm = gcp.CreateVPNSecretsMock()
	auth.saveCertificates()

	// fake the auth
	auth.mutex.Lock()
	auth.Token = "token"
	auth.Email = "email"
	auth.mutex.Unlock()
	// fake connect to the VPN
	os.WriteFile(auth.GetLogFile(), []byte("Initialization Sequence Completed"), 0777)

	vpnConfig := newVPNConfigFake(t, tmpDir)
	vpnConfig.Auth = auth
	err := vpnConfig.Connect()
	assert.NoError(t, err)
}
