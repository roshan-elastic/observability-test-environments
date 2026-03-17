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

package maps

import (
	"encoding/json"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
)

// Converts a interface, map[interface{}]interface, or []interface[]
// to a interface type we can marshall to JSON.
// https://github.com/go-yaml/yaml/issues/139
func InterfaceToMapInterface(v interface{}) (r interface{}) {
	switch v2 := v.(type) {
	case []interface{}:
		for i, e := range v2 {
			v2[i] = InterfaceToMapInterface(e)
		}
		r = []interface{}(v2)
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{}, len(v2))
		for k, e := range v2 {
			newMap[k.(string)] = InterfaceToMapInterface(e)
		}
		r = map[string]interface{}(newMap)
	default:
		r = v2
	}
	return
}

// convertToMap convert a interface{} to map[string]interface{}
func ConvertToMap(object interface{}) (objMap map[string]interface{}) {
	logger.Debugf("Validate Recipe.")
	jsonString, err := json.Marshal(InterfaceToMapInterface(object))
	cobra.CheckErr(err)
	err = json.Unmarshal(jsonString, &objMap)
	cobra.CheckErr(err)
	return objMap
}
