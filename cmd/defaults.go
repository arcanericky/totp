package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/arcanericky/totp"
)

const defaultBaseCollectionFile = "totp-config.json"

var collectionFile struct {
	filename string
	useStdio bool
	loader   func() (*totp.Collection, error)
}

func loadCollectionFromStdin() (*totp.Collection, error) {
	c, err := totp.NewCollectionWithReader(os.Stdin)
	c.SetWriter(os.Stdout)

	return c, err
}

func loadCollectionFromDefaultFile() (*totp.Collection, error) {
	return totp.NewCollectionWithFile(collectionFile.filename)
}

func setCollectionFile(goos string) {
	if goos == "windows" {
		collectionFile.filename = filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile)
	} else {
		collectionFile.filename = filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile)
	}
}

var reservedCommands = []string{configCmd.Use, versionCmd.Use}

func isReservedCommand(name string) bool {
	for _, c := range reservedCommands {
		if name == c {
			return true
		}
	}

	return false
}

func init() {
	setCollectionFile(runtime.GOOS)
	collectionFile.loader = loadCollectionFromDefaultFile
}
