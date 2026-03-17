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

package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Repository(t *testing.T) {
	t.Run("Git Protocol", func(t *testing.T) {
		repo := Repository{
			Branch:  "main",
			Owner:   "elastic",
			Name:    "foo",
			SSHMode: true,
		}

		assert.Equal(t, "git@github.com:elastic/foo.git", repo.String())
	})

	t.Run("HTTPS Protocol", func(t *testing.T) {
		repo := Repository{
			Branch:  "main",
			Owner:   "elastic",
			Name:    "foo",
			SSHMode: false,
		}

		assert.Equal(t, "https://github.com/elastic/foo.git", repo.String())
	})
}

func Test_LabelsClusterType(t *testing.T) {
	t.Run("Should return an empty list", func(t *testing.T) {
		labels := LabelsClusterType(nil)
		assert.True(t, len(labels) == 0)
	})

	t.Run("Should return on entry if duplicated states", func(t *testing.T) {
		labels := LabelsClusterType([]string{"M\tenvironments/users/foo/cluster.yml", "M\tenvironments/users/bar/config-cluster.yml"})
		assert.True(t, len(labels) == 1)
	})

	t.Run("Should return an destroy label when deleted", func(t *testing.T) {
		labels := LabelsClusterType([]string{"D\tenvironments/users/foo/deleted.yml"})
		assert.Equal(t, "cluster:destroy", labels[0])
	})

	t.Run("Should return an update label when changed", func(t *testing.T) {
		labels := LabelsClusterType([]string{"M\tenvironments/users/foo/updated-cluster.yml"})
		assert.Equal(t, "cluster:update", labels[0])
	})

	t.Run("Should return an update label when renamed", func(t *testing.T) {
		labels := LabelsClusterType([]string{"R\tenvironments/users/foo/ccs.yml\tenvironments/users/foo/edge-ccs.yml"})
		assert.Equal(t, "cluster:update", labels[0])
	})

	t.Run("Should return an create label", func(t *testing.T) {
		labels := LabelsClusterType([]string{"A\tenvironments/users/foo/create.yml"})
		assert.Equal(t, "cluster:create", labels[0])
	})

	t.Run("Should return an empty labels if unmatched path", func(t *testing.T) {
		labels := LabelsClusterType([]string{"A\tconfig-cluster.yml"})
		assert.True(t, len(labels) == 0)
	})

	t.Run("Should return an empty labels if unmatched extension", func(t *testing.T) {
		labels := LabelsClusterType([]string{"A\tnew-file"})
		assert.True(t, len(labels) == 0)
	})
}

func Test_getLabelsGivenACluster(t *testing.T) {
	tmpDir := t.TempDir()
	t.Run("Should return an empty list if empty value", func(t *testing.T) {
		labels := getLabelsGivenACluster("")
		assert.True(t, len(labels) == 0)
	})
	t.Run("Should return an empty list if unknown file", func(t *testing.T) {
		labels := getLabelsGivenACluster("unknown")
		assert.True(t, len(labels) == 0)
	})
	t.Run("Should return a list with some labels if a template", func(t *testing.T) {
		ymlFilePath := filepath.Join(tmpDir, "template.yml")
		ymlContent := `
template_name: serverless
stack:
  mode: "serverless"
  template: "observability"
  target: "qa"
`
		err := os.WriteFile(ymlFilePath, []byte(ymlContent), 0644)
		assert.Nil(t, err, "failed to write yaml file")
		labels := getLabelsGivenACluster(ymlFilePath)
		assert.Equal(t, 4, len(labels))
	})
	t.Run("Should return a list with some labels if no a template", func(t *testing.T) {
		ymlFilePath := filepath.Join(tmpDir, "no-template.yml")
		ymlContent := `
stack:
  mode: "serverless"
  template: "observability"
  target: "qa"
`
		err := os.WriteFile(ymlFilePath, []byte(ymlContent), 0644)
		assert.Nil(t, err, "failed to write yaml file")
		labels := getLabelsGivenACluster(ymlFilePath)
		assert.Equal(t, 3, len(labels))
	})
}
