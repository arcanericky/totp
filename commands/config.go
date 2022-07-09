package commands

import (
	"github.com/spf13/cobra"
)

func getConfigCmd(rootCmd *cobra.Command) *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:   cmdConfig,
		Short: "Configure totp",
		Long:  `Configure totp`,
	}

	cobraCmd.AddCommand(getConfigListCmd())
	cobraCmd.AddCommand(getConfigRenameCmd(rootCmd))
	cobraCmd.AddCommand(getConfigUpdateCmd(rootCmd))
	cobraCmd.AddCommand(getConfigDeleteCmd())
	cobraCmd.AddCommand(getConfigResetCmd())

	return cobraCmd
}
