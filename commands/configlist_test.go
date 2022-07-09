package commands

import (
	"os"
	"testing"
)

func TestConfigList(t *testing.T) {
	collectionFile.filename = "testcollection.json"

	createTestData(t)

	configListCmd := getConfigListCmd()

	configListCmd.Run(configListCmd, []string{})

	// names only
	_ = configListCmd.Flags().Set("names", "true")
	configListCmd.Run(configListCmd, []string{})

	// all
	configListCmd = getConfigListCmd()
	_ = configListCmd.Flags().Set("all", "true")
	configListCmd.Run(configListCmd, []string{})

	// names and all
	configListCmd = getConfigListCmd()
	_ = configListCmd.Flags().Set("all", "true")
	_ = configListCmd.Flags().Set("names", "true")
	configListCmd.Run(configListCmd, []string{})

	configListCmd.ResetFlags()
	configListCmd.Flags().Int64P("names", "n", 0, "")
	configListCmd.Run(configListCmd, []string{})

	// No collections file
	os.Remove(collectionFile.filename)
	configListCmd.Run(configListCmd, []string{})
}
