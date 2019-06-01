package cmd

import (
	"fmt"
	"os"

	"github.com/arcanericky/totp"
	"github.com/spf13/cobra"
)

var configDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "erase", "rm", "del"},
	Short:   "Delete a key",
	Long:    `Delete a key`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Must provide a key name to delete.")
			return
		}

		deleteKey(args[0])
	},
}

func deleteKey(name string) {
	s, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading settings", err)
	} else {
		_, err := s.DeleteKey(name)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting key", err)
		} else {
			s.Save()
			fmt.Printf("Deleted key %s\n", name)
		}
	}
}

func init() {
	configCmd.AddCommand(configDeleteCmd)
}
