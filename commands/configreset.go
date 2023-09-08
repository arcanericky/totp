package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func configReset(filename string) error {
	if err := os.Remove(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing collection file %s: %s\n", filename, err)
		return err
	}

	fmt.Printf("Collection file %s removed\n", filename)
	return nil
}

func getConfigResetCmd() *cobra.Command {
	var (
		confirmAll bool
		cobraCmd   = &cobra.Command{
			Use:   "reset",
			Short: "Reset the TOTP colllection",
			Long:  "Reset the TOTP colllection",
			Run: func(_ *cobra.Command, _ []string) {
				if !confirmAll {
					confirm, err := userConfirm(bufio.NewReader(os.Stdin), "This will remove all secrets.")
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error getting response:", err)
						exitVal = 1
						return
					}

					if !confirm {
						fmt.Println("Skipping reset")
						return
					}
				}

				if err := configReset(collectionFile.filename); err != nil {
					exitVal = 1
				}
			},
		}
	)

	cobraCmd.Flags().BoolVarP(&confirmAll, optionYes, "y", false, "confirm all prompts")

	return cobraCmd
}
