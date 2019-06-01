package cmd

import (
	"os"
	"testing"
)

func TestConfigReset(t *testing.T) {
	defaultCollectionFile = "testcollection"

	createTestData(t)

	configResetCmd.Run(nil, []string{})

	_, err := os.Stat(defaultCollectionFile)
	if !os.IsNotExist(err) {
		t.Error("Failed to remove the collection file")
	}
}
