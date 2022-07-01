package cmd

import (
	"os"
	"testing"
	"time"

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
	Execute()
	collectionFile.filename = "testcollection.json"

	secretList := createTestData(t)

	rootCmd := getRootCmd()

	// No parameters
	rootCmd.Run(rootCmd, []string{})

	// Valid entry and secret
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Non-existing entry
	rootCmd.Run(rootCmd, []string{"invalidsecret"})

	// Test follow condition
	savedGenerateCodesService := generateCodesService
	generateCodesService = func(time.Duration, time.Duration, time.Duration, func(time.Duration), string, string) {}
	_ = rootCmd.Flags().Lookup(optionFollow).Value.Set("true")
	rootCmd.Run(rootCmd, []string{"name0"})
	generateCodesService = savedGenerateCodesService
	_ = rootCmd.Flags().Lookup(optionFollow).Value.Set("false")

	// Completion
	rootCmd.ValidArgsFunction(rootCmd, []string{}, "na")

	// Completion with args
	rootCmd.ValidArgsFunction(rootCmd, []string{"secret"}, "na")

	// No collections file
	os.Remove(collectionFile.filename)
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Completion without collections
	rootCmd.ValidArgsFunction(rootCmd, []string{}, "na")

	// Excessive args
	rootCmd.Run(rootCmd, []string{"secretname", "extraarg"})

	// Provide secret option
	_ = rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{})

	// Provide invalid secret option
	_ = rootCmd.Flags().Set(optionSecret, "seed1")
	rootCmd.Run(rootCmd, []string{})

	// File option
	_ = rootCmd.Flags().Set(optionFile, collectionFile.filename)
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	// Stdio option
	_ = rootCmd.Flags().Set(optionStdio, "true")
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})
	collectionFile.loader = loadCollectionFromDefaultFile
	collectionFile.useStdio = false
	_ = rootCmd.Flags().Set(optionStdio, "")

	// Time option
	_ = rootCmd.Flags().Set(optionTime, "2019-06-01T20:00:00-05:00")
	rootCmd.Run(rootCmd, []string{})

	// Give secret and secret name
	_ = rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{"secretname"})
	_ = rootCmd.Flags().Set(optionSecret, "")

	// Invalid time option
	_ = rootCmd.Flags().Set(optionTime, "invalidtime")
	rootCmd.Run(rootCmd, []string{})
	_ = rootCmd.Flags().Set(optionTime, "")

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

	// optionFollow error
	f = rootCmd.Flags().Lookup(optionFollow)
	savedFlagValue = f.Value
	f.Value = new(flagValue)
	rootCmd.Run(rootCmd, []string{})
	f.Value = savedFlagValue
}

func TestExecOnInterval(t *testing.T) {
	execCount := 0
	preAndExecNormal := func() bool { return false }
	preAndExecEarlyExit := func() bool { execCount++; return execCount == 2 }

	// Normal execution
	execCount = 0
	callOnInterval(2*time.Millisecond, 1*time.Millisecond, preAndExecNormal)

	// Exit at preExec
	execCount = 1
	callOnInterval(2*time.Millisecond, 1*time.Millisecond, preAndExecNormal)

	// Exit at top exec
	execCount = 1
	callOnInterval(2*time.Millisecond, 1*time.Millisecond, preAndExecEarlyExit)

	// Exit at loop exec
	execCount = 0
	callOnInterval(2*time.Millisecond, 1*time.Millisecond, preAndExecEarlyExit)

	// No callbacks
	callOnInterval(2*time.Millisecond, 1*time.Millisecond, nil)
}

func TestDurationToNextInterval(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2019-06-23T20:00:01-05:00")
	expectedResult := time.Duration(29 * time.Second)
	actualResult := durationToNextInterval(now)
	if expectedResult != actualResult {
		t.Errorf("durationToNextInterval(%s) expected %s but returned %s", now, expectedResult, actualResult)
	}

	now, _ = time.Parse(time.RFC3339, "2019-06-23T20:00:31-05:00")
	expectedResult = time.Duration(29 * time.Second)
	actualResult = durationToNextInterval(now)
	if expectedResult != actualResult {
		t.Errorf("durationToNextInterval(%s) expected %s but returned %s", now, expectedResult, actualResult)
	}

	now, _ = time.Parse(time.RFC3339, "2019-06-23T20:00:31.001-05:00")
	expectedResult = time.Duration(28*time.Second + 999*time.Millisecond)
	actualResult = durationToNextInterval(now)
	if expectedResult != actualResult {
		t.Errorf("durationToNextInterval(%s) expected %s but returned %s", now, expectedResult, actualResult)
	}
}

func TestGenerateCodes(t *testing.T) {
	var d time.Duration
	generateCodes(d, 2*time.Millisecond, 1*time.Millisecond,
		func(d time.Duration) {}, "", "seed")
	generateCodes(d, 2*time.Millisecond, 1*time.Millisecond,
		func(d time.Duration) {}, "", "invalidseed")
}
