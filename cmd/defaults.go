package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

var defaultCollectionFile string
var reservedCommands = []string{"config", "version"}

func isReservedCommand(name string) bool {
	for _, c := range reservedCommands {
		if name == c {
			return true
		}
	}

	return false
}

func setCollectionFile(goos string) {
	if goos == "windows" {
		defaultCollectionFile = filepath.Join(os.Getenv("LOCALAPPDATA"), "totp-config.json")
	} else {
		defaultCollectionFile = filepath.Join(os.Getenv("HOME"), ".totp-config.json")
	}
}

func init() {
	setCollectionFile(runtime.GOOS)
}
