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

// Package prompt It contains the functions to interact with the user in the console.
package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGithubToken(t *testing.T) {

	assert.True(t, IsGithubToken("ghp_123472Yl3MVsM2XPrYqVi2h5m30BxF02O2ab"))

	assert.False(t, IsGithubToken("my_passwords"))

	assert.True(t, IsGithubToken("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"))
}
