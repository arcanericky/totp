package cmd

import (
	"os"
	"testing"
)

func TestConfigReset(t *testing.T) {
	collectionFile.filename = "testcollection"

	createTestData(t)

	configResetCmd.Run(nil, []string{})

	_, err := os.Stat(collectionFile.filename)
	if !os.IsNotExist(err) {
		t.Error("Failed to remove the collection file")
	}
}
