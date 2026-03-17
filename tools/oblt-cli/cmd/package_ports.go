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
	"fmt"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/internal/apm"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters/k8s"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// PackagesPortsCmd represents the packages ports command
var PackagesPortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Command to forward the ports to local ports from the containers deployed in the cluster.",
	Long:  `Command to forward the ports to local ports from the containers deployed in the cluster.`,
	Run:   runPackagesPorts,
}

func init() {
	PackagesCmd.AddCommand(PackagesPortsCmd)

	PackagesPortsCmd.Flags().String(config.ClusterNameFlag, "", "Name of the cluster. (Required)")
	PackagesPortsCmd.Flags().StringSlice(config.PortsFlag, []string{}, "Ports to forward [LOCAL_PORT:]REMOTE_PORT[,...[LOCAL_PORT_N:]REMOTE_PORT_N] e.g. 8080:80,1234 (Required)")

	cobra.MarkFlagRequired(PackagesPortsCmd.Flags(), config.ClusterNameFlag)
	cobra.MarkFlagRequired(PackagesPortsCmd.Flags(), config.PortsFlag)
}

// runPorts forward the ports to local ports from the containers deployed in the cluster
func runPackagesPorts(cmd *cobra.Command, args []string) {
	userConfig, err := newObltConfig(cmd)
	config.CheckErr(err)
	checkOpperationLock(userConfig)

	clusterName, _ := cmd.Flags().GetString(config.ClusterNameFlag)
	ports, _ := cmd.Flags().GetStringSlice(config.PortsFlag)

	tx, ctx := apm.StartTransaction("runPackagesPorts", "request", []apm.Label{labelVersion}, userConfig)
	defer apm.Flush(tx)

	obltTestEnvironments, err := clusters.BootstrapRepository(dryRun)
	apm.CobraCheckErr(err, tx, ctx)
	config, err := obltTestEnvironments.FindClusterConfig(clusterName)
	apm.CobraCheckErr(err, tx, ctx)

	script := `oblt-fordward-ports`
	for _, port := range ports {
		script += fmt.Sprintf(" %s", port)
	}
	logger.Infof("Forwarding ports %s", script)
	k8sShell := k8s.NewK8sShell(config, script, userConfig.GetDir(), dryRun)
	err = k8sShell.ExecPackagesScript()
	apm.CobraCheckErr(err, tx, ctx)
}
