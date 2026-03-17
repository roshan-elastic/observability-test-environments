package certs

import (
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	key, keyPem, err := CreateKey()
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotEmpty(t, keyPem)

	// Check if the key PEM can be parsed
	block, _ := pem.Decode([]byte(keyPem))
	assert.NotNil(t, block)
	assert.Equal(t, "PRIVATE KEY", block.Type)
}

func TestCreateCertificate(t *testing.T) {
	key, _, _ := CreateKey()
	pemCert, err := CreateCertificate(key)
	assert.NoError(t, err)
	assert.NotEmpty(t, pemCert)

	// Check if the certificate PEM can be parsed
	block, _ := pem.Decode([]byte(pemCert))
	assert.NotNil(t, block)
	assert.Equal(t, "CERTIFICATE", block.Type)
}
