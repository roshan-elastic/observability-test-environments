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

package cmd

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/vpn"
	"github.com/spf13/cobra"
)

// ServicesVpnCmd represents the create command
var ServicesVpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Command to connect to the Elastic VPN.",
	Long: `Command to connect to the Elastic VPN.

The following example connect to the production VPN:

	oblt-cli services vpn`,
	Args: validateEnvironmentFlag,
	Run:  runServicesVpn,
}

func init() {
	ServicesCmd.AddCommand(ServicesVpnCmd)
}

/*
runServicesVpn Open a VPN connection to the Elastic VPN,
*/
func runServicesVpn(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	tx, ctx := apm.StartTransaction("runServicesVpn", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)
	checkOpperationLock(userConfig)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	initialChecks(obltTestEnvironments, tx, ctx)

	vpnClient := vpn.NewVPNConfig(userConfig.GetDir(), obltTestEnvironments.GetPath(), dryRun)
	err = vpnClient.Connect()
	apm.CobraCheckErr(err, tx, ctx)
}
