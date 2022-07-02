package commands

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigUpdate(t *testing.T) {
	collectionFile.filename = "testcollection.json"

	createTestData(t)

	configUpdateCmd := getConfigUpdateCmd(getRootCmd())

	// Valid parameters
	secretName := "testsecret"
	configUpdateCmd.Run(nil, []string{secretName, "seed"})
	c, err := totp.NewCollectionWithFile(collectionFile.filename)
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
	c, err = totp.NewCollectionWithFile(collectionFile.filename)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	secret, err := c.GetSecret(secretName)
	if err != nil || secret.Value != newSecret {
		t.Error("Secret not updated", secret)
	}

	// Test using secret named 'config'
	secretName = "config"
	configUpdateCmd.Run(nil, []string{secretName, "seed"})
	c, err = totp.NewCollectionWithFile(collectionFile.filename)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	secret, err = c.GetSecret(secretName)
	if err == nil {
		t.Error("Secret named \"" + secretName + "\" should not have been saved")
	}

	// No parameters passed
	configUpdateCmd.Run(nil, []string{})

	// Invalid secret value
	configUpdateCmd.Run(nil, []string{"testsecret", "seed1"})

	// No collections file
	os.Remove(collectionFile.filename)
	configUpdateCmd.Run(nil, []string{"testsecret", "seed"})
}
