package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arcanericky/totp"
	"github.com/spf13/cobra"
)

var configUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"add"},
	Short:   "Add or update a key",
	Long:    `Add or update a key`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Must provide name and seed.")
			return
		}

		updateKey(args[0], args[1])
	},
}

func updateKey(name, seed string) {
	if isReservedCommand(name) {
		fmt.Fprintln(os.Stderr, "The name \""+name+"\" is reserved for the "+name+" command")
		return
	}

	s, _ := totp.NewCollectionWithFile(defaultCollectionFile)
	key, err := s.UpdateKey(name, seed)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error updating key:", err)
	} else {
		s.Save()

		action := "Updated"
		if key.DateAdded == key.DateModified {
			action = "Added"
		}

		fmt.Printf("%s key %s\n", action, name)
	}
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
	configUpdateCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [seed name] [seed value]", 1))
}
