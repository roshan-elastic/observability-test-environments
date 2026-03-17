// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package interactions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateClustersData(t *testing.T) {
	t.Run("Should run every 1 minute", func(t *testing.T) {
		assert.True(t, UpdateClustersData())
		time.Sleep(10 * time.Second)
		assert.False(t, UpdateClustersData(), "Should not update the data until 1 minute has passed")
		time.Sleep(1 * time.Minute)
		assert.True(t, UpdateClustersData(), "Should update the data after 1 minute has passed")
	})
}

func TestUpdateClusters(t *testing.T) {
	t.Run("Should find clusters", func(t *testing.T) {
		obltTestEnvironments, err := OnMemoryUserConf(t.Name(), "test", true)
		assert.NoError(t, err)
		UpdateClusters(*obltTestEnvironments)
		assert.True(t, len(Clusters.files) > 0, "Should find clusters")
		assert.True(t, len(ClustersExtra.data) > 0, "Should populate the extra data")
	})
}

func TestUpdateUsers(t *testing.T) {
	t.Run("Should find users", func(t *testing.T) {
		obltTestEnvironments, err := OnMemoryUserConf(t.Name(), "test", true)
		assert.NoError(t, err)
		UpdateClusters(*obltTestEnvironments)
		UpdateSlackUsers()
		assert.True(t, len(Users.data) > 0, "Should find users")
	})
}

func TestUpdateTemplates(t *testing.T) {
	t.Run("Should find templates", func(t *testing.T) {
		obltTestEnvironments, err := OnMemoryUserConf(t.Name(), "test", true)
		assert.NoError(t, err)
		UpdateTemplates(*obltTestEnvironments)
		assert.True(t, len(Templates.files) > 0, "Should find templates")
	})
}

func TestGoldenClustesr(t *testing.T) {
	t.Run("Should find golden clusters", func(t *testing.T) {
		obltTestEnvironments, err := OnMemoryUserConf(t.Name(), "test", true)
		assert.NoError(t, err)
		UpdateClusters(*obltTestEnvironments)
		goldenClusters := GoldenClusters()
		assert.True(t, len(goldenClusters) > 0, "Should find golden clusters")
	})
}
