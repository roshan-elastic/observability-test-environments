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

	"go.elastic.co/apm/v2"
)

type Label struct {
	Key   string
	Value string
}

// StartTransaction returns a new Transaction with the specified
// name, type and labels and also the transaction context.
func StartTransaction(name, transactionType string, labels []Label) (*apm.Transaction, context.Context) {
	tx := apm.DefaultTracer().StartTransaction(name, transactionType)
	for _, label := range labels {
		tx.Context.SetLabel(label.Key, label.Value)
	}
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	return tx, ctx
}

// StartTransactionForm returns a new Transaction with the specified
// name, type and labels and also the transaction context.
func StartTransactionForm(name, transactionType string, triggerID string, user string) (*apm.Transaction, context.Context) {
	labelUser := Label{Key: "slack-user", Value: user}
	labelTriggerID := Label{Key: "slack-trigger", Value: triggerID}
	tx, ctx := StartTransaction(name, transactionType, []Label{labelUser, labelTriggerID})
	return tx, ctx
}
