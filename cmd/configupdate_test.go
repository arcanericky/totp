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
	keyName := "testkey"
	configUpdateCmd.Run(nil, []string{keyName, "seed"})
	c, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	_, err = c.GetKey(keyName)
	if err != nil {
		t.Error("Key not added")
	}

	// Test update seed
	newSeed := "seedseed"
	configUpdateCmd.Run(nil, []string{keyName, newSeed})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	key, err := c.GetKey(keyName)
	if err != nil || key.Seed != newSeed {
		t.Error("Key not updated", key)
	}

	// Test using seed named 'config'
	keyName = configName
	configUpdateCmd.Run(nil, []string{keyName, "seed"})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for update test from file")
	}

	key, err = c.GetKey(keyName)
	if err == nil {
		t.Error("Key named \"" + configName + "\" should not have been saved")
	}

	// No parameters passed
	configUpdateCmd.Run(nil, []string{})

	// Invalid seed value
	configUpdateCmd.Run(nil, []string{"testkey", "seed1"})

	// No collections file
	os.Remove(defaultCollectionFile)
	configUpdateCmd.Run(nil, []string{"testkey", "seed"})
}
