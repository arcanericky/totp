package cmd

import (
	"runtime"
	"testing"
)

func TestDefaults(t *testing.T) {
	setCollectionFile("windows")
	setCollectionFile(runtime.GOOS)

	for _, c := range reservedCommands {
		if isReservedCommand(c) != true {
			t.Error("Error checking valid reserved commands")
		}
	}

	if isReservedCommand("validkeyname") == true {
		t.Error("Error checking invalid reserved command")
	}
}
