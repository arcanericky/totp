package cmd

import (
	"os"
	"testing"
)

func TestRoot(t *testing.T) {
	defaultCollectionFile = "testcollection"

	secretList := createTestData(t)

	// No parameters
	rootCmd.Run(rootCmd, []string{})

	// Valid entry and secret
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Non-existing entry
	rootCmd.Run(rootCmd, []string{"invalidsecret"})

	// No collections file
	os.Remove(defaultCollectionFile)
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Provide secret option
	rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{})

	// Provide invalid secret option
	rootCmd.Flags().Set(optionSecret, "seed1")
	rootCmd.Run(rootCmd, []string{})

	// File option
	rootCmd.Flags().Set(optionFile, defaultCollectionFile)
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	rootCmd.ResetFlags()

	// Test error on secret option
	rootCmd.Run(rootCmd, []string{})

	// Test error on secret option
	rootCmd.Flags().Int64P(optionFile, "f", 0, "")
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	Execute()

	savedArgs := os.Args
	os.Args = []string{"totp", "--invalidoption"}
	Execute()
	os.Args = savedArgs
}
