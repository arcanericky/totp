package commands

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/arcanericky/totp"
)

const (
	defaultBaseCollectionFile = "totp-config.json"

	cmdVersion    = "version"
	cmdConfig     = "config"
	cmdCompletion = "completion"
)

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
	if totpFile := os.Getenv("TOTP_CONFIG"); totpFile != "" {
		collectionFile.filename = totpFile
		return
	}

	if goos == "windows" {
		collectionFile.filename = filepath.Join(os.Getenv("LOCALAPPDATA"), defaultBaseCollectionFile)
		return
	}

	collectionFile.filename = filepath.Join(os.Getenv("HOME"), "."+defaultBaseCollectionFile)
}

var reservedCommands = []string{cmdConfig, cmdVersion, cmdCompletion}

func isReservedCommand(name string) bool {
	for _, c := range reservedCommands {
		if name == c {
			return true
		}
	}

	return false
}

func defaults() {
	setCollectionFile(runtime.GOOS)
	collectionFile.loader = loadCollectionFromDefaultFile
}
