package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the TOTP colllection",
	Long:  "Reset the TOTP colllection",
	Run: func(cmd *cobra.Command, args []string) {
		configReset()
	},
}

func configReset() {
	os.Remove(defaultCollectionFile)
}

func init() {
	configCmd.AddCommand(configResetCmd)
}
