package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	// Doesn't set and check exactly under on Linux but good enough for test
	setCollectionFile("windows")
	if collectionFile.filename != filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile) {
		t.Error("Windows collection file not set properly")
	}

	setCollectionFile("linux")
	if collectionFile.filename != filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile) {
		t.Error("Runtime OS collection file not set properly")
	}

	// Not sure how to unit test but at least run it for now
	loadCollectionFromStdin()

	for _, c := range reservedCommands {
		if isReservedCommand(c) != true {
			t.Error("Error checking valid reserved commands")
		}
	}

	if isReservedCommand("validsecretname") == true {
		t.Error("Error checking invalid reserved command")
	}
}
