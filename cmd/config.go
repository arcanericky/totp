package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure totp",
	Long:  `Configure totp`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
