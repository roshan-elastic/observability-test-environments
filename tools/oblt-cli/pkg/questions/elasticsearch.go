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

// Package questions It contains the functions to interact with Elasticsearch
package questions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

const (
	INDEX       string = "faq"
	SEARCH_SIZE int    = 5
)

// Answer defines an answer entry defined in the FAQ index
type Answer struct {
	Title  string
	Answer string
	Url    string
}

// FAQRecord defines the FAQ entry in the FAQ index
type FAQRecord struct {
	Question string
	Answer   string
	Url      string
	User     string
}

type Elasticsearch struct {
	client *elasticsearch.Client
}

// NewClient initializes a new Elasticsearch client.
func (v *Elasticsearch) NewClient() error {
	logger.Debugf("Elasticsearch NewClient")

	es, err := newClient()
	v.client = es
	if err != nil {
		return err
	}

	res, err := v.client.Info()
	if err != nil {
		return err
	}

	defer res.Body.Close()
	logger.Infof("%v", res)
	return nil
}

// AddAnswer add an answer with the question
func (v *Elasticsearch) AddAnswer(question string, answer string, url string, user string) error {

	// Create Elasticsearch client
	v.NewClient()

	// Build the request body.
	data, err := json.Marshal(FAQRecord{
		Question: question,
		Answer:   answer,
		User:     user,
		Url:      url,
	})
	if err != nil {
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   INDEX,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), v.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	logger.Infof("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	return nil
}

// SearchAnswers searches for all the answers for the given question
func (v *Elasticsearch) SearchAnswers(question string) ([]Answer, error) {
	var ret []Answer

	query := map[string]interface{}{
		"size": SEARCH_SIZE,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"question": question,
			},
		},
	}
	values, err := search(INDEX, query)
	if err != nil {
		return nil, err
	}

	// If there are values then parse them as Answers
	if values != nil {
		for i, hit := range values["hits"].(map[string]interface{})["hits"].([]interface{}) {
			score := hit.(map[string]interface{})["_score"].(float64)
			if score > 1.0 {
				question := hit.(map[string]interface{})["_source"].(map[string]interface{})["question"]
				answer := hit.(map[string]interface{})["_source"].(map[string]interface{})["answer"]
				url := hit.(map[string]interface{})["_source"].(map[string]interface{})["url"]
				answerType := Answer{
					Title:  fmt.Sprintf("%d - %s", i, question),
					Answer: fmt.Sprintf("%s", answer),
					Url:    fmt.Sprintf("%s", url),
				}
				ret = append(ret, answerType)
			}
		}
	}

	return ret, nil
}

func search(index string, query map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	var r map[string]interface{}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return r, err
	}

	// Create Elasticsearch client
	client, err := newClient()
	if err != nil {
		return r, err
	}

	// Perform the search request.
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return r, err
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return r, err
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}

	return r, nil
}

func newClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
		Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
	}

	es, err := elasticsearch.NewClient(cfg)
	return es, err
}
