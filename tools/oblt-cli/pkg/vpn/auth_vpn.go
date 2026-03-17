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

// inspired by https://github.com/christophgysin/avpnc/blob/master/avpnc
package vpn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/cmd"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

const (
	connectionTimeOut = 90
)

// AuthVPN is the object to store the information from the OpenVPN config file.
type AuthVPN struct {
	Token      string `json:"Token"`
	Email      string `json:"Email"`
	OvpnFile   string `json:"OvpnFile"`
	AuthUrl    string `json:"Url"`
	ListenPort string `json:"ListenPort"`
	ScriptFile string `json:"scriptFile"`
	AuthFile   string `json:"authFile"`
	DryRun     bool   `json:"dryRun"`
	gcsm       *gcp.VpnSecrets
	mutex      sync.Mutex
}

type checkCondition func() bool

// NewAuthVPN creates a new AuthVPN object from a OpenVPN config file.
func NewAuthVPN(ovpnFile string) (auth *AuthVPN) {
	auth = &AuthVPN{}
	auth.mutex.Lock()
	defer auth.mutex.Unlock()
	logger.Debugf("Loading Open VPN config: %s", ovpnFile)
	data, err := os.ReadFile(ovpnFile)
	if err == nil {
		config := string(data)
		lines := strings.Split(config, "\n")
		auth.OvpnFile = ovpnFile
		for _, line := range lines {
			if strings.HasPrefix(line, "#AviatrixController") {
				auth.AuthUrl = strings.TrimSpace(strings.Split(line, " ")[1])
			}
			if strings.HasPrefix(line, "#ClientLocalPort") {
				auth.ListenPort = strings.TrimSpace(strings.Split(line, " ")[1])
			}
		}
	}
	if err = auth.createVpnDir(); err != nil {
		logger.Fatal(err)
	}
	os.Create(auth.GetLogFile())
	os.OpenFile(auth.GetPidFile(), os.O_RDWR|os.O_CREATE, 0666)
	return auth
}

func NewAuthVPNForEnv(env string) (auth *AuthVPN, err error) {
	logger.Debugf("Loading Open VPN config for %s", env)
	v := AuthVPN{}
	err = errors.Join(v.saveCertificates(),
		v.saveVPNConfig(env))
	return NewAuthVPN(v.GetOvpnFile()), err
}

// getRoot Process the response from the browser after authenticate.
func (v *AuthVPN) getRoot(w http.ResponseWriter, r *http.Request) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	logger.Debugf("got / request: %v\n", r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, "SuccessAviatrix")
	// remove the first character which is /
	data := r.URL.Path[1:]
	json.Unmarshal([]byte(data), v)
}

func (v *AuthVPN) createVpnDir() (err error) {
	err = os.MkdirAll(v.getVPNConfigFolder(), 0755)
	return err
}

// createAuthFile Create the file with the VPN credentials
func (v *AuthVPN) createAuthFile() (tmpFile *os.File, err error) {
	authFile := filepath.Join(v.getVPNConfigFolder(), "vpn-auth")
	if tmpFile, err = os.OpenFile(authFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err == nil {
		defer tmpFile.Close()
		io.WriteString(tmpFile, v.Email+"\n"+v.Token+"\n")
		v.AuthFile = tmpFile.Name()
	}
	return tmpFile, err
}

// getVPNConfigFolder return the Path to the folder used for VPN files.
func (v *AuthVPN) getVPNConfigFolder() (vpnHome string) {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".oblt-cli", "vpn")
}

// getPidFile return the path a file to save PID of the OpenVPN process.
func (v *AuthVPN) GetPidFile() (pidFile string) {
	return filepath.Join(v.getVPNConfigFolder(), "openvpn.pid")
}

// getLogFile return the path a file to save logs of the OpenVPN process.
func (v *AuthVPN) GetLogFile() (logFile string) {
	return filepath.Join(v.getVPNConfigFolder(), "openvpn.log")
}

// getSvcCertPath return the path to the VPN service certificate.
func (v *AuthVPN) getSvcCertPath() (svcCertPath string) {
	return filepath.Join(v.getVPNConfigFolder(), "tls.crt")
}

