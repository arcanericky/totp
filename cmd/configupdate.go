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
	Short:   "Add or update a secret",
	Long:    `Add or update a secret`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Must provide name and secret")
			return
		}

		updateSecret(args[0], args[1])
	},
}

func updateSecret(name, value string) {
	if isReservedCommand(name) {
		fmt.Fprintln(os.Stderr, "The name \""+name+"\" is reserved for the "+name+" command")
		return
	}

	s, _ := totp.NewCollectionWithFile(defaultCollectionFile)
	secret, err := s.UpdateSecret(name, value)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error updating secret:", err)
	} else {
		s.Save()

		action := "Updated"
		if secret.DateAdded == secret.DateModified {
			action = "Added"
		}

		printResultf("%s secret %s\n", action, name)
	}
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
	configUpdateCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name] [secret value]", 1))
}
