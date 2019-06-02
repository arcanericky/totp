package cmd

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigRename(t *testing.T) {
	defaultCollectionFile = "testcollection"

	keys := createTestData(t)

	configRenameCmd.Run(nil, []string{})

	// Valid parameters
	newName := "newName"
	configRenameCmd.Run(nil, []string{keys[0].name, newName})
	c, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for rename test from file")
	}

	_, err = c.GetKey(newName)
	if err != nil {
		t.Error("Key not renamed")
	}

	// Test rename to config
	configRenameCmd.Run(nil, []string{newName, configCmd.Use})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for rename test from file")
	}

	_, err = c.GetKey(configCmd.Use)
	if err == nil {
		t.Error("Key should not have been renamed to \"" + configCmd.Use + "\"")
	}

	// No collections file
	os.Remove(defaultCollectionFile)
	configRenameCmd.Run(nil, []string{keys[0].name, "newname"})
}
