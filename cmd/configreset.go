package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func configReset() {
	if err := os.Remove(collectionFile.filename); err != nil {
		fmt.Println("Error removing collection file:", err)
		return
	}

	fmt.Println("Collection file removed")
}

func getConfigResetCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset the TOTP colllection",
		Long:  "Reset the TOTP colllection",
		Run: func(_ *cobra.Command, _ []string) {
			configReset()
		},
	}

	return cobraCmd
}
