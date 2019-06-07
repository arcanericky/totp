package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionText = "unspecified"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show topt version",
	Long:  "Show topt version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("totp version %s %s/%s\n", versionText, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
