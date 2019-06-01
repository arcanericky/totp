package cmd

import (
	"testing"
)

func TestConfigCompletion(t *testing.T) {
	completionCmd.Run(nil, []string{})
}
