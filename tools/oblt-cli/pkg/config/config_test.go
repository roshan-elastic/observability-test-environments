package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("File does exist", func(t *testing.T) {
		f, err := os.CreateTemp(tmpDir, "i-am-alive")
		if err != nil {
			t.Fatal("could not create temp file")
		}

		assert.True(t, fileExists(f.Name()))
	})

	t.Run("File does not exist", func(t *testing.T) {
		assert.False(t, fileExists(filepath.Join(tmpDir, "ghost")))
	})
}

func TestForUser(t *testing.T) {
	cfgFile := ForUser("foo")

	expectedSuffix := filepath.Join(configDir, "foo", configFile)

	// we do not control User's dir at test time, that's why we check against everything else
	assert.True(t, strings.HasSuffix(cfgFile, expectedSuffix))
}

func TestInitialise(t *testing.T) {
	tmpDir := t.TempDir()
	parentDir := "foo"
	parentPath := filepath.Join(tmpDir, parentDir)
	fallbackCfgFile := filepath.Join(parentPath, configFile)

	t.Run("Passing a fixed configuration file", func(t *testing.T) {
		originalCfgFile := CfgFile
		defer func() {
			CfgFile = originalCfgFile
		}()

		CfgFile = filepath.Join(tmpDir, "another-cfg-file.yml")
		Initialise(fallbackCfgFile)

		assert.Equal(t, CfgFile, viper.GetViper().ConfigFileUsed())
	})

	t.Run("Using fallback config file", func(t *testing.T) {
		Initialise(fallbackCfgFile)

		// should create default config parent dir
		assert.True(t, fileExists(parentPath))
		assert.Equal(t, fallbackCfgFile, viper.GetViper().ConfigFileUsed())
	})
}

func TestValidateLength(t *testing.T) {
	// Test a valid input
	err := ValidateLength("hello", 5, 10)
	assert.NoError(t, err)

	// Test an input that is too short
	err = ValidateLength("hi", 5, 10)
	assert.Error(t, err)

	// Test an input that is too long
	err = ValidateLength("this is a really long string", 5, 10)
	assert.Error(t, err)
}

func TestValidateSlackChannel(t *testing.T) {
	// Test a valid input
	err := ValidateSlackChannel("#my-channel")
	assert.NoError(t, err)

	// Test an input that is missing the pound sign
	err = ValidateSlackChannel("my-channel")
	assert.Error(t, err)

	// Test an input that contains invalid characters
	err = ValidateSlackChannel("#my_channel")
	assert.Error(t, err)

	// Test an input that is too long
	err = ValidateSlackChannel("#this-channel-name-is-way-too-long-for-slack-oooooooooooooooooooooooooooooooooooo")
	assert.Error(t, err)

	// Test an input that is a Slack User ID
	err = ValidateSlackChannel("@U12345678")
	assert.NoError(t, err)
}

func TestValidatePrefix(t *testing.T) {
	// Test a valid input
	err := ValidatePrefix("my-prefix")
	assert.NoError(t, err)

	// Test an input that is too long
	err = ValidatePrefix("this-prefix-is-way-too-long-for-the-config")
	assert.Error(t, err)

	// Test an input that contains invalid characters
	err = ValidatePrefix("my_prefix")
	assert.Error(t, err)
}

func TestValidateDockerImage(t *testing.T) {
	// Test a valid input
	err := ValidateDockerImage("my-image:latest")
	assert.NoError(t, err)

	// Test an input with an invalid character
	err = ValidateDockerImage("my_image:latest")
	assert.Error(t, err)

	// Test an input with an invalid tag
	err = ValidateDockerImage("my-image:invalid_tag!")
	assert.Error(t, err)

	// Test an input with an invalid name
	err = ValidateDockerImage("my/image:latest")
	assert.NoError(t, err)

	// Test an input with an valid name
	err = ValidateDockerImage("docker.elastic.co/elasticsearch/elasticsearch:7.14.0-SNAPSHOT")
	assert.NoError(t, err)
}

func TestValidateNames(t *testing.T) {
	// Test a valid input
	err := ValidateNames("myname")
	assert.NoError(t, err)

	// Test an input with an invalid character
	err = ValidateNames("my_name")
	assert.Error(t, err)

	// Test an input with an uppercase character
	err = ValidateNames("MyName")
	assert.Error(t, err)

	// Test an input with a hyphen
	err = ValidateNames("my-name")
	assert.NoError(t, err)

	// Test an input with a period
	err = ValidateNames("my.name")
	assert.NoError(t, err)

	// Test an input with a space
	err = ValidateNames("my name")
	assert.Error(t, err)
}

func TestValidateAlphanumeric(t *testing.T) {
	// Test a valid input
	err := ValidateAlphanumeric("myvalue123")
	assert.NoError(t, err)

	// Test an input with an uppercase character
	err = ValidateAlphanumeric("MyValue123")
	assert.Error(t, err)

	// Test an input with a special character
	err = ValidateAlphanumeric("my_value_123")
	assert.Error(t, err)

	// Test an input with a space
	err = ValidateAlphanumeric("my value 123")
	assert.Error(t, err)
}

func TestValidateSemVer(t *testing.T) {
	// Test a valid input
	err := ValidateSemVer("1.2.3")
	assert.NoError(t, err)

	// Test a valid input
	err = ValidateSemVer("1.2.3-SNAPSHOT")
	assert.NoError(t, err)

	// Test a valid input
	err = ValidateSemVer("1.2.3-abcde-SNAPSHOT")
	assert.NoError(t, err)

	// Test an input with a non-numeric character
	err = ValidateSemVer("1.2.a")
	assert.Error(t, err)

	// Test an input with too many version components
	err = ValidateSemVer("1.2.3.4")
	assert.Error(t, err)

	// Test an input with too few version components
	err = ValidateSemVer("1.2")
	assert.Error(t, err)

	// Test an input with a negative version component
	err = ValidateSemVer("1.2.-3")
	assert.Error(t, err)
}
