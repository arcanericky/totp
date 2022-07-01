package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	setCollectionFile("windows")
	if collectionFile.filename != filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile) {
		t.Error("Windows collection file not set properly")
	}

	setCollectionFile("linux")
	if collectionFile.filename != filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile) {
		t.Error("Runtime OS collection file not set properly")
	}

	os.Setenv("TOTP_FILE", "testcollectionfile.json")
	setCollectionFile("windows")
	if collectionFile.filename != os.Getenv("TOTP_FILE") {
		t.Error("Collection file not set properly with environment variable")
	}

	// Not sure how to unit test but at least run it for now
	_, _ = loadCollectionFromStdin()

	for _, c := range reservedCommands {
		if isReservedCommand(c) != true {
			t.Error("Error checking valid reserved commands")
		}
	}

	if isReservedCommand("validsecretname") == true {
		t.Error("Error checking invalid reserved command")
	}
}
