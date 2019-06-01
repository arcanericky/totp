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
	configRenameCmd.Run(nil, []string{newName, configName})
	c, err = totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for rename test from file")
	}

	_, err = c.GetKey(configName)
	if err == nil {
		t.Error("Key should not have been renamed to \"" + configName + "\"")
	}

	// No collections file
	os.Remove(defaultCollectionFile)
	configRenameCmd.Run(nil, []string{keys[0].name, "newname"})
}