// getSvcKeyPath return the path to the VPN service key.
func (v *AuthVPN) getSvcKeyPath() (svcKeyPath string) {
	return filepath.Join(v.getVPNConfigFolder(), "tls.key")
}

// getCACertPath return the path to the VPN CA certificate.
func (v *AuthVPN) getCACertPath() (svcCertPath string) {
	return filepath.Join(v.getVPNConfigFolder(), "rootCA.pem")
}

// getCAKeyPath return the path to the VPN CA key.
func (v *AuthVPN) getCAKeyPath() (svcKeyPath string) {
	return filepath.Join(v.getVPNConfigFolder(), "rootCA.key")
}

// OpenVPN starts the OpenVPN client
func (v *AuthVPN) OpenVPN() {
	logger.Infof("Starting OpenVPN")
	logger.Warnf("We need elevated privileges to run OpenVPN, so you will be asked for your password.")
	var err error
	if _, err = v.createAuthFile(); err == nil {
		err = cmd.RunBashScript(v.scriptRunOpenVPN(), v.DryRun)
	}
	if err != nil {
		logOutput, _ := os.ReadFile(v.GetLogFile())
		logger.Infof(string(logOutput))
		logger.Errorf("OpenVPN failed to start")
		logger.Fatal(err)
	}
}

// Authenticate opens the Authentication URL in a browser.
func (v *AuthVPN) Authenticate() (err error) {
	logger.Debugf("browser.open: %s", v.AuthUrl)
	err = cmd.OpenBrowser(v.AuthUrl, v.DryRun)
	if err != nil {
		logger.Fatal(err)
	}
	err = v.WaitForCredentials()
	return err
}

// RunServer starts the server to listen for the response from the browser
func (v *AuthVPN) RunServer() {
	logger.Debugf("Starting Local Server : https://localhost:%s/", v.ListenPort)
	crt := v.getSvcCertPath()
	key := v.getSvcKeyPath()
	go http.HandleFunc("/", v.getRoot)
	err := http.ListenAndServeTLS(":"+v.ListenPort, crt, key, nil)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}

// getGcsm get an instance of VpnSecrets
func (v *AuthVPN) getGcsm() (vpnSecrets *gcp.VpnSecrets, err error) {
	if v.gcsm == nil {
		v.gcsm, err = gcp.NewVpnSecrets()
	}
	return v.gcsm, err
}

// mkDirVpnConfig creates the VPN config folder if it s needed
func (v *AuthVPN) mkDirVpnConfig() error {
	return os.MkdirAll(v.getVPNConfigFolder(), 0755)
}

// existPath check if a file or directory exits
func (v *AuthVPN) existPath(pathToCheck string) (bool, error) {
	_, err := os.Stat(pathToCheck)
	existsPath := !os.IsNotExist(err)
	if !existsPath {
		err = v.mkDirVpnConfig()
	}
	return existsPath, err
}

// ReadVPNSvcCertificates wrapper to instantiate the VPN secrets manager and retrieve the callback server certificates
func (v *AuthVPN) ReadVPNSvcCertificates() (certs gcp.VPNSvcCert, err error) {
	var gcsm *gcp.VpnSecrets
	if gcsm, err = v.getGcsm(); err == nil {
		certs, err = gcsm.ReadVPNSvcCertificates()
	}
	return certs, err
}

// ReadVPNConfig wrapper to instantiate the VPN secrets manager and retrieve the VPN config file for the given environment
func (v *AuthVPN) ReadVPNConfig() (config string, err error) {
	var gcsm *gcp.VpnSecrets
	if gcsm, err = v.getGcsm(); err == nil {
		config, err = gcsm.ReadVPNConfig("production")
	}
	return config, err
}

// saveCertificates saves the VPN certificates and keys.
func (v *AuthVPN) saveCertificates() (err error) {
	var exist bool
	var certs gcp.VPNSvcCert
	if exist, err = v.existPath(v.getCACertPath()); err == nil && !exist {
		if certs, err = v.ReadVPNSvcCertificates(); err == nil {
			err = errors.Join(
				os.WriteFile(v.getSvcCertPath(), []byte(certs.ServiceCert), 0644),
				os.WriteFile(v.getSvcKeyPath(), []byte(certs.ServiceKey), 0644),
				os.WriteFile(v.getCACertPath(), []byte(certs.CaCert), 0644),
				os.WriteFile(v.getCAKeyPath(), []byte(certs.CaKey), 0644),
				v.importCACert(),
			)
		}
	}
	return err
}

