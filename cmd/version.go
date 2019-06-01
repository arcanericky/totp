package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionText = "unspecified"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show topt version",
	Long:  "Show topt version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("totp version", versionText)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
