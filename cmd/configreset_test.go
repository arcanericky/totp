package cmd

import (
	"os"
	"testing"
)

func TestConfigReset(t *testing.T) {
	collectionFile.filename = "testcollection.json"

	createTestData(t)
	configResetCmd := getConfigResetCmd()
	_ = configResetCmd.Flags().Set(optionYes, "true")
	configResetCmd.Run(nil, []string{})

	_, err := os.Stat(collectionFile.filename)
	if !os.IsNotExist(err) {
		t.Error("Failed to remove the collection file")
	}

	configResetCmd = getConfigResetCmd()
	configResetCmd.Run(nil, []string{})

	if err := configReset("nosuchfilename.example"); err == nil {
		t.Error("Failed to generate error removing invalid collection file")
	}
}
