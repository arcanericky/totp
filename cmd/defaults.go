package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

const defaultBaseCollectionFile = "totp-config.json"

var defaultCollectionFile string
var reservedCommands = []string{configCmd.Use, versionCmd.Use}

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
		defaultCollectionFile = filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile)
	} else {
		defaultCollectionFile = filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile)
	}
}

func init() {
	setCollectionFile(runtime.GOOS)
}
