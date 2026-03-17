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
package apm

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"go.elastic.co/apm/v2"
)

type Label struct {
	Key   string
	Value string
}

// StartTransaction returns a new Transaction with the specified
// name, type and labels and also the transaction context.
func StartTransaction(name, transactionType string, labels []Label, config config.ObltConfiguration) (*apm.Transaction, context.Context) {
	tx := apm.DefaultTracer().StartTransaction(name, transactionType)
	for _, label := range labels {
		tx.Context.SetLabel(label.Key, label.Value)
	}

	tx.Context.SetLabel("slack-channel", config.SlackChannel)
	tx.Context.SetLabel("user", config.Username)

	ctx := apm.ContextWithTransaction(context.Background(), tx)
	return tx, ctx
}

// Flush ensures the data reaches the APM server.
func Flush(tx *apm.Transaction) {
	tx.End()
	apm.DefaultTracer().Flush(nil)
}

// ReportError reports the error and send APM errors.
func ReportError(tx *apm.Transaction, message string) {
	e := apm.DefaultTracer().NewError(errors.New(message))
	e.SetTransaction(tx)
	e.Send()
	Flush(tx)
	log.Fatalf(message)
}

// CobraCheckErr reports the error if any with cobra.CheckErr and send APM errors.
func CobraCheckErr(err error, tx *apm.Transaction, ctx context.Context) {
	if err != nil {
		e := apm.CaptureError(ctx, err)
		e.Send()
		Flush(tx)
	}
	if err != nil {
		logger.Errorf("%v", err)
		os.Exit(1)
	}
}
