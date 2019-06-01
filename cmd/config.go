package cmd

import (
	"github.com/spf13/cobra"
)

const configName = "config"

var configCmd = &cobra.Command{
	Use:   configName,
	Short: "Configure totp",
	Long:  `Configure totp`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
