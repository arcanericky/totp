package cmd

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigRename(t *testing.T) {
	collectionFile.filename = "testcollection.json"

	secrets := createTestData(t)

	configRenameCmd.Run(nil, []string{})

	// Valid parameters
	newName := "newName"
	configRenameCmd.Run(nil, []string{secrets[0].name, newName})
	c, err := totp.NewCollectionWithFile(collectionFile.filename)
	if err != nil {
		t.Error("Could not load collection for rename test from file")
	}

	_, err = c.GetSecret(newName)
	if err != nil {
		t.Error("Secret not renamed")
	}

	// Test rename to config
	configRenameCmd.Run(nil, []string{newName, configCmd.Use})
	c, err = totp.NewCollectionWithFile(collectionFile.filename)
	if err != nil {
		t.Error("Could not load collection for rename test from file")
	}

	_, err = c.GetSecret(configCmd.Use)
	if err == nil {
		t.Error("Secret should not have been renamed to \"" + configCmd.Use + "\"")
	}

	// No collections file
	os.Remove(collectionFile.filename)
	configRenameCmd.Run(nil, []string{secrets[0].name, "newname"})
}