// saveVPNConfig saves the VPN configuration file for the environment.
func (v *AuthVPN) saveVPNConfig(env string) (err error) {
	var exist bool
	var config string
	ovpnFile := v.GetOvpnFile()
	if exist, err = v.existPath(ovpnFile); err == nil && !exist {
		if config, err = v.ReadVPNConfig(); err == nil {
			err = os.WriteFile(ovpnFile, []byte(config), 0644)
		}
	}
	return err
}

// GetOvpnFile return the path to the VPN configuration file for the environment.
func (v *AuthVPN) GetOvpnFile() (ovpnFile string) {
	return filepath.Join(v.getVPNConfigFolder(), "production.ovpn")
}

// importCACert imports the CA certificate to the keychain.
func (v *AuthVPN) importCACert() (err error) {
	logger.Warnf("We will import a new CA certificate to your keychain.")
	err = cmd.RunBashScript(v.scriptImportCACert(), v.DryRun)
	if !v.DryRun && err == nil {
		err = v.waitForCertImported()
	}
	return err
}

// waitForCertImported waits for the CA certificate to be imported.
func (v *AuthVPN) waitForCertImported() (err error) {
	for {
		if _, err = os.Stat(v.getConfigureMark()); os.IsNotExist(err) {
			logger.Infof("Waiting for CA certificate to be imported.")
		} else {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return err
}

// WaitForCredentials waits for the credentials to be available.
func (v *AuthVPN) WaitForCredentials() (err error) {
	return v.timeout(v.checkCredentialsLoaded, "VPN credentials timeout")
}

// checkCredentialsLoaded checks if the credentials are available.
func (v *AuthVPN) checkCredentialsLoaded() (ret bool) {
	ret = false
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if v.Token != "" && v.Email != "" {
		logger.Infof("VPN credentials received")
		ret = true
	}
	return ret
}

// timeout waits for a condition to be true or timeout.
func (v *AuthVPN) timeout(condition checkCondition, timeOutMessage string) (err error) {
	timeOutTime := connectionTimeOut * time.Second
	if v.DryRun {
		timeOutTime = 10 * time.Second
	}
	for iteration := 0 * time.Second; ; iteration++ {
		currentTime := iteration * time.Second
		if condition() {
			break
		}
		if currentTime > timeOutTime {
			err = errors.New(timeOutMessage)
			break
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

// WaitForConnection waits for the VPN connection to be established.
func (v *AuthVPN) WaitForConnection() (err error) {
	return v.timeout(v.checkConnectionStablished, "VPN connection timeout")
}

// checkConnectionStablished checks if the VPN connection is established.
func (v *AuthVPN) checkConnectionStablished() (ret bool) {
	ret = false
	logFile := v.GetLogFile()
	if _, err := os.Stat(logFile); err == nil {
		logContent, _ := os.ReadFile(logFile)
		if strings.Contains(string(logContent), "Initialization Sequence Completed") {
			logger.Infof("VPN connection established")
			ret = true
		}
	}
	return ret
}

// scriptImportCACert return the script to import the CA certificate to the keychain.
func (v *AuthVPN) scriptImportCACert() (script string) {
	return fmt.Sprintf(`#!/usr/bin/env bash
		. "${OBLT_ELASTIC_SCRIPT}"
		elastic::vpn-import-certs-macOS %s %s
		`, v.getCACertPath(), v.getConfigureMark())
}

// scriptRunOpenVPN return the script to run OpenVPN.
func (v *AuthVPN) scriptRunOpenVPN() (script string) {
	return fmt.Sprintf(`#!/usr/bin/env bash
		. "${OBLT_ELASTIC_SCRIPT}"
		elastic::vpn-run %s %s %s %s
		`, v.OvpnFile, v.AuthFile, v.GetLogFile(), v.GetPidFile())
}

// getConfigureMark return the path to the file used to mark the CA certificate as imported.
func (v *AuthVPN) getConfigureMark() (configureMark string) {
	return filepath.Join(v.getVPNConfigFolder(), "configure.mark")
}
