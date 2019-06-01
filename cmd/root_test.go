package cmd

import (
	"os"
	"testing"
)

func TestRoot(t *testing.T) {
	defaultCollectionFile = "testcollection"

	seedList := createTestData(t)

	// No parameters
	rootCmd.Run(rootCmd, []string{})

	// Valid entry and seed
	rootCmd.Run(rootCmd, []string{seedList[0].name})

	// Non-existing entry
	rootCmd.Run(rootCmd, []string{"invalidkey"})

	// No collections file
	os.Remove(defaultCollectionFile)
	rootCmd.Run(rootCmd, []string{seedList[0].name})

	// Provide seed option
	rootCmd.Flags().Set(optionSeed, "seed")
	rootCmd.Run(rootCmd, []string{})

	// Provide invalid seed option
	rootCmd.Flags().Set(optionSeed, "seed1")
	rootCmd.Run(rootCmd, []string{})

	// File option
	rootCmd.Flags().Set(optionFile, defaultCollectionFile)
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"key"})

	// Error when parsing seed option
	rootCmd.ResetFlags()
	rootCmd.Flags().Int64P(optionSeed, "s", 0, "")
	rootCmd.Run(rootCmd, []string{})

	// Error when parsing file option
	rootCmd.ResetFlags()
	rootCmd.Flags().Int64P(optionFile, "f", 0, "")
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"key"})

	Execute()

	savedArgs := os.Args
	os.Args = []string{"totp", "--invalidoption"}
	Execute()
	os.Args = savedArgs
}
