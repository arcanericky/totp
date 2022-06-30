package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func deleteSecret(name string) {
	s, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading settings:", err)
		return
	}

	if _, err := s.DeleteSecret(name); err != nil {
		fmt.Fprintln(os.Stderr, "Error deleting secret:", err)
		return
	}

	if err := s.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving settings:", err)
		return
	}

	if _, err := printResultf("Deleted secret %s\n", name); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func getConfigDeleteCmd() *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"remove", "erase", "rm", "del"},
		Short:   "Delete a secret",
		Long:    `Delete a secret`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "Must provide a secret name to delete.")
				return
			}

			deleteSecret(args[0])
		},
	}

	cobraCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
	return cobraCmd
}
