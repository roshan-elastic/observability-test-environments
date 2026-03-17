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

// Package certs contains functions to create self-signed certificates
package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// CreateKey creates a new RSA private key and returns it as a string
func CreateKey() (key *rsa.PrivateKey, keyPem string, err error) {
	var buf *bytes.Buffer
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err == nil {
		buf = bytes.NewBufferString("")
		pem.Encode(buf, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
	}
	return key, buf.String(), err
}

// CreateCertificate creates a new self-signed certificate using the given RSA private key
func CreateCertificate(key *rsa.PrivateKey) (pemCert string, err error) {
	var buf *bytes.Buffer
	serialNumber, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err == nil {
		now := time.Now()
		template := x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				CommonName:   "localhost",
				Country:      []string{"US"},
				Organization: []string{"Internet Widgits Pty Ltd"},
				Province:     []string{"Some-State"},
			},
			Issuer: pkix.Name{
				CommonName:   "localhost",
				Country:      []string{"AU"},
				Organization: []string{"Internet Widgits Pty Ltd"},
				Province:     []string{"Some-State"},
			},
			NotBefore: now,
			NotAfter:  now.AddDate(0, 0, 1),
			//PublicKeyAlgorithm: x509.RSA,
			//SignatureAlgorithm: x509.RSA,
		}
		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
		if err == nil {
			buf = bytes.NewBufferString("")
			pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		}
	}
	return buf.String(), err
}
