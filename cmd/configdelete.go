package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "erase", "rm", "del"},
	Short:   "Delete a secret",
	Long:    `Delete a secret`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Must provide a secret name to delete.")
			return
		}

		deleteSecret(args[0])
	},
}

func deleteSecret(name string) {
	s, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading settings", err)
	} else {
		_, err := s.DeleteSecret(name)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting secret", err)
		} else {
			s.Save()
			printResultf("Deleted secret %s\n", name)
		}
	}
}

func init() {
	configCmd.AddCommand(configDeleteCmd)
	configDeleteCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
}
