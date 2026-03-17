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

// Package releases It contains the functions to interact with the Unified Releases.
package releases

import (
	"fmt"
)

// GetSnapshots searches for all the current snapshots in the Unified Release
// and parse them
func GetSnapshots() ([]Builds, error) {

	aliases, err := getAliases()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the aliases: %w", err)
	}

	snapshots, err := getSnapshotVersions(aliases)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the snapshot: %w", err)
	}

	return snapshots, nil
}

func getSnapshotVersions(versions []string) ([]Builds, error) {
	var ret []Builds
	for _, version := range versions {
		if isSnapshot(version) {
			snapshot, _ := getLatestBuildVersion(version)
			ret = append(ret, snapshot)
		}
	}
	return ret, nil
}
