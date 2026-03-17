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
	"os"
	"path/filepath"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/cmd"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

const (
	scriptHeader = `#!/usr/bin/env bash
. "${OBLT_INTALL_SCRIPT}"
. "${OBLT_ELASTIC_SCRIPT}"
install::vpn
`
	scriptPlaceholder = "#!/usr/bin/env bash\nfunction install::vpn(){\ntrue\n};\nfunction elastic::vpn(){\ntrue\n};"
	installLib        = ".ci/scripts/lib/install.sh"
	elasticLib        = ".ci/scripts/lib/elastic.sh"
)

// VPNConfig is the object to store the information to make the VPN connection.
type VPNConfig struct {
	environment string
	Script      string
	repoPath    string
	obltCliHome string
	dryRun      bool
	Auth        *AuthVPN
}

// NewVPNConfig creates a new VPNConfig object.
func NewVPNConfig(obltCliHome, repoPath string, dryRun bool) *VPNConfig {
	return &VPNConfig{
		Script:      "",
		repoPath:    repoPath,
		obltCliHome: obltCliHome,
		dryRun:      dryRun,
	}
}

// prerequisites runs the prerequisites for the VPN.
func (c *VPNConfig) prerequisites() (err error) {
	os.Setenv("OBLT_INTALL_SCRIPT", filepath.Join(c.repoPath, installLib))
	os.Setenv("OBLT_ELASTIC_SCRIPT", filepath.Join(c.repoPath, elasticLib))
	os.Setenv("OBLT_CLI_HOME", c.obltCliHome)
	os.Setenv("OBLT_VPN_ENVIRONMENT", c.environment)
	if logger.Verbose {
		os.Setenv("DEBUG", "true")
	}
	cmd.RunBashScript(scriptHeader+c.Script, c.dryRun)
	if c.Auth == nil {
		c.Auth, err = NewAuthVPNForEnv(c.environment)
	}
	return err
}

// Connect connects to the VPN.
func (c *VPNConfig) Connect() (err error) {
	if err = c.prerequisites(); err == nil {
		logger.Infof("Starting")
		go c.Auth.RunServer()
		if err = c.Auth.Authenticate(); err == nil {
			c.Auth.OpenVPN()
			if err = c.Auth.WaitForConnection(); err == nil {
				logger.Infof("Done")
			}
		}
	}
	return err
}
