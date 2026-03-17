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
package clusters

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowBanner(t *testing.T) {
	yaml := `- msg: zero
  level: INFO
- msg: one
  level: WARN
- msg: two
  level: DEBUG`

	tmp := t.TempDir()
	bannersPath := filepath.Join(tmp, "banners.yml")

	err := os.WriteFile(bannersPath, []byte(yaml), os.ModePerm)
	assert.Nil(t, err)

	err = os.Chdir(tmp)
	assert.Nil(t, err)

	messages := ShowBanner(tmp)

	assert.Equal(t, 3, len(messages))

	msg0 := messages[0]
	assert.Equal(t, "zero", msg0.Text)
	assert.Equal(t, "INFO", msg0.Level)

	msg1 := messages[1]
	assert.Equal(t, "one", msg1.Text)
	assert.Equal(t, "WARN", msg1.Level)

	msg2 := messages[2]
	assert.Equal(t, "two", msg2.Text)
	assert.Equal(t, "DEBUG", msg2.Level)
}

func TestShowBannerFileNotFound(t *testing.T) {
	tmp := t.TempDir()

	err := os.Chdir(tmp)
	assert.Nil(t, err)

	messages := ShowBanner(tmp)

	assert.Equal(t, 0, len(messages))
}
