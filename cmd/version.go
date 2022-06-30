package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionText = "unspecified"

func getVersionCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:   cmdVersion,
		Short: "Show totp version",
		Long:  "Show totp version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("totp version %s %s/%s\n", versionText, runtime.GOOS, runtime.GOARCH)
		},
	}

	return cobraCmd
}
