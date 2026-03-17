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
	"os"
	"path/filepath"
	"testing"

	"github.com/blang/semver"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/clusters"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	tempDir string
}

func NewMockRepository(tempDir string) clusters.Repository {
	return &mockRepository{tempDir: tempDir}
}
func (m *mockRepository) GetPath() string {
	return m.tempDir
}
func (m *mockRepository) GetVersion() (repoVersion semver.Version, err error) {
	return semver.MustParse("1.0.0"), nil
}

func Test_verifyCompatibility(t *testing.T) {
	version := semver.MustParse("1.5.0")

	currentVersion = semver.MustParse("1.5.0")
	err := verifyCompatibility(version)
	assert.NoError(t, err)

	currentVersion = semver.MustParse("1.5.1")
	err = verifyCompatibility(version)
	assert.NoError(t, err)

	currentVersion = semver.MustParse("1.6.0")
	err = verifyCompatibility(version)
	assert.NoError(t, err)

	currentVersion = semver.MustParse("2.0.0")
	err = verifyCompatibility(version)
	assert.Error(t, err)
}

func Test_Banner(t *testing.T) {
	dir := t.TempDir()
	obltRepo := NewMockRepository(dir)
	bannerYAML := `---
- msg: test info message
  level: INFO
- msg: test warn message
  level: WARN
- msg: test error message
  level: ERROR
`
	err := os.WriteFile(filepath.Join(dir, "banners.yml"), []byte(bannerYAML), 0644)
	assert.NoError(t, err)

	disableBanner = false
	ret := showBanner(obltRepo)
	assert.Equal(t, 3, len(ret))

	disableBanner = true
	ret = showBanner(obltRepo)
	assert.Equal(t, 0, len(ret))
}
