package gcp

import (
	"os"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Skipping test that requires GCP credentials")
	}
	auth := NewGcpAuth()
	err := auth.Authenticate()
	assert.NoError(t, err, "Error should be nil")
}

func TestAuthenticateInteractive(t *testing.T) {
	t.Skip("This test is interactive and should be run manually")
	seed := config.Seed(5)
	secretPath := "projects/8560181848/secrets/oblt-cli-test-" + seed
	auth := NewGcpAuth()
	auth.AuthenticateInteractive()
	gcsm, err := NewClientWithOpts(auth.Options...)
	assert.NoError(t, err, "failed to create the client: %v", err)
	defer gcsm.Close()
	err = gcsm.CreateSecret(secretPath, "foo")
	assert.NoError(t, err, "Error should be nil")
	secretValue, err := gcsm.GetSecret(secretPath)
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, "foo", secretValue, "Secret value should match")
	err = gcsm.DeleteSecret(secretPath)
	assert.NoError(t, err, "Error should be nil")
}
