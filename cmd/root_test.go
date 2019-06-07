package cmd

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

type flagValue struct{}

func (f flagValue) Set(s string) error {
	return nil
}

func (f flagValue) Type() string {
	return ""
}

func (f flagValue) String() string {
	return ""
}

func TestRoot(t *testing.T) {
	collectionFile.filename = "testcollection"

	secretList := createTestData(t)

	// No parameters
	rootCmd.Run(rootCmd, []string{})

	// Valid entry and secret
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Non-existing entry
	rootCmd.Run(rootCmd, []string{"invalidsecret"})

	// No collections file
	os.Remove(collectionFile.filename)
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Excessive args
	rootCmd.Run(rootCmd, []string{"secretname", "extraarg"})

	// Provide secret option
	rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{})

	// Provide invalid secret option
	rootCmd.Flags().Set(optionSecret, "seed1")
	rootCmd.Run(rootCmd, []string{})

	// File option
	rootCmd.Flags().Set(optionFile, collectionFile.filename)
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	// Stdio option
	rootCmd.Flags().Set(optionStdio, "true")
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	// Time option
	rootCmd.Flags().Set(optionTime, "2019-06-01T20:00:00-05:00")
	rootCmd.Run(rootCmd, []string{})

	// Give secret and secret name
	rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{"secretname"})

	// Invalid time option
	rootCmd.Flags().Set(optionTime, "invalidtime")
	rootCmd.Run(rootCmd, []string{})

	var f *pflag.Flag
	var savedFlagValue pflag.Value

	// optionFile error
	f = rootCmd.Flags().Lookup(optionFile)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	f.Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})
	f.Value = savedFlagValue

	// optionStdio error
	f = rootCmd.Flags().Lookup(optionStdio)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})
	f.Value = savedFlagValue

	// optionSecret error
	f = rootCmd.Flags().Lookup(optionSecret)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.Run(rootCmd, []string{})
	f.Value = savedFlagValue

	// optionBackward error
	f = rootCmd.Flags().Lookup(optionBackward)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.Run(rootCmd, []string{})
	f.Value = savedFlagValue

	// optionForward error
	f = rootCmd.Flags().Lookup(optionForward)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.Run(rootCmd, []string{})
	f.Value = savedFlagValue

	// optionTime error
	f = rootCmd.Flags().Lookup(optionTime)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.Run(rootCmd, []string{})
	f.Value = savedFlagValue

	Execute()
	savedArgs := os.Args
	os.Args = []string{"totp", "--invalidoption"}
	Execute()
	os.Args = savedArgs
}
