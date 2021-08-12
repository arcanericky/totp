package cmd

import (
	"os"
	"testing"
)

func TestConfigList(t *testing.T) {
	collectionFile.filename = "testcollection.json"

	createTestData(t)

	configListCmd.Run(configListCmd, []string{})

	configListCmd.Flags().Set("names", "true")
	configListCmd.Run(configListCmd, []string{})

	configListCmd.ResetFlags()
	configListCmd.Flags().Int64P("names", "n", 0, "")
	configListCmd.Run(configListCmd, []string{})

	// No collections file
	os.Remove(collectionFile.filename)
	configListCmd.Run(configListCmd, []string{})
}
