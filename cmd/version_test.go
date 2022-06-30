package cmd

import (
	"testing"
)

func TestVersion(t *testing.T) {
	versionCmd := getVersionCmd()
	versionCmd.Run(nil, []string{})
}
