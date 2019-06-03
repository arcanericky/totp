package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaults(t *testing.T) {
	// Doesn't set and check exactly under on Linux but good enough for test
	setCollectionFile("windows")
	if defaultCollectionFile != filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile) {
		t.Error("Windows collection file not set properly")
	}

	setCollectionFile(runtime.GOOS)
	if defaultCollectionFile != filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile) {
		t.Error("Runtime OS collection file not set properly")
	}

	for _, c := range reservedCommands {
		if isReservedCommand(c) != true {
			t.Error("Error checking valid reserved commands")
		}
	}

	if isReservedCommand("validsecretname") == true {
		t.Error("Error checking invalid reserved command")
	}
}
