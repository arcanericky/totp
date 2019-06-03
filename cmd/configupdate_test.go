package cmd

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigUpdate(t *testing.T) {
	defaultCollectionFile = "testcollection"

	createTestData(t)

	// Valid parameters
	secretName := "testsecret"
	configUpdateCmd.Run(nil, []string{secretName, "seed"})
	c, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	_, err = c.GetSecret(secretName)
	if err != nil {
		t.Error("Secret not added")
	}

	// Test update secret
	newSecret := "seedseed"
	configUpdateCmd.Run(nil, []string{secretName, newSecret})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	secret, err := c.GetSecret(secretName)
	if err != nil || secret.Value != newSecret {
		t.Error("Secret not updated", secret)
	}

	// Test using secret named 'config'
	secretName = configCmd.Use
	configUpdateCmd.Run(nil, []string{secretName, "seed"})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	secret, err = c.GetSecret(secretName)
	if err == nil {
		t.Error("Secret named \"" + configCmd.Use + "\" should not have been saved")
	}

	// No parameters passed
	configUpdateCmd.Run(nil, []string{})

	// Invalid secret value
	configUpdateCmd.Run(nil, []string{"testsecret", "seed1"})

	// No collections file
	os.Remove(defaultCollectionFile)
	configUpdateCmd.Run(nil, []string{"testsecret", "seed"})
}
