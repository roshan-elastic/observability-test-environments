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
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-robot/cmd"
	"github.com/joho/godotenv"
)

var (
	Version string = "unknown"
	Build   string = "unknown"
	Commit  string = "unknown"
)

func main() {
	// Load Env variables from .dot file
	godotenv.Load(".env")

	// show the file name and line number in the logs
	log.SetFlags(log.LstdFlags | log.Llongfile)
	os.Setenv("OBLT_ROBOT_VERSION", Version)
	os.Setenv("OBLT_ROBOT_BUILD", Build)
	os.Setenv("ELASTIC_APM_SERVICE_VERSION", fmt.Sprintf("%s-%s", Version, Commit))
	os.Setenv("ELASTIC_APM_ENVIRONMENT", "production")

	log.Printf("Version: %v (%v) [%v]\n", os.Getenv("OBLT_ROBOT_VERSION"), os.Getenv("OBLT_ROBOT_BUILD"), Commit)
	cmd.Execute()
}
