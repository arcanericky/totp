package cmd

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigDelete(t *testing.T) {
	defaultCollectionFile = "testcollection"

	seedList := createTestData(t)

	// Key does not exit
	configDeleteCmd.Run(nil, []string{"key"})

	// No key provided
	configDeleteCmd.Run(nil, []string{})

	// Successful delete
	configDeleteCmd.Run(nil, []string{seedList[3].name})
	c, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		t.Error("Could not load collection for delete test from file")
	}

	_, err = c.GetKey(seedList[3].name)
	if err == nil {
		t.Error("Key not deleted")
	}

	// No collections file
	os.Remove(defaultCollectionFile)
	deleteKey(seedList[3].name)
}
