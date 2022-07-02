package commands

import (
	"testing"
)

func TestConfigCompletion(t *testing.T) {
	completionCmd := getConfigCompletionCmd(getRootCmd())
	completionCmd.Run(nil, []string{})
}
