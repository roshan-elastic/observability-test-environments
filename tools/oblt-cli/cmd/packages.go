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
	"github.com/spf13/cobra"
)

// ClusterCmd represents the cluster command
var PackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Command to manage integration packages.",
	Long: `Command to manage integration packages.
	With this command you can create, update and destroy oblt clusters running integrations.
	You can access to the component of the clusters for debugging purposes.`,
}

func init() {
	RootCmd.AddCommand(PackagesCmd)
}
